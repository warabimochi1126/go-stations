package router

import (
	"database/sql"
	"net/http"

	"github.com/TechBowl-japan/go-stations/handler"
	"github.com/TechBowl-japan/go-stations/handler/middleware"
	"github.com/TechBowl-japan/go-stations/service"
)

func NewRouter(todoDB *sql.DB) *http.ServeMux {
	// register routes
	mux := http.NewServeMux()

	healthzHandler := handler.NewHealthzHandler()
	todoHandler := handler.NewTODOHandler(service.NewTODOService(todoDB))

	mux.Handle("/healthz", healthzHandler)
	mux.Handle("/todos", todoHandler)

	// Go基礎編Station1
	doPanicHandler := handler.NewDoPanicHandler()
	mux.Handle("/do-panic", middleware.Recovery(doPanicHandler))

	// Go基礎編Station2
	getOsHandler := handler.NewGetOsHandler()
	mux.Handle("/get-os", middleware.GetOS(getOsHandler))

	// Go基礎編Station3
	latencyHandler := handler.NewLatencyHandler()
	mux.Handle("/latency", middleware.GetOS(middleware.GetTime(latencyHandler)))

	return mux
}
