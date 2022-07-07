package service

import (
	server "allincecup-server"
	"allincecup-server/pkg/repository"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"math/rand"
	"time"
)

const (
	salt       = "dsadkasdi212312mdmacmxz00"
	tokenTTL   = 15 * time.Minute
	signingKey = "das345=FF@!a;212&&dsDFCwW12e112d%#d$c"
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
func (s *AuthService) CreateUser(user server.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
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
	b := make([]byte, 32)

	src := rand.NewSource(time.Now().Unix())
	r := rand.New(src)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
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

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
