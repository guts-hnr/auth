package handler

import (
	"auth/middleware"
	"auth/model"
	"auth/service"
	"encoding/json"
	"net/http"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input model.UserRegister

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = h.userService.Register(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte("user created"))
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var request model.UserLogin

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	response, err := h.userService.Login(request.Identifier, request.IdentifierType, request.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)

	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) Profile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.userService.GetProfile(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := model.ProfileResponse{
		Name:     user.Name,
		Age:      user.Age,
		Username: user.Username,
		Email:    user.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
}
