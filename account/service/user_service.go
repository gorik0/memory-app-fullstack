package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"log"
	"memory-app/account/models"
	"memory-app/account/models/apprerrors"
	"mime/multipart"
	"strings"
)

type UserService struct {
	UserRepository  models.UserRepositoryI
	ImageRepository models.ImageRepository
}

func (u *UserService) ClearProfileImage(ctx context.Context, uid uuid.UUID) error {

	//	::GET USER
	user, err := u.UserRepository.GetById(ctx, uid)
	if err != nil {
		fmt.Println("Error getting user id:", err.Error())
		return err
	}
	//	::: GET IMAGE NAME
	imageNameToDelete := user.ImageURL
	if imageNameToDelete == "" {
		return nil
	}
	//	:: DELETE ON IMAGE REPO
	err = u.ImageRepository.DeleteProfile(ctx, imageNameToDelete)
	if err != nil {
		fmt.Println("Failed to delete user id:", err.Error())
		return err

	}
	//	:: UPDATE USER REPO
	user, err = u.UserRepository.UpdateImage(ctx, uid, "")

	if err != nil {

		fmt.Println("failed to update image in repository:", err.Error())
		return err

	}

	return nil
}

func (u *UserService) SetProfileImage(ctx context.Context, uid uuid.UUID, imageFileHeader *multipart.FileHeader) (*models.User, error) {

	//	:::GET USER

	user, err := u.UserRepository.GetById(ctx, uid)
	if err != nil {
		fmt.Println("Error getting user id:", err.Error())
		return nil, err
	}
	//	:::CREATE IMAGE NAME (NEW OR EXISTED BASED ON)

	fromUrl, err := objNameFromUrl(user.ImageURL)
	if err != nil {
		return nil, err
	}
	imageExtension := strings.Split(imageFileHeader.Filename, ".")[1]
	imageFileNameWithExtension := fmt.Sprintf("%s.%s", fromUrl, imageExtension)
	//	:::OPEN FILE from imageFileHeader
	file, err := imageFileHeader.Open()

	if err != nil {
		fmt.Println("Failed to open file:", err.Error())
		return nil, err
	}
	//	::: Upload to IMAGE REPO
	imageName, err := u.ImageRepository.UpdateProfile(ctx, imageFileNameWithExtension, file)
	if err != nil {
		fmt.Println("failed to update image:", err.Error())
		return nil, err

	}
	_ = imageName
	//	::: Upload to USER REPO
	user, err = u.UserRepository.UpdateImage(ctx, uid, imageName)
	if err != nil {

		fmt.Println("failed to update image in repository:", err.Error())
		return nil, err

	}

	return user, nil
}

func objNameFromUrl(urlName string) (string, error) {
	//if urlName == "" {
	//	return uuid.New().String(), nil
	//
	//}
	//parsed, err := url.Parse(urlName)
	//if err != nil {
	//	fmt.Println("Error parsing url:", err.Error())
	//	return "", err
	//
	//}
	//
	//return path.Base(parsed.Path), nil
	return uuid.New().String(), nil
}

func (u *UserService) UpdateDetail(ctx context.Context, user *models.User) error {

	return u.UserRepository.Update(ctx, user)
}

func (u *UserService) Signin(ctx context.Context, user *models.User) error {
	userGetted, err := u.UserRepository.GetByEmail(ctx, user.Email)
	if err != nil {
		log.Printf("Unnable get user  : %v\n", user)
		return err
	}

	match, err := Compare(userGetted.Password, user.Password)
	if err != nil {

		e := apprerrors.NewAuthorization("Invalid passsword")
		log.Printf("Invalid password!!!:::%s", err)
		return e
	}
	if !match {
		e := apprerrors.NewAuthorization("Invalid passsword (doesn't match")
		log.Printf("Passwords doesn't match!!!")
		return e
	}

	*user = *userGetted

	return nil

}

func (u *UserService) Signup(context context.Context, user *models.User) error {
	hash, err := generateHashPassword(user.Password)
	if err != nil {
		e := fmt.Errorf("Unable generatePassword for user : %v\n", user)
		log.Printf(e.Error())
		return e

	}
	user.Password = hash

	err = u.UserRepository.Create(context, user)
	if err != nil {
		return err
	}
	return nil
}

type UserServiceConfig struct {
	UserRepo  models.UserRepositoryI
	ImageRepo models.ImageRepository
}

func NewUserService(cfg *UserServiceConfig) models.UserServiceI {
	return &UserService{
		UserRepository:  cfg.UserRepo,
		ImageRepository: cfg.ImageRepo,
	}
}

func (u *UserService) Get(ctx context.Context, uid uuid.UUID) (*models.User, error) {
	return u.UserRepository.GetById(ctx, uid)

}

var _ models.UserServiceI = &UserService{}
