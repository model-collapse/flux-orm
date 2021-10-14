package florm

import (
	"context"
	"fmt"
	"strings"
	"time"
)

const keyReduceFilter = `(r, accumulator) => ({_value:r._value + "," + accumulator._value})`
const keyReduceIdentity = `{_value:""}`
const fieldReductFilter = `(r, accumulator) => ({_value:r._field + "," + accumulator._value})`

type schemaPerm struct {
	Perm string `florm:"k,_value"`
}

type SchemaHelper struct {
	ss *FluxSession
}

func NewSchemaHelper(ss *FluxSession) SchemaHelper {
	return SchemaHelper{ss: ss}
}

func (s *SchemaHelper) GetTagKeysPerm(bucket, measurement string, start, stop time.Time) (ret [][]string, reterr error) {
	var perms []schemaPerm
	fn := fmt.Sprintf(`(r) => (r._measurement == "%s")`, measurement)
	s.ss.From(bucket).Range(start, stop).Filter(fn, "drop").Keys().Reduce(keyReduceFilter, keyReduceIdentity).Keep([]string{"_value"}).Distinct().Yield(&perms)

	if reterr = s.ss.ExecuteQuery(context.Background()); reterr != nil {
		return
	}

	for _, p := range perms {
		tags := strings.Split(p.Perm, ",")
		ret = append(ret, tags)
	}

	return
}

func (s *SchemaHelper) GetTagKeysDistinct(bucket, measurement string, start, stop time.Time) (ret []string, reterr error) {
	var perms []schemaPerm
	fn := fmt.Sprintf(`(r) => (r._measurement == "%s")`, measurement)
	s.ss.From(bucket).Range(start, stop).Filter(fn, "drop").Keys().Keep([]string{"_value"}).Distinct().Yield(&perms)

	if reterr = s.ss.ExecuteQuery(context.Background()); reterr != nil {
		return
	}

	for _, p := range perms {
		ret = append(ret, p.Perm)
	}

	return
}

func (s *SchemaHelper) GetFieldsPerm(bucket, measurement string, start, stop time.Time) (ret [][]string, reterr error) {
	var perms []schemaPerm
	fn := fmt.Sprintf(`(r) => (r._measurement == "%s")`, measurement)
	s.ss.From(bucket).Range(start, stop).Filter(fn, "drop").Drop([]string{"_time", "_start", "_stop", "_value"}).GroupExcept([]string{"_field"}).Reduce(fieldReductFilter, keyReduceIdentity).Keep([]string{"_value"}).Distinct().Yield(&perms)

	if reterr = s.ss.ExecuteQuery(context.Background()); reterr != nil {
		return
	}

	for _, p := range perms {
		tags := strings.Split(p.Perm, ",")
		var ntg []string
		for _, t := range tags {
			if len(t) > 0 {
				ntg = append(ntg, t)
			}
		}

		ret = append(ret, ntg)
	}

	return
}

func (s *SchemaHelper) GetFieldsDistinct(bucket, measurement string, start, stop time.Time) (ret []string, reterr error) {
	var perms []schemaPerm
	fn := fmt.Sprintf(`(r) => (r._measurement == "%s")`, measurement)
	s.ss.From(bucket).Range(start, stop).Filter(fn, "drop").Keep([]string{"_field"}).DistinctCol("_field").Yield(&perms)

	if reterr = s.ss.ExecuteQuery(context.Background()); reterr != nil {
		return
	}

	for _, p := range perms {
		ret = append(ret, p.Perm)
	}

	return
}
