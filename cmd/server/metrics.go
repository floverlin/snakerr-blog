package server

import (
	"runtime"
	"runtime/metrics"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func Prom() {
	metrGorutines := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "blog",
			Name:      "gorutines",
			Help:      "количество горутин",
		},
	)
	metrMemory := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "blog",
			Name:      "memory",
			Help:      "используемая память",
		},
	)
	prometheus.MustRegister(metrGorutines)
	prometheus.MustRegister(metrMemory)
	nameGorutines := "/sched/goroutines:goroutines"
	nameMemory := "/memory/classes/heap/free:bytes"
	go func() {
		for {
			getMetric := make([]metrics.Sample, 2)
			getMetric[0].Name = nameGorutines
			getMetric[1].Name = nameMemory
			metrics.Read(getMetric)

			runtime.GC()

			metrGorutines.Set(float64(getMetric[0].Value.Uint64()))
			metrMemory.Set(float64(getMetric[1].Value.Uint64()))
			time.Sleep(5 * time.Second)
		}
	}()
}
