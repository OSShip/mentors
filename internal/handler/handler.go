package handler

import (
	"encoding/json"
	"net/http"

	"github.com/OSShip/mentors/internal/events"
	"github.com/OSShip/mentors/internal/github"
	"github.com/OSShip/mentors/internal/store"
)

type Handler struct {
	Store  *store.Store
	Events *events.Publisher
	Github *github.Client
}

func WriteJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
