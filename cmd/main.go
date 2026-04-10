package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	pb "github.com/Danil-Ivonin/GrinexRates/gen/rates"
	"github.com/Danil-Ivonin/GrinexRates/internal/config"
	"github.com/Danil-Ivonin/GrinexRates/internal/handler"
	"github.com/Danil-Ivonin/GrinexRates/internal/http/client"
	"github.com/Danil-Ivonin/GrinexRates/internal/observability"
	"github.com/Danil-Ivonin/GrinexRates/internal/services"
	"github.com/Danil-Ivonin/GrinexRates/internal/storage"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Initialize logger
	logger, err := observability.NewLogger()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync() //nolint:errcheck

	// Read configs
	err = config.Load()
	if err != nil {
		logger.Fatal("error loading config", zap.Error(err))
	}

	// Initialize db
	initTimeout := viper.GetDuration("postgres.initTimeout")
	initCtx, initCancel := context.WithTimeout(context.Background(), initTimeout)
	defer initCancel()

	pool, err := storage.New(initCtx, config.DSN())
	if err != nil {
		logger.Fatal("startup: storage init failed", zap.Error(err))
	}
	defer pool.Close()

	repo := storage.NewRatesRepository(pool)

	// Initialize Grinex client
	url := viper.GetString("grinex.url")
	timeout := viper.GetDuration("grinex.tumeout")
	cl := client.New(url, timeout)

	// Service and handler
	svc := services.NewRatesService(cl, repo, config.AvgNMPrecision())
	h := handler.NewRatesHandler(svc)

	// Server startup
	s := grpc.NewServer(grpc.StatsHandler(otelgrpc.NewServerHandler()))
	pb.RegisterRatesServiceServer(s, h)
	handler.RegisterHealth(s)
	reflection.Register(s)

	port := viper.GetString("grpc-port")
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		logger.Fatal("startup: listen failed", zap.String("port", port), zap.Error(err))
	}

	// Start serving in a goroutine
	serveErr := make(chan error, 1)
	go func() {
		logger.Info("startup: gRPC server listening", zap.String("port", port))
		serveErr <- s.Serve(lis)
	}()

	// Block until SIGTERM or SIGINT is received.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	select {
	case sig := <-quit:
		logger.Info("shutdown: received signal", zap.String("signal", sig.String()))
	case err := <-serveErr:
		logger.Error("shutdown: server exited unexpectedly", zap.Error(err))
	}

	// GracefulStop finishes in-flight RPCs then stops the server
	// Force-stop after 10 seconds
	stopped := make(chan struct{})
	go func() {
		s.GracefulStop()
		close(stopped)
	}()

	forceCtx, forceCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer forceCancel()

	select {
	case <-stopped:
		logger.Info("shutdown: graceful stop complete")
	case <-forceCtx.Done():
		logger.Warn("shutdown: graceful stop timed out, forcing stop")
		s.Stop()
	}
}
