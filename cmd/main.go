package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/ncostamagna/passit-back/pkg/grpc"
)

func main() {

	dotenv.Load()
	grpcServer := grpc.New(grpc.Configs{
		Api:  grpcapi.New(app),
		Host: os.Getenv("GRPC_HOST"),
		Addr: os.Getenv("GRPC_PORT"),
	})

	shutdown := make(chan struct{}, 1)

	go func() {
		if err := s.GrpcServer.Serve(); err != nil {
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