package user

import (
	"net/http"
)

type Handler struct {
	store interface {
		CreateUser(string) error
		Follow(string, string) error
	}
}

func NewHandler(store interface {
		CreateUser(string) error
		Follow(string, string) error
	}) *Handler {
		return &Handler{store: store}
	}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	h.store.CreateUser(id)
	w.Write([]byte("user created"))
}

func (h *Handler) Follow(w http.ResponseWriter, r *http.Request) {
	u := r.URL.Query().Get("u")
	v := r.URL.Query().Get("v")

	if u == "" || v == "" {
		http.Error(w, "missing u or v", http.StatusBadRequest)
		return
	}
	h.store.Follow(u, v)
	w.Write([]byte("followed"))
}