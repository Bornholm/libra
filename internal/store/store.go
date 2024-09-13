package store

import (
	"context"
	"io"
	"time"
)

type Entry struct {
	Name    string
	ModTime time.Time
	Size    int64
}

type Store interface {
	Writer(ctx context.Context, name string) (io.WriteCloser, error)
	Reader(ctx context.Context, name string) (io.ReadCloser, error)
	Remove(ctx context.Context, name string) error
	Exists(ctx context.Context, name string) (bool, error)
	Range(ctx context.Context) (chan Entry, error)
}
