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

	"github.com/influxdata/influxdb-client-go/v2/api/query"
)

const YieldBufSize = 100

type YieldFunc func(r *query.FluxRecord)

func (ss *FluxSession) buildYields() string {
	buf := bytes.NewBuffer(make([]byte, YieldBufSize))

	for i, out := range ss.outputs {
		name := out.stream.Name()
		fmt.Fprintf(buf, "%s\n|> yield(name:\"%d\")\n\n", name, i)
	}

	return buf.String()
}

const updatePivotClause = `|> pivot(columnKey: ["_field"], rowKey: ["_time"], valueColumn: "_value")`
const updateToClause = `|> to(bucket:%s, tagColumns: %s, fieldFn: %s)`

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
					case float32, float64, int32, int64, int:
						nf.floatValue, reterr = interfaceToFloat(v)
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

	buf := bytes.NewBuffer(make([]byte, YieldBufSize))
	fmt.Fprintf(buf, "%s\n", name)
	fmt.Fprintf(buf, updatePivotClause)
	fmt.Fprintf(buf, "\n")

	flds, err := getUpdateFields(ss.update.src, ss.update.vals, replace)
	if err != nil {
		return "", err
	}

	var tagCols []string
	var valCols []string
	for _, f := range flds {
		if f.tp == ValueField {
			if f.strValue == "" {
				fmt.Fprintf(buf, "|> set(\"%s\", %f)\n", f.name, f.floatValue)
			} else {
				fmt.Fprintf(buf, "|> set(\"%s\", \"%s\")\n", f.name, f.strValue)
			}
			valCols = append(valCols, f.name)
		} else if f.tp == KeyField {
			tagCols = append(tagCols, f.name)
		}
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
	buf := bytes.NewBuffer(make([]byte, YieldBufSize))

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
		v := reflect.ValueOf(o.output)
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
			}
		} else if v.Kind() == reflect.Struct {
			fnc = func(r *query.FluxRecord) {
				assignRecordToStruct(r, v)
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
		script += "\n\n" + ss.buildYields()
	}

	if ss.update != nil {
		q, err := ss.buildUpdates()
		if err != nil {
			return err
		}

		script += "\n\n" + q
	}

	log.Printf(script)

	if ss.dbg {
		return nil
	}

	qApi := ss.mgr.QueryAPI(ss.buckets)
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

func (ss *FluxSession) Insert(m InfluxModel) error {
	script, err := modelToInsertString(m)
	log.Printf(script)
	if err != nil {
		return err
	}

	api := ss.mgr.WriteAPI(m.Bucket())
	defer api.Flush()
	api.WriteRecord(script)
	return nil
}

var defaultAPIManager APIManager

func RegisterDefaultAPIManager(m APIManager) {
	defaultAPIManager = m
}
