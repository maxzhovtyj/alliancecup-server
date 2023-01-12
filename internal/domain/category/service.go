package category

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	minioPkg "github.com/zh0vtyj/allincecup-server/pkg/client/minio"
)

type Service interface {
	GetAll() ([]Category, error)
	Update(category Category) (int, error)
	Create(dto CreateDTO) (int, error)
	Delete(id int) error
	DeleteFiltration(id int) error
	GetFiltration(fkName string, id int) ([]Filtration, error)
	GetFiltrationItems() ([]Filtration, error)
	AddFiltration(dto CreateFiltrationDTO) (int, error)
}

type service struct {
	repo        Storage
	fileStorage *minio.Client
}

func NewCategoryService(repo Storage, fileStorage *minio.Client) Service {
	return &service{
		repo:        repo,
		fileStorage: fileStorage,
	}
}

func (s *service) GetAll() ([]Category, error) {
	return s.repo.GetAll()
}

func (s *service) Update(category Category) (int, error) {
	return s.repo.Update(category)
}

func (s *service) Create(dto CreateDTO) (int, error) {
	var imgUUIDPtr *uuid.UUID

	if dto.Img != nil {
		imgUUID := uuid.New()
		imgUUIDPtr = &imgUUID
	}

	category := Category{
		CategoryTitle: dto.CategoryTitle,
		ImgUrl:        dto.ImgUrl,
		ImgUUID:       imgUUIDPtr,
		Description:   dto.CategoryDescription,
	}

	id, err := s.repo.Create(category)
	if err != nil {
		return 0, err
	}

	if imgUUIDPtr != nil {
		exists, errBucketExists := s.fileStorage.BucketExists(context.Background(), minioPkg.ImagesBucket)
		if errBucketExists != nil || !exists {
			err := s.fileStorage.MakeBucket(context.Background(), "images", minio.MakeBucketOptions{})
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
					"Name": dto.CategoryTitle,
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

func (s *service) Delete(id int) error {
	category, err := s.repo.GetById(id)
	if err != nil {
		return err
	}

	err = s.repo.Delete(id)
	if err != nil {
		return fmt.Errorf("failed to delete category due to %v", err)
	}

	if category.ImgUUID != nil {
		err = s.fileStorage.RemoveObject(
			context.Background(),
			minioPkg.ImagesBucket,
			category.ImgUUID.String(),
			minio.RemoveObjectOptions{},
		)
		if err != nil {
			return fmt.Errorf("failed to remove image object due to %v", err)
		}
	}

	return err
}

func (s *service) DeleteFiltration(id int) error {
	item, err := s.repo.GetFiltrationItem(id)
	if err != nil {
		return err
	}

	err = s.repo.DeleteFiltration(id)
	if err != nil {
		return err
	}

	if item.ImgUUID != nil {
		err = s.fileStorage.RemoveObject(
			context.Background(),
			minioPkg.ImagesBucket,
			item.ImgUUID.String(),
			minio.RemoveObjectOptions{},
		)
		if err != nil {
			return err
		}
	}

	return err
}

func (s *service) AddFiltration(dto CreateFiltrationDTO) (int, error) {
	var imgUUIDPtr *uuid.UUID
	if dto.Img != nil {
		imgUUID := uuid.New()
		imgUUIDPtr = &imgUUID
	}

	id, err := s.repo.AddFiltration(Filtration{
		CategoryId:            dto.CategoryId,
		ImgUrl:                dto.ImgUrl,
		ImgUUID:               imgUUIDPtr,
		SearchKey:             dto.SearchKey,
		SearchCharacteristic:  dto.SearchCharacteristic,
		FiltrationTitle:       dto.FiltrationTitle,
		FiltrationDescription: dto.FiltrationDescription,
		FiltrationListId:      dto.FiltrationListId,
	})
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
					"Name": dto.FiltrationTitle,
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

func (s *service) GetFiltration(fkName string, id int) ([]Filtration, error) {
	return s.repo.GetFiltration(fkName, id)
}

func (s *service) GetFiltrationItems() ([]Filtration, error) {
	return s.repo.GetFiltrationItems()
}
