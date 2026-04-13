package routes

import (
	"backend/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupApiRoutes(group *gin.RouterGroup, h *handler.SubscriptionHandler) {
	subscriptions := group.Group("/subscriptions")
	{
		subscriptions.POST("", h.HandleCreateSubscription)
		subscriptions.GET("", h.HandleGetSubscriptions)
		subscriptions.GET("/:id", h.HandleGetSubscription)
		subscriptions.PUT("/:id", h.HandleUpdateSubscription)
		subscriptions.DELETE("/:id", h.HandleDeleteSubscription)
	}

	group.GET("/calculate_total_price", h.HandleCalculateTotalPrice)
}
