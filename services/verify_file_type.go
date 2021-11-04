package services

import (
	"mime/multipart"

	"github.com/anti-lgbt/medusa/types"
	"github.com/h2non/filetype"
)

func VerifyFileType(file_header *multipart.FileHeader, file_type types.FileType) bool {
	file, _ := file_header.Open()
	defer file.Close()
	buf := make([]byte, 512)
	file.Read(buf)

	switch file_type {
	case types.FileTypeImage:
		return filetype.IsImage(buf)
	case types.FileTypeAudio:
		return filetype.IsAudio(buf)
	}

	return false
}
