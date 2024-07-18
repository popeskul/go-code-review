package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"coupon_service/internal/api"
	"coupon_service/internal/config"
	"coupon_service/internal/repository/memdb"
	"coupon_service/internal/service"
	_ "go.uber.org/automaxprocs"
)

func main() {
	cfg := config.New()
	repo := memdb.New()
	svc := service.New(repo)
	apiServer := api.New(cfg.API, svc)

	go func() {
		fmt.Println("Starting Coupon service server")
		if err := apiServer.Start(); err != nil {
			fmt.Printf("Error starting API server: %v\n", err)
			os.Exit(1)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	fmt.Println("Shutting down Coupon service server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := apiServer.Close(ctx); err != nil {
		fmt.Printf("Server forced to shutdown: %v\n", err)
	}
}
