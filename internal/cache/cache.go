package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Cache defines the operations the auth service needs from the token store.
// Implementations are expected to be backed by Redis.
type Cache interface {
	// StoreRefreshToken persists a refresh token associated with a user, expiring after ttl.
	StoreRefreshToken(ctx context.Context, userID, token string, ttl time.Duration) error

	// DeleteRefreshToken removes a refresh token (e.g. on logout or rotation).
	DeleteRefreshToken(ctx context.Context, token string) error

	// IsTokenBlacklisted reports whether a token has been explicitly revoked.
	IsTokenBlacklisted(ctx context.Context, token string) (bool, error)
}

// redisCache implements Cache using a Redis backend.
type redisCache struct {
	// TODO: add a real Redis client here once the dependency is wired.
	client *redis.Client
}

// NewRedisCache returns a Cache backed by Redis.
// Pass the Redis client once the connection is initialised in main.
//
//	rdb := redis.NewClient(&redis.Options{Addr: cfg.RedisURL})
//	cacheSvc := cache.NewRedisCache(rdb)
func NewRedisCache(client *redis.Client) Cache {
	return &redisCache{
		// client: client,
	}
}

// StoreRefreshToken stores the token in Redis under the key "refresh:<token>"
// with the given TTL so it expires automatically.
func (c *redisCache) StoreRefreshToken(ctx context.Context, userID, token string, ttl time.Duration) error {
	// TODO: implement
	// key := "refresh:" + token
	// return c.client.Set(ctx, key, userID, ttl).Err()
	return nil
}

// DeleteRefreshToken removes the key "refresh:<token>" from Redis.
func (c *redisCache) DeleteRefreshToken(ctx context.Context, token string) error {
	// TODO: implement
	// key := "refresh:" + token
	// return c.client.Del(ctx, key).Err()
	return nil
}

// IsTokenBlacklisted checks for the presence of the key "blacklist:<token>".
// A token is blacklisted when it is revoked before its natural expiry (e.g. logout).
func (c *redisCache) IsTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	// TODO: implement
	// key := "blacklist:" + token
	// n, err := c.client.Exists(ctx, key).Result()
	// if err != nil {
	// 	return false, fmt.Errorf("check blacklist: %w", err)
	// }
	// return n > 0, nil
	return false, nil
}
