package repository

type Authorization interface {
}

type ShopItemCup interface {
}

type ShopList interface {
}

type Category interface {
}

type Repository struct {
	Authorization
	ShopItemCup
	ShopList
	Category
}

func NewRepository() *Repository {
	return &Repository{}
}
