package user

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/zh0vtyj/allincecup-server/internal/domain/models"
	"github.com/zh0vtyj/allincecup-server/pkg/client/postgres"
)

type Storage interface {
	CreateUser(user User, code string) (int, string, error)
	GetUser(email string, password string) (User, error)
	NewSession(session models.Session) (models.Session, error)
	GetSessionByRefresh(refresh string) (models.Session, error)
	DeleteSessionByRefresh(refresh string) error
	DeleteSessionByUserId(id int) error
	UpdateRefreshToken(userId int, newRefreshToken string) error
	GetUserPasswordHash(userId int) (string, error)
	UpdatePassword(userId int, newPassword string) error
	UserExists(email string) (int, string, error)
	SelectUserInfo(id int) (InfoDTO, error)
	UpdatePersonalInfo(user InfoDTO, id int) error
	GetModerators(createdAt string, roleCode string) (moderators []User, err error)
	Delete(id int) error
}

type storage struct {
	db *sqlx.DB
	qb squirrel.StatementBuilderType
}

func NewAuthPostgres(db *sqlx.DB, qb squirrel.StatementBuilderType) *storage {
	return &storage{db: db, qb: qb}
}

func (s *storage) CreateUser(user User, code string) (int, string, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return 0, "", err
	}

	var userId int

	query := fmt.Sprintf(
		`
		INSERT INTO %s 
		(role_id, email, lastname, firstname, middle_name, password_hash, phone_number) 
		values ((SELECT id FROM roles WHERE code = $1), $2, $3, $4, $5, $6, $7) 
		RETURNING id
		`,
		postgres.UsersTable,
	)

	row := tx.QueryRow(
		query,
		code,
		user.Email,
		user.Lastname,
		user.Firstname,
		user.MiddleName,
		user.Password,
		user.PhoneNumber,
	)

	if err = row.Scan(&userId); err != nil {
		_ = tx.Rollback() // db rollback in error case
		return 0, "", err
	}

	// new user's cart query
	query = fmt.Sprintf("INSERT INTO %s (user_id) values ($1)", postgres.CartsTable)
	_, err = tx.Exec(query, userId)
	if err != nil {
		_ = tx.Rollback() // db rollback in error case
		return 0, "", err
	}

	// return id and transaction commit
	return userId, code, tx.Commit()
}

func (s *storage) GetUser(email, password string) (User, error) {
	var user User
	query := fmt.Sprintf(
		`
		SELECT users.id, roles.code as role_code
		FROM %s 
		JOIN %s ON users.role_id = roles.id 
		WHERE email = $1 AND password_hash = $2
		`,
		postgres.UsersTable,
		postgres.RolesTable,
	)

	err := s.db.Get(&user, query, email, password)

	return user, err
}

func (s *storage) NewSession(session models.Session) (models.Session, error) {
	queryDeleteOldSession := fmt.Sprintf("DELETE FROM %s WHERE user_id = $1", postgres.SessionsTable)
	_, err := s.db.Exec(queryDeleteOldSession, session.UserId)
	if err != nil {
		return models.Session{}, err
	}

	var newSession models.Session
	query := fmt.Sprintf(
		`
		INSERT INTO %s 
		(user_id, role_code, refresh_token, client_ip, user_agent, expires_at) 
		VALUES 
		($1, $2, $3, $4, $5, $6) 
		RETURNING id, user_id, role_code, refresh_token, client_ip, user_agent, is_blocked, expires_at, created_at
		`,
		postgres.SessionsTable,
	)

	row := s.db.QueryRow(
		query,
		session.UserId,
		session.RoleCode,
		session.RefreshToken,
		session.ClientIp,
		session.UserAgent,
		session.ExpiresAt,
	)

	if err = row.Scan(
		&newSession.Id,
		&newSession.UserId,
		&newSession.RoleCode,
		&newSession.RefreshToken,
		&newSession.ClientIp,
		&newSession.UserAgent,
		&newSession.IsBlocked,
		&newSession.ExpiresAt,
		&newSession.CreatedAt,
	); err != nil {
		return models.Session{}, err
	}

	return newSession, nil
}

func (s *storage) DeleteSessionByRefresh(refresh string) error {
	var id int
	query := fmt.Sprintf("DELETE from %s WHERE refresh_token = $1 RETURNING id", postgres.SessionsTable)
	row := s.db.QueryRow(query, refresh)

	err := row.Scan(&id)
	if err != nil {
		return err
	}
	return nil
}

func (s *storage) GetSessionByRefresh(refresh string) (models.Session, error) {
	var session models.Session
	queryGetSession := fmt.Sprintf(`
		SELECT 
			sessions.id,
			sessions.user_id,
			(SELECT code FROM roles WHERE users.role_id = roles.id) as role_code,
			sessions.refresh_token,
			sessions.client_ip,
			sessions.user_agent,
			sessions.is_blocked,
			sessions.expires_at,
			sessions.created_at
		FROM %s
		LEFT JOIN %s ON sessions.user_id = users.id 
		WHERE refresh_token = $1 LIMIT 1
		`,
		postgres.SessionsTable,
		postgres.UsersTable,
	)
	err := s.db.Get(&session, queryGetSession, refresh)

	if err != nil {
		return models.Session{}, fmt.Errorf("session wasn't found by refresh=%s, due to: %v", refresh, err)
	}

	return session, nil
}

func (s *storage) DeleteSessionByUserId(id int) error {
	queryDeleteSession := fmt.Sprintf("DELETE FROM %s WHERE user_id=$1", postgres.SessionsTable)
	_, err := s.db.Exec(queryDeleteSession, id)
	return err
}

func (s *storage) UpdateRefreshToken(userId int, newRefreshToken string) error {
	queryUpdateRefreshToken := fmt.Sprintf("UPDATE %s SET refresh_token=$1 WHERE user_id=$2", postgres.SessionsTable)
	_, err := s.db.Exec(queryUpdateRefreshToken, newRefreshToken, userId)
	if err != nil {
		return err
	}
	return nil
}

func (s *storage) GetUserPasswordHash(userId int) (string, error) {
	var hash string
	queryGetHash := fmt.Sprintf("SELECT password_hash FROM %s WHERE id=$1", postgres.UsersTable)

	err := s.db.Get(&hash, queryGetHash, userId)
	if err != nil {
		return "", fmt.Errorf("failed to get password hash due to: %v", err)
	}

	return hash, nil
}

func (s *storage) UpdatePassword(userId int, newPassword string) error {
	queryUpdatePassword := fmt.Sprintf("UPDATE %s SET password_hash=$1 WHERE id=$2", postgres.UsersTable)

	_, err := s.db.Exec(queryUpdatePassword, newPassword, userId)
	if err != nil {
		return fmt.Errorf("failed to update password in db due to: %v", err)
	}

	return nil
}

func (s *storage) UserExists(email string) (int, string, error) {
	var userId int
	var userRoleCode string

	queryGetUserId := fmt.Sprintf(
		`SELECT users.id, roles.code as role_code FROM %s JOIN %s ON users.role_id = roles.id WHERE email = $1`,
		postgres.UsersTable,
		postgres.RolesTable,
	)

	row := s.db.QueryRow(queryGetUserId, email)
	if err := row.Scan(&userId, userRoleCode); err != nil {
		return 0, "", err
	}

	return userId, userRoleCode, nil
}

func (s *storage) SelectUserInfo(id int) (user InfoDTO, err error) {
	querySelectUserInfo := fmt.Sprintf(
		`
		SELECT email, lastname, firstname, middle_name, phone_number 
		FROM %s 
		WHERE id = $1
		`,
		postgres.UsersTable,
	)

	err = s.db.Get(&user, querySelectUserInfo, id)
	if err != nil {
		return InfoDTO{}, fmt.Errorf("failed to select user info due to %v", err)
	}

	return user, err
}

func (s *storage) UpdatePersonalInfo(user InfoDTO, id int) error {
	queryUpdateUser := fmt.Sprintf(
		`
		UPDATE %s
		SET lastname = $1,
			firstname = $2,
			middle_name = $3,
			phone_number = $4
		WHERE id = $5
		`,
		postgres.UsersTable,
	)

	_, err := s.db.Exec(
		queryUpdateUser,
		user.Lastname,
		user.Firstname,
		user.MiddleName,
		user.PhoneNumber,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to update user info due to %v", err)
	}

	return nil
}

func (s *storage) GetModerators(createdAt string, roleCode string) (moderators []User, err error) {
	var moderatorsColumnsSelect = []string{
		"users.id",
		"users.email",
		"users.lastname",
		"users.firstname",
		"users.middle_name",
		"users.phone_number",
		"users.created_at",
	}

	querySelectModerators := s.qb.
		Select(moderatorsColumnsSelect...).
		From(postgres.UsersTable).
		LeftJoin(postgres.RolesTable + " ON users.role_id = roles.id").
		Where(squirrel.Eq{"roles.code": roleCode})

	if createdAt != "" {
		querySelectModerators = querySelectModerators.Where(squirrel.Lt{"users.created_at": createdAt})
	}

	querySelectModeratorsSql, args, err := querySelectModerators.
		OrderBy("users.created_at").
		Limit(12).
		ToSql()

	err = s.db.Select(&moderators, querySelectModeratorsSql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to select moderators due to %v", err)
	}

	return moderators, err
}

func (s *storage) Delete(id int) (err error) {
	queryDeleteUser := fmt.Sprintf("DELETE FROM %s WHERE id = $1", postgres.UsersTable)

	_, err = s.db.Exec(queryDeleteUser, id)
	if err != nil {
		return fmt.Errorf("failed to delete user %v", err)
	}

	return err
}
