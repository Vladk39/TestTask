package api

import (
	userservice "TestTask/pkg/userService"
	usersrepository "TestTask/pkg/usersRepository"
	"encoding/json"
	"net/http"
	"strconv"
)

type Handler struct {
	UserService *userservice.UserService
}

func NewHandler(userService *userservice.UserService) *Handler {
	return &Handler{
		UserService: userService,
	}
}

func (h *Handler) AddUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "неверный метод", http.StatusMethodNotAllowed)
		return
	}

	var req userservice.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "неверный запрос", http.StatusBadRequest)
		return
	}

	err := h.UserService.AddUserService(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Пользователь добавлен"))
}

func (h *Handler) GetUserByFilterServiceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "неверный метод", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()

	national := query.Get("national")
	gender := query.Get("gender")

	limit := 1
	offset := 0

	if lim := query.Get("limit"); lim != "" {
		if parsedLimit, err := strconv.Atoi(lim); err == nil {
			limit = parsedLimit
		}
	}

	if offs := query.Get("offset"); offs != "" {
		if parsedOffset, err := strconv.Atoi(offs); err == nil {
			offset = parsedOffset
		}
	}

	users, err := h.UserService.GetUserByFilterService(gender, national, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "ошибка отправки JSON", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "неверный метод", http.StatusMethodNotAllowed)
		return
	}

	allusers, err := h.UserService.GetAllUserService()
	if err != nil {
		http.Error(w, "ошибка получения пользователей", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(allusers); err != nil {
		http.Error(w, "ошибка отправки JSON", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "неверный метод", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "id обязателен", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "id должен быть числом", http.StatusBadRequest)
		return
	}

	err = h.UserService.DeleteUserService(id)
	if err != nil {
		http.Error(w, "ошибка удаления пользователя", http.StatusBadRequest)
		return
	}

	w.Write([]byte("Пользователь удалён"))
}

func (h *Handler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "неверный метод", http.StatusMethodNotAllowed)
		return
	}

	var req *usersrepository.DBUser
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "неверное тело запроса", http.StatusBadRequest)
		return
	}

	err := h.UserService.UpdateUserService(req.ID, req)
	if err != nil {
		http.Error(w, "ошибка обновления пользователя", http.StatusBadRequest)
		return
	}

	w.Write([]byte("пользователь обновлен"))
}
