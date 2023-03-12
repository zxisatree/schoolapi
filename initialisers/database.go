package initialisers

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var MockDB *sql.DB
var Mock sqlmock.Sqlmock

func ConnectToDB() {
	var err error
	dsn := os.Getenv("MYSQL_DB_URL")
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database")
	}
	fmt.Println(DB, err)
}

func InitialiseMockDb(t *testing.T) {
	var err error
	MockDB, Mock, err = sqlmock.New()
	if err != nil {
		t.Fatal("Failed to open a stub database connection with error: ", err)
	}

	dialector := mysql.New(mysql.Config{
		DSN:                       "sqlmock_school_api_db",
		DriverName:                "mysql",
		Conn:                      MockDB,
		SkipInitializeWithVersion: true,
	})

	DB, err = gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Error("Failed to connect to SQLMock database with error: ", err)
	}
}
