package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

type (
	MetricsHandler struct {
		requestChan  chan uint64
		responseChan chan uint64
	}
)

func NewMetricsHandler(requestChan chan uint64, responseChan chan uint64) http.Handler {
	return &MetricsHandler{
		requestChan:  requestChan,
		responseChan: responseChan,
	}
}

func (m *MetricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")

	if contentType != "" {
		contentType = "application/json"
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	//w.Header().Set("Access-Control-Allow-Headers", "origin, content-type, accept, authorization")
	//w.Header().Set("Access-Control-Allow-Credentials", "true")
	//w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, HEAD")
	//w.Header().Set("Content-Disposition", "attachment; filename=\"results.json\"")

	m.requestChan <- 1

	err := json.NewEncoder(w).Encode(struct {
		Count uint64 `json:"count"`
	}{
		Count: <-m.responseChan,
	})

	if err != nil {
		log.Fatal(err)
	}
}
