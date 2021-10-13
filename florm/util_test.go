package florm

import (
	"fmt"
	"math"
	"reflect"
	"testing"

	"github.com/influxdata/influxdb-client-go/v2/api/query"
)

//--------------------Format Convert---------------------
func TestStrDeDedup(t *testing.T) {
	cases := [][]string{
		{"a", "a", "g", "k"},
		{"b", "c", "b", "c"},
		{"a", "b", "c"},
		{},
	}

	groundTruth := [][]string{
		{"a", "g", "k"},
		{"b", "c"},
		{"a", "b", "c"},
		{},
	}

	for i, c := range cases {
		r := strDedup(c)
		if !cmpStrArray(r, groundTruth[i]) {
			t.Errorf("inconsistent on case @%d", i)
		}
	}

	return
}

func TestFKVs(t *testing.T) {
	if fkvS("a", "b") != "a:b" {
		t.Errorf("[bad] fkvS")
	}

	if fkvSq("a", "b") != "a:\"b\"" {
		t.Errorf("[bad] fkvSq")
	}

	if fkvF("a", 0.6) != fmt.Sprintf("a:%f", 0.6) {
		t.Errorf("[bad] fkvF")
	}

	if fkvI("a", 1) != "a:1" {
		t.Errorf("[bad] fkvI")
	}

	if fkvB("a", true) != "a:true" {
		t.Errorf("[bad] fkvB (%s)", fkvB("a", true))
	}

	if fkvB("a", false) != "a:false" {
		t.Errorf("[bad] fkvB (%s)", fkvB("a", false))
	}

	if v := fkvSJ("a", []string{"a", "b", "c"}); v != "a:[\"a\",\"b\",\"c\"]" {
		t.Errorf("[bad] fkvJ %v", v)
	}
}

func TestInterfaceToFloat(t *testing.T) {
	goodCases := []interface{}{int(1), int32(1), int64(1), float32(1.0), float64(1.0), uint8(1)}
	badCases := []interface{}{"haha", []byte{10, 3, 4}, t}

	for i, g := range goodCases {
		r, err := interfaceToFloat(g)
		if err != nil {
			t.Errorf("incorrect good case @%d", i)
		}

		if math.Abs(r-1.0) > 0.00001 {
			t.Errorf("inaccurate good case @%d", i)
		}
	}

	for i, g := range badCases {
		_, err := interfaceToFloat(g)
		if err == nil {
			t.Errorf("incorrect bad case @%d", i)
		}
	}
}

func TestParseParam(t *testing.T) {
	c1 := "s1:s2"
	k, v, err := parseParam(c1)
	if err != nil {
		t.Fatal(err)
	}

	if k != "s1" {
		t.Error("error k")
	}

	if v != "s2" {
		t.Error("error v")
	}

	c2 := "s1:s3:s4"
	_, _, err = parseParam(c2)
	if err == nil {
		t.Fatal("nil error for incorrect param")
	}
}

//--------------------Field Process----------------------

type Person struct {
	STable
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

func TestParseTag(t *testing.T) {
	t1 := "s"
	_, err := parseTag(t1)
	if err == nil {
		t.Fatal("nil err for illegal tag")
	}

	t2 := "ss,f,v"
	_, err = parseTag(t2)
	if err == nil {
		t.Fatal("nil err for illegal tag")
	}

	t3 := "k,value"
	f, err := parseTag(t3)
	if err != nil {
		t.Fatal(err)
	}

	if f.tp != KeyField {
		t.Errorf("type is wrong")
	}

	if f.name != "value" {
		t.Errorf("name is wrong")
	}

	t4 := "v,value"
	f, err = parseTag(t4)
	if err != nil {
		t.Fatal(err)
	}

	if f.tp != ValueField {
		t.Errorf("type is wrong")
	}

	if f.name != "value" {
		t.Errorf("name is wrong")
	}
}

func assert(t *testing.T, name string, a, b interface{}) {
	equal := reflect.DeepEqual(a, b)

	if !equal {
		switch av := a.(type) {
		case float64:
			bv := b.(float64)
			if math.Abs(float64(av-bv)) < 0.000001 {
				equal = true
			}

		case float32:
			bv := b.(float32)
			if math.Abs(float64(av-bv)) < 0.000001 {
				equal = true
			}
		}
	}

	if !equal {
		t.Errorf("Field %s is inconsistent [%v / %v]", name, a, b)
	}
}

func TestRecurseFluxFlieds(t *testing.T) {
	s := &Student{
		Person: Person{
			STable: STable{ID: 19},
			Name:   "agnes",
			Age:    10,
			Sex:    "female",
		},
		Grade: 0.9,
		Class: "star",
	}

	flds, err := recurseFluxFlieds(reflect.ValueOf(s))
	if err != nil {
		t.Fatal(err)
	}

	for _, fld := range flds {
		switch fld.name {
		case "name":
			assert(t, "name", fld.strValue, "agnes")
			assert(t, "name[type]", KeyField, fld.tp)
		case "age":
			assert(t, "age", fld.intValue, int64(10))
			assert(t, "age[type]", ValueField, fld.tp)
		case "grade":
			assert(t, "grade", fld.floatValue, float64(0.9))
			assert(t, "grade[type]", ValueField, fld.tp)
		case "primary":
			assert(t, "primary", fld.uintValue, uint64(19))
			assert(t, "primary[type]", KeyField, fld.tp)
		}
	}
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

	if s.Name != "agnes" {
		t.Error("incorrect name")
	}

	if s.Sex != "male" {
		t.Error("incorrect sex")
	}

	if s.Class != "xiao" {
		t.Error("incorrect class")
	}

	if s.Grade != 0.2 {
		t.Error("incorrect grade")
	}

	if s.Age != 5 {
		t.Error("incorrect age")
	}
}
