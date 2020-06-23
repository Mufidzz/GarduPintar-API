package main

import (
	"./config"
	"./controllers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	db := config.DBInit()
	inDB := &controllers.InDB{DB: db}

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "HEAD", "LAST", "AUTH"},
	}))

	user := router.Group("/user")
	{
		user.GET("/", inDB.GetUsers)
		user.GET("/:id", inDB.GetUser)

		user.POST("/", inDB.CreateUser)
		user.PUT("/:id", inDB.UpdateUser)
		user.DELETE("/:id", inDB.DeleteUser)

		user.Handle("AUTH", "/:credential/:password", inDB.AuthorizeUser)
	}

	log := router.Group("/log")
	{

		log.GET("/", inDB.GetLogs)
		log.GET("/filter/date/:date", inDB.GetLogsWithSpecificDate)
		log.GET("/list", inDB.GetLogDateLists)

		log.POST("/", inDB.CreateLog)
		log.PUT("/:id", inDB.UpdateLog)
		log.DELETE("/:id", inDB.DeleteLog)

		log.Handle("LAST", "/", inDB.GetLastLog)
	}

	router.Run(":1022")
}
