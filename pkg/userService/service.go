package userservice

import (
	"TestTask/pkg/config"
	userclient "TestTask/pkg/userClient"
	usersrepository "TestTask/pkg/usersRepository"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	// "gorm.io/gorm/logger"
)

type UserService struct {
	repo   *usersrepository.Repository
	uc     *userclient.UserClient
	c      *config.Config
	logger *zap.Logger
}

type UserRequest struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
}

type DeleteUserRequest struct {
	ID int `json:"id"`
}

func NewUserService(repo *usersrepository.Repository, uc *userclient.UserClient, c *config.Config, logger *zap.Logger) *UserService {

	return &UserService{repo: repo, uc: uc, c: c, logger: logger}
}

func (s *UserService) AddUserService(req UserRequest) error {
	nationalCh := make(chan string, 1)
	ageCh := make(chan int, 1)
	genderCh := make(chan string, 1)

	go func() {
		nationalCh <- s.uc.GetNational(req.Name)
	}()

	go func() {
		ageCh <- s.uc.GetAge(req.Name)
	}()

	go func() {
		genderCh <- s.uc.GetGender(req.Name)
	}()

	err := s.repo.AddUser(req.Name, req.Surname, <-nationalCh, <-genderCh, <-ageCh)
	if err != nil {
		s.logger.Error("Ошибка добавления пользователя", zap.Error(err))
		return errors.Wrap(err, "ошибка добавления пользователя")
	}

	return nil
}

func (s *UserService) SearchUserService(req UserRequest) (bool, error) {
	exist, err := s.repo.SearchUser(req.Name, req.Surname)
	if err != nil {
		s.logger.Error("Ошибка поиска пользователя", zap.Error(err))
		return false, errors.Wrap(err, "ошибка поиска пользователя")
	}

	return exist, nil
}

func (s *UserService) GetUserByFilterService(gender, national string, limit, offset int) ([]usersrepository.DBUser, error) {

	return s.repo.GetUserByFilter(gender, national, limit, offset)
}

func (s *UserService) GetAllUserService() ([]usersrepository.DBUser, error) {

	return s.repo.GetAllusers()
}

func (s *UserService) DeleteUserService(id int) error {
	err := s.repo.DeleteUser(id)
	if err != nil {
		s.logger.Error("Ошибка  удаления пользователя", zap.Error(err))
		return err
	}

	return nil
}

func (s *UserService) UpdateUserService(id int, user *usersrepository.DBUser) error {
	err := s.repo.UpdateUser(id, user)
	if err != nil {
		s.logger.Error("Ошибка обновления пользователя", zap.Error(err))
		return err
	}
	return nil
}
