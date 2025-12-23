package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"ms-parcel-core/internal/infrastructure/http/middleware"
	httpRouter "ms-parcel-core/internal/infrastructure/http/router"
)

func main() {
	// Gin base (manténlo simple por ahora)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.DevClaimsMiddleware())
	r.Use(middleware.ErrorMiddleware())

	// Registrar rutas del monolito
	httpRouter.RegisterRoutes(r)

	// Puerto
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("listening on :" + port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
