package usersrepository

import (
	"TestTask/pkg/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/pressly/goose"
)

const GetAllusers = `select * from users`

const GetByIdSQL = `select * from users where id = $1`

const UpdateUser = `update users
set name = $1, surname = $2, age =$3, gender = $4, national = $5
where id = $6`

const DeleteUser = `DELETE FROM users WHERE id = $1`

const GetWithFilterPaginated = `SELECT * FROM users WHERE gender = $1 AND national = $2 ORDER BY id LIMIT $3 OFFSET $4`

const AddUser = `insert into users (name, surname, national, gender, age) values ($1, $2, $3, $4, $5)`

type Repository struct {
	DB *sqlx.DB
}

type DBUser struct {
	ID       int    `db:"id" json:"id"`
	Name     string `db:"name" json:"name"`
	Surname  string `db:"surname" json:"surname"`
	Age      int    `db:"age" json:"age"`
	Gender   string `db:"gender" json:"gender"`
	Natioanl string `db:"national" json:"national"`
}

func connectDB(c *config.DBConfig) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", c.DBconn)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка подключения к БД")
	}
	return db, nil
}

func NewRepository(c *config.DBConfig) (*Repository, error) {
	db, err := connectDB(c)
	if err != nil {
		return nil, err
	}

	if err := applyMigrations(db); err != nil {
		return nil, errors.Wrap(err, "ошибка применения миграций")
	}

	return &Repository{
		DB: db,
	}, nil
}

func applyMigrations(DB *sqlx.DB) error {
	migrationsDir := "../migrations"
	if err := goose.Up(DB.DB, migrationsDir); err != nil {
		return errors.Wrap(err, "ошибка миграции")
	}
	return nil
}

func (repo *Repository) GetUserByFilter(gender, national string, limit, offset int) ([]DBUser, error) {
	var users []DBUser
	err := repo.DB.Select(&users, GetWithFilterPaginated, gender, national, limit, offset)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка поиска по базу данных")
	}

	return users, nil
}

func (repo *Repository) DeleteUser(id int) error {
	_, err := repo.DB.Exec(DeleteUser, id)
	if err != nil {
		return errors.Wrap(err, "ошибка удаления пользователя")
	}

	return nil
}

func (repo *Repository) GetAllusers() ([]DBUser, error) {
	var AllUsers []DBUser
	err := repo.DB.Select(&AllUsers, GetAllusers)
	if err != nil {
		return nil, err
	}

	return AllUsers, err
}

func (repo *Repository) GetUserById(id int) (*DBUser, error) {
	var User DBUser
	err := repo.DB.Get(&User, GetByIdSQL, id)

	return &User, err
}

func (repo *Repository) UpdateUser(id int, user *DBUser) error {
	_, err := repo.DB.Exec(UpdateUser, user.Name, user.Surname, user.Age, user.Gender, user.Natioanl, id)
	if err != nil {
		return errors.Wrap(err, "ошибка обновления пользователя")
	}

	return nil
}

func (repo *Repository) AddUser(name, surname, national, gender string, age int) error {

	_, err := repo.DB.Exec(AddUser, name, surname, national, gender, age)
	if err != nil {
		return errors.Wrap(err, "ошибка добавления пользователя")
	}
	return nil
}
