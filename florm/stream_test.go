package florm

import (
	"context"
	"testing"
	"time"
)

type TOPT struct {
	F string `florm:"k,f"`
	S string `florm:"v,g"`
}

func (t *TOPT) Measurement() string {
	return "m"
}

func (t *TOPT) Bucket() string {
	return "b"
}

func TestFluxVars(t *testing.T) {
	ss := NewFluxSession()
	a := ss.From("testBucket").
		Range(time.Date(2020, 10, 1, 12, 0, 0, 0, time.Local), time.Date(2020, 11, 1, 12, 0, 0, 0, time.Local)).
		Filter("(r) => (r.s > 0)", "drop").
		Window("1d").
		StateCount("(r) => (r.s > 1)", "cnt")

	b := ss.From("testBucket2").
		Range(time.Date(2020, 10, 1, 12, 0, 0, 0, time.Local), time.Date(2020, 11, 1, 12, 0, 0, 0, time.Local)).
		Map("(r) => ({r with share = 100})").
		Window("1d").
		Mean().
		Difference(false, []string{"c1", "c2"}, true)

	c := a.Join(b, []string{"c3", "c4"}, "a_", "b_", "inner").
		Filter("(r) => (r.c4 > 0)", "drop").
		Sum()

	r := &TOPT{}
	c.Yield(r)
	t.Log(ss.QueryString())
	ss.ExecuteQuery(context.Background())
}

func TestFluxVars2(t *testing.T) {
	ss := NewFluxSession()
	a := ss.From("testBucket").
		Range(time.Date(2020, 10, 1, 12, 0, 0, 0, time.Local), time.Date(2020, 11, 1, 12, 0, 0, 0, time.Local)).
		Filter("(r) => (r.s > 0)", "drop").
		Window("1d")

	b := a.StateCount("(r) => (r.s > 1)", "cnt")

	c := a.
		Filter("(r) => (r.c4 > 0)", "drop").
		Sum()

	r := &TOPT{}
	c.Yield(r)
	b.Yield(r)
	t.Log(ss.QueryString())
}

func TestFluxVarsUpdate(t *testing.T) {
	ss := NewFluxSession()

	r := &TOPT{}
	ss.From("testBucket").Static().Filter("(r) => (r.c2 > 0)", "drop").Update(r)
	if err := ss.ExecuteQuery(context.Background()); err != nil {
		t.Fatal(err)
	}

	t.Log(ss.QueryString())
}
