package product

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	server "github.com/zh0vtyj/allincecup-server/internal/domain/shopping"
	minioPkg "github.com/zh0vtyj/allincecup-server/pkg/client/minio"
)

type Service interface {
	Search(searchInput string) ([]Product, error)
	GetWithParams(params server.SearchParams) ([]Product, error)
	GetProductById(id int) (Product, error)
	Add(product CreateDTO) (int, error)
	GetFavourites(userId int) ([]Product, error)
	Update(product Product) (int, error)
	UpdateImage(dto UpdateImageDTO) (int, error)
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

func (s *service) Add(dto CreateDTO) (int, error) {
	var imgUUIDPtr *uuid.UUID
	if dto.Img != nil {
		imgUUID := uuid.New()
		imgUUIDPtr = &imgUUID
	}

	product := Product{
		Article:         dto.Article,
		CategoryTitle:   &dto.CategoryTitle,
		ProductTitle:    dto.ProductTitle,
		AmountInStock:   dto.AmountInStock,
		ImgUUID:         imgUUIDPtr,
		Price:           dto.Price,
		Characteristics: dto.Characteristics,
		Packaging:       dto.Packaging,
		Description:     dto.Description,
	}

	id, err := s.repo.Create(product)
	if err != nil {
		return 0, err
	}

	if imgUUIDPtr != nil {
		exists, errBucketExists := s.fileStorage.BucketExists(context.Background(), minioPkg.ImagesBucket)
		if errBucketExists != nil || !exists {
			err = s.fileStorage.MakeBucket(context.Background(), "images", minio.MakeBucketOptions{})
			if err != nil {
				return 0, fmt.Errorf("failed to create new bucket. err: %w", err)
			}
		}

		_, err = s.fileStorage.PutObject(
			context.Background(),
			minioPkg.ImagesBucket,
			imgUUIDPtr.String(),
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
	}

	return id, err
}

func (s *service) GetFavourites(userId int) ([]Product, error) {
	return s.repo.GetFavourites(userId)
}

func (s *service) Update(product Product) (int, error) {
	return s.repo.Update(product)
}

func (s *service) UpdateImage(dto UpdateImageDTO) (int, error) {
	var imgUUIDPtr *uuid.UUID
	if dto.Img != nil {
		imgUUID := uuid.New()
		imgUUIDPtr = &imgUUID
	} else {
		return 0, nil
	}

	oldProduct, err := s.repo.GetProductById(dto.Id)
	if err != nil {
		return 0, err
	}

	_, err = s.fileStorage.PutObject(
		context.Background(),
		minioPkg.ImagesBucket,
		imgUUIDPtr.String(),
		dto.Img.Reader,
		dto.Img.Size,
		minio.PutObjectOptions{},
	)
	if err != nil {
		return 0, err
	}

	if oldProduct.ImgUUID != nil {
		err = s.fileStorage.RemoveObject(
			context.Background(),
			minioPkg.ImagesBucket,
			oldProduct.ImgUUID.String(),
			minio.RemoveObjectOptions{},
		)
		if err != nil {
			return 0, fmt.Errorf("failed to remove old product image due to %v", err)
		}
	}

	id, err := s.repo.UpdateImage(Product{
		Id:      dto.Id,
		ImgUUID: imgUUIDPtr,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to update product image uuid")
	}

	return id, err
}

func (s *service) Delete(productId int) error {
	product, err := s.repo.GetProductById(productId)
	if err != nil {
		return err
	}

	err = s.repo.Delete(productId)
	if err != nil {
		return fmt.Errorf("failed to delete product %d due to %v", productId, err)
	}

	if product.ImgUUID != nil {
		err = s.fileStorage.RemoveObject(
			context.Background(),
			minioPkg.ImagesBucket,
			product.ImgUUID.String(),
			minio.RemoveObjectOptions{},
		)
		if err != nil {
			return fmt.Errorf("failed to remove object")
		}
	}

	return err
}

func (s *service) GetProductById(id int) (Product, error) {
	return s.repo.GetProductById(id)
}
