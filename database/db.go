package database

import (
	"github.com/go-redis/redis"
	"github.com/siddontang/ledisdb/config"
	"github.com/siddontang/ledisdb/ledis"
)

type Database interface {
	Get(string) (string, error)
	Set(string, string) error
	Close() error
}

// GetConnection connects to the Redis/Ledis database (or creates it if it does not exist, in the case of Ledis),
// and creates a space for saving the repos, also only if it does not exist.
func GetConnection(redisHost string) (Database, error) {
	if redisHost == "" {
		return newLedisDatabase()
	}
	return newRedisDatabase(redisHost)
}

type ledisDatabase struct {
	connection *ledis.DB
}

func (l *ledisDatabase) Get(k string) (string, error) {
	b, _ := l.connection.Get([]byte(k))
	return string(b), nil
}

func (l *ledisDatabase) Set(k, v string) error {
	return l.connection.Set([]byte(k), []byte(v))
}

func (l *ledisDatabase) Close() error {
	return nil
}

func newLedisDatabase() (*ledisDatabase, error) {
	db, err := connectLedis()
	if err != nil {
		return nil, err
	}
	return &ledisDatabase{connection: db}, nil
}

func connectLedis() (*ledis.DB, error) {
	// connect to ledis
	l, err := ledis.Open(config.NewConfigDefault())
	if err != nil {
		return nil, err
	}
	db, err := l.Select(0)
	if err != nil {
		return nil, err
	}
	return db, nil
}

type redisDatabase struct {
	connection *redis.Client
}

func newRedisDatabase(redisHost string) (*redisDatabase, error) {
	db := connectRedis(redisHost)
	err := db.Echo("test").Err()
	return &redisDatabase{connection: db}, err
}

func (r *redisDatabase) Get(k string) (string, error) {
	return r.connection.Get(k).Result()
}

func (r *redisDatabase) Set(k, v string) error {
	return r.connection.Set(k, v, 0).Err()
}

func (r *redisDatabase) Close() error {
	return r.connection.Close()
}

func connectRedis(redisHost string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return client
}
