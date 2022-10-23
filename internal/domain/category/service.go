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
	GetFiltration(fkName string, id int) ([]Filtration, error)
	Update(category Category) (int, error)
	Create(dto CreateDTO) (int, error)
	Delete(id int) error
	AddFiltration(filtration Filtration) (int, error)
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
		minioPkg.ImagesBucket,
		imgUUID.String(),
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

	category := Category{
		CategoryTitle:       dto.CategoryTitle,
		ImgUrl:              dto.ImgUrl,
		ImgUUID:             &imgUUID,
		CategoryDescription: dto.CategoryDescription,
	}

	id, err := s.repo.Create(category)
	if err != nil {
		return 0, err
	}

	return id, err
}

func (s *service) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *service) AddFiltration(filtration Filtration) (int, error) {
	return s.repo.AddFiltration(filtration)
}

func (s *service) GetFiltration(fkName string, id int) ([]Filtration, error) {
	return s.repo.GetFiltration(fkName, id)
}
