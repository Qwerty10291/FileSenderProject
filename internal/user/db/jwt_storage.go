package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type Jwt struct {
	UserId int `db:"user_id"`
	Token string `db:"token"`
	RefreshToken string `db:"refresh_token"`
	Expires time.Time `db:"expires"`
}

type JwtStorage struct{
	db *sqlx.DB
}

func (s *JwtStorage) Create(jwt Jwt) error{
	_, err := s.db.Exec("INSERT INTO jwt_tokens (user_id, token, refresh_token, expires) VALUES ($1, $2, $3, $4)")
	return err
}

func (s *JwtStorage) Delete(userId int) (bool, error){
	result, err := s.db.Exec("DELETE FROM jwt_tokens WHERE user_id = $1", userId)
	if err != nil{
		return false, fmt.Errorf("error when remove jwt token of %d: %s", userId, err)
	}
	numAffected, err := result.RowsAffected()
	if err != nil{
		return false, fmt.Errorf("error when getting affected rows of deleted jwt: %s", err)
	}
	return numAffected > 0, nil
}

func (s *JwtStorage) FromUser(user *User) (*Jwt, error) {
	jwt := new(Jwt)
	err := s.db.Get(jwt, "SELECT * FROM jwt_tokens WHERE user_id = $1", user.Id)
	if err != nil{
		if err == sql.ErrNoRows{
			return nil, nil
		}
		return nil, fmt.Errorf("error when getting jwt from user %d: %s", user.Id, err)
	}
	return jwt, nil
}

func (s *JwtStorage) FromUsername(name string) (*Jwt, error) {
	jwt := new(Jwt)
	err := s.db.Get(jwt, "SELECT * FROM jwt_tokens WHERE user_id = (SELECT id from users WHERE login = $1)", name)
	if err != nil{
		if err == sql.ErrNoRows{
			return nil, nil
		}
		return nil, err
	}
	return jwt, nil
}