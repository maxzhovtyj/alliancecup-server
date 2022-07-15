package service

import (
	server "allincecup-server"
	"allincecup-server/internal/domain"
	"allincecup-server/pkg/repository"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
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

type tokenClaims struct {
	jwt.StandardClaims
	UserId     int `json:"user_id"`
	UserRoleId int `json:"user_role_id"`
}

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}
func (s *AuthService) CreateUser(user server.User) (int, int, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user, clientRole)
}

func (s *AuthService) CreateModerator(user server.User) (int, int, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user, moderatorRole)
}

func (s *AuthService) GenerateTokens(email, password string) (string, string, error) {
	user, err := s.repo.GetUser(email, generatePasswordHash(password))
	if err != nil {
		return "", "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
		user.RoleId,
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
	},
	)

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

func (s *AuthService) RefreshAccessToken(refreshToken string) (string, error) {
	session, err := s.repo.GetSessionByRefresh(refreshToken)

	if err != nil {
		return "", err
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
		return "", err
	}

	return accessToken, nil
}

func (s *AuthService) CreateNewSession(session *domain.Session) (*domain.Session, error) {
	newSession, err := s.repo.NewSession(*session)
	if err != nil {
		return nil, err
	}
	return newSession, err
}
