package florm

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"

	"github.com/influxdata/influxdb-client-go/v2/api/query"
)

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

func parseParam(p string) (k, v string) {
	eles := strings.Split(p, ":")
	k, v = eles[0], eles[1]
	return
}

func fillStrToValue(s string, v reflect.Value) {
	return
}

// tags
const (
	KeyField = iota
	ValueField
	TimeField
)

type fieldTag struct {
	tp   int
	name string
}

func praseTag(tag string) (ret fieldTag, err error) {
	eles := strings.Split(tag, ",")
	if len(eles) != 2 {
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
	return
}

type fluxField struct {
	fieldTag
	floatValue float64
	strValue   string
}

func recurseFluxFlieds(v reflect.Value) (ret []fluxField, reterr error) {
	if v.Kind() != reflect.Struct {
		reterr = errors.New("invalid value kind, need to be struct!")
		//panic(reterr)
		return
	}

	tp := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fld := v.Field(i)
		tf := tp.Field(i)
		if tf.Anonymous && fld.Kind() == reflect.Struct {
			fds, err := recurseFluxFlieds(fld)
			if err != nil {
				reterr = err
				return
			}

			ret = append(ret, fds...)
		}

		if !tf.Anonymous {
			ft, err := praseTag(string(tf.Tag.Get("florm")))
			if err != nil {
				reterr = err
				return
			}

			f := fluxField{fieldTag: ft}
			switch fld.Kind() {
			case reflect.Float32, reflect.Float64:
				f.floatValue = fld.Float()
			case reflect.Int, reflect.Int32, reflect.Int64:
				f.floatValue = float64(fld.Int())
			case reflect.String:
				f.strValue = fld.String()
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

func assignRecordToStruct(r *query.FluxRecord, v reflect.Value) {
	if v.Kind() != reflect.Struct {
		return
	}

	tp := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fld := tp.Field(i)
		if fld.Anonymous && fld.Type.Kind() == reflect.Struct {
			assignRecordToStruct(r, v.Field(i))
			continue
		}

		tag := fld.Tag.Get("florm")
		log.Printf("tag is %s", tag)

		if tag == "" {
			continue
		}

		tg, err := praseTag(tag)
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
					log.Printf("WTF?")
				}

				if vf, suc := val.(float64); suc {
					if v.Field(i).Kind() == reflect.Int32 ||
						v.Field(i).Kind() == reflect.Int64 ||
						v.Field(i).Kind() == reflect.Int {
						v.Field(i).SetInt(int64(vf))
					} else {
						v.Field(i).SetFloat(vf)
					}
				} else if vs, suc := val.(string); suc {
					v.Field(i).SetString(vs)
				}
			}()
		}
	}
}
