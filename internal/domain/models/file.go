package models

import "mime/multipart"

type FileDTO struct {
	Name   string
	Size   int64
	Reader multipart.File
}
