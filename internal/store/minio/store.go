package minio

import (
	"bytes"
	"context"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/oc-docker/libra/internal/store"
	"github.com/pkg/errors"
)

type Store struct {
	client *minio.Client
	bucket string
}

// Range implements store.Store.
func (s *Store) Range(ctx context.Context) (chan store.Entry, error) {
	entries := make(chan store.Entry, 0)

	go func() {
		defer close(entries)
		objects := s.client.ListObjects(ctx, s.bucket, minio.ListObjectsOptions{})
		for obj := range objects {
			entries <- store.Entry{
				Name:    obj.Key,
				ModTime: obj.LastModified,
				Size:    obj.Size,
			}
		}
	}()

	return entries, nil
}

// Exists implements store.Store.
func (s *Store) Exists(ctx context.Context, name string) (bool, error) {
	_, err := s.client.GetObjectAttributes(ctx, s.bucket, name, minio.ObjectAttributesOptions{})
	if err != nil {
		return false, errors.WithStack(err)
	}

	return true, nil
}

// Reader implements store.Store.
func (s *Store) Reader(ctx context.Context, name string) (io.ReadCloser, error) {
	obj, err := s.client.GetObject(ctx, s.bucket, name, minio.GetObjectOptions{})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return obj, nil
}

// Remove implements store.Store.
func (s *Store) Remove(ctx context.Context, name string) error {
	err := s.client.RemoveObject(ctx, s.bucket, name, minio.RemoveObjectOptions{})
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Writer implements store.Store.
func (s *Store) Writer(ctx context.Context, name string) (io.WriteCloser, error) {
	return &WriterCloser{
		name:   name,
		client: s.client,
		bucket: s.bucket,
	}, nil
}

func NewStore(client *minio.Client, bucket string) *Store {
	return &Store{client, bucket}
}

var _ store.Store = &Store{}

type WriterCloser struct {
	client *minio.Client
	bucket string
	name   string
	buff   bytes.Buffer
}

// Close implements io.WriteCloser.
func (w *WriterCloser) Close() error {
	_, err := w.client.PutObject(context.Background(), w.bucket, w.name, &w.buff, -1, minio.PutObjectOptions{})
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Write implements io.WriteCloser.
func (w *WriterCloser) Write(p []byte) (n int, err error) {
	return w.buff.Write(p)
}

var _ io.WriteCloser = &WriterCloser{}
