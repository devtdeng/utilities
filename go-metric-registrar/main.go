package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var simpleCounter = promauto.NewCounter(prometheus.CounterOpts{
	Name: "simple",
	Help: "Increments when /simple handler is called",
})

var customGauge = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "custom",
	Help: "Custom gauge used to scale app instances",
})

var requestDuration = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "request_duration_seconds",
		Help: "A histogram of latencies for requests.",
	},
	[]string{"code", "method", "handler"},
)

func main() {
	router := mux.NewRouter()

	// endpoint worker
	router.Handle("/metrics", promhttp.Handler())

	router.HandleFunc("/simple", instrument(simple, "app_NAME.metrics_endpoint.simple_counter")).Methods(http.MethodGet)
	router.HandleFunc("/high_latency", instrument(highLatency, "high_latency")).Methods(http.MethodGet)
	router.HandleFunc("/custom_metric", instrument(customMetric, "custom_metric"))

	// log worker
	router.HandleFunc("/log_metric_dogstatsd", instrument(logMetricDogStatsD, "log_metric_dogstatsd"))
	router.HandleFunc("/log_metric_json", instrument(logMetricJSON, "log_metric_json"))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(":"+port, router))
}

func instrument(handlerFunc http.HandlerFunc, name string) http.HandlerFunc {
	handlerDuration, err := requestDuration.CurryWith(prometheus.Labels{
		"handler": name,
	})

	if err != nil {
		panic(err)
	}

	return promhttp.InstrumentHandlerDuration(handlerDuration, handlerFunc)
}

func simple(w http.ResponseWriter, _ *http.Request) {
	simpleCounter.Inc()
	w.Write([]byte("{}"))
}

func highLatency(w http.ResponseWriter, _ *http.Request) {
	time.Sleep(2 * time.Second)
	w.Write([]byte("{}"))
}

func customMetric(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("inc") != "" {
		customGauge.Inc()
	} else {
		customGauge.Dec()
	}

	w.Write([]byte("{}"))
}

func logMetricDogStatsD(w http.ResponseWriter, r *http.Request) {
	log.Printf("metric_registrar.go-metric-registra.count:%d|c\n", 1)
	w.Write([]byte("{}"))
}

func logMetricJSON(w http.ResponseWriter, r *http.Request) {
	log.Printf("{\"type\": \"counter\",\"name\": \"metric_registrar.go-metric-registra.count\",\"delta\": %d}\n", 1)
	w.Write([]byte("{}"))
}
