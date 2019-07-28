package database

import (
	"math"
	"time"

	"github.com/go-redis/redis"
	"github.com/siddontang/ledisdb/config"
	"github.com/siddontang/ledisdb/ledis"
)

// Database represents the main functions that should be made available by a GRC database
// implementation
type Database interface {
	GetRepo(string) (string, error)
	SetRepo(repo string, value string) error
	SetScore(repo string, score int) error
	GetHighScores(n int) ([]string, error)
	SetRecentlyViewed(repo string) error
	GetMostRecentlyViewed(n int) ([]string, error)
	GetRepoCount() (int, error)
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

func (l *ledisDatabase) GetRepo(k string) (string, error) {
	b, _ := l.connection.Get([]byte(k))
	return string(b), nil
}

func (l *ledisDatabase) SetRepo(k, v string) error {
	return l.connection.Set([]byte(k), []byte(v))
}

func (l *ledisDatabase) SetScore(repo string, score int) error {
	pair := ledis.ScorePair{Score: int64(score), Member: []byte(repo)}
	_, err := l.connection.ZAdd([]byte("scores"), pair)
	return err
}
func (l *ledisDatabase) GetHighScores(n int) ([]string, error) {
	pairs, err := l.connection.ZRevRange([]byte("scores"), 0, n)
	if err != nil {
		return nil, err
	}

	s := make([]string, len(pairs))
	for i := range pairs {
		s[i] = string(pairs[i].Member)
	}
	return s, nil
}

func (l *ledisDatabase) SetRecentlyViewed(repo string) error {
	pair := ledis.ScorePair{Score: time.Now().UnixNano(), Member: []byte(repo)}
	_, err := l.connection.ZAdd([]byte("last_viewed"), pair)
	return err
}

func (l *ledisDatabase) GetMostRecentlyViewed(n int) ([]string, error) {
	pairs, err := l.connection.ZRevRange([]byte("last_viewed"), 0, n-1)
	if err != nil {
		return nil, err
	}

	s := make([]string, len(pairs))
	for i := range pairs {
		s[i] = string(pairs[i].Member)
	}
	return s, nil
}

func (l *ledisDatabase) GetRepoCount() (int, error) {
	n, err := l.connection.ZCount([]byte("last_viewed"), math.MinInt64, math.MaxInt64)
	return int(n), err
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

func (r *redisDatabase) GetRepo(k string) (string, error) {
	return r.connection.Get(k).Result()
}

func (r *redisDatabase) SetRepo(k, v string) error {
	return r.connection.Set(k, v, 0).Err()
}

func (r *redisDatabase) SetScore(repo string, score int) error {
	return r.connection.ZAdd("scores", &redis.Z{Member: repo, Score: float64(score)}).Err()
}

func (r *redisDatabase) GetHighScores(n int) ([]string, error) {
	return r.connection.ZRevRange("scores", 0, int64(n)).Result()
}

func (r *redisDatabase) SetRecentlyViewed(repo string) error {
	return r.connection.ZAdd("last_viewed", &redis.Z{Member: repo, Score: float64(time.Now().UnixNano())}).Err()
}

func (r *redisDatabase) GetMostRecentlyViewed(n int) ([]string, error) {
	return r.connection.ZRevRange("last_viewed", 0, int64(n)).Result()
}

func (r *redisDatabase) GetRepoCount() (int, error) {
	n, err := r.connection.ZCount("last_viewed", "-inf", "inf").Result()
	return int(n), err
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
