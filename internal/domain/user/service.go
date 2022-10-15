package user

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/zh0vtyj/allincecup-server/internal/domain/models"
	"time"
)

const (
	salt = "dsadkasdi212312mdmacmxz00"
	//tokenTTL          = 30 * time.Minute // RELEASE VERSION
	tokenTTL          = 120 * time.Minute
	signingKey        = "das345=FF@!a;212&&dsDFCwW12e112d%#d$c"
	refreshTokenTTL   = 1440 * time.Hour
	refreshSigningKey = "Sepasd213*99921@@#dsad+-=SXxassd@lLL;"
	clientRole        = "CLIENT"
	moderatorRole     = "MODERATOR"
)

type AuthorizationService interface {
	CreateUser(user User) (int, int, error)
	CreateModerator(user User) (int, int, error)
	GenerateTokens(email string, password string) (string, string, error)
	ParseToken(token string) (int, int, error)
	ParseRefreshToken(refreshToken string) error
	RefreshTokens(refreshToken, clientIp, userAgent string) (string, string, int, int, error)
	CreateNewSession(session *models.Session) (*models.Session, error)
	Logout(id int) error
	ChangePassword(userId int, oldPassword, newPassword string) error
	UserForgotPassword(email string) error
}

type tokenClaims struct {
	jwt.StandardClaims
	UserId     int `json:"user_id"`
	UserRoleId int `json:"user_role_id"`
}

type AuthService struct {
	repo AuthorizationStorage
}

func NewAuthService(repo AuthorizationStorage) AuthorizationService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user User) (int, int, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user, clientRole)
}

func (s *AuthService) CreateModerator(user User) (int, int, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user, moderatorRole)
}

func (s *AuthService) GenerateTokens(email, password string) (string, string, error) {
	selectedUser, err := s.repo.GetUser(email, generatePasswordHash(password))
	if err != nil {
		return "", "", fmt.Errorf("user are not found")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		selectedUser.Id,
		selectedUser.RoleId,
	})

	refreshToken, err := s.GenerateRefreshToken()
	if err != nil {
		return "", "", err
	}

	accessToken, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) GenerateRefreshToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(refreshTokenTTL).Unix(),
		IssuedAt:  time.Now().Unix(),
	})

	refreshToken, err := token.SignedString([]byte(refreshSigningKey))
	if err != nil {
		return "", err
	}
	return refreshToken, err
}

func (s *AuthService) ParseToken(accessToken string) (int, int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, claims.UserRoleId, nil
}

func (s *AuthService) ParseRefreshToken(refreshToken string) error {
	token, err := jwt.ParseWithClaims(refreshToken, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(refreshSigningKey), nil
	})
	if err != nil {
		return err
	}

	_, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		return errors.New("token claims are not of type *StandardClaims")
	}

	return nil
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func (s *AuthService) RefreshTokens(refreshToken, clientIp, userAgent string) (string, string, int, int, error) {
	// get user session by old refresh token
	session, err := s.repo.GetSessionByRefresh(refreshToken)
	if err != nil {
		return "", "", 0, 0, err
	}

	// validation if client IP or user agent is not the same
	if session.ClientIp != clientIp || session.UserAgent != userAgent {
		err = s.repo.DeleteSessionByRefresh(session.RefreshToken)
		if err != nil {
			return "", "", 0, 0, fmt.Errorf("cannot delete session: " + err.Error())
		}
		return "", "", 0, 0, fmt.Errorf("invalid meta data")
	}

	// validation if refresh token is expired
	if time.Now().After(session.ExpiresAt) {
		err = s.repo.DeleteSessionByRefresh(session.RefreshToken)
		if err != nil {
			return "", "", 0, 0, fmt.Errorf("cannot delete session: " + err.Error())
		}
		return "", "", 0, 0, fmt.Errorf("refresh expired token, session deleted from db")
	}

	// validation if refresh token is blocked
	if session.IsBlocked {
		err = s.repo.DeleteSessionByRefresh(session.RefreshToken)
		if err != nil {
			return "", "", 0, 0, fmt.Errorf("cannot delete session: " + err.Error())
		}
		return "", "", 0, 0, fmt.Errorf("session is blocked, session deleted from db")
	}

	newRefreshToken, err := s.GenerateRefreshToken()

	err = s.repo.UpdateRefreshToken(session.UserId, newRefreshToken)
	if err != nil {
		return "", "", 0, 0, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		session.UserId,
		session.RoleId,
	})
	accessToken, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return "", "", 0, 0, err
	}

	return accessToken, newRefreshToken, session.UserId, session.RoleId, err
}

func (s *AuthService) CreateNewSession(session *models.Session) (*models.Session, error) {
	newSession, err := s.repo.NewSession(*session)
	if err != nil {
		return nil, err
	}
	return newSession, err
}

func (s *AuthService) Logout(id int) error {
	return s.repo.DeleteSessionByUserId(id)
}

func (s *AuthService) ChangePassword(userId int, oldPassword, newPassword string) error {
	hash, err := s.repo.GetUserPasswordHash(userId)
	if err != nil {
		return err
	}

	passwordHash := generatePasswordHash(oldPassword)

	if hash != passwordHash {
		return fmt.Errorf("invalid input password")
	}

	newPasswordHash := generatePasswordHash(newPassword)

	err = s.repo.UpdatePassword(userId, newPasswordHash)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) UserForgotPassword(email string) error {
	// check whether user with such email exists
	userId, userRoleId, err := s.repo.UserExists(email)
	if err != nil {
		return fmt.Errorf("failed to get user with email %s due to %v", email, err)
	}

	// generate a token for changing a password
	_ = jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		userId,
		userRoleId,
	})

	// TODO send a letter to an email

	return nil
}
