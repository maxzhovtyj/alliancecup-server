package service

import (
	server "allincecup-server"
	"allincecup-server/pkg/repository"
)

type Authorization interface {
	CreateUser(user server.User) (int, error)
	GenerateToken(email string, password string) (string, error)
}

type ShopItemCup interface {
}

type ShopList interface {
}

type Category interface {
}

type Service struct {
	Authorization
	ShopItemCup
	ShopList
	Category
}

func NewService(repos *repository.Repository) *Service {
	return &Service{Authorization: NewAuthService(repos.Authorization)}
}
