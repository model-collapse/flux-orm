package florm

import (
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

	ss := NewFluxSession()
	ss.From("misc2")
}

func TestUpdate(t *testing.T) {

}
