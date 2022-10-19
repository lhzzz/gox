package metric

import (
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/armon/go-metrics"
	prometheussink "github.com/armon/go-metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

type metricsFunc func(key []string, val float32, labels []metrics.Label)

// InitPrometheusMetrics 初始化Prometheus监控
func InitPrometheusMetrics(serviceName, metricAddr string) {

	prometheusOpts := prometheussink.PrometheusOpts{
		Expiration: 60 * time.Second,
	}

	if viper.GetBool("SCRAPE_SLOW") {
		prometheusOpts.Expiration = 10 * time.Minute
	}

	sink, _ := prometheussink.NewPrometheusSinkFrom(prometheusOpts)
	config := metrics.DefaultConfig(serviceName)
	config.EnableHostname = false
	config.EnableHostnameLabel = true
	config.EnableRuntimeMetrics = false
	metrics.NewGlobal(config, sink)

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Infof("Beginning to serve %s metrics on %s, expiration: %s", serviceName, metricAddr, prometheusOpts.Expiration)
		log.Fatal(http.ListenAndServe(metricAddr, nil))
	}()
}

// MapToMetricsLables map转换成Labels
func MapToMetricsLables(maps map[string]interface{}) []metrics.Label {

	labels := make([]metrics.Label, len(maps))
	index := 0
	for k, v := range maps {
		labels[index].Name = k
		sv, err := cast.ToStringE(v)
		if err != nil {
			labels[index].Value = fmt.Sprintf("%v", v)
		} else {
			labels[index].Value = sv
		}
		index++
	}

	sort.Slice(labels[:], func(i, j int) bool {
		return labels[i].Name < labels[j].Name
	})

	return labels
}

// PairsToMetricsLables 通过传入的kv返回一组标签
// len(kv)必须为偶数，传入参数格式为key, value交替出现，且key类型必须为string
func PairsToMetricsLables(kv ...interface{}) []metrics.Label {
	if len(kv)%2 == 1 {
		panic(fmt.Sprintf("metrics: Pairs got the odd number of input pairs for metrics: %d", len(kv)))
	}
	maps := map[string]interface{}{}
	var key string
	for i, s := range kv {
		if i%2 == 0 {
			key = s.(string)
			continue
		}
		maps[key] = s
	}
	return MapToMetricsLables(maps)
}
