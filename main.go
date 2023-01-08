package main

import (
	"log"
	"os"
	"time"

	"github.com/CobaKauPikirkan/aplikasi-ladangku/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.SetTrustedProxies(nil)
	router.Use(gin.Logger())
	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowCredentials: true,
		MaxAge: 12 * time.Hour,
	}))

	routes.Routes(router)

	log.Fatal(router.Run(":"+port))
}