// Package api implements the HTTP layer for the Natter API.
// This file configures the HTTP router with all Natter API endpoints.
package api

import (
	"net/http"
)

func NewRouter(h *Handlers) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", h.Health)
	mux.HandleFunc("POST /spaces", h.CreateSpace)
	mux.HandleFunc("POST /spaces/{spaceId}/messages", h.AddMessage)
	mux.HandleFunc("GET /spaces/{spaceId}/messages", h.ListMessages)
	mux.HandleFunc("GET /spaces/{spaceId}/messages/{messageId}", h.GetMessage)

	return mux
}
