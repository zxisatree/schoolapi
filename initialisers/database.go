package initialisers

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() {
	var err error
	dsn := os.Getenv("MYSQL_DB_URL")
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database")
	}
	fmt.Println(DB, err)
}
