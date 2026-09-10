package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/polar-bear-cu/sgt-subscription-service/config"
	"github.com/polar-bear-cu/sgt-subscription-service/controllers"
	"github.com/polar-bear-cu/sgt-subscription-service/repositories"
	"github.com/polar-bear-cu/sgt-subscription-service/routes"
	"github.com/polar-bear-cu/sgt-subscription-service/usecases"
)

func main() {
	config, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	repo := repositories.NewInMemorySubscription()
	uc := usecases.NewSubscription(repo)
	subCtrl := controllers.NewSubscription(uc)
	routes.Register(r, subCtrl)

	log.Println("listening :" + config.Port)
	if err := r.Run(":" + config.Port); err != nil {
		log.Fatal(err)
	}
}

