package users

import (
	"encoding/json"
	"fazz/app/internal/domain"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	baseHandler
	users Service
}

func NewUsersHandler(users Service) *Handler {
	return &Handler{users: users}
}

func (h *Handler) Register(router *mux.Router) {
	router.HandleFunc("/api/v1/users", h.listUsers).Methods("GET")
	router.HandleFunc("/api/v1/users", h.registrationUser).Methods("POST")
	router.HandleFunc("/api/v1/users/auth", h.authorizationUser).Methods("POST")
}

func (h *Handler) listUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.users.GetListUsers(r.Context())
	if err != nil {
		h.ResponseErrorJson(w, "", http.StatusBadRequest)
		return
	}
	h.ResponseJson(w, users, 200)
}

func (h *Handler) registrationUser(w http.ResponseWriter, r *http.Request) {
	var user domain.CreateUser
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		h.ResponseErrorJson(w, "wrong data", http.StatusBadRequest)
		return
	}
	if err := user.Validate(); err != nil {
		h.ResponseErrorJson(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	newUser, err := h.users.CreateUser(r.Context(), user)
	if err != nil {
		h.ResponseErrorJson(w, "user is exists", http.StatusBadRequest)
		return
	}

	h.ResponseJson(w, newUser, http.StatusCreated)
}

func (h *Handler) authorizationUser(w http.ResponseWriter, r *http.Request) {
	var rUser domain.LoginUser

	if err := json.NewDecoder(r.Body).Decode(&rUser); err != nil {
		h.ResponseErrorJson(w, "wrong data", http.StatusBadRequest)
		return
	}

	token, err := h.users.AuthorizationUser(r.Context(), rUser)
	if err != nil {
		h.ResponseErrorJson(w, fmt.Sprint(err), http.StatusBadRequest)
		return
	}
	h.ResponseJson(w, map[string]string{"token": token}, 200)
}
