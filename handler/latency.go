package handler

import (
	"fmt"
	"net/http"
)

type LatencyHandler struct{}

func NewLatencyHandler() *LatencyHandler {
	return &LatencyHandler{}
}

func (l *LatencyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("latencyハンドラーが発火した")
}
