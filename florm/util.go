package florm

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"

	"github.com/influxdata/influxdb-client-go/v2/api/query"
)

var ErrUnsupportedType = errors.New("unsupported data type in data model")
var ErrParamFormat = errors.New("invalid parameter format")
var ErrAssignTypeInconsistent = errors.New("type is inconsistent between value to assign and to be assigned")

func strDedup(s []string) []string {
	ret := make([]string, 0, len(s))
	sort.Strings(s)
	past := ""
	for _, ss := range s {
		if ss != past {
			ret = append(ret, ss)
		}

		past = ss
	}

	return ret
}

func cmpStrArray(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i, aa := range a {
		bb := b[i]
		if aa != bb {
			return false
		}
	}

	return true
}

func strInSlice(s string, b []string) bool {
	for _, ss := range b {
		if ss == s {
			return true
		}
	}

	return false
}

func isSubSet(sub, sup []string) bool {
	for _, s := range sub {
		if !strInSlice(s, sup) {
			return false
		}
	}

	return true
}

func fkvS(k, v string) string {
	return fmt.Sprintf("%s:%s", k, v)
}

func fkvSq(k, v string) string {
	return fmt.Sprintf("%s:\"%s\"", k, v)
}

func fkvF(k string, v float64) string {
	return fmt.Sprintf("%s:%f", k, v)
}

func fkvI(k string, v int) string {
	return fmt.Sprintf("%s:%d", k, v)
}

func fkvB(k string, b bool) string {
	s := "true"
	if !b {
		s = "false"
	}

	return fmt.Sprintf("%s:%s", k, s)
}

func fkvSJ(k string, s interface{}) string {
	ss, _ := json.Marshal(s)

	return fmt.Sprintf("%s:%s", k, ss)
}

func interfaceToFloat(v interface{}) (ret float64, reterr error) {
	switch f := v.(type) {
	case int:
		ret = float64(f)
	case int8:
		ret = float64(f)
	case int16:
		ret = float64(f)
	case uint8:
		ret = float64(f)
	case uint16:
		ret = float64(f)
	case uint32:
		ret = float64(f)
	case uint64:
		ret = float64(f)
	case int32:
		ret = float64(f)
	case int64:
		ret = float64(f)
	case float32:
		ret = float64(f)
	case float64:
		ret = f
	default:
		reterr = fmt.Errorf("invalid type to be converted float")
	}

	return
}

func parseParam(p string) (k, v string, reterr error) {
	eles := strings.Split(p, ":")
	if len(eles) != 2 {
		reterr = ErrParamFormat
	}

	k, v = eles[0], eles[1]
	return
}

// tags
const (
	UnknownField = iota
	KeyField
	ValueField
	TimeField
)

const (
	FieldTypeUnknown = iota
	FieldTypeString
	FieldTypeInt
	FieldTypeUint
	FieldTypeFloat
)

type fieldTag struct {
	tp        int
	name      string
	omitEmpty bool
}

func parseTag(tag string) (ret fieldTag, err error) {
	eles := strings.Split(tag, ",")
	if len(eles) < 2 {
		err = errors.New("invalid tag format, we need 2 fields splittedd by comma (,)")
		return
	}

	if k := eles[0]; k == "k" || k == "key" {
		ret.tp = KeyField
	} else if k == "v" || k == "value" {
		ret.tp = ValueField
	} else if k == "t" || k == "time" {
		ret.tp = TimeField
	} else {
		err = errors.New("invalid tag format, prefix need to be \"k(key)\", \"v(value)\", \"t(time)\"")
		return
	}

	ret.name = eles[1]
	if len(eles) > 2 {
		if eles[2] == "omitempty" {
			ret.omitEmpty = true
		}
	}

	return
}

type fluxField struct {
	fieldTag
	valType    int
	floatValue float64
	uintValue  uint64
	intValue   int64
	strValue   string
}

func (f *fluxField) ValueString(withQuote bool) (ret string) {
	switch f.valType {
	case FieldTypeString:
		if withQuote {
			ret = fmt.Sprintf("\"%s\"", f.strValue)
		} else {
			ret = f.strValue
		}
	case FieldTypeFloat:
		ret = fmt.Sprintf("%f", f.floatValue)
	case FieldTypeInt:
		ret = fmt.Sprintf("%d", f.intValue)
	case FieldTypeUint:
		ret = fmt.Sprintf("%d", f.uintValue)
	}

	return
}

func (f *fluxField) CheckAndAssign(v interface{}) (ret string, reterr error) {
	var nf fluxField
	fillFieldValue(&nf, reflect.ValueOf(v))
	if nf.valType != f.valType {
		reterr = ErrAssignTypeInconsistent
		return
	}

	nf.name = f.name
	nf.fieldTag = f.fieldTag
	nf.omitEmpty = f.omitEmpty

	*f = nf
	return
}

func recurseFluxFlieds(v reflect.Value) (ret []fluxField, reterr error) {
	defer func() {
		if e := recover(); e != nil {
			reterr = e.(error)
		}
	}()

	return recurseFluxFliedsImpl(v)
}

func fillFieldValue(f *fluxField, v reflect.Value) (reterr error) {
	switch v.Kind() {
	case reflect.Float32, reflect.Float64:
		f.floatValue = v.Float()
		f.valType = FieldTypeFloat
	case reflect.Int, reflect.Int32, reflect.Int8, reflect.Int64:
		f.intValue = v.Int()
		f.valType = FieldTypeInt
	case reflect.Uint64, reflect.Uint8, reflect.Uint16:
		f.uintValue = v.Uint()
		f.valType = FieldTypeUint
	case reflect.Bool:
		if v.Bool() {
			f.intValue = 1
		} else {
			f.intValue = 0
		}
		f.valType = FieldTypeInt
	case reflect.String:
		f.strValue = v.String()
		f.valType = FieldTypeString
	default:
		reterr = ErrUnsupportedType
	}

	return
}

func recurseFluxFliedsImpl(v reflect.Value) (ret []fluxField, reterr error) {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		reterr = errors.New("invalid value kind, need to be struct!")
		return
	}

	tp := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fld := v.Field(i)
		tf := tp.Field(i)
		if tf.Anonymous && fld.Kind() == reflect.Struct {
			fds, err := recurseFluxFliedsImpl(fld)
			if err != nil {
				reterr = err
				return
			}

			ret = append(ret, fds...)
		}

		if !tf.Anonymous {
			ft, err := parseTag(string(tf.Tag.Get("florm")))
			if err != nil {
				reterr = err
				return
			}

			f := fluxField{fieldTag: ft}
			if fld.Kind() == reflect.Struct {
				f.strValue, reterr = EncodeMeta(fld.Interface())
				f.tp = FieldTypeString
			} else if fld.Kind() == reflect.Ptr {
				if fld.Elem().Kind() == reflect.Struct {
					f.strValue, err = EncodeMeta(fld.Interface())
					f.tp = FieldTypeString
				} else {
					err = ErrUnsupportedType
				}
			} else {
				err = fillFieldValue(&f, fld)
			}

			if err != nil {
				reterr = err
				return
			}

			ret = append(ret, f)
		}
	}

	return
}

func getYieldIDFromMeta(m *query.FluxTableMetadata) string {
	cols := m.Columns()

	for _, c := range cols {
		if c.Name() == "result" {
			return c.DefaultValue()
		}
	}

	return ""
}

func assignRecordToStruct(r *query.FluxRecord, v reflect.Value) (ret error) {
	defer func() {
		if e := recover(); e != nil {
			ret = e.(error)
		}
	}()

	assignRecordToStructImpl(r, v)
	return
}

func assignRecordToStructImpl(r *query.FluxRecord, v reflect.Value) {
	if v.Kind() != reflect.Struct {
		return
	}

	tp := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fld := tp.Field(i)
		if fld.Anonymous && fld.Type.Kind() == reflect.Struct {
			assignRecordToStructImpl(r, v.Field(i))
			continue
		}

		tag := fld.Tag.Get("florm")

		if tag == "" {
			continue
		}

		tg, err := parseTag(tag)
		if err != nil {
			panic(err)
		}

		val := r.ValueByKey(tg.name)
		if val != nil {
			func() {
				defer func() {
					if e := recover(); e != nil {
						log.Printf("inconsistent type of field assign: [%s / %s], e = %v", v.Field(i).Kind().String(), reflect.TypeOf(val).Kind().String(), e)
					}
				}()

				if !v.Field(i).CanSet() {
					panic(fmt.Errorf("field cannot be set: %s | %s | %s", tg.name, fld.Name, fld.Tag))
				}

				if vf, suc := val.(float64); suc {
					switch v.Field(i).Kind() {
					case reflect.Int32, reflect.Int64, reflect.Int, reflect.Int16, reflect.Int8:
						v.Field(i).SetInt(int64(vf))
					case reflect.Uint32, reflect.Uint64, reflect.Uint, reflect.Uint16, reflect.Uint8:
						v.Field(i).SetUint(uint64(vf))
					case reflect.Float32, reflect.Float64:
						v.Field(i).SetFloat(vf)
					case reflect.Bool:
						v.Field(i).SetBool(vf != 0)
					}
				} else if vs, suc := val.(string); suc {
					if v.Field(i).Kind() == reflect.Struct {
						if err := DecodeMeta(vs, v.Field(i).Interface()); err != nil {
							panic(err)
						}
					} else if v.Field(i).Kind() == reflect.Ptr && v.Field(i).Elem().Kind() == reflect.Struct {
						if err := DecodeMeta(vs, v.Field(i).Interface()); err != nil {
							panic(err)
						}
					} else if v.Field(i).Kind() == reflect.String {
						v.Field(i).SetString(vs)
					}
				}
			}()
		}
	}
}

func idSetToFilter(measurement string, ids []uint64) string {
	buf := bytes.NewBuffer(nil)
	for i, id := range ids {
		fmt.Fprintf(buf, "r.key==\"%d\"", id)
		if i < len(ids)-1 {
			fmt.Fprintf(buf, " || ")
		}
	}

	return fmt.Sprintf("(r) => (r._measurement == \"%s\" && (%s))", measurement, buf)
}

func EncodeMeta(v interface{}) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", nil
	}

	return base64.RawStdEncoding.EncodeToString(data), nil
}

func DecodeMeta(src string, v interface{}) error {
	data, err := base64.RawStdEncoding.DecodeString(src)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, v)
}

func extractTagAndValueCols(v interface{}) (tagCols, valCols []string, reterr error) {
	flds, err := recurseFluxFlieds(reflect.ValueOf(v))
	if err != nil {
		reterr = err
		return
	}

	for _, f := range flds {
		if f.tp == ValueField {
			valCols = append(valCols, f.name)
		} else if f.tp == KeyField {
			tagCols = append(tagCols, f.name)
		}
	}

	return
}
