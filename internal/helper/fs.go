package helper

import (
	"os"

	"github.com/pkg/errors"
)

var (
	errPathNotExists = errors.New("path not exists")
)

// exists returns whether the given file or directory exists or not
// from http://stackoverflow.com/a/10510783
func Exists(path string) (ok bool, err error) {
	if _, err = os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, errPathNotExists
		}

		return false, errors.Wrap(err, "fs.Exists failed to os.Stat")
	}

	return true, err
}

// EnsurePath make sure the path has been exists.
// it will create if path not exists
func EnsurePath(path string) (err error) {
	var ok bool
	if ok, err = Exists(path); ok {
		return nil
	}

	if err == errPathNotExists {
		// true: not exists then make dirs
		err = os.MkdirAll(path, 0777)
		if err != nil {
			return errors.Wrap(err, "fs.EnsurePath failed to mkdir")
		}

		return nil
	}

	return errors.Wrap(err, "fs.EnsurePath failed to check path")
}
