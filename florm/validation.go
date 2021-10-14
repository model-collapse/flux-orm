package florm

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"time"
)

const (
	// In this check mode, all historical permutation of tags/fields should be
	// consistent and equal with the corresponding model in the code.
	CheckModeStrict = iota

	// In this check mode, all historical permutation of tags/fields should be
	// a subset of corresponding model in the code.
	// [Warning] in this mode, all tag/field inconsistency need be manually controlled.
	CheckModeCompatible

	// In this mode, we only check whether Table model have primary key in the DB.
	// All data inconsitency should be managed manually.
	CheckModeMinimum
)

var schemaCheckOffset = -time.Hour * 24 * 30
var ErrNoSchemaData = errors.New("no schema data found in the time period")

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

func CheckSchema(m InfluxModel, mgr APIManager, mode int) error {
	start := time.Unix(0, 0)
	end := time.Unix(0, 1)

	if m.IsSeries() {
		start = time.Now().Add(schemaCheckOffset)
		end = time.Now()
	}

	return checkSchema(m, mgr, mode, start, end)
}

func checkSchema(m InfluxModel, mgr APIManager, mode int, start, stop time.Time) error {
	if mgr == nil {
		mgr = defaultAPIManager
	}

	ss := NewFluxSessionCustomAPI(mgr)
	shp := NewSchemaHelper(ss)

	tagPerms, err := shp.GetTagKeysPerm(m.Bucket(), m.Measurement(), start, stop)
	if err != nil {
		return err
	}

	if len(tagPerms) == 0 {
		return ErrNoSchemaData
	}

	for i, t := range tagPerms {
		tagPerms[i] = strDedup(append(t, "_time"))
	}

	tagCols, valCols, err := extractTagAndValueCols(m)
	if err != nil {
		return err
	}

	tagCols = strDedup(append(tagCols, "_field", "_measurement", "_time", "_start", "_stop"))
	if mode == CheckModeStrict {
		if len(tagPerms) > 1 {
			return fmt.Errorf("[strict] historical tags permutations are not consistent, versions are %v", tagPerms)
		}

		if !cmpStrArray(tagCols, tagPerms[0]) {
			return fmt.Errorf("[strict] historical tags permutations != model tags, [%v/%v]", tagCols, tagPerms[0])
		}
	} else if mode == CheckModeCompatible {
		for _, p := range tagPerms {
			if !isSubSet(p, tagCols) {
				return fmt.Errorf("[compatible] historical tags permutations is not compatible with model, [%v / %v]", p, tagCols)
			}
		}
	} else if mode == CheckModeMinimum {
		for _, p := range tagPerms {
			if strInSlice("primary", tagCols) {
				if !strInSlice("primary", p) {
					return fmt.Errorf("[minimum] no primary in table model")
				}
			}
		}
	}

	ss = NewFluxSessionCustomAPI(mgr)
	shp = NewSchemaHelper(ss)

	fieldPerms, err := shp.GetFieldsPerm(m.Bucket(), m.Measurement(), start, stop)
	if err != nil {
		return err
	}

	for _, t := range fieldPerms {
		sort.Strings(t)
	}

	if mode == CheckModeStrict {
		if len(fieldPerms) > 1 {
			return fmt.Errorf("[strict] historical field permutations are not consistent, versions are: %v", fieldPerms)
		}

		if !cmpStrArray(valCols, fieldPerms[0]) {
			return fmt.Errorf("[strict] historical field permutations != model field, [%v/%v]", valCols, fieldPerms[0])
		}
	} else if mode == CheckModeCompatible {
		for _, p := range fieldPerms {
			if !isSubSet(p, valCols) {
				return fmt.Errorf("[compatible] historical field permutations is not compatible with model, [%v / %v]", p, valCols)
			}
		}
	}

	return nil
}
