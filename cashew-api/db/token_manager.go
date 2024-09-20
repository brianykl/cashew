package db

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type TokenManager interface {
	StoreToken(userId, accessToken string, expiration time.Duration) error
	GetTokens(userId string) ([]string, error)
	DeleteTokens(userId string) error
}

type redisTokenManager struct {
	client *redis.Client
}

func NewTokenManager(addr string) (TokenManager, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return &redisTokenManager{client: client}, nil
}

func (rtm *redisTokenManager) DeleteTokens(userId string) error {
	ctx := context.Background()
	key := fmt.Sprintf("plaid_tokens:%s", userId)

	_, err := rtm.client.SRem(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("failed to delete token: %v", err)
	}
	return nil
}

func (rtm *redisTokenManager) GetTokens(userId string) ([]string, error) {
	ctx := context.Background()
	key := fmt.Sprintf("plaid_tokens:%s", userId)

	tokens, err := rtm.client.SMembers(ctx, key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("token not found for user %s", userId)
	} else if err != nil {
		return nil, fmt.Errorf("failed to retrieve token: %v", err)
	}

	return tokens, nil
}

func (rtm *redisTokenManager) StoreToken(userId string, accessToken string, expiration time.Duration) error {
	ctx := context.Background()
	key := fmt.Sprintf("plaid_tokens:%s", userId)

	err := rtm.client.SAdd(ctx, key, accessToken).Err()
	if err != nil {
		return fmt.Errorf("failed to store token: %v", err)
	}

	return nil
}
