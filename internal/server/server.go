package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"order-pack-calculator/internal/database"
	"order-pack-calculator/internal/domain/repositories"
	"order-pack-calculator/internal/domain/services"
)

type Server struct {
	port int

	dbService       database.Service
	packSizeService services.PackSizeService
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	dbService := database.New()
	packSizeRepository := repositories.NewPackSizeRepository(dbService.GetDB())
	packSizeService := services.NewPackSizeService(packSizeRepository)
	NewServer := &Server{
		port:      port,
		dbService: dbService,

		packSizeService: packSizeService,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
