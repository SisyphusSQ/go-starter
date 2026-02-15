package common_srv

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	promModel "github.com/prometheus/common/model"
	"github.com/spf13/cast"

	"go-starter/config"
	"go-starter/internal/lib/log"
)

type PrometheusService interface {
	QueryMemUsage(address string) (int, error)

	queryVector(promql string) (promModel.Vector, error)
}

type PrometheusServiceImpl struct {
	ctxTimeout time.Duration

	url    string
	client api.Client
	v1api  v1.API
}

func NewPrometheusService(config config.Config) (PrometheusService, error) {
	client, err := api.NewClient(api.Config{
		Address: config.Prometheus.URL,
	})
	if err != nil {
		log.Logger.Errorf("[PrometheusService.NewPrometheusService] url[%s] NewClient error: %v", config.Prometheus.URL, err)
		return nil, err
	}

	log.Logger.Infof("[PrometheusService.NewPrometheusService] url[%s] init success", config.Prometheus.URL)
	return &PrometheusServiceImpl{
		ctxTimeout: 5 * time.Second,
		url:        config.Prometheus.URL,
		client:     client,
		v1api:      v1.NewAPI(client),
	}, nil
}

func (p *PrometheusServiceImpl) QueryMemUsage(address string) (int, error) {
	log.Logger.Debugf("[PrometheusService.QueryMemUsage] address[%s] start", address)
	promqlFmt := `java_lang_Memory_HeapMemoryUsage_used{instance="%s"}/java_lang_Memory_HeapMemoryUsage_max{instance="%s"} * 100`
	promql := fmt.Sprintf(promqlFmt, address, address)

	vec, err := p.queryVector(promql)
	if err != nil {
		log.Logger.Errorf("[PrometheusService.QueryMemUsage] address[%s] queryVector error: %v", address, err)
		return 0, err
	}

	// vector must be one
	if len(vec) == 0 {
		log.Logger.Warnf("[PrometheusService.QueryMemUsage] address[%s] empty vector", address)
		return 0, errors.New("empty vector")
	}
	usage := cast.ToInt(vec[0].Value)
	log.Logger.Infof("[PrometheusService.QueryMemUsage] address[%s] usage[%d%%] success", address, usage)
	return usage, nil
}

func (p *PrometheusServiceImpl) queryVector(promql string) (promModel.Vector, error) {
	ctx, cancel := context.WithTimeout(context.Background(), p.ctxTimeout)
	defer cancel()

	res, _, err := p.v1api.Query(ctx, promql, time.Now())
	if err != nil {
		log.Logger.Errorf("[PrometheusService.queryVector] promql[%s] query error: %v", promql, err)
		return nil, err
	}

	vector, ok := res.(promModel.Vector)
	if !ok {
		log.Logger.Errorf("[PrometheusService.queryVector] result type[%v] not Vector", reflect.TypeOf(res))
		return nil, fmt.Errorf("query Vector Error: %v", reflect.TypeOf(res))
	}
	return vector, nil
}
