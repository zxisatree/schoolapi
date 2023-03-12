package main

import (
	"schoolapi/controllers"
	"schoolapi/initialisers"

	"github.com/gin-gonic/gin"
)

func init() {
	initialisers.LoadEnvVariables()
	initialisers.ConnectToDB()
}

func main() {
	r := gin.Default()
	r.GET("/api/commonstudents", controllers.GetCommonStudents)
	r.POST("/api/register", controllers.RegisterStudents)
	r.POST("/api/suspend", controllers.SuspendStudent)
	r.POST("/api/retrievefornotifications", controllers.GetNotifiableStudents)
	r.Run() // listen and serve on localhost:PORT, where PORT is defined in the environment variable file .env
}
