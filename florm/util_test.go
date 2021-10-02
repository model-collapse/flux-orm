package florm

import (
	"reflect"
	"testing"

	"github.com/influxdata/influxdb-client-go/v2/api/query"
)

type Person struct {
	Model
	Name string `florm:"k,name"`
	Age  int    `florm:"v,age"`
	Sex  string `florm:"v,sex"`
}

type Student struct {
	Person
	Class string  `florm:"k,class"`
	Grade float32 `florm:"v,grade"`
}

func (s *Student) Measurement() string {
	return "student"
}

func (s *Student) Bucket() string {
	return "misc2"
}

func TestAssignRecordToStruct(t *testing.T) {
	rec := query.NewFluxRecord(0, map[string]interface{}{
		"name":  "agnes",
		"age":   5.0,
		"sex":   "male",
		"class": "xiao",
		"grade": 0.2,
	})

	s := &Student{}
	assignRecordToStruct(rec, reflect.ValueOf(s).Elem())

	t.Log(s)
}
