package internal

import (
	"encoding/json"
	"net/http"
)

type DefaultRouteHandler struct{}

func NewDefaultRouteHandler() *DefaultRouteHandler {
	return &DefaultRouteHandler{}
}

func (h *DefaultRouteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	response := map[string]string{"message": "hello, world"}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
