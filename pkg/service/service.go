package service

import "allincecup-server/pkg/repository"

type Authorization interface {
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

func NewService(repository *repository.Repository) *Service {
	return &Service{}
}
