package florm

import (
	"errors"
	"reflect"
	"time"
)

var schemaCheckOffset = -time.Hour * 24 * 30

func checkYieldReceiver(v interface{}) error {
	vv := reflect.ValueOf(v)
	if vv.Kind() != reflect.Ptr {
		return errors.New("yield receiver should be a pointer")
	}

	tp := vv.Type()
	tpEle := tp.Elem()

	if tpEle.Kind() == reflect.Slice {
		tpEE := tpEle.Elem()
		if tpEE.Kind() == reflect.Ptr {
			tpEE = tpEE.Elem()
		}

		if tpEE.Kind() != reflect.Struct {
			return errors.New("if receiver should be a pointer to struct or a slice of structs")
		}
	} else if tpEle.Kind() != reflect.Struct {
		return errors.New("yield receiver should be a pointer to struct or slice")
	}

	return nil
}

func checkSchema(m InfluxModel, mgr APIManager) error {
	/*if mgr == nil {
		mgr = defaultAPIManager
	}

	start := time.Unix(0, 0)
	end := time.Unix(0, 2)

	if m.IsSeries() {
		start := time.Now().Add(schemaCheckOffset)
		end := time.Now()
	}

	ss := NewFluxSessionCustomAPI(mgr)
	ss.From(m.Bucket()).Range(start, end).Filter(fmt.Sprintf(`(r) => (r._measurement=%s)`, m.Measurement()), "drop").Keys()*/

	return nil
}
