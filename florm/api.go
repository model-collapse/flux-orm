package florm

import (
	"log"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	http2 "github.com/influxdata/influxdb-client-go/v2/api/http"
)

type APIManager interface {
	QueryAPI(buckets []string) api.QueryAPI
	UpdateAPI(bucket string) api.QueryAPI
	WriteAPI(bucket string) api.WriteAPI
	DeleteAPI(bucket string) api.DeleteAPI
	Org() string

	// When NeedInstantFlush is true, any writeAPI operation (e.g. Insert)
	// will be flushed after WriteRecord is called. This is required when
	// LordAPIManager.WriteAPI instantly create the api when called and
	// maintain no cache, because the api will be deprecated after the insert call.
	NeedInstantFlush() bool
}

type LordAPIManager struct {
	client influxdb2.Client
	qAPI   api.QueryAPI
	dAPI   api.DeleteAPI
	org    string
}

func NewLordAPIManager(client influxdb2.Client, org string) *LordAPIManager {
	ret := &LordAPIManager{
		client: client,
		qAPI:   client.QueryAPI(org),
		dAPI:   client.DeleteAPI(),
		org:    org,
	}

	return ret
}

func (l *LordAPIManager) QueryAPI(buckets []string) api.QueryAPI {
	return l.qAPI
}

func (l *LordAPIManager) UpdateAPI(bucket string) api.QueryAPI {
	return l.qAPI
}

func (l *LordAPIManager) WriteAPI(bucket string) api.WriteAPI {
	ret := l.client.WriteAPI(l.org, bucket)
	ret.SetWriteFailedCallback(
		func(batch string, err http2.Error, retryAttempts uint) bool {
			log.Printf("Fail to write record, %v", err)
			return true
		},
	)
	return ret
}

func (l *LordAPIManager) NeedInstantFlush() bool {
	return true
}

func (l *LordAPIManager) DeleteAPI(bucket string) api.DeleteAPI {
	return l.dAPI
}

func (l *LordAPIManager) Org() string {
	return l.org
}
