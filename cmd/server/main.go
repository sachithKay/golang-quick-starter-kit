package main

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/sachithKay/ghost/internal/config"
	"github.com/sachithKay/ghost/internal/handler"
	"github.com/sachithKay/ghost/internal/repository"
	"github.com/sachithKay/ghost/internal/service"

	// 1. Point to your generated code
	orderv1 "github.com/sachithKay/ghost/gen/go/v1"
)

func main() {
	// 1. Standard Logging & Config
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize Database Connection Pool
	dbPool, err := pgxpool.New(ctx, cfg.DB.URL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	// Dependency Injection Wiring
	orderRepo := repository.NewPostgresOrderRepository(dbPool)
	orderSvc := service.NewOrderService(orderRepo)
	orderHandler := handler.NewOrderHandler(cfg, orderSvc)

	// ---------------------------------------------------------
	// SERVER 1: gRPC (Internal Port - 50051)
	// ---------------------------------------------------------
	grpcServer := grpc.NewServer()
	// This is the "Plug-in" moment we discussed!
	orderv1.RegisterOrderServiceServer(grpcServer, orderHandler)

	grpcAddr := "0.0.0.0:50051"
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		slog.Error("failed to listen for gRPC", "error", err)
		os.Exit(1)
	}

	go func() {
		slog.Info("gRPC server starting", "addr", grpcAddr)
		if err := grpcServer.Serve(lis); err != nil {
			slog.Error("gRPC serve error", "error", err)
		}
	}()

	// ---------------------------------------------------------
	// SERVER 2: REST Gateway (External Port - e.g. 8080)
	// ---------------------------------------------------------

	// This is the "Translator" mux
	mux := runtime.NewServeMux()

	// We tell the Gateway how to connect to the gRPC server we just started
	// Since they are in the same binary, we use localhost:50051
	dialOpts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err = orderv1.RegisterOrderServiceHandlerFromEndpoint(ctx, mux, grpcAddr, dialOpts)
	if err != nil {
		slog.Error("failed to register REST Gateway", "error", err)
		os.Exit(1)
	}

	httpSrv := &http.Server{
		Addr:         ":" + cfg.Port, // This is your 8080 from config
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		slog.Info("REST Gateway starting", "port", cfg.Port)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("REST Gateway error", "error", err)
		}
	}()

	// ---------------------------------------------------------
	// GRACEFUL SHUTDOWN
	// ---------------------------------------------------------
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop // Wait here for Ctrl+C or kill signal
	slog.Info("shutting down servers...")

	// 1. Stop gRPC first (stops accepting new binary calls)
	grpcServer.GracefulStop()

	// 2. Stop HTTP Gateway with a timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()

	if err := httpSrv.Shutdown(shutdownCtx); err != nil {
		slog.Error("forced gateway shutdown", "error", err)
	}

	slog.Info("server exited gracefully")
}
