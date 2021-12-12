package prometheus_exporter

import (
	"log"
	"net/http"
	"time"

	"github.com/kalgurn/github-rate-limits-prometheus-exporter/internal/github_client"
	"github.com/kalgurn/github-rate-limits-prometheus-exporter/internal/utils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	github_app = utils.GetOSVar("GITHUB_APP_NAME")
)

func newLimitsCollector() *LimitsCollector {
	return &LimitsCollector{
		LimitTotal: prometheus.NewDesc(prometheus.BuildFQName(github_app, "", "limit_total"),
			"Total limit of requests for the installation",
			nil, nil),
		LimitRemaining: prometheus.NewDesc(prometheus.BuildFQName(github_app, "", "limit_remaining"),
			"Amount of remaining requests for the installation",
			nil, nil),
		LimitUsed: prometheus.NewDesc(prometheus.BuildFQName(github_app, "", "limit_used"),
			"Amount of used requests for the installation",
			nil, nil),
	}
}

func (collector *LimitsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.LimitTotal
	ch <- collector.LimitRemaining
	ch <- collector.LimitUsed
}

func (collector *LimitsCollector) Collect(ch chan<- prometheus.Metric) {

	auth := github_client.InitConfig()
	limits := github_client.GetRemainingLimits(auth.InitClient())
	log.Printf("Collected metrics for %s", github_app)
	log.Printf("Limit: %d | Used: %d | Remaining: %d", limits.Limit, limits.Used, limits.Remaining)
	//Write latest value for each metric in the prometheus metric channel.
	//Note that you can pass CounterValue, GaugeValue, or UntypedValue types here.
	m1 := prometheus.MustNewConstMetric(collector.LimitTotal, prometheus.GaugeValue, float64(limits.Limit))
	m2 := prometheus.MustNewConstMetric(collector.LimitRemaining, prometheus.GaugeValue, float64(limits.Remaining))
	m3 := prometheus.MustNewConstMetric(collector.LimitUsed, prometheus.GaugeValue, float64(limits.Used))
	m1 = prometheus.NewMetricWithTimestamp(time.Now(), m1)
	m2 = prometheus.NewMetricWithTimestamp(time.Now(), m2)
	m3 = prometheus.NewMetricWithTimestamp(time.Now(), m3)
	ch <- m1
	ch <- m2
	ch <- m3
}

func Run() {
	limit := newLimitsCollector()
	prometheus.MustRegister(limit)

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
