package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

const (
	infuraEndpoint = "https://mainnet.infura.io/v3/3ceb75252d2f4137b099deda23a7b6e9"
	ankrEndpoint   = "https://rpc.ankr.com/eth/1566f3525a809c0e61d27de6e9412f0218267051e6e833e15904421c340f40fd"
)

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

func main() {
	http.HandleFunc("/", handleScrape)
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		// Add Prometheus metrics here if needed
	})

	http.ListenAndServe(":8080", nil)
}
