package store

import "io"

type Store interface {
	Writer(name string) (io.WriteCloser, error)
	Reader(name string) (io.ReadCloser, error)
	Remove(name string) error
	Exists(name string) (bool, error)
}
