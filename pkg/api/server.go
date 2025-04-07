package api

import (
	"TestTask/pkg/config"
	"net/http"
)

func StartServer(handler *Handler, c *config.ServerConfig) error {
	http.HandleFunc("/add-user", handler.AddUserHandler)
	http.HandleFunc("/get-users", handler.GetUserByFilterServiceHandler)
	http.HandleFunc("/get-all-users", handler.GetAllUsersHandler)
	http.HandleFunc("/delete-user", handler.DeleteUserHandler)
	http.HandleFunc("/update-user", handler.UpdateUserHandler)

	return http.ListenAndServe(c.Port, nil)
}
