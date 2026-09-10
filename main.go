package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/polar-bear-cu/sgt-subscription-service/config"
	"github.com/polar-bear-cu/sgt-subscription-service/controllers"
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

	routes.Register(r, subCtrl)

	log.Println("listening :" + cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
