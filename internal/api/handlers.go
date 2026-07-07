// Package api implements the HTTP layer for the Natter API.
// It provides request parsing, JSON responses, and error mapping from service errors to HTTP status codes.
package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/mrl00/natter/internal/model"
	"github.com/mrl00/natter/internal/service"
)

// Handlers holds the service dependency and exposes HTTP handler methods for each endpoint.
type Handlers struct {
	svc service.Service
}

func NewHandlers(s service.Service) *Handlers {
	return &Handlers{svc: s}
}

func jsonResponse(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func jsonError(w http.ResponseWriter, status int, msg string) {
	jsonResponse(w, status, map[string]string{"error": msg})
}

// Health godoc
// @Summary      Health check
// @Description  Returns the service status
// @Tags         health
// @Produce      json
// @Success      200 {object} map[string]string
// @Router       /health [get]
func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, http.StatusOK, map[string]string{"status": "ok"})
}

// CreateSpace godoc
// @Summary      Create a space
// @Description  Creates a new social space. The user that performs this request becomes the owner.
// @Tags         spaces
// @Accept       json
// @Produce      json
// @Param        request body model.CreateSpaceRequest true "Space data"
// @Success      201 {object} model.CreateSpaceResponse
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /spaces [post]
func (h *Handlers) CreateSpace(w http.ResponseWriter, r *http.Request) {
	var req model.CreateSpaceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" || req.Owner == "" {
		jsonError(w, http.StatusBadRequest, "name and owner are required")
		return
	}

	space, err := h.svc.CreateSpace(req.Name, req.Owner)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to create space")
		return
	}

	jsonResponse(w, http.StatusCreated, model.CreateSpaceResponse{
		Name: space.Name,
		URI:  fmt.Sprintf("/spaces/%s", space.ID),
	})
}

// AddMessage godoc
// @Summary      Add a message
// @Description  Adds a message to a social space
// @Tags         messages
// @Accept       json
// @Produce      json
// @Param        spaceId path string true "Space ID"
// @Param        request body model.CreateMessageRequest true "Message data"
// @Success      201 {object} model.Message
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /spaces/{spaceId}/messages [post]
func (h *Handlers) AddMessage(w http.ResponseWriter, r *http.Request) {
	spaceID := r.PathValue("spaceId")

	var req model.CreateMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Author == "" || req.Content == "" {
		jsonError(w, http.StatusBadRequest, "author and content are required")
		return
	}

	msg, err := h.svc.AddMessage(spaceID, req.Author, req.Content)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			jsonError(w, http.StatusNotFound, "space not found")
			return
		}
		jsonError(w, http.StatusInternalServerError, "failed to add message")
		return
	}

	jsonResponse(w, http.StatusCreated, msg)
}

// ListMessages godoc
// @Summary      List messages
// @Description  Returns all messages in a space, optionally filtered by a since timestamp
// @Tags         messages
// @Produce      json
// @Param        spaceId path string true "Space ID"
// @Param        since query string false "RFC3339 timestamp to filter messages"
// @Success      200 {array} model.Message
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /spaces/{spaceId}/messages [get]
func (h *Handlers) ListMessages(w http.ResponseWriter, r *http.Request) {
	spaceID := r.PathValue("spaceId")

	var since time.Time
	if sinceStr := r.URL.Query().Get("since"); sinceStr != "" {
		var err error
		since, err = time.Parse(time.RFC3339, sinceStr)
		if err != nil {
			jsonError(w, http.StatusBadRequest, "invalid since timestamp (use RFC3339)")
			return
		}
	}

	msgs, err := h.svc.ListMessages(spaceID, since)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			jsonError(w, http.StatusNotFound, "space not found")
			return
		}
		jsonError(w, http.StatusInternalServerError, "failed to list messages")
		return
	}

	jsonResponse(w, http.StatusOK, msgs)
}

// GetMessage godoc
// @Summary      Get a message
// @Description  Returns the details of a single message
// @Tags         messages
// @Produce      json
// @Param        spaceId path string true "Space ID"
// @Param        messageId path string true "Message ID"
// @Success      200 {object} model.Message
// @Failure      404 {object} map[string]string
// @Router       /spaces/{spaceId}/messages/{messageId} [get]
func (h *Handlers) GetMessage(w http.ResponseWriter, r *http.Request) {
	spaceID := r.PathValue("spaceId")
	messageID := r.PathValue("messageId")

	msg, err := h.svc.GetMessage(spaceID, messageID)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			jsonError(w, http.StatusNotFound, "space or message not found")
			return
		}
		jsonError(w, http.StatusInternalServerError, "failed to get message")
		return
	}

	jsonResponse(w, http.StatusOK, msg)
}
