package client

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/gopaytech/patroni_exporter/opts"
	options "github.com/gopaytech/patroni_exporter/opts"
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
	options.PatroniOpts
}

func (p *patroniClient) GetMetrics() (patroniJSONResp, error) {
	resp, err := p.resty.R().EnableTrace().Get("/patroni")
	if err != nil {
		return patroniJSONResp{}, err
	}

	var objmap patroniJSONResp

	err = json.Unmarshal(resp.Body(), &objmap)
	if err != nil {
		return patroniJSONResp{}, err
	}

	fmt.Println(objmap)

	return objmap, nil
}

func NewPatroniClient(opts opts.PatroniOpts) PatroniClient {
	r := resty.New()
	r.SetHostURL(fmt.Sprintf("%s:%s", opts.Host, opts.Host))
	return &patroniClient{
		resty: r,
	}
}
