package user

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/zh0vtyj/allincecup-server/internal/db"
	"github.com/zh0vtyj/allincecup-server/pkg/models"
)

type AuthorizationStorage interface {
	CreateUser(user User, role string) (int, int, error)
	GetUser(email string, password string) (User, error)
	NewSession(session models.Session) (*models.Session, error)
	GetSessionByRefresh(refresh string) (*models.Session, error)
	DeleteSessionByRefresh(refresh string) error
	DeleteSessionByUserId(id int) error
	UpdateRefreshToken(userId int, newRefreshToken string) error
	GetUserPasswordHash(userId int) (string, error)
	UpdatePassword(userId int, newPassword string) error
}

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (a *AuthPostgres) CreateUser(user User, role string) (int, int, error) {
	// Transaction begin
	tx, err := a.db.Begin()
	if err != nil {
		return 0, 0, err
	}

	var clientRoleId int
	query := fmt.Sprintf("SELECT id FROM %s WHERE role_title=$1", db.RolesTable)
	row := tx.QueryRow(query, role)
	if err = row.Scan(&clientRoleId); err != nil {
		_ = tx.Rollback()
		return 0, 0, err
	}

	var id int // variable for user's id
	var userRoleId int

	query = fmt.Sprintf("INSERT INTO %s (role_id, name, email, password_hash, phone_number) values ($1, $2, $3, $4, $5) RETURNING id, role_id", db.UsersTable)
	row = tx.QueryRow(query, clientRoleId, user.Name, user.Email, user.Password, user.PhoneNumber)
	if err = row.Scan(&id, &userRoleId); err != nil {
		_ = tx.Rollback() // db rollback in error case
		return 0, 0, err
	}

	// new user's cart query
	query = fmt.Sprintf("INSERT INTO %s (user_id) values ($1)", db.CartsTable)
	_, err = tx.Exec(query, id)
	if err != nil {
		_ = tx.Rollback() // db rollback in error case
		return 0, 0, err
	}

	// return id and transaction commit
	return id, userRoleId, tx.Commit()
}

func (a *AuthPostgres) GetUser(email, password string) (User, error) {
	var user User
	query := fmt.Sprintf("SELECT id, role_id FROM %s WHERE email=$1 AND password_hash=$2", db.UsersTable)
	err := a.db.Get(&user, query, email, password)

	return user, err
}

func (a *AuthPostgres) NewSession(session models.Session) (*models.Session, error) {
	queryDeleteOldSession := fmt.Sprintf("DELETE FROM %s WHERE user_id=$1", db.SessionsTable)
	_, err := a.db.Exec(queryDeleteOldSession, session.UserId)
	if err != nil {
		return nil, err
	}

	var newSession models.Session
	query := fmt.Sprintf(
		"INSERT INTO %s (user_id, role_id, refresh_token, client_ip, user_agent, expires_at) values ($1, $2, $3, $4, $5, $6) RETURNING *",
		db.SessionsTable,
	)
	row := a.db.QueryRow(query, session.UserId, session.RoleId, session.RefreshToken, session.ClientIp, session.UserAgent, session.ExpiresAt)

	if err = row.Scan(
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
	query := fmt.Sprintf("DELETE from %s WHERE refresh_token=$1 RETURNING id", db.SessionsTable)
	row := a.db.QueryRow(query, refresh)

	err := row.Scan(&id)
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthPostgres) GetSessionByRefresh(refresh string) (*models.Session, error) {
	var session models.Session
	queryGetSession := fmt.Sprintf("SELECT * from %s WHERE refresh_token=$1 LIMIT 1", db.SessionsTable)
	err := a.db.Get(&session, queryGetSession, refresh)

	if err != nil {
		return nil, fmt.Errorf("session wasn't found by refresh=%s, due to: %v", refresh, err)
	}

	return &session, nil
}

func (a *AuthPostgres) DeleteSessionByUserId(id int) error {
	queryDeleteSession := fmt.Sprintf("DELETE FROM %s WHERE user_id=$1", db.SessionsTable)
	_, err := a.db.Exec(queryDeleteSession, id)
	return err
}

func (a *AuthPostgres) UpdateRefreshToken(userId int, newRefreshToken string) error {
	queryUpdateRefreshToken := fmt.Sprintf("UPDATE %s SET refresh_token=$1 WHERE user_id=$2", db.SessionsTable)
	_, err := a.db.Exec(queryUpdateRefreshToken, newRefreshToken, userId)
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthPostgres) GetUserPasswordHash(userId int) (string, error) {
	var hash string
	queryGetHash := fmt.Sprintf("SELECT password_hash FROM %s WHERE id=$1", db.UsersTable)

	err := a.db.Get(&hash, queryGetHash, userId)
	if err != nil {
		return "", fmt.Errorf("failed to get password hash due to: %v", err)
	}

	return hash, nil
}

func (a *AuthPostgres) UpdatePassword(userId int, newPassword string) error {
	queryUpdatePassword := fmt.Sprintf("UPDATE %s SET password_hash=$1 WHERE id=$2", db.UsersTable)

	_, err := a.db.Exec(queryUpdatePassword, newPassword, userId)
	if err != nil {
		return fmt.Errorf("failed to update password in db due to: %v", err)
	}

	return nil
}
