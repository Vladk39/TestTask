package api

import (
	"TestTask/pkg/config"
	"net/http"

	"TestTask/pkg/api/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

func StartServer(handler *Handler, c *config.ServerConfig) error {
	docs.SwaggerInfo.BasePath = "/"

	http.HandleFunc("/add-user", handler.AddUserHandler)
	http.HandleFunc("/get-users", handler.GetUserByFilterServiceHandler)
	http.HandleFunc("/get-all-users", handler.GetAllUsersHandler)
	http.HandleFunc("/delete-user", handler.DeleteUserHandler)
	http.HandleFunc("/update-user", handler.UpdateUserHandler)

	http.Handle("/swagger/", httpSwagger.WrapHandler)

	return http.ListenAndServe(c.Port, nil)
}
