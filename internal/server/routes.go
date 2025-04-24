package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "order-pack-calculator/docs"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	// Swagger route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")
	api.GET("/health", s.healthHandler)
	v1 := api.Group("/v1")

	packsizes := v1.Group("/packsizes")
	packsizes.POST("/", s.CreatePackSizeHandler)
	packsizes.PATCH("/", s.UpdatePackSizeHandler)

	orders := v1.Group("/orders")
	orders.POST("/calculate", s.CalculatePackSizeHandler)

	return r
}
