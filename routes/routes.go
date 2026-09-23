package routes

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"

	"github.com/polar-bear-cu/sgt-subscription-service/controllers"
	"github.com/polar-bear-cu/sgt-subscription-service/middlewares"
)

func Register(r *gin.Engine, sub *controllers.SubscriptionController, jwtSecret string, swaggerEnabled bool) {
	r.GET("/health", controllers.GetHealth)

	if swaggerEnabled {
		r.GET("/swagger/*any", ginswagger.WrapHandler(swaggerfiles.Handler))
	}

	v1 := r.Group("/api/v1")
	v1.Use(middlewares.RequireAuth(jwtSecret))

	subs := v1.Group("/subscriptions")
	subs.POST("", sub.Create)
	subs.GET("/:id", sub.GetByID)
	subs.PUT("/:id", sub.Update)
	subs.PATCH("/:id", sub.UpdateStatus)
	subs.DELETE("/:id", sub.Delete)
}
