package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"fmt"

	"github.com/joho/godotenv"
	"github.com/ncostamagna/passit-back/adapters/cache"
	"github.com/ncostamagna/passit-back/pkg/config"
	"github.com/ncostamagna/passit-back/pkg/grpc"
	"github.com/ncostamagna/passit-back/pkg/instance"
	"github.com/ncostamagna/passit-back/pkg/log"
	"github.com/ncostamagna/passit-back/transport/grpcapi"
)

type Config struct {
	Redis struct {
		Addr string `mapstructure:"addr"`
		Pass string `mapstructure:"pass"`
		Db   int    `mapstructure:"db"`
	} `mapstructure:"redis"`
	API struct {
		GRPC struct {
			Host string `mapstructure:"host"`
			Port string `mapstructure:"port"`
		} `mapstructure:"grpc"`
		SuperToken string `mapstructure:"super_token"`
	} `mapstructure:"api"`
	Token string `mapstructure:"token"`
}

const fileConfid = ".config.yaml"

func main() {

	_ = godotenv.Load()
	cfg := &Config{}
	if err := config.Load(cfg, fileConfid); err != nil {
		slog.Error("Error initializing config manager", "err", err)
		os.Exit(1)
	}

	ctx := context.Background()
	log := log.New(log.Config{
		AppName:   "passit-back",
		Level:     "info",
		AddSource: true,
	})

	log.Info("config file")
	log.Info(fmt.Sprint(cfg))
	log.Info(cfg.Redis.Addr)

	redisDb, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	cache := cache.NewCache(os.Getenv("REDIS_ADDR"), os.Getenv("REDIS_PASS"), redisDb)
	srv := instance.NewSecretService(cache, log)

	grpcConfig := grpc.Configs{
		API:  grpcapi.New(srv),
		Host: os.Getenv("GRPC_HOST"),
		Addr: os.Getenv("GRPC_PORT"),
	}

	grpcServer := grpc.New(ctx, grpcConfig)

	shutdown := make(chan struct{}, 1)

	go func() {
		slog.Info("Starting grpc server", "host", grpcConfig.Host, "addr", grpcConfig.Addr)
		if err := grpcServer.Serve(); err != nil {
			slog.Error("Error starting grpc server", "err", err)
			shutdown <- struct{}{}
		}
	}()

	kill := make(chan os.Signal, 1)
	signal.Notify(kill, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-kill
		shutdown <- struct{}{}
	}()

	<-shutdown
	grpcServer.Shutdown()
}
