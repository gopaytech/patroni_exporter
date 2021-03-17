package client_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"

	"github.com/gopaytech/patroni_exporter/client"
	"github.com/gopaytech/patroni_exporter/opts"
)

const PATRONI_HOST = "http://localhost"
const PATRONI_PORT = "8008"
const PATRONI_RESPONSE = `
{
  "state": "running",
  "postmaster_start_time": "2021-02-10 09:07:37.860 UTC",
  "role": "master",
  "server_version": 110010,
  "cluster_unlocked": false,
  "xlog": {
    "location": 1247250352
  },
  "timeline": 23,
  "replication": [
    {
      "usename": "repluser",
      "application_name": "cluster-02",
      "client_addr": "10.1.1.2",
      "state": "streaming",
      "sync_state": "sync",
      "sync_priority": 1
    }
  ],
  "database_system_identifier": "3922806929873881258",
  "patroni": {
    "version": "2.0.1",
    "scope": "cluster"
  }
}
`

type clientTestContext struct {
	httpClient    *resty.Client
	client        client.PatroniClient
	getMetricsUrl string
}

func (context *clientTestContext) setUp(t *testing.T) {
	context.httpClient = resty.New()
	httpmock.ActivateNonDefault(context.httpClient.GetClient())
	options := opts.PatroniOpts{
		Host: PATRONI_HOST,
		Port: PATRONI_PORT,
	}
	context.getMetricsUrl = fmt.Sprintf("%s:%s/patroni", PATRONI_HOST, PATRONI_PORT)
	context.client = client.NewPatroniClient(context.httpClient, options)
}

func (context *clientTestContext) tearDown() {
	httpmock.DeactivateAndReset()
}

func TestPatroniClient_GetMetricsSuccess(t *testing.T) {
	context := clientTestContext{}
	context.setUp(t)
	defer context.tearDown()

	httpmock.RegisterResponder(
		"GET",
		context.getMetricsUrl,
		func(request *http.Request) (response *http.Response, err error) {
			return httpmock.NewStringResponse(200, PATRONI_RESPONSE), nil
		},
	)

	resp, err := context.client.GetMetrics()
	assert.NoError(t, err)
	assert.Equal(t, "running", resp.State)
	assert.Equal(t, "master", resp.Role)
	assert.Equal(t, "cluster", resp.Patroni.Scope)
	assert.Equal(t, "2.0.1", resp.Patroni.Version)
}

func TestPatroniClient_GetMetricsFailed(t *testing.T) {
	context := clientTestContext{}
	context.setUp(t)
	defer context.tearDown()

	httpmock.RegisterResponder(
		"GET",
		context.getMetricsUrl,
		func(request *http.Request) (response *http.Response, err error) {
			return httpmock.NewStringResponse(503, PATRONI_RESPONSE), nil
		},
	)

	resp, err := context.client.GetMetrics()
	assert.Error(t, err)
	assert.Empty(t, resp.State)
	assert.Empty(t, resp.Role)
	assert.Empty(t, resp.Patroni.Scope)
	assert.Empty(t, resp.Patroni.Version)
}
