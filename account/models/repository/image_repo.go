package repository

import (
	"context"
	"io"
	"memory-app/account/models"
	"mime/multipart"
	"os"
)

type ImageRepository struct {
	PathToFile string
}

func (i *ImageRepository) DeleteProfile(ctx context.Context, objName string) error {

	err := os.Remove(i.PathToFile + "/" + objName)
	if err != nil {

		return err
	}
	return nil

}

func (i *ImageRepository) UpdateProfile(ctx context.Context, objName string, imageFile multipart.File) (string, error) {

	pathToFile := i.PathToFile + objName
	fileToWrite, err := os.OpenFile(pathToFile, os.O_WRONLY|os.O_CREATE, 0666)
	println("!!!!!!!", pathToFile)
	if err != nil {
		println("*while open*")
		return "", err

	}
	_, err = io.Copy(fileToWrite, imageFile)
	if err != nil {
		println("*while copy*")
		return "", err
	}

	return objName, nil

}

var _ models.ImageRepository = &ImageRepository{}
