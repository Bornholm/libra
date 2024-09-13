package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/oc-docker/libra/internal/store"
	aferostore "github.com/oc-docker/libra/internal/store/afero"
	miniostore "github.com/oc-docker/libra/internal/store/minio"
	"github.com/spf13/afero"
)

func createStoreFromDSN(dsn string) (store.Store, error) {
	url, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}

	var store store.Store

	switch url.Scheme {
	case "file":
		store, err = createLocalStore(url)
		if err != nil {
			return nil, err
		}

	case "memory":
		fs := afero.NewMemMapFs()
		store = aferostore.NewStore(fs)

	case "s3":
		store, err = createS3Store(url)
		if err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("could not find store implementation for scheme '%s'", url.Scheme)
	}

	go runCleanJob(store)

	return store, nil
}

func createLocalStore(dsn *url.URL) (store.Store, error) {
	path, err := filepath.Abs(filepath.Clean(dsn.Host + dsn.Path))
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(path, 0750); err != nil {
		return nil, err
	}

	fs := afero.NewBasePathFs(afero.NewOsFs(), path)

	return aferostore.NewStore(fs), nil
}

func createS3Store(dsn *url.URL) (store.Store, error) {
	region := dsn.Query().Get("region")
	if region == "" {
		region = "us-east-1"
	}

	var (
		id     string
		secret string
	)

	if dsn.User != nil {
		id = dsn.User.Username()
		secret, _ = dsn.User.Password()
	}

	token := dsn.Query().Get("token")

	client, err := minio.New(dsn.Host, &minio.Options{
		BucketLookup: minio.BucketLookupAuto,
		Creds:        credentials.NewStaticV4(id, secret, token),
		Region:       region,
	})
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	bucket := filepath.Base(dsn.Path)

	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		return nil, err
	}

	if !exists {
		if err := client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{Region: region}); err != nil {
			return nil, err
		}
	}

	return miniostore.NewStore(client, bucket), nil
}
