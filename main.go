package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	grpclib "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	subscriptionv1 "github.com/polar-bear-cu/sgt-proto/gen/go/subscription/v1"
	"github.com/polar-bear-cu/sgt-subscription-service/config"
	"github.com/polar-bear-cu/sgt-subscription-service/controllers"
	_ "github.com/polar-bear-cu/sgt-subscription-service/docs"
	grpcserver "github.com/polar-bear-cu/sgt-subscription-service/grpc"
	"github.com/polar-bear-cu/sgt-subscription-service/repositories"
	"github.com/polar-bear-cu/sgt-subscription-service/routes"
	"github.com/polar-bear-cu/sgt-subscription-service/usecases"
)

// @title        Subscription Service API
// @version      1.0
// @description  REST API for the Subglutee subscription service
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pool, err := config.ConnectPostgres(ctx, cfg.DB.DSN())
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	repo := repositories.NewSubscriptionPostgres(pool)
	uc := usecases.NewSubscription(repo)
	subCtrl := controllers.NewSubscriptionController(uc)

	gs, lis, err := newGRPCServer(ctx, cfg.GRPCPort, uc)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		log.Println("gRPC :" + cfg.GRPCPort)
		if err := gs.Serve(lis); err != nil {
			log.Printf("grpc serve: %v", err)
		}
	}()

	r := gin.Default()
	routes.Register(r, subCtrl, cfg.SwaggerEnabled)
	srv := &http.Server{Addr: ":" + cfg.Port, Handler: r}

	go func() {
		log.Println("listening :" + cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
	stop()
	log.Println("shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("http shutdown: %v", err)
	}
	gs.GracefulStop()
}

func newGRPCServer(
	ctx context.Context,
	port string,
	uc *usecases.SubscriptionUsecase,
) (*grpclib.Server, net.Listener, error) {
	var lc net.ListenConfig
	lis, err := lc.Listen(ctx, "tcp", ":"+port)
	if err != nil {
		return nil, nil, err
	}

	gs := grpclib.NewServer()
	subscriptionv1.RegisterSubscriptionServiceServer(gs, grpcserver.NewSubscriptionServer(uc))
	reflection.Register(gs)
	return gs, lis, nil
}
