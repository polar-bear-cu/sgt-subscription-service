package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/polar-bear-cu/sgt-subscription-service/controllers"
	"github.com/polar-bear-cu/sgt-subscription-service/repositories"
	"github.com/polar-bear-cu/sgt-subscription-service/routes"
	"github.com/polar-bear-cu/sgt-subscription-service/usecases"
)

func main() {
	r := gin.Default()
	repo := repositories.NewInMemorySubscription()
	uc := usecases.NewSubscription(repo)
	subCtrl := controllers.NewSubscription(uc)
	routes.Register(r, subCtrl)

	log.Println("listening :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
