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

// AddUserHandler добавляет пользователя
// @Summary Добавить пользователя
// @Description Добавляет пользователя с именем, фамилией, и автоматически подтягивает пол, возраст, нацию
// @Tags users
// @Accept json
// @Produce json
// @Param user body userservice.UserRequest true "Данные пользователя"
// @Success 200 {string} string "Пользователь добавлен"
// @Failure 500 {string} string
// @Router /add-user [post]
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

	exists, err := h.UserService.SearchUserService(req)
	if err != nil {
		http.Error(w, "Ошибка поиска пользователя: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if exists {
		http.Error(w, "Пользователь уже существует", http.StatusConflict) // Статус 409 - конфликт
		return
	}

	err = h.UserService.AddUserService(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Пользователь добавлен"))
}

// GetUserByFilterServiceHandler получает пользователей по фильтрам
// @Summary Получить пользователей по фильтру
// @Description Возвращает пользователей по полу, национальности и с пагинацией
// @Tags users
// @Accept json
// @Produce json
// @Param gender query string true "Пол пользователя (например: male, female)"
// @Param national query string true "Национальность (например: US, RU)"
// @Param limit query int false "Максимальное количество пользователей (по умолчанию 10)"
// @Param offset query int false "Смещение для пагинации (по умолчанию 0)"
// @Success 200 {array} usersrepository.DBUser
// @Failure 400 {string} string "Неверный запрос"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /get-users [get]
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

// @Tags users
// @Summary Получить всех пользователей
// @Description Эндпоинт для получения списка всех пользователей
// @Accept json
// @Produce json
// @Success 200 {array} usersrepository.DBUser "Список пользователей"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /get-all-users [get]
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

// @Tags users
// @Summary Удалить пользователя
// @Description Эндпоинт для удаления пользователя по его ID
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя для удаления"
// @Success 200 {string} string "Успешное удаление пользователя"
// @Failure 400 {string} string "Неверный запрос"
// @Failure 404 {string} string "Пользователь не найден"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /delete-user/{id} [delete]
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

// @Tags users
// @Summary Обновить информацию о пользователе
// @Description Эндпоинт для обновления данных пользователя по его ID
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя, данные которого нужно обновить"
// @Param user body userservice.UserRequest true "Данные пользователя для обновления"
// @Success 200 {string} string "Информация о пользователе обновлена"
// @Failure 400 {string} string "Неверный запрос"
// @Failure 404 {string} string "Пользователь не найден"
// @Failure 500 {string} string "Ошибка сервера"
// @Router /update-user/ [post]
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
