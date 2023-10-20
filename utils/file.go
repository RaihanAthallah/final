package utils

import (
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

func UploadFile(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {

	if !SizeValidator(*fileHeader) {
		return "", errors.New("File size too large")
	}

	fileName := fileHeader.Filename

	if !ExtensionValidator(fileName) {
		return "", errors.New("Invalid file extension")
	}

	uploadFolder := "client/uploads"
	err := os.MkdirAll(uploadFolder, os.ModePerm)
	if err != nil {
		return "", errors.New("Failed to create upload folder")
	}

	var filename = RandString(32) + filepath.Ext(fileHeader.Filename)
	imagePath := filepath.Join(uploadFolder, filename)
	imageFile, err := os.Create(imagePath)
	if err != nil {
		return "", errors.New("Failed to create image file")
	}

	defer imageFile.Close()

	_, err = io.Copy(imageFile, file)
	if err != nil {
		return "", errors.New("Failed to copy image")
	}

	return filename, nil
}

func UploadFilePDF(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {

	if !SizeValidator(*fileHeader) {
		return "", errors.New("File size too large")
	}

	fileName := fileHeader.Filename

	if !ExtensionValidatorPDF(fileName) {
		return "", errors.New("Invalid file extension")
	}

	uploadFolder := "client/uploads"
	err := os.MkdirAll(uploadFolder, os.ModePerm)
	if err != nil {
		return "", errors.New("Failed to create upload folder")
	}

	var filename = RandString(32) + filepath.Ext(fileHeader.Filename)
	imagePath := filepath.Join(uploadFolder, filename)
	imageFile, err := os.Create(imagePath)
	if err != nil {
		return "", errors.New("Failed to create image file")
	}

	defer imageFile.Close()

	_, err = io.Copy(imageFile, file)
	if err != nil {
		return "", errors.New("Failed to copy image")
	}

	return filename, nil
}

func ExtensionValidator(filename string) bool {
	ext := filepath.Ext(filename)
	allowedExtensions := []string{".jpg", ".jpeg", ".png"}
	validExtension := false
	for _, allowedExt := range allowedExtensions {
		if ext == allowedExt {
			validExtension = true
			break
		}
	}

	return validExtension
}

func ExtensionValidatorPDF(filename string) bool {
	ext := filepath.Ext(filename)
	allowedExtensions := []string{".pdf"}
	validExtension := false
	for _, allowedExt := range allowedExtensions {
		if ext == allowedExt {
			validExtension = true
			break
		}
	}

	return validExtension
}

func SizeValidator(file multipart.FileHeader) bool {
	size := file.Size
	maxSize := int64(1024 * 1024 * 2) // 2MB

	if size > maxSize {
		return false
	}

	return true
}
