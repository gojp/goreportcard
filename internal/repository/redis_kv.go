package repository

import (
	"github.com/go-redis/redis"
)

type redisRepo struct {
	*redis.Client
}

// NewRedisRepo .
func NewRedisRepo() (IRepository, error) {
	opt := &redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	}
	db := redis.NewClient(opt)
	_, err := db.Ping().Result()
	if err != nil {
		return nil, err
	}

	return redisRepo{Client: db}, nil
}

func (rd redisRepo) Get(key []byte) ([]byte, error) {
	out, err := rd.Client.Get(string(key)).Bytes()
	if err != nil && err == redis.Nil {
		// TODO: handler not exists error
		return nil, ErrKeyNotFound
	}

	return out, err
}

func (rd redisRepo) Update(key, value []byte) error {
	return rd.Client.Set(string(key), value, 0).Err()
}

func (rd redisRepo) Close() {
	rd.Client.Close()
}
