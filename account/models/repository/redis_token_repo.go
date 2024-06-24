package repository

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"memory-app/account/models"
	"memory-app/account/models/apprerrors"
	"time"
)

type tokenRedisRepository struct {
	Redis *redis.Client
}

func (t tokenRedisRepository) SetRefreshToken(ctx context.Context, userId string, tokenId string, expiTime time.Duration) error {
	key := fmt.Sprintf("%s:%s", userId, tokenId)
	if err := t.Redis.Set(ctx, key, 0, expiTime).Err(); err != nil {
		log.Printf("Could not SET refresh token to redis for userID/prevTokenId: %s/%s: %v\n", userId, tokenId, err)
		return apprerrors.NewInternal()
	}
	return nil
}

func (t tokenRedisRepository) DeleteRefreshToken(ctx context.Context, userId string, prevTokenId string) error {
	key := fmt.Sprintf("%s:%s", userId, prevTokenId)
	if err := t.Redis.Del(ctx, key).Err(); err != nil {
		log.Printf("Could not delete refresh token to redis for userId/prevTokenId: %s/%s: %v\n", userId, prevTokenId, err)
		return apprerrors.NewInternal()
	}

	return nil
}

func NewTokenRepo(redis *redis.Client) models.TokenRepository {
	return &tokenRedisRepository{
		Redis: redis,
	}
}
