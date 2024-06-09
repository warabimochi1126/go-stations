package handler

import (
	"fmt"
	"net/http"
)

type GetOsHandler struct{}

func NewGetOsHandler() *GetOsHandler {
	return &GetOsHandler{}
}

func (n *GetOsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	osName := r.Context().Value("osName")
	fmt.Println("osName:", osName)
}
