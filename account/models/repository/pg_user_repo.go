package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"log"
	"memory-app/account/models"
	"memory-app/account/models/apprerrors"
)

type PgUSerRepo struct {
	DB *sqlx.DB
}

func (p PgUSerRepo) UpdateImage(ctx context.Context, uid uuid.UUID, imageURL string) (*models.User, error) {
	query := `UPDATE users SET image_url = $1 WHERE uid = $2 returning *`

	user := &models.User{}
	err := p.DB.GetContext(ctx, user, query, imageURL, uid)
	if err != nil {
		fmt.Println("Failed to update user", err.Error())
		return nil, apprerrors.NewInternal()

	}

	return user, nil

}

func (p PgUSerRepo) Update(ctx context.Context, u *models.User) error {
	query := `UPDATE users SET name=:name, email=:email where uid=:uid returning *`
	fmt.Println("!!!!!!!", u)

	nstmt, err := p.DB.PrepareNamedContext(ctx, query)
	if err != nil {
		log.Printf("Error preparing query: %v", err)
		return apprerrors.NewInternal()
	}
	err = nstmt.GetContext(ctx, u, u)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return apprerrors.NewInternal()
	}
	fmt.Println("222!!!!!!!", u)
	return nil
}

func (p PgUSerRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user = &models.User{}

	query := `SELECT * from users where email=$1`
	println(email)
	err := p.DB.GetContext(ctx, user, query, email)
	if err != nil {

		log.Printf("Errro while selecting useer by email ::: %s", err.Error())

		return nil, apprerrors.NewNotFound("email", email)

	}
	return user, nil

}

func (p PgUSerRepo) GetById(ctx context.Context, uid uuid.UUID) (*models.User, error) {
	query := `SELECT * from   users where uid = $1 `

	var user = new(models.User)

	if err := p.DB.GetContext(ctx, user, query, uid); err != nil {
		log.Printf("Couldn't get user for id  ::: %v with error ::: %v", uid, err.Error())

		return nil, apprerrors.NewNotFound("uid", uid.String())
	}
	return user, nil

}

func (p PgUSerRepo) Create(ctx context.Context, u *models.User) error {
	query := `INSERT into users  (email,password) values ($1,$2) returning *`

	err := p.DB.GetContext(ctx, u, query, u.Email, u.Password)
	if err != nil {
		if pqError, ok := err.(*pq.Error); ok {
			if pqError.Code.Name() == "unique_violation" {
				log.Printf("Couldn't create user, already exist ::: %v with error ::: %v", u, err.Error())
				return apprerrors.NewConflict("email", u.Email)
			}
		}
		log.Printf("Couldn't create user with error ::: %v", err.Error())

		return apprerrors.NewInternal()
	}
	return nil
}

var _ models.UserRepositoryI = &PgUSerRepo{}
