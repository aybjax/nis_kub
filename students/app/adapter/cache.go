package adapter

import (
	"fmt"
	"time"

	"github.com/aybjax/nis_lib/cmntypes"
	"github.com/go-redis/redis"
)

const _CACHE_TIME = time.Minute * 5

func NewCacheEngine() (cmntypes.AppCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", ENV.REDIS_URL, ENV.REDIS_PORT),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := client.Ping().Result()

	if err != nil {
		return nil, err
	}

	return &redisCache{client}, nil
}

type redisCache struct {
	redis *redis.Client
}

func (r *redisCache) Get(key string) ([]byte, error) {
	return r.redis.Get(key).Bytes()
}

func (r *redisCache) Set(key string, data []byte) error {
	return r.redis.Set(key, data, _CACHE_TIME).Err()
}

func (r *redisCache) Delete(key string) error {
	return r.redis.Del(key).Err()
}
