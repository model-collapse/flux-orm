package florm

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
)

type InfluxModel interface {
	Measurement() string
	Bucket() string
}

func isInfluxModel(v interface{}) bool {
	_, suc := v.(InfluxModel)
	return suc
}

type Model struct {
	Time  time.Time `florm:"t,time"`
	Start time.Time `florm:"t,start"`
	Stop  time.Time `florm:"t,stop"`
}

func (m *Model) Measurement() string {
	panic(errors.New("Cannot use raw model struct in IO"))
}

func (m *Model) Buckets() string {
	panic(errors.New("Cannot use raw model struct in IO"))
}

func (m *Model) SetTime(t time.Time) {
	m.Time = t
}

func getModelObj(val reflect.Value, tp reflect.Type) (ret Model, suc bool) {
	for i := 0; i < tp.NumField(); i++ {
		fld := val.Field(i)
		tfld := tp.Field(i)

		if tfld.Anonymous && tfld.Type.Kind() == reflect.Struct {
			if tfld.Type == reflect.TypeOf(Model{}) {
				ret = fld.Interface().(Model)
				suc = true
			} else {
				ret, suc = getModelObj(fld, tfld.Type)
				if suc {
					return
				}
			}
		}
	}

	return
}

func modelToInsertString(m InfluxModel) (ret string, reterr error) {
	val := reflect.ValueOf(m)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	tp := val.Type()

	if tp.Kind() != reflect.Struct {
		reterr = errors.New("the input type should be a struct or a pointer to struct")
		return
	}

	md, suc := getModelObj(val, tp)
	if !suc {
		reterr = errors.New("input should be a composed with a model")
		return
	}

	flds, err := recurseFluxFlieds(val)
	if err != nil {
		reterr = err
		return
	}

	var keys []string
	var vals []string
	for _, f := range flds {
		if f.tp == KeyField {
			keys = append(keys, fmt.Sprintf("%s=%s", f.name, f.strValue))
		} else if f.tp == ValueField {
			if f.strValue == "" {
				vals = append(vals, fmt.Sprintf("%s=%f", f.name, f.floatValue))
			} else {
				vals = append(vals, fmt.Sprintf("%s=\"%s\"", f.name, f.strValue))
			}
		}
	}

	measurement := m.Measurement()
	timeStamp := md.Time.UnixNano()
	if timeStamp < 0 {
		timeStamp = 0
	}

	ret = fmt.Sprintf("%s,%s %s %d", measurement, strings.Join(keys, ","), strings.Join(vals, ","), timeStamp)
	return
}
