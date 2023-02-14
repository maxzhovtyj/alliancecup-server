package handler

import (
	"fmt"
	"github.com/zh0vtyj/alliancecup-server/internal/domain/models"
	"mime/multipart"
)

func parseFile(file map[string][]*multipart.FileHeader) (*models.FileDTO, error) {
	files, ok := file["file"]

	if len(files) == 0 {
		return nil, ErrEmptyFile
	}

	if !ok {
		return nil, fmt.Errorf("something wrong with file you provided")
	}

	fileInfo := files[0]
	fileReader, err := fileInfo.Open()
	if err != nil {
		return nil, err
	}

	if fileInfo.Size > fileMaxSize {
		return nil, fmt.Errorf("file size exceeded, max size is 5MB")
	}

	return &models.FileDTO{
		Name:   fileInfo.Filename,
		Size:   fileInfo.Size,
		Reader: fileReader,
	}, nil

}
