package product

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	server "github.com/zh0vtyj/allincecup-server/internal/domain/shopping"
)

type Service interface {
	Search(searchInput string) ([]Product, error)
	GetWithParams(params server.SearchParams) ([]Product, error)
	GetProductById(id int) (Product, error)
	AddProduct(product CreateProductDTO) (int, error)
	GetFavourites(userId int) ([]Product, error)
	Update(product Product) (int, error)
	Delete(productId int) error
}

type service struct {
	repo        Storage
	fileStorage *minio.Client
}

func NewProductsService(repo Storage, fileStorage *minio.Client) Service {
	return &service{
		repo:        repo,
		fileStorage: fileStorage,
	}
}

func (s *service) Search(searchInput string) ([]Product, error) {
	searchInput = "%" + searchInput + "%"
	products, err := s.repo.Search(searchInput)
	if err != nil {
		return nil, err
	}
	return products, err
}

func (s *service) GetWithParams(params server.SearchParams) ([]Product, error) {
	return s.repo.GetWithParams(params)
}

func (s *service) AddProduct(dto CreateProductDTO) (int, error) {
	imgUUID := uuid.New()

	exists, errBucketExists := s.fileStorage.BucketExists(context.Background(), "images")
	if errBucketExists != nil || !exists {
		err := s.fileStorage.MakeBucket(context.Background(), "images", minio.MakeBucketOptions{})
		if err != nil {
			return 0, fmt.Errorf("failed to create new bucket. err: %w", err)
		}
	}

	_, err := s.fileStorage.PutObject(
		context.Background(),
		"images",
		imgUUID.String(),
		dto.Img.Reader,
		dto.Img.Size,
		minio.PutObjectOptions{
			UserMetadata: map[string]string{
				"Name": dto.ProductTitle,
			},
			ContentType: "application/octet-stream",
		},
	)
	if err != nil {
		return 0, err
	}

	product := Product{
		Article:         dto.Article,
		CategoryTitle:   dto.CategoryTitle,
		ProductTitle:    dto.ProductTitle,
		AmountInStock:   dto.AmountInStock,
		ImgUUID:         imgUUID,
		Price:           dto.Price,
		Characteristics: dto.Characteristics,
		Packaging:       dto.Packaging,
		Description:     dto.Description,
	}

	return s.repo.Create(product)
}

func (s *service) GetFavourites(userId int) ([]Product, error) {
	return s.repo.GetFavourites(userId)
}

func (s *service) Update(product Product) (int, error) {
	return s.repo.Update(product)
}

func (s *service) Delete(productId int) error {
	return s.repo.Delete(productId)
}

func (s *service) GetProductById(id int) (Product, error) {
	return s.repo.GetProductById(id)
}
