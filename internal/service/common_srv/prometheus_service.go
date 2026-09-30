package common_srv

import (
	"context"
	"fmt"
	"github.com/SisyphusSQ/go-starter/v2/config"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	promModel "github.com/prometheus/common/model"
	"net/http"
	"time"
)

// PrometheusService 只封装查询客户端，业务指标和 PromQL 由调用方定义。
type PrometheusService interface {
	Query(context.Context, string) (promModel.Value, v1.Warnings, error)
}
type PrometheusServiceImpl struct {
	timeout time.Duration
	api     v1.API
}

func NewPrometheusService(c config.Config) (PrometheusService, error) {
	transport := http.DefaultTransport
	if c.Prometheus.Token != "" {
		transport = bearerTransport{base: transport, token: c.Prometheus.Token}
	}
	client, err := api.NewClient(api.Config{Address: c.Prometheus.URL, Client: &http.Client{Transport: transport, Timeout: c.Prometheus.Timeout, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}})
	if err != nil {
		return nil, fmt.Errorf("configure prometheus client: %w", err)
	}
	return &PrometheusServiceImpl{timeout: c.Prometheus.Timeout, api: v1.NewAPI(client)}, nil
}
func (p *PrometheusServiceImpl) Query(ctx context.Context, query string) (promModel.Value, v1.Warnings, error) {
	queryCtx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()
	return p.api.Query(queryCtx, query, time.Now())
}

type bearerTransport struct {
	base  http.RoundTripper
	token string
}

func (t bearerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	clone := req.Clone(req.Context())
	clone.Header.Set("Authorization", "Bearer "+t.token)
	return t.base.RoundTrip(clone)
}
