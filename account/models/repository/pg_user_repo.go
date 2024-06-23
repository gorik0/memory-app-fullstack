package repository

import (
	"context"
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

func (p PgUSerRepo) GetById(ctx context.Context, uid uuid.UUID) (*models.User, error) {
	query := `SELECT * from   users where id = $1 `

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
