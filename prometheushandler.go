package prometheushandler

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Count of all HTTP requests",
		},
		[]string{"code", "method", "handler"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "A histogram of latencies for requests.",
			Buckets: []float64{.25, .5, 1, 2.5, 5, 10},
		},
		[]string{"code", "method", "handler"},
	)

	httpRequestSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_size_bytes",
			Help:    "A histogram of response sizes for requests.",
			Buckets: []float64{200, 500, 900, 1500},
		},
		[]string{"code", "method", "handler"},
	)

	httpResponseSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "A histogram of response sizes for responses.",
			Buckets: []float64{200, 500, 900, 1500},
		},
		[]string{"code", "method", "handler"},
	)
)

func init() {
	// prometheus.MustRegister(httpRequestsTotal, httpRequestDuration, httpRequestSize, httpResponseSize)
	prometheus.MustRegister()
}

// PromRequestHandler return a handler Func
func PromRequestHandler(name string, handler func(http.ResponseWriter, *http.Request)) (string, http.Handler) {
	return name, promhttp.InstrumentHandlerDuration(httpRequestDuration.MustCurryWith(prometheus.Labels{"handler": name}),
		promhttp.InstrumentHandlerCounter(httpRequestsTotal.MustCurryWith(prometheus.Labels{"handler": name}),
			promhttp.InstrumentHandlerRequestSize(httpRequestSize.MustCurryWith(prometheus.Labels{"handler": name}),
				promhttp.InstrumentHandlerResponseSize(httpResponseSize.MustCurryWith(prometheus.Labels{"handler": name}),
					http.HandlerFunc(handler)),
			),
		),
	)
}
