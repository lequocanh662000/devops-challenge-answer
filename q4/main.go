package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	infuraEndpoint = "https://mainnet.infura.io/v3/3ceb75252d2f4137b099deda23a7b6e9"
	ankrEndpoint   = "https://rpc.ankr.com/eth/1566f3525a809c0e61d27de6e9412f0218267051e6e833e15904421c340f40fd"
)

type metrics struct {
	cpuTemp    prometheus.Gauge
	hdFailures *prometheus.CounterVec
	apiRequest *prometheus.CounterVec
}

func getBlockNumber(endpoint string) (int64, error) {
	reqBody := []byte(`{"jsonrpc": "2.0", "method": "eth_blockNumber", "params": [], "id": 1}`)
	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}

	result, ok := data["result"].(string)
	if !ok {
		return 0, fmt.Errorf("Invalid response format")
	}

	blockNumber, err := strconv.ParseInt(result, 0, 64)
	if err != nil {
		return 0, err
	}

	return blockNumber, nil
}

func checkBlockNumberDifference() (bool, error) {
	ankrBlockNumber, err := getBlockNumber(ankrEndpoint)
	if err != nil {
		return false, err
	}

	infuraBlockNumber, err := getBlockNumber(infuraEndpoint)
	if err != nil {
		return false, err
	}

	difference := ankrBlockNumber - infuraBlockNumber
	return difference < 5, nil
}

func handleScrape(w http.ResponseWriter, r *http.Request) {
	success, err := checkBlockNumberDifference()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if success {
		w.Write([]byte("succes"))
	} else {
		w.Write([]byte("fail"))
	}
}

func NewMetrics(reg prometheus.Registerer) *metrics {
	m := &metrics{
		cpuTemp: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "cpu_temperature_celsius",
			Help: "Current temperature of the CPU.",
		}),
		hdFailures: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "hd_errors_total",
				Help: "Number of hard-disk errors.",
			},
			[]string{"device"},
		),
		apiRequest: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Number of http requests received.",
			},
			[]string{"method", "path"},
		),
	}
	reg.MustRegister(m.cpuTemp)
	reg.MustRegister(m.hdFailures)
	reg.MustRegister(m.apiRequest)
	return m
}

func main() {
	// Create a non-global registry.
	reg := prometheus.NewRegistry()

	// Create new metrics and register them using the custom registry.
	m := NewMetrics(reg)

	// Set values for the new created metrics.
	m.cpuTemp.Set(65.3)
	m.hdFailures.With(prometheus.Labels{"device": "/dev/sda"}).Inc()

	http.HandleFunc("/", handleScrape)
	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		// Set values for total /api/* requests
		m.apiRequest.WithLabelValues(r.Method, r.URL.Path).Inc()
	})
	// Expose metrics and custom registry via an HTTP server
	// using the HandleFor function. "/metrics" is the usual endpoint for that.
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
