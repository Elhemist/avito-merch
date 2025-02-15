package handler

import (
	"merch-test/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {

	router := gin.New()

	api := router.Group("/api")
	{
		api.GET("/auth", h.authorization)

		api.Use(h.JWTMiddleware())
		api.GET("/info", h.getInfo)
		api.GET("/sendCoin", h.sendCoin)
		api.GET("/buy/:item", h.buyItem)
	}
	return router
}
