package routers

import (
	"siencefestival/api/controllers"
	"github.com/gin-gonic/gin"
)

func ShowRouter() *gin.Engine {
	router := gin.Default()

	authorGroup := router.Group("/show")
	{
		authorGroup.POST("/add", controllers.AddShow)
		authorGroup.GET("/get", controllers.GetAllShows)
		authorGroup.GET("/get/:id", controllers.GetShowById)
		authorGroup.GET("/get-accepted", controllers.GetAcceptedShows)
		authorGroup.GET("/get-unaccepted", controllers.GetUnacceptedShows)
	}

	return router

}