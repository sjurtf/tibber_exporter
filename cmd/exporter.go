package cmd

import (
	"github.com/prometheus/client_golang/prometheus"
	promcollectors "github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sjurtf/tibber_exporter/tibber"
	"net/http"
)

type Exporter struct {
	Tibber   *tibber.Tibber
	Gatherer *prometheus.Registry
}

func NewExporter(tibber *tibber.Tibber) *Exporter {
	return &Exporter{
		Tibber:   tibber,
		Gatherer: prometheus.NewRegistry(),
	}
}

func (e *Exporter) StartGather() {
	e.Gatherer.MustRegister(promcollectors.NewGoCollector())
	e.Gatherer.MustRegister(NewTibberCollector(e.Tibber))
}

func (e *Exporter) Listen() {
	metricsPath := "/metrics"

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
<html>
<head><title>Tibber Exporter</title></head>
<body>
<h1>Tibber Exporter</h1>
<p><a href=` + metricsPath + `>Metrics</a></p>
</body>
</html>`))
	})

	http.Handle(metricsPath, promhttp.HandlerFor(
		e.Gatherer,
		promhttp.HandlerOpts{
			EnableOpenMetrics: true,
		},
	))
	_ = http.ListenAndServe(":8080", nil)
}
