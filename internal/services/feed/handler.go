package feed

import (
	"encoding/json"
	"net/http"
)

type FeedStore interface {
	GetFeed(userID string, limit int) ([]string, error)
}

type Handler struct {
	store FeedStore
}

func NewHandler(store FeedStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) GetFeed(w http.ResponseWriter, r *http.Request) {
	user := r.URL.Query().Get("user")

	if user == "" {
		http.Error(w, "missing user", http.StatusBadRequest)
		return
	}

	posts, err := h.store.GetFeed(user, 20)
	if err != nil {
		http.Error(w, "failed to fetch feed", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]any{
		"user": user,
		"posts": posts,
	})
}