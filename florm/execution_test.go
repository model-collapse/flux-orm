package florm

import (
	"context"
	"encoding/json"
	"sync"
	"testing"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

var initialAPIOnce sync.Once

func initializeAPI() {
	client := influxdb2.NewClient("https://us-west-2-1.aws.cloud2.influxdata.com", "hDHUGj34tYfx1HEwVIqiW_uZhMUxSiKpFLRCr__csNY_1QgVuKmPi8zrqsO-QEtzA-Rc8WsAKH_cZSCZS1nhaA==")
	mgr := NewLordAPIManager(client, "charlie.yang.nju@gmail.com")

	RegisterDefaultAPIManager(mgr)
}

func TestStaticGet(t *testing.T) {
	initialAPIOnce.Do(initializeAPI)

	var res []Student
	ss := NewFluxSession()
	ss.From("misc2").Static().Pivot().Yield(&res)
	if err := ss.ExecuteQuery(context.Background()); err != nil {
		t.Fatal(err)
	}

	if len(res) == 0 {
		t.Fatalf("zero result received")
	}

	d, _ := json.Marshal(&res)
	t.Log(string(d))
}

func TestUpdate(t *testing.T) {

}

func TestInsert(t *testing.T) {
	initialAPIOnce.Do(initializeAPI)

	ss := NewFluxSession()

	st := &Student{
		Person: Person{
			Table: Table{
				ID: 1,
			},
			Name: "sean",
			Age:  4,
			Sex:  "male",
		},
		Class: "5",
		Grade: 0.93,
	}

	err := ss.Insert(st)

	if err != nil {
		t.Fatal(err)
	}
}
