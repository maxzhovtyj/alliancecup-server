package minioPkg

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/zh0vtyj/allincecup-server/internal/config"
)

func NewClient(cfg config.MinIO) (*minio.Client, error) {
	client, err := minio.New(
		cfg.Endpoint,
		&minio.Options{
			Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
			Secure: false,
		},
	)
	if err != nil {
		return nil, err
	}

	return client, err
}
