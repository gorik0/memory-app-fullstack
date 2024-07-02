package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"log"
	"memory-app/account/models"
	"memory-app/account/models/apprerrors"
)

type UserService struct {
	UserRepository models.UserRepositoryI
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
	UserRepo models.UserRepositoryI
}

func NewUserService(cfg *UserServiceConfig) models.UserServiceI {
	return &UserService{
		UserRepository: cfg.UserRepo,
	}
}

func (u *UserService) Get(ctx context.Context, uid uuid.UUID) (*models.User, error) {
	return u.UserRepository.GetById(ctx, uid)

}

var _ models.UserServiceI = &UserService{}
