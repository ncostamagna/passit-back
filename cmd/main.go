package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"fmt"

	"github.com/ncostamagna/passit-back/adapters/cache"
	"github.com/ncostamagna/passit-back/pkg/config"
	"github.com/ncostamagna/passit-back/pkg/grpc"
	"github.com/ncostamagna/passit-back/pkg/instance"
	"github.com/ncostamagna/passit-back/pkg/log"
	"github.com/ncostamagna/passit-back/transport/grpcapi"
	"github.com/ncostamagna/passit-back/transport/httpapi"
)

type Config struct {
	LogLevel string `mapstructure:"log_level"`
	Redis    struct {
		Addr string `mapstructure:"addr"`
		Pass string `mapstructure:"pass"`
		Db   int    `mapstructure:"db"`
	} `mapstructure:"redis"`
	API struct {
		GRPC struct {
			Host string `mapstructure:"host"`
			Port string `mapstructure:"port"`
		} `mapstructure:"grpc"`
		HTTP struct {
			Host string `mapstructure:"host"`
			Port string `mapstructure:"port"`
		} `mapstructure:"http"`
	} `mapstructure:"api"`
	Token string `mapstructure:"token"`
}

const fileConfig = ".config.yaml"

func main() {

	cfg := &Config{}
	if err := config.Load(cfg, fileConfig); err != nil {
		slog.Error("Error initializing config manager", "err", err)
		os.Exit(1)
	}

	ctx := context.Background()
	log := log.New(log.Config{
		AppName:   "passit-back",
		Level:     cfg.LogLevel,
		AddSource: true,
	})

	cache := cache.NewCache(cfg.Redis.Addr, cfg.Redis.Pass, cfg.Redis.Db)
	srv := instance.NewSecretService(cache, log)

	grpcConfig := grpc.Configs{
		API:  grpcapi.New(srv),
		Host: cfg.API.GRPC.Host,
		Addr: cfg.API.GRPC.Port,
	}

	grpcServer := grpc.New(ctx, grpcConfig)
	apiServer := httpapi.NewHTTPAPI()

	shutdown := make(chan struct{}, 1)
	errs := make(chan error, 1)

	go func() {
		slog.Info("Starting grpc server", "host", grpcConfig.Host, "addr", grpcConfig.Addr)
		if err := grpcServer.Serve(); err != nil {
			slog.Error("Error starting grpc server", "err", err)
			shutdown <- struct{}{}
		}
	}()

	go func() {
		url := fmt.Sprintf("%s:%s", cfg.API.HTTP.Host, cfg.API.HTTP.Port)
		slog.Info("Listening", "url", url)
		errs <- apiServer.Listen(url)
	}()

	kill := make(chan os.Signal, 1)
	signal.Notify(kill, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-kill
		shutdown <- struct{}{}
		errs <- errors.New("received signal to shutdown")
	}()

	<-shutdown
	<-errs
	grpcServer.Shutdown()
}
