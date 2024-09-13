package afero

import (
	"context"
	"log/slog"

	"io"
	"io/fs"
	"os"

	"github.com/oc-docker/libra/internal/store"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

type Store struct {
	fs afero.Fs
}

// Range implements store.Store.
func (s *Store) Range(ctx context.Context) (chan store.Entry, error) {
	entries := make(chan store.Entry, 0)

	go func() {
		defer close(entries)

		err := afero.Walk(s.fs, "", func(path string, info fs.FileInfo, err error) error {
			select {
			case <-ctx.Done():
				if err := ctx.Err(); err != nil {
					return errors.WithStack(err)
				}
			default:
				if err != nil || info == nil || info.IsDir() {
					return nil
				}

				entries <- store.Entry{
					Name:    info.Name(),
					ModTime: info.ModTime(),
					Size:    info.Size(),
				}
			}

			return nil
		})
		if err != nil {
			slog.ErrorContext(ctx, "could not walk filesystem", slog.Any("error", errors.WithStack(err)))
		}
	}()

	return entries, nil
}

// Remove implements store.Store.
func (s *Store) Remove(ctx context.Context, name string) error {
	return s.fs.Remove(name)
}

// Exists implements store.Store.
func (s *Store) Exists(ctx context.Context, name string) (bool, error) {
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
func (s *Store) Reader(ctx context.Context, name string) (io.ReadCloser, error) {
	return s.fs.Open(name)
}

// Writer implements store.Store.
func (s *Store) Writer(ctx context.Context, name string) (io.WriteCloser, error) {
	return s.fs.OpenFile(name, os.O_CREATE|os.O_WRONLY, 0640)
}

func NewStore(fs afero.Fs) *Store {
	return &Store{fs}
}

var _ store.Store = &Store{}
