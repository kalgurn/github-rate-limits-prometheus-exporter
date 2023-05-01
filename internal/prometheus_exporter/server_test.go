package prometheus_exporter

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type FakeCollector struct {
	LimitTotal     *prometheus.Desc
	LimitRemaining *prometheus.Desc
	LimitUsed      *prometheus.Desc
	SecondsLeft    *prometheus.Desc
}

func newFakeCollector() *LimitsCollector {
	return &LimitsCollector{
		LimitTotal: prometheus.NewDesc(prometheus.BuildFQName(githubAccount, "", "limit_total"),
			"Total limit of requests for the installation",
			nil, nil),
		LimitRemaining: prometheus.NewDesc(prometheus.BuildFQName(githubAccount, "", "limit_remaining"),
			"Amount of remaining requests for the installation",
			nil, nil),
		LimitUsed: prometheus.NewDesc(prometheus.BuildFQName(githubAccount, "", "limit_used"),
			"Amount of used requests for the installation",
			nil, nil),
		SecondsLeft: prometheus.NewDesc(prometheus.BuildFQName(githubAccount, "", "seconds_left"),
			"Time left in seconds until limit is reset for the installation",
			nil, nil),
	}
}

func (collector *FakeCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.LimitTotal
	ch <- collector.LimitRemaining
	ch <- collector.LimitUsed
	ch <- collector.SecondsLeft
}

func (collector *FakeCollector) Collect(ch chan<- prometheus.Metric) {

	log.Printf("Collecting metrics for %s", githubAccount)
	//Write latest value for each metric in the prometheus metric channel.
	//Note that you can pass CounterValue, GaugeValue, or UntypedValue types here.
	m1 := prometheus.MustNewConstMetric(collector.LimitTotal, prometheus.GaugeValue, float64(10))
	m2 := prometheus.MustNewConstMetric(collector.LimitRemaining, prometheus.GaugeValue, float64(6))
	m3 := prometheus.MustNewConstMetric(collector.LimitUsed, prometheus.GaugeValue, float64(4))
	m4 := prometheus.MustNewConstMetric(collector.SecondsLeft, prometheus.GaugeValue, time.Duration(time.Second*30).Seconds())
	m1 = prometheus.NewMetricWithTimestamp(time.Now().Add(-time.Hour), m1)
	m2 = prometheus.NewMetricWithTimestamp(time.Now(), m2)
	m3 = prometheus.NewMetricWithTimestamp(time.Now(), m3)
	m4 = prometheus.NewMetricWithTimestamp(time.Now(), m4)
	ch <- m1
	ch <- m2
	ch <- m3
	ch <- m4
}

func TestNewLimitsCollector(t *testing.T) {
	newCollector := newFakeCollector()
	prometheus.MustRegister(newCollector)

	mux := http.NewServeMux()

	mux.Handle("/limits", promhttp.Handler())

	ts := httptest.NewServer(mux)
	defer ts.Close()

	resp, err := http.Get("0.0.0.0:2112/limits")
	if err != nil {
		log.Print(err)
	}
	fmt.Println(resp)

}
