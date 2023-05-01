package prometheus_exporter

import "github.com/prometheus/client_golang/prometheus"

type LimitsCollector struct {
	LimitTotal     *prometheus.Desc
	LimitRemaining *prometheus.Desc
	LimitUsed      *prometheus.Desc
	SecondsLeft    *prometheus.Desc
}
