package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/DmitriyKomarovCoder/banner-api/config"
	routerInit "github.com/DmitriyKomarovCoder/banner-api/internal/app/router"
	deliveryBanner "github.com/DmitriyKomarovCoder/banner-api/internal/banner/delivery/http"
	repositoryBanner "github.com/DmitriyKomarovCoder/banner-api/internal/banner/repository"
	usecaseBanner "github.com/DmitriyKomarovCoder/banner-api/internal/banner/usecase"
	"github.com/DmitriyKomarovCoder/banner-api/pkg/closer"
	"github.com/DmitriyKomarovCoder/banner-api/pkg/logger"
	"github.com/DmitriyKomarovCoder/banner-api/pkg/postgres"
	"github.com/DmitriyKomarovCoder/banner-api/pkg/redis"
)

func Run(cfg *config.Config) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	l, err := logger.NewLogger(cfg.Log.Path)
	if err != nil {
		log.Fatalf("Logger initialisation error %s", err)
	}

	pg, err := postgres.New(cfg.PG.URL, cfg.PG.PoolMax)
	if err != nil {
		l.Fatal(fmt.Errorf("error: postgres.New: %w", err))
	}

	rd := redis.NewRedisRepository(cfg.Redis.Address, cfg.Redis.DB, *l)

	if err := rd.Connect(); err != nil {
		l.Errorf("Error Initializing connect: %v", err)
		return
	}

	l.Info("Db Connect successfully")
	cache := repositoryBanner.NewCache(rd.Client)
	repBanner := repositoryBanner.NewRepository(pg.Pool)
	useBanner := usecaseBanner.NewUsecase(repBanner, cache)
	handlerBanner := deliveryBanner.NewHandler(useBanner, *l)
	router := *routerInit.NewRouter(handlerBanner, l)

	httpServer := &http.Server{
		Addr:         cfg.Http.Host + ":" + cfg.Http.Port,
		Handler:      &router,
		ReadTimeout:  cfg.Http.ReadTimeout,
		WriteTimeout: cfg.Http.WriteTimeout,
	}

	c := &closer.Closer{}
	c.Add(httpServer.Shutdown)
	c.Add(rd.Close)
	c.Add(pg.Close)

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			l.Fatalf("Erorr starting server: %v", err)
		}
	}()
	l.Infof("server start in port: %v", cfg.Http.Port)

	<-ctx.Done()
	l.Info("shutting down server gracefully")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := c.Close(shutdownCtx); err != nil {
		l.Fatalf("closer: %v", err)
	}

	l.Info("Service close without error")
}
