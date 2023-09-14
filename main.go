package main

import (
	"APT-Exporter/cuda_on_time"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

func main() {
	cuda_on_time.Start()
	log.Println("http server init")
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("HTTP Server init failed:", err)
	}
	// cuda use time collect service
}
