package grpc

import (
	"context"
	"log/slog"
	"net"

	"google.golang.org/grpc"

	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type Configs struct {
	Host             string
	Addr             string
	Api              serverRegister
	EnableReflection bool
}

type serverRegister interface {
	Register(grpc.ServiceRegistrar)
}

type Grpc struct {
	instance *grpc.Server
	cfgs     Configs
}

func New(cfgs Configs) *Grpc {
	return &Grpc{
		cfgs: cfgs,
	}
}

func (g *Grpc) Serve() error {
	address := g.cfgs.Host + ":" + g.cfgs.Addr
	lis, err := net.Listen("tcp", address)
	if err != nil {
		slog.Error("failed to listen", "error", err)
		return err
	}
	g.instance = grpc.NewServer()
	if g.cfgs.EnableReflection {
		reflection.Register(g.instance)
	}
	healthcheck := health.NewServer()
	healthgrpc.RegisterHealthServer(g.instance, healthcheck)

	g.cfgs.Api.Register(g.instance)

	err = g.instance.Serve(lis)
	if err != nil {
		slog.Error("failed to serve gRPC", "error", err)
		return err
	}

	return nil
}

func (g *Grpc) Shutdown() {
	if g.instance == nil {
		return
	}
	g.instance.GracefulStop()
}

func (g *Grpc) Name() string {
	return "gRPC"
}

func (g *Grpc) Healthy() bool {
	cred := insecure.NewCredentials()
	address := g.cfgs.Host + ":" + g.cfgs.Addr
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(cred))
	if err != nil {
		slog.Error("error starting grpc client", "error", err)
		return false
	}
	defer conn.Close()

	client := healthgrpc.NewHealthClient(conn)
	resp, err := client.Check(context.Background(), &healthgrpc.HealthCheckRequest{})
	if err != nil {
		slog.Error("error checking server health", "error", err)
		return false
	}
	if resp.Status != healthgrpc.HealthCheckResponse_SERVING {
		slog.Error("server is not serving", "status", resp.Status)
		return false
	}
	return true
}