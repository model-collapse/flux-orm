package florm

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	sfk "github.com/godruoyi/go-snowflake"
)

type InfluxModel interface {
	Measurement() string
	Bucket() string
	IsSeries() bool
}

func isInfluxModel(v interface{}) bool {
	_, suc := v.(InfluxModel)
	return suc
}

type Series struct {
	Time  time.Time `florm:"t,_time"`
	Start time.Time `florm:"t,_start"`
	Stop  time.Time `florm:"t,_stop"`
}

func (m *Series) Measurement() string {
	panic(errors.New("Cannot use raw model struct in IO"))
}

func (m *Series) Buckets() string {
	panic(errors.New("Cannot use raw model struct in IO"))
}

func (m *Series) IsSeries() bool {
	return true
}

func (m *Series) SetTime(t time.Time) {
	m.Time = t
}

func findSeriesObj(val reflect.Value, tp reflect.Type) (ret Series, suc bool) {
	for i := 0; i < tp.NumField(); i++ {
		fld := val.Field(i)
		tfld := tp.Field(i)

		if tfld.Anonymous && tfld.Type.Kind() == reflect.Struct {
			if tfld.Type == reflect.TypeOf(Series{}) {
				ret = fld.Interface().(Series)
				suc = true
			} else {
				ret, suc = findSeriesObj(fld, tfld.Type)
				if suc {
					return
				}
			}
		}
	}

	return
}

type STable struct {
	ID uint64 `florm:"k,primary"`
}

func (m *STable) Measurement() string {
	panic(errors.New("Cannot use raw model struct in IO"))
}

func (m *STable) Bucket() string {
	panic(errors.New("Cannot use raw model struct in IO"))
}

func (m *STable) IsSeries() bool {
	return false
}

func (m *STable) Snowflake() uint64 {
	m.ID = sfk.ID()

	return m.ID
}

func findSTableObj(val reflect.Value, tp reflect.Type) (ret STable, suc bool) {
	for i := 0; i < tp.NumField(); i++ {
		fld := val.Field(i)
		tfld := tp.Field(i)

		if tfld.Anonymous && tfld.Type.Kind() == reflect.Struct {
			if tfld.Type == reflect.TypeOf(STable{}) {
				ret = fld.Interface().(STable)
				suc = true
			} else {
				ret, suc = findSTableObj(fld, tfld.Type)
				if suc {
					return
				}
			}
		}
	}

	return
}

type DTable struct {
	Series
	ID uint64 `florm:"k,primary"`
}

func (m *DTable) Snowflake() uint64 {
	m.ID = sfk.ID()
	sid := sfk.ParseID(m.ID)
	m.Time = sid.GenerateTime()

	return m.ID
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

	var timeStamp int64
	md, suc := findSeriesObj(val, tp)
	if !suc {
		if _, suc := findSTableObj(val, tp); suc {
			timeStamp = 0
		} else {
			reterr = errors.New("input should be a composed with a model")
			return
		}
	} else {
		timeStamp = md.Time.UnixNano()
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
			keys = append(keys, fmt.Sprintf("%s=%s", f.name, f.ValueString(false)))
		} else if f.tp == ValueField {
			vals = append(vals, fmt.Sprintf("%s=%s", f.name, f.ValueString(true)))
		}
	}

	measurement := m.Measurement()
	if timeStamp < 0 {
		timeStamp = 0
	}

	ret = fmt.Sprintf("%s,%s %s %d", measurement, strings.Join(keys, ","), strings.Join(vals, ","), timeStamp)
	return
}
