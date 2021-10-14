package florm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/query"
)

var ErrInvalidInput = errors.New("the input should be InfluxModel or a slice / channle of InfluxModel")
var ErrSeriesNotSupported = errors.New("delete by id can not be used in series data model, please use other delete api")
var ErrSeriesNotFound = errors.New("with IsSeries=true, we cannot find any presence of Series struct in the data model")
var ErrEmptySlice = errors.New("slice for insertion is empty")
var ErrNothingToUpdate = errors.New("nothing to update")

type YieldFunc func(r *query.FluxRecord)

func (ss *FluxSession) buildYields() string {
	buf := bytes.NewBuffer(nil)

	for i, out := range ss.outputs {
		name := out.stream.Name()
		fmt.Fprintf(buf, "%s\n|> yield(name:\"%d\")\n\n", name, i)
	}

	return buf.String()
}

const updatePivotClause = `|> pivot(columnKey: ["_field"], rowKey: ["_time"], valueColumn: "_value")`
const updateToClause = `|> to(bucket:"%s", tagColumns: %s, fieldFn: %s)`

func getUpdateFields(v interface{}, mp map[string]interface{}, replaceValue bool) (ret []fluxField, reterr error) {
	vv := reflect.ValueOf(v)
	if vv.Kind() == reflect.Ptr {
		vv = vv.Elem()
	}

	flds, err := recurseFluxFlieds(vv)
	if err != nil {
		return nil, err
	}

	if mp != nil {
		passAll := false
		if _, suc := mp["*"]; suc {
			passAll = true
		}

		for _, f := range flds {
			if v, suc := mp[f.name]; suc || passAll {
				nf := f

				if replaceValue {
					switch ff := v.(type) {
					case float32, float64:
						nf.floatValue = reflect.ValueOf(v).Float()
					case int8, int32, int64, int, uint, uint8, uint32, uint16:
						nf.intValue = reflect.ValueOf(v).Int()
					case string:
						nf.strValue = ff
					default:
						reterr = fmt.Errorf("invalid type for field %s", f.name)
						return
					}
				}

				if reterr != nil {
					return
				}

				ret = append(ret, nf)
			}
		}
	}

	return
}

func (ss *FluxSession) buildUpdateClause() (string, error) {
	replace := false
	if ss.update.src == nil && ss.model != nil {
		ss.update.src = ss.model
		replace = true
	} else if ss.update.src != nil {
		ss.update.vals = make(map[string]interface{})
		for _, s := range ss.update.sel {
			ss.update.vals[s] = nil
		}
	}

	name := ss.update.stream.Name()

	buf := bytes.NewBuffer(nil)
	fmt.Fprintf(buf, "%s\n", name)

	flds, err := getUpdateFields(ss.update.src, ss.update.vals, replace)
	if err != nil {
		return "", err
	}

	if len(flds) == 0 {
		return "", ErrNothingToUpdate
	}

	ss.imports = append(ss.imports, "experimental")

	var setParams []string
	for _, f := range flds {
		if f.tp == ValueField {
			p := fmt.Sprintf("%s:%s", f.name, f.ValueString(true))
			setParams = append(setParams, p)
		}
	}

	fmt.Fprintf(buf, "|> experimental.set(o:{%s})\n", strings.Join(setParams, ","))

	tagCols, valCols, err := extractTagAndValueCols(ss.update.src)
	if err != nil {
		return "", err
	}

	tgnc, _ := json.Marshal(tagCols)

	var fnEles []string
	for _, vv := range valCols {
		fnEles = append(fnEles, fmt.Sprintf("\"%s\":r.%s", vv, vv))
	}

	bucket := ss.update.src.(InfluxModel).Bucket()
	tagFn := fmt.Sprintf("(r) => ({%s})", strings.Join(fnEles, ", "))
	fmt.Fprintf(buf, updateToClause, bucket, string(tgnc), tagFn)
	fmt.Fprintf(buf, "\n")

	return buf.String(), nil
}

func (ss *FluxSession) buildUpdates() (string, error) {
	buf := bytes.NewBuffer(nil)

	if up := ss.update; up != nil {
		cl, err := ss.buildUpdateClause()
		if err != nil {
			return "", err
		}

		fmt.Fprintf(buf, "%s\n", cl)
	}

	return buf.String(), nil
}

func (ss *FluxSession) buildYieldProgram() (ret map[string]YieldFunc) {
	ret = make(map[string]YieldFunc)

	for i, o := range ss.outputs {
		outInterface := o.output
		v := reflect.ValueOf(outInterface)
		if v.Kind() != reflect.Ptr {
			panic(errors.New("the yield receiver should be a ptr"))
		}

		var fnc YieldFunc
		v = v.Elem()
		if v.Kind() == reflect.Slice {
			etp := v.Type().Elem()
			eleIsPtr := etp.Kind() == reflect.Ptr
			if eleIsPtr {
				etp = etp.Elem()
			}

			fnc = func(r *query.FluxRecord) {
				n := reflect.New(etp)
				assignRecordToStruct(r, n.Elem())

				if eleIsPtr {
					v = reflect.Append(v, n)
				} else {
					v = reflect.Append(v, n.Elem())
				}

				reflect.ValueOf(outInterface).Elem().Set(v)
			}
		} else if v.Kind() == reflect.Struct {
			etp := v.Type()
			fnc = func(r *query.FluxRecord) {
				vv := reflect.New(etp)
				assignRecordToStruct(r, vv.Elem())
				reflect.ValueOf(outInterface).Set(vv)
			}
		}

		ret[fmt.Sprintf("%d", i)] = fnc
	}

	return
}

// Will block util query finish
func (ss *FluxSession) ExecuteQuery(ctx context.Context) error {
	script := ss.QueryString()

	if len(ss.outputs) > 0 {
		script = fmt.Sprintf("%s\n\n%s", script, ss.buildYields())
	}

	if ss.update != nil {
		q, err := ss.buildUpdates()
		if err != nil {
			return err
		}

		script += "\n\n" + q
	}

	if len(ss.imports) > 0 {
		script = ss.buildImports() + script
	}

	if ss.dbg {
		return nil
	}

	log.Printf(script)
	var qApi api.QueryAPI

	if ss.update != nil {
		qApi = ss.mgr.UpdateAPI(ss.buckets, ss.update.src.(InfluxModel).Bucket())
	} else {
		qApi = ss.mgr.QueryAPI(ss.buckets)
	}

	result, err := qApi.Query(ctx, script)
	if err != nil {
		return err
	}

	if len(ss.outputs) == 0 {
		return nil
	}

	fncs := ss.buildYieldProgram()
	var fn YieldFunc
	yid := ""
	for result.Next() {
		if result.TableChanged() {
			yid = getYieldIDFromMeta(result.TableMetadata())
			fn = fncs[yid]
		}

		fn(result.Record())
	}

	return nil
}

func (ss *FluxSession) Insert(v interface{}) error {
	var target []InfluxModel
	var bucket string
	modelStub := reflect.TypeOf((*InfluxModel)(nil)).Elem()

	userChan := false
	var chanIn reflect.Value

	if m, suc := v.(InfluxModel); suc {
		target = []InfluxModel{m}
		bucket = target[0].Bucket()
	} else {
		vv := reflect.ValueOf(v)
		if vv.Kind() == reflect.Ptr {
			vv = vv.Elem()
		}

		if vv.Kind() == reflect.Slice {
			target = make([]InfluxModel, 0, vv.Len())
			for i := 0; i < vv.Len(); i++ {
				vvv := vv.Index(i)
				if mm, suc := vvv.Interface().(InfluxModel); suc {
					target = append(target, mm)
				} else {
					return ErrInvalidInput
				}
			}

			if vv.Len() == 0 {
				return ErrEmptySlice
			} else {
				bucket = target[0].Bucket()
			}
		} else if vv.Kind() == reflect.Chan {
			if !vv.Type().Elem().ConvertibleTo(modelStub) {
				return ErrInvalidInput
			}

			userChan = true
			chanIn = vv
		} else {
			return ErrInvalidInput
		}
	}

	var wApi api.WriteAPI
	if bucket != "" {
		wApi = ss.mgr.WriteAPI(bucket)
	}

	// when instant flush needed, the api will be flushed after record written
	if ss.mgr.NeedInstantFlush() {
		defer wApi.Flush()
	}

	if userChan {
		m, suc := chanIn.Recv()
		for suc {
			mm := m.Interface().(InfluxModel)
			if mm == nil {
				continue
			}

			if bucket == "" {
				bucket = mm.Bucket()
				wApi = ss.mgr.WriteAPI(bucket)
			}

			script, err := modelToInsertString(mm)
			log.Printf(script)
			if err != nil {
				return err
			}

			wApi.WriteRecord(script)
			m, suc = chanIn.Recv()
		}
	} else {
		for _, m := range target {
			script, err := modelToInsertString(m)
			log.Printf(script)
			if err != nil {
				return err
			}

			wApi.WriteRecord(script)
		}
	}

	return nil
}

func (ss *FluxSession) Delete(ctx context.Context, m InfluxModel, ids ...uint64) error {
	dapi := ss.mgr.DeleteAPI(m.Bucket())

	if m.IsSeries() {
		return ErrSeriesNotSupported
	}

	dapi.DeleteWithName(ctx, ss.mgr.Org(), m.Bucket(), time.Unix(0, 0), time.Unix(0, 1), idSetToFilter(m.Measurement(), ids))
	return nil
}

func (ss *FluxSession) DeleteWithFilter(ctx context.Context, m InfluxModel, fn string) error {
	dapi := ss.mgr.DeleteAPI(m.Bucket())

	if m.IsSeries() {
		srs, suc := findSeriesObj(reflect.ValueOf(m), reflect.TypeOf(m))
		if !suc {
			return ErrSeriesNotFound
		}

		dapi.DeleteWithName(ctx, ss.mgr.Org(), m.Bucket(), srs.Start, srs.Stop, fn)
	} else {
		dapi.DeleteWithName(ctx, ss.mgr.Org(), m.Bucket(), time.Unix(0, 0), time.Unix(0, 1), fn)
	}

	return nil
}
