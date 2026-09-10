package routes

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"

	"github.com/polar-bear-cu/sgt-subscription-service/controllers"
)

func Register(r *gin.Engine, sub *controllers.SubscriptionController, swaggerEnabled bool) {
	r.GET("/health", controllers.GetHealth)

	if swaggerEnabled {
		r.GET("/swagger/*any", ginswagger.WrapHandler(swaggerfiles.Handler))
	}

	v1 := r.Group("/api/v1")
	v1.POST("/subscriptions", sub.CreateSubscription)
}
