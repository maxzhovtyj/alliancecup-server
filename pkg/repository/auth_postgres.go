package repository

import (
	server "allincecup-server"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (a *AuthPostgres) CreateUser(user server.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (role_id, name, email, password_hash, phone_number) values ($1, $2, $3, $4, $5) RETURNING id", usersTable)
	row := a.db.QueryRow(query, user.RoleId, user.Name, user.Email, user.Password, user.PhoneNumber)

	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (a *AuthPostgres) GetUser(email, password string) (server.User, error) {
	var user server.User
	query := fmt.Sprintf("SELECT id, role_id FROM %s WHERE email=$1 AND password_hash=$2", usersTable)
	err := a.db.Get(&user, query, email, password)

	return user, err
}
