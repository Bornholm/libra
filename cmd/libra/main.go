package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/oc-docker/libra/internal/server"
)

var (
	address            string        = ":8080"
	rawLogLevel        string        = slog.LevelInfo.String()
	maxUploadSize      int64         = 32 << 20
	baseURL            string        = "http://localhost:8080"
	storeDSN           string        = "file://uploads"
	fileTTL            time.Duration = 24 * time.Hour
	cleanupJobInterval time.Duration = time.Hour
)

func init() {
	flag.StringVar(&address, "address", envString("LIBRA_ADDRESS", address), "server listening address (env: LIBRA_ADDRESS)")
	flag.StringVar(&rawLogLevel, "log-level", envString("LIBRA_LOG_LEVEL", rawLogLevel), "logging level (env: LIBRA_LOG_LEVEL)")
	flag.StringVar(&baseURL, "base-url", envString("LIBRA_BASE_URL", baseURL), "base url (env: LIBRA_BASE_URL)")
	flag.Int64Var(&maxUploadSize, "max-upload-size", envInt64("LIBRA_MAX_UPLOAD_SIZE", maxUploadSize), "max upload size (env: LIBRA_MAX_UPLOAD_SIZE)")
	flag.StringVar(&storeDSN, "store-dsn", envString("LIBRA_STORE_DSN", storeDSN), "store dsn (env: LIBRA_STORE_DSN)")
	flag.DurationVar(&fileTTL, "file-ttl", envDuration("LIBRA_FILE_TTL", fileTTL), "file ttl (env: LIBRA_FILE_TTL)")
	flag.DurationVar(&cleanupJobInterval, "cleanup-job-interval", envDuration("LIBRA_CLEANUP_JOB_INTERVAL", fileTTL), "cleanup job interval (env: LIBRA_CLEANUP_JOB_INTERVAL)")
}

func main() {
	flag.Parse()

	var logLevel slog.Level
	if err := logLevel.UnmarshalText([]byte(rawLogLevel)); err != nil {
		slog.Error("could not unmarshal log level", slog.Any("error", err))
		os.Exit(1)
	}

	slog.Info("setting log level", slog.Any("level", logLevel.String()))
	slog.SetLogLoggerLevel(logLevel)

	slog.Info("using store dsn", slog.Any("dsn", storeDSN))
	store, err := createStoreFromDSN(storeDSN)
	if err != nil {
		slog.Error("could not create store", slog.Any("error", err))
		os.Exit(1)
	}

	opts := &server.ServerOptions{
		MaxUploadSize: maxUploadSize,
		BaseURL:       baseURL,
		Store:         store,
	}

	server := server.New(opts)

	slog.Info("listening", slog.Any("address", address))

	if err := http.ListenAndServe(address, server); err != nil {
		slog.Error("could not listen", slog.Any("error", err))
		os.Exit(1)
	}
}
