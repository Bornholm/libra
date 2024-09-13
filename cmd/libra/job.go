package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/oc-docker/libra/internal/store"
	"github.com/pkg/errors"
)

func runCleanJob(store store.Store) {
	ticker := time.NewTicker(cleanupJobInterval)
	for {
		slog.Info("starting cleanup job", slog.Any("interval", cleanupJobInterval), slog.Any("ttl", fileTTL))

		treshold := time.Now().Add(-fileTTL)

		func() {
			ctx, cancel := context.WithTimeout(context.Background(), cleanupJobInterval*10)
			defer cancel()

			entries, err := store.Range(ctx)
			if err != nil {
				slog.Error("could not range over entries", slog.Any("error", errors.WithStack(err)))
			}

			for e := range entries {
				slog.Debug("checking entry", slog.Any("name", e.Name), slog.Any("modtime", e.ModTime))

				if e.ModTime.After(treshold) {
					continue
				}

				slog.Info("deleting entry", slog.Any("name", e.Name), slog.Any("modtime", e.ModTime))

				if err := store.Remove(ctx, e.Name); err != nil {
					slog.Error("could not remove entry", slog.Any("name", e.Name))
				}
			}
		}()

		slog.Info("cleanup job finished", slog.Any("next", time.Now().Add(cleanupJobInterval)))

		<-ticker.C
	}
}
