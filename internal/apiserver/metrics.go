package apiserver

import "github.com/prometheus/client_golang/prometheus"

type serverMetrics struct {
	requestErrors   prometheus.Counter
	requestCounter  prometheus.Counter
	requestDuration prometheus.Histogram
}

func (s *APIServer) registerMetrics() {
	s.metrics.requestErrors = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "blog",
			Name:      "request_error_total",
			Help:      "количество ошибок",
		},
	)
	s.metrics.requestCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "blog",
			Name:      "request_total",
			Help:      "запросов в минуту",
		},
	)
	s.metrics.requestDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: "blog",
			Name:      "requst_duration_second",
			Help:      "длительность запросов",
			Buckets:   prometheus.DefBuckets,
		},
	)
	prometheus.MustRegister(s.metrics.requestErrors)
	prometheus.MustRegister(s.metrics.requestCounter)
	prometheus.MustRegister(s.metrics.requestDuration)
}
