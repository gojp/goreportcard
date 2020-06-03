package repository

import (
	"github.com/dgraph-io/badger"
	"github.com/pkg/errors"
	"github.com/yeqown/log"
)

type badgerRepo struct {
	*badger.DB
}

// NewBadgerRepo .
func NewBadgerRepo(dir string) (IRepository, error) {
	db, err := badger.Open(badger.DefaultOptions("./.badger"))
	if err != nil {
		log.Errorf("repository.NewBadgerRepo failed to load db, err=%v", err)
		return nil, errors.Wrap(err, "failed to load db")
	}

	return badgerRepo{
		DB: db,
	}, nil
}

func (br badgerRepo) Get(key []byte) (out []byte, err error) {
	if err = br.DB.View(func(txn *badger.Txn) error {
		// try to load and handle error, if key not found
		item, err := txn.Get(key)
		if err != nil && err != badger.ErrKeyNotFound {
			return err
		} else if err == badger.ErrKeyNotFound || item == nil {
			return ErrKeyNotFound
		}

		// copy value to `out`s
		out = make([]byte, item.ValueSize())
		_, err = item.ValueCopy(out)
		if err != nil {
			return errors.Wrap(err, "badgerRepo.Get.ValueCopy")
		}

		return nil
	}); err != nil {
		return nil, errors.Wrap(err, "badgerRepo.Get key="+string(key))
	}

	log.Debugf("badgerRepo.Get(key=%s) v=%s", key, out)
	return out, nil
}

func (br badgerRepo) Update(key, value []byte) (err error) {
	if err = br.DB.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	}); err != nil {
		log.Debugf("BR update error, err=%v", err)
		return errors.Wrap(err, "badgerRepo.Update")
	}

	return nil
}

func (br badgerRepo) Close() {
	br.DB.Close()
}
