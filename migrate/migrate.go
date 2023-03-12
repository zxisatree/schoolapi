package main

import (
	"schoolapi/initialisers"
	"schoolapi/models"
)

func init() {
	initialisers.LoadEnvVariables()
	initialisers.ConnectToDB()
}

func main() {
	// Migrate the schema
	initialisers.DB.AutoMigrate(&models.Student{})
	initialisers.DB.AutoMigrate(&models.Teacher{})
}
