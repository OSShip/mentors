package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/OSShip/mentors/internal/model"
)

func (h *Handler) Apply(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-Id")
	githubUser := r.Header.Get("X-Github-Username")
	if userID == "" {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
		return
	}
	if githubUser == "" {
		var req struct {
			GithubUsername string `json:"github_username"`
		}
		_ = json.NewDecoder(r.Body).Decode(&req)
		githubUser = req.GithubUsername
	}
	if githubUser == "" {
		http.Error(w, `{"error":"github_username required"}`, http.StatusBadRequest)
		return
	}

	pending, err := h.Store.HasPendingApplication(r.Context(), userID)
	if err != nil {
		http.Error(w, `{"error":"internal"}`, http.StatusInternalServerError)
		return
	}
	if pending {
		http.Error(w, `{"error":"application already pending"}`, http.StatusConflict)
		return
	}

	githubData := h.Github.FetchContributions(githubUser)
	app, err := h.Store.CreateApplication(r.Context(), userID, githubData)
	if err != nil {
		http.Error(w, `{"error":"application already exists"}`, http.StatusConflict)
		return
	}
	WriteJSON(w, http.StatusCreated, app)
}

func (h *Handler) ListApplications(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-User-Role") != "admin" {
		http.Error(w, `{"error":"admin required"}`, http.StatusForbidden)
		return
	}
	list, err := h.Store.ListApplications(r.Context(), r.URL.Query().Get("status"))
	if err != nil {
		http.Error(w, `{"error":"internal"}`, http.StatusInternalServerError)
		return
	}
	if list == nil {
		list = []model.Application{}
	}
	WriteJSON(w, http.StatusOK, list)
}

func (h *Handler) ReviewApplication(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-User-Role") != "admin" {
		http.Error(w, `{"error":"admin required"}`, http.StatusForbidden)
		return
	}
	id := chi.URLParam(r, "id")
	adminID := r.Header.Get("X-User-Id")
	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || (req.Status != "approved" && req.Status != "rejected") {
		http.Error(w, `{"error":"status must be approved or rejected"}`, http.StatusBadRequest)
		return
	}
	userID, err := h.Store.ReviewApplication(r.Context(), id, req.Status, adminID)
	if err != nil {
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
		return
	}
	if req.Status == "approved" {
		_ = h.Store.PromoteToMentor(r.Context(), userID)
		email, _ := h.Store.GetUserEmail(r.Context(), userID)
		_ = h.Events.PublishApproved(r.Context(), userID, email)
	}
	WriteJSON(w, http.StatusOK, map[string]string{"status": req.Status})
}
