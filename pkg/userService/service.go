package userservice

import (
	"TestTask/pkg/config"
	userclient "TestTask/pkg/userClient"
	usersrepository "TestTask/pkg/usersRepository"

	"github.com/pkg/errors"
)

type UserService struct {
	repo *usersrepository.Repository
	uc   *userclient.UserClient
	c    *config.Config
}

type UserRequest struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
}

type DeleteUserRequest struct {
	ID int `json:"id"`
}

func NewUserService(repo *usersrepository.Repository, uc *userclient.UserClient, c *config.Config) *UserService {

	return &UserService{repo: repo, uc: uc, c: c}
}

func (s *UserService) AddUserService(req UserRequest) error {
	national := s.uc.GetNational(req.Name)
	age := s.uc.GetAge(req.Name)
	gender := s.uc.GetGender(req.Name)

	err := s.repo.AddUser(req.Name, req.Surname, national, gender, age)
	if err != nil {
		return errors.Wrap(err, "ошибка добавления пользователя")
	}

	return nil
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
		return err
	}

	return nil
}

func (s *UserService) UpdateUserService(id int, user *usersrepository.DBUser) error {
	err := s.repo.UpdateUser(id, user)
	if err != nil {
		return err
	}
	return nil
}
