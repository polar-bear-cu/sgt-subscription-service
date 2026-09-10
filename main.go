package main

import (
	"context"
	"log"
	"net"

	"github.com/gin-gonic/gin"
	grpclib "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	subscriptionv1 "github.com/polar-bear-cu/sgt-proto/gen/go/subscription/v1"
	"github.com/polar-bear-cu/sgt-subscription-service/config"
	"github.com/polar-bear-cu/sgt-subscription-service/controllers"
	grpcserver "github.com/polar-bear-cu/sgt-subscription-service/grpc"
	"github.com/polar-bear-cu/sgt-subscription-service/repositories"
	"github.com/polar-bear-cu/sgt-subscription-service/routes"
	"github.com/polar-bear-cu/sgt-subscription-service/usecases"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	ctx := context.Background()
	pool, err := config.ConnectPostgres(ctx, cfg.DB.DSN())
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	repo := repositories.NewSubscriptionPostgres(pool)
	uc := usecases.NewSubscription(repo)
	subCtrl := controllers.NewSubscriptionController(uc)

	go startGRPC(ctx, cfg.GRPCPort, uc)

	routes.Register(r, subCtrl)

	log.Println("listening :" + cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}

func startGRPC(ctx context.Context, port string, uc *usecases.SubscriptionUsecase) {
	var lc net.ListenConfig
	lis, err := lc.Listen(ctx, "tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}

	gs := grpclib.NewServer()
	subscriptionv1.RegisterSubscriptionServiceServer(gs, grpcserver.NewSubscriptionServer(uc))
	reflection.Register(gs)

	log.Println("gRPC :" + port)
	if err := gs.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
