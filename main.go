package main

import (
	"fmt"
	"log"
	"net/http"
	"siencefestival/api/controllers"
	"siencefestival/api/routers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	healthcheck "github.com/tavsec/gin-healthcheck"
	"github.com/tavsec/gin-healthcheck/checks"
	"github.com/tavsec/gin-healthcheck/config"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	fmt.Println("Server is getting started...")

	r := gin.Default()

	showRouter := routers.ShowRouter()

	showGroup := r.Group("/show")
	showGroup.Any("/*path", gin.WrapH(showRouter))

	healthcheck.New(r, config.DefaultConfig(), []checks.Check{})

	go controllers.StartMessageConsumer()

	log.Fatal(http.ListenAndServe(":8000", r))
	fmt.Println("Listening on port 8000")
}
