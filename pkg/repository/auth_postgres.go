package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	server "github.com/zh0vtyj/allincecup-server"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (a *AuthPostgres) CreateUser(user server.User, role string) (int, int, error) {
	// Transaction begin
	tx, err := a.db.Begin()
	if err != nil {
		return 0, 0, err
	}

	var clientRoleId int
	query := fmt.Sprintf("SELECT id FROM %s WHERE role_title=$1", rolesTable)
	row := tx.QueryRow(query, role)
	if err = row.Scan(&clientRoleId); err != nil {
		_ = tx.Rollback()
		return 0, 0, err
	}

	var id int // variable for user's id
	var userRoleId int

	query = fmt.Sprintf("INSERT INTO %s (role_id, name, email, password_hash, phone_number) values ($1, $2, $3, $4, $5) RETURNING id, role_id", usersTable)
	row = tx.QueryRow(query, clientRoleId, user.Name, user.Email, user.Password, user.PhoneNumber)
	if err = row.Scan(&id, &userRoleId); err != nil {
		_ = tx.Rollback() // db rollback in error case
		return 0, 0, err
	}

	// new user's cart query
	query = fmt.Sprintf("INSERT INTO %s (user_id) values ($1)", cartsTable)
	_, err = tx.Exec(query, id)
	if err != nil {
		_ = tx.Rollback() // db rollback in error case
		return 0, 0, err
	}

	// return id and transaction commit
	return id, userRoleId, tx.Commit()
}

func (a *AuthPostgres) GetUser(email, password string) (server.User, error) {
	var user server.User
	query := fmt.Sprintf("SELECT id, role_id FROM %s WHERE email=$1 AND password_hash=$2", usersTable)
	err := a.db.Get(&user, query, email, password)

	return user, err
}

func (a *AuthPostgres) NewSession(session server.Session) (*server.Session, error) {
	queryDeleteOldSession := fmt.Sprintf("DELETE FROM %s WHERE user_id=$1", sessionsTable)
	_, err := a.db.Exec(queryDeleteOldSession, session.UserId)
	if err != nil {
		return nil, err
	}

	var newSession server.Session
	query := fmt.Sprintf(
		"INSERT INTO %s (user_id, role_id, refresh_token, client_ip, user_agent, expires_at) values ($1, $2, $3, $4, $5, $6) RETURNING *",
		sessionsTable)
	row := a.db.QueryRow(query, session.UserId, session.RoleId, session.RefreshToken, session.ClientIp, session.UserAgent, session.ExpiresAt)

	if err := row.Scan(
		&newSession.Id,
		&newSession.UserId,
		&newSession.RoleId,
		&newSession.RefreshToken,
		&newSession.ClientIp,
		&newSession.UserAgent,
		&newSession.IsBlocked,
		&newSession.ExpiresAt,
		&newSession.CreatedAt,
	); err != nil {
		return nil, err
	}

	return &newSession, nil
}

func (a *AuthPostgres) DeleteSessionByRefresh(refresh string) error {
	var id int
	query := fmt.Sprintf("DELETE from %s WHERE refresh_token=$1 RETURNING id", sessionsTable)
	row := a.db.QueryRow(query, refresh)

	err := row.Scan(&id)
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthPostgres) GetSessionByRefresh(refresh string) (*server.Session, error) {
	var session server.Session
	query := fmt.Sprintf("SELECT * from %s WHERE refresh_token=$1 LIMIT 1", sessionsTable)
	row := a.db.QueryRow(query, refresh)

	if err := row.Scan(
		&session.Id,
		&session.UserId,
		&session.RoleId,
		&session.RefreshToken,
		&session.ClientIp,
		&session.UserAgent,
		&session.IsBlocked,
		&session.ExpiresAt,
		&session.CreatedAt,
	); err != nil {
		return nil, err
	}

	return &session, nil
}

func (a *AuthPostgres) DeleteSessionByUserId(id int) error {
	queryDeleteSession := fmt.Sprintf("DELETE FROM %s WHERE user_id=$1", sessionsTable)
	_, err := a.db.Exec(queryDeleteSession, id)
	return err
}

func (a *AuthPostgres) UpdateRefreshToken(userId int, newRefreshToken string) error {
	queryUpdateRefreshToken := fmt.Sprintf("UPDATE %s SET refresh_token=$1 WHERE user_id=$2", sessionsTable)
	_, err := a.db.Exec(queryUpdateRefreshToken, newRefreshToken, userId)
	return err
}
