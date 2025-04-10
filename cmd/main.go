package main

import (
	"TestTask/pkg/api"
	_ "TestTask/pkg/api/docs"
	"TestTask/pkg/config"
	userclient "TestTask/pkg/userClient"
	userservice "TestTask/pkg/userService"
	usersrepository "TestTask/pkg/usersRepository"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		errors.Wrap(err, "ошибка запуска логера")
		return
	}
	defer logger.Sync()

	logger.Info("Приложение запущено")

	config, err := config.GetConfig()
	if err != nil {
		logger.Error("не удалось получить конфигурации", zap.Error(err))
	}

	repo, err := usersrepository.NewRepository(&config.DBConfig)
	if err != nil {
		logger.Error("Не удалось создать репозиторий", zap.Error(err))
		return
	}
	logger.Info("Репозиторий создан")

	userClient := userclient.NewUserClient(config, logger)

	userService := userservice.NewUserService(repo, userClient, config, logger)

	handler := api.NewHandler(userService)

	logger.Info("Запуск HTTP сервера", zap.String("port", config.Port))
	if err := api.StartServer(handler, &config.ServerConfig); err != nil {
		logger.Error("ошибка сервер", zap.Error(err))
	}
}
