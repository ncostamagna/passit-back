package main

import (
	"os"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/ncostamagna/passit-back/pkg/grpc"
	"github.com/ncostamagna/passit-back/transport/grpcapi"
)

func main() {

	_ = godotenv.Load()
	host := os.Getenv("GRPC_HOST")
	addr := os.Getenv("GRPC_PORT")
	grpcServer := grpc.New(grpc.Configs{
		Api:  grpcapi.New(nil),
		Host: host,
		Addr: addr,
	})

	shutdown := make(chan struct{}, 1)

	go func() {
		slog.Info("Starting grpc server", "host", host, "addr", addr)
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