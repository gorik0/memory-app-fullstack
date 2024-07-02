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

func (t tokenRedisRepository) DeleteUserRefreshToken(ctx context.Context, userId string) error {

	pattern := fmt.Sprintf("%s*", userId)
	println(pattern)
	iter := t.Redis.Scan(ctx, 0, pattern, 5).Iterator()

	failCount := 0
	for iter.Next(ctx) {
		if err := t.Redis.Del(ctx, iter.Val()).Err(); err != nil {
			fmt.Println("Error while deleting user refresh token ::: ", err)
			failCount++
		}
	}

	if err := iter.Err(); err != nil {
		fmt.Println("Error while deleting user refresh token ::: ", err)
		return err

	}
	if failCount > 0 {
		return apprerrors.NewInternal()
	}
	return nil

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
	_ = key
	println("!!!!!! DELETE")
	var err error
	var res int64
	if res, err = t.Redis.Del(ctx, key).Result(); res == 0 || err != nil {
		log.Printf("Could not delete refresh token to redis for userId/prevTokenId: %s/%s: %v\n", userId, prevTokenId, err)
		return apprerrors.NewInternal()
	}
	println("!!!!  ", res)
	return nil
}

func NewTokenRepo(redis *redis.Client) models.TokenRepository {
	return &tokenRedisRepository{
		Redis: redis,
	}
}
