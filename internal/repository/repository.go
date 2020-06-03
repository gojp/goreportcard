package repository

import "github.com/pkg/errors"

// IRepository .
type IRepository interface {
	Get(key []byte) ([]byte, error)

	Update(key, value []byte) error

	Close()
}

var (
	_repo IRepository

	// ErrKeyNotFound .
	ErrKeyNotFound = errors.New("key not found")
)

// Init .
func Init() (err error) {
	// _repo, err = NewBadgerRepo("./.badger")
	_repo, err = NewRedisRepo()
	if err != nil {
		return err
	}

	return nil
}

// GetRepo .
func GetRepo() IRepository {
	return _repo
}
