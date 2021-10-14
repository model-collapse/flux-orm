package florm

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

var testStartTime = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)

var scriptDTableWithoutID = []string{
	`carTypeD,brand=GM,branch=1 tp="gas",tpid=1 %d`,
	`carTypeD,brand=GM,branch=2 tp="ev",tpid=2 %d`,
	`carTypeD,brand=BMW,branch=1 tp="gas",tpid=1 %d`,
	`carTypeD,brand=BMW,branch=2 tp="ev",tpid=2 %d`,
	`carTypeD,brand=BMW,branch=3 tp="ev",tpid=2 %d`,
}

var scriptSTableWithoutID = []string{
	`carTypeSS,brand=GM,branch=1 tp="gas",tpid=1 0`,
	`carTypeSS,brand=GM,branch=2 tp="ev",tpid=2 0`,
	`carTypeSS,brand=BMW,branch=1 tp="gas",tpid=1 0`,
	`carTypeSS,brand=BMW,branch=2 tp="ev",tpid=2 0`,
	`carTypeSS,brand=BMW,branch=3 tp="ev",tpid=2 0`,
}

var scriptDTableWithID = []string{
	`carTypeD,primary=%d,brand=GM,branch=1 tp="gas",tpid=1 %d`,
	`carTypeD,primary=%d,brand=GM,branch=2 tp="ev",tpid=2 %d`,
	`carTypeD,primary=%d,brand=BMW,branch=1 tp="gas",tpid=1 %d`,
	`carTypeD,primary=%d,brand=BMW,branch=2 tp="ev",tpid=2 %d`,
	`carTypeD,primary=%d,brand=BMW tp="ev",tpid=2 %d`,
	`carTypeD,primary=%d,brand=BMW,branch=4,ditto=h tp="ev",tpid=2 %d`,
}

var scriptSTableWithID = []string{
	`carTypeS,primary=1,brand=GM,branch=1 tp="gas",tpid=1 0`,
	`carTypeS,primary=2,brand=GM,branch=2 tp="ev",tpid=2 0`,
	`carTypeS,primary=3,brand=BMW,branch=1 tp="gas",tpid=1 0`,
	`carTypeS,primary=4,brand=BMW,branch=2 tp="ev",tpid=2 0`,
	`carTypeS,primary=5,brand=BMW,branch=3 tp="ev",tpid=2 0`,
}

func refineScripts() {
	for i := 0; i < 5; i++ {
		t := testStartTime.Add(time.Second * time.Duration(i))
		ts := t.UnixNano()
		scriptDTableWithoutID[i] = fmt.Sprintf(scriptDTableWithoutID[i], ts)

	}

	for i := 0; i < 6; i++ {
		t := testStartTime.Add(time.Second * time.Duration(i))
		ts := t.UnixNano()
		pid := SID{Sequence: uint64(i), MachineID: 0}
		pid.SetTime(t)
		scriptDTableWithID[i] = fmt.Sprintf(scriptDTableWithID[i], pid.CalculateID(), ts)
	}
}

var fillMockDataM3Once sync.Once
var refineScriptsOnce sync.Once

func fillMockDataM3() {
	initialAPIOnce.Do(initializeAPI)
	refineScriptsOnce.Do(refineScripts)

	dapi := defaultAPIManager.DeleteAPI("misc3")
	if err := dapi.DeleteWithName(context.Background(), testOrg, "misc3", time.Unix(0, 0), time.Now(), ""); err != nil {
		panic(fmt.Errorf("cannot clean the unit test bucket, %v", err))
	}

	wapi := defaultAPIManager.WriteAPI("misc3")
	for _, scripts := range [][]string{scriptDTableWithoutID, scriptSTableWithoutID, scriptDTableWithID, scriptSTableWithID} {
		for _, s := range scripts {
			log.Println(s)
			wapi.WriteRecord(s)
		}
	}

	wapi.Flush()
}

type CarTypeS struct {
	STable
	Brand  string  `florm:"k,brand"`
	Branch int     `florm:"k,branch"`
	Type   string  `florm:"v,tp"`
	Tid    float64 `florm:"v,tpid"`
}

func (s *CarTypeS) Bucket() string {
	return "misc3"
}

func (s *CarTypeS) Measurement() string {
	return "carTypeS"
}

type CarTypeSS struct {
	STable
	Brand  string  `florm:"k,brand"`
	Branch int     `florm:"k,branch"`
	Type   string  `florm:"v,tp"`
	Tid    float64 `florm:"v,tpid"`
}

func (s *CarTypeSS) Bucket() string {
	return "misc3"
}

func (s *CarTypeSS) Measurement() string {
	return "carTypeSS"
}

type CarTypeD struct {
	DTable
	Brand  string  `florm:"k,brand"`
	Branch int     `florm:"k,branch"`
	Type   string  `florm:"v,tp"`
	Tid    float64 `florm:"v,tpid"`
}

func (s *CarTypeD) Bucket() string {
	return "misc3"
}

func (s *CarTypeD) Measurement() string {
	return "carTypeD"
}
