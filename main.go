package main

import (
	"fmt"
	"log"
	"net/http"
	"siencefestival/api/routers"
	"siencefestival/api/controllers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

	go controllers.StartMessageConsumer()

	log.Fatal(http.ListenAndServe(":8000", r))
	fmt.Println("Listening on port 8000")
}
