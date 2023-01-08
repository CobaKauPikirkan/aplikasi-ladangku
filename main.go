package main

import (
	"log"
	"os"

	"github.com/CobaKauPikirkan/aplikasi-ladangku/routes"
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

	routes.Routes(router)

	log.Fatal(router.Run(":"+port))
}