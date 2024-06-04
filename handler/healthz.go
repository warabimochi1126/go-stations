package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/TechBowl-japan/go-stations/model"
)

// A HealthzHandler implements health check endpoint.
type HealthzHandler struct{}

// NewHealthzHandler returns HealthzHandler based http.Handler.
func NewHealthzHandler() *HealthzHandler {
	return &HealthzHandler{}
}

// ServeHTTP implements http.Handler interface.
func (h *HealthzHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_ = &model.HealthzResponse{}

	var healthzResponse model.HealthzResponse
	healthzResponse.Message = "OK"

	err := json.NewEncoder(w).Encode(healthzResponse)

	if err != nil {
		log.Println(err)
	}

}
