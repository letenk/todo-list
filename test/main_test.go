package test

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/letenk/todo-list/config"
	"github.com/letenk/todo-list/router"
	"gorm.io/gorm"
)

var ConnTest *gorm.DB
var Route *gin.Engine

func TestMain(m *testing.M) {
	// Set env
	os.Setenv("MYSQL_USER", "root")
	os.Setenv("MYSQL_PASSWORD", "root")
	os.Setenv("MYSQL_HOST", "127.0.0.1")
	os.Setenv("MYSQL_PORT", "3306")
	os.Setenv("MYSQL_DBNAME", "todo4")

	// Open connection
	db := config.SetupDB()
	ConnTest = db

	// Setup router
	Route = router.SetupRouter(db)

	m.Run()
}
