package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os/signal"
	"syscall"
	"time"
	"url-shortener/internal/api"
	"url-shortener/internal/api/getoriginal"
	"url-shortener/internal/api/redirect"
	"url-shortener/internal/api/shorten"
	"url-shortener/internal/domain/shortgen"
	"url-shortener/internal/infrastructure/config"
	"url-shortener/internal/infrastructure/db/postgres"
	"url-shortener/internal/infrastructure/logger"
	"url-shortener/internal/infrastructure/persistence"
	dbPersistence "url-shortener/internal/infrastructure/persistence/db/shortenedurls"
	inmemoryPersistence "url-shortener/internal/infrastructure/persistence/inmemory/shortenedurls"
	"url-shortener/internal/middleware"
	getoriginalsvc "url-shortener/internal/service/getoriginal"
	shortensvc "url-shortener/internal/service/shorten"

	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
)

type Storage struct {
	urlRepo persistence.UrlRepository
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger.Init()

	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	storage, cleanup, err := initStorage(ctx, cfg)
	if err != nil {
		panic(fmt.Errorf("failed to initialize storage: %w", err))
	}

	defer cleanup()

	generator := shortgen.NewGenerator(rand.New(rand.NewSource(time.Now().UnixNano())))

	shortenService := shortensvc.New(storage.urlRepo, generator)
	getOriginalService := getoriginalsvc.New(storage.urlRepo)

	shortenHandler := shorten.New(shortenService)
	redirectHandler := redirect.New(getOriginalService)
	getOriginalHandler := getoriginal.New(getOriginalService)

	e := echo.New()
	e.Use(echomw.RequestID())
	e.Use(middleware.LoggingMiddleware())

	api := api.NewAPI(
		shortenHandler,
		redirectHandler,
		getOriginalHandler,
	)
	api.InitRoutes(e)

	serverAddress := ":" + cfg.Port
	go func() {
		logger.Log.Info("Starting HTTP server", "addr", serverAddress)
		if err := e.Start(serverAddress); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Log.Error("Failed to start server", "error", err)
			panic(err)
		}
	}()

	<-ctx.Done()

	logger.Log.Info("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := e.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Failed to shutdown server: %v", err)
	} else {
		logger.Log.Info("Server gracefully stopped")
	}
}

func initStorage(ctx context.Context, cfg *config.Config) (*Storage, func(), error) {
	switch cfg.StorageType {
	case config.StorageTypePostgres:
		db, err := postgres.NewFromConfig(ctx)
		if err != nil {
			log.Fatalf("failed to connect to database: %v", err)
		}
		return &Storage{
				urlRepo: dbPersistence.NewRepository(db),
			}, func() {
				db.Close()
			}, nil
	case config.StorageTypeMemory:
		return &Storage{
			urlRepo: inmemoryPersistence.NewRepository(),
		}, func() {}, nil
	default:
		return nil, nil, errors.New("invalid storage type")
	}
}
