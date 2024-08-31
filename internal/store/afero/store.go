package afero

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/oc-docker/libra/internal/store"
	"github.com/spf13/afero"
)

type Store struct {
	fs afero.Fs
}

// Remove implements store.Store.
func (s *Store) Remove(name string) error {
	return s.fs.Remove(name)
}

// Exists implements store.Store.
func (s *Store) Exists(name string) (bool, error) {
	_, err := s.fs.Stat(name)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// Reader implements store.Store.
func (s *Store) Reader(name string) (io.ReadCloser, error) {
	return s.fs.Open(name)
}

// Writer implements store.Store.
func (s *Store) Writer(name string) (io.WriteCloser, error) {
	dirname := filepath.Dir(name)
	if err := s.fs.MkdirAll(dirname, 0750); err != nil {
		return nil, err
	}

	return s.fs.OpenFile(name, os.O_CREATE|os.O_WRONLY, 0640)
}

func NewStore(fs afero.Fs) *Store {
	return &Store{fs}
}

var _ store.Store = &Store{}
