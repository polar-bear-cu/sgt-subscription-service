package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/polar-bear-cu/sgt-subscription-service/controllers"
)

func Register(r *gin.Engine, sub *controllers.SubscriptionController) {
	r.GET("/health", controllers.Health)

	v1 := r.Group("/api/v1")
	v1.POST("/subscriptions", sub.Create)
}
