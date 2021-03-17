package client

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"

	"github.com/gopaytech/patroni_exporter/opts"
)

type PatroniClient interface {
	GetMetrics() (patroniJSONResp, error)
}

type patroniField struct {
	Scope   string `json:"scope"`
	Version string `json:"version"`
}

type patroniJSONResp struct {
	State   string       `json:"state"`
	Role    string       `json:"role"`
	Patroni patroniField `json:"patroni"`
}

type patroniClient struct {
	resty *resty.Client
}

func (p *patroniClient) GetMetrics() (patroniJSONResp, error) {
	resp, err := p.resty.R().Get("/patroni")
	if err != nil {
		return patroniJSONResp{}, err
	}
	if resp.IsError() {
		return patroniJSONResp{}, fmt.Errorf("got response status %d", resp.StatusCode())
	}

	var objmap patroniJSONResp
	err = json.Unmarshal(resp.Body(), &objmap)
	if err != nil {
		return patroniJSONResp{}, err
	}

	return objmap, nil
}

func NewPatroniClient(httpClient *resty.Client, opts opts.PatroniOpts) *patroniClient {
	httpClient.SetHostURL(fmt.Sprintf("%s:%s", opts.Host, opts.Port))
	return &patroniClient{
		resty: httpClient,
	}
}
