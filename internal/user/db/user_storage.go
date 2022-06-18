package db

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type User struct{
	Id int `json:"id" db:"id"`
	Login string `json:"login" db:"login"`
	Password string `db:"password"`
}

type UserStorage struct {
	db *sqlx.DB
}

func NewUserStorage(db *sqlx.DB) *UserStorage {
	return &UserStorage{
		db: db,
	}
}

func (s *UserStorage) Delete(id int) (bool, error) {
	result, err := s.db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil{
		return false, fmt.Errorf("error when deleting user:%s", err)
	}
	count, err :=  result.RowsAffected()
	if err != nil{
		return false, fmt.Errorf("cannot get affected rows count: %s", err)
	}
	if count == 0{
		return false, nil
	}
	return true, nil
}

func (s *UserStorage) GetById(id int) (*User, error) {
	user := new(User)
	err := s.db.Get(user, "SELECT * FROM users WHERE id = $1", id)
	if err != nil{
		if err == sql.ErrNoRows{
			return nil, nil
		}
		return nil, fmt.Errorf("error when fetching user by id: %s", err)
	}
	return user, nil
}

func (s *UserStorage) GetByLogin(login string) (*User, error) {
	user := new(User)
	err := s.db.Get(user, "SELECT * FROM users WHERE login = $1", login)
	if err != nil{
		return nil, fmt.Errorf("error when fetching user by login: %s", err)
	}
	return user, nil
}

func (s *UserStorage) Create(login string, password string) (*User, error) {
	result, err := s.db.Exec("INSERT INTO users (login, password) VALUES ($1, $2)")
	if err != nil{
		return nil, fmt.Errorf("error when creating user: %s", err)
	}
	id, err := result.LastInsertId()
	if err != nil{
		return nil, fmt.Errorf("error when getting id of created user: %s", err)
	}
	user := new(User)
	user.Id = int(id)
	user.Login = login
	user.Password = password
	return user, nil
}