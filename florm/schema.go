package florm

import (
	"context"
	"fmt"
	"strings"
	"time"
)

const keyReduceFilter = `(r, accumulator) => ({_value:r._value + "," + accumulator._value})`
const keyReduceIdentity = `{_value:""}`
const fieldReductFilter = `reduce(fn: (r, accumulator) => ({_value:r._field + "," + accumulator._value}), identity: {_value:""})`

type SchemaPerm struct {
	perm string `florm:"k,_value"`
}

type SchemaHelper struct {
	ss *FluxSession
}

func (s *SchemaHelper) GetTagKeysPerm(bucket, measurement string, start, stop time.Time) (ret [][]string, reterr error) {
	var perms []SchemaPerm
	fn := fmt.Sprintf(`(r) => (r._measurement == "%s")`, measurement)
	s.ss.From(bucket).Range(start, stop).Filter(fn, "drop").Keys().Reduce(keyReduceFilter, keyReduceIdentity).Keep([]string{"_value"}).Distinct().Yield(&perms)

	if reterr = s.ss.ExecuteQuery(context.Background()); reterr != nil {
		return
	}

	for _, p := range perms {
		tags := strings.Split(p.perm, ",")
		ret = append(ret, tags)
	}

	return
}

func (s *SchemaHelper) GetTagKeysDistinct(bucket, measurement string, start, stop time.Time) (ret []string, reterr error) {
	var perms []SchemaPerm
	fn := fmt.Sprintf(`(r) => (r._measurement == "%s")`, measurement)
	s.ss.From(bucket).Range(start, stop).Filter(fn, "drop").Keys().Keep([]string{"_value"}).Distinct().Yield(&perms)

	if reterr = s.ss.ExecuteQuery(context.Background()); reterr != nil {
		return
	}

	for _, p := range perms {
		ret = append(ret, p.perm)
	}

	return
}

func (s *SchemaHelper) GetFieldsPerm(bucket, measurement string, start, stop time.Time) (ret [][]string, reterr error) {
	var perms []SchemaPerm
	fn := fmt.Sprintf(`(r) => (r._measurement == "%s")`, measurement)
	s.ss.From(bucket).Range(start, stop).Filter(fn, "drop").Drop([]string{"_time", "_start", "_stop", "_value"}).GroupExcept([]string{"field"}).Reduce(fieldReductFilter, keyReduceIdentity).Keep([]string{"_value"}).Distinct().Yield(&perms)

	if reterr = s.ss.ExecuteQuery(context.Background()); reterr != nil {
		return
	}

	for _, p := range perms {
		tags := strings.Split(p.perm, ",")
		ret = append(ret, tags)
	}

	return
}

func (s *SchemaHelper) GetFieldsDistinct(bucket, measurement string, start, stop time.Time) (ret []string, reterr error) {
	var perms []SchemaPerm
	fn := fmt.Sprintf(`(r) => (r._measurement == "%s")`, measurement)
	s.ss.From(bucket).Range(start, stop).Filter(fn, "drop").Keep([]string{"_field"}).Distinct().Yield(&perms)

	if reterr = s.ss.ExecuteQuery(context.Background()); reterr != nil {
		return
	}

	for _, p := range perms {
		ret = append(ret, p.perm)
	}

	return
}
