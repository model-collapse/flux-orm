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
}

type LordAPIManager struct {
	client influxdb2.Client
	qAPI   api.QueryAPI
	org    string
}

func NewLordAPIManager(client influxdb2.Client, org string) *LordAPIManager {
	ret := &LordAPIManager{
		client: client,
		qAPI:   client.QueryAPI(org),
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
