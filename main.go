package main

import (
	"fmt"

	"github.com/letenk/todo-list/config"
	"github.com/letenk/todo-list/router"
)

func main() {
	fmt.Println("App is starting...")
	db := config.SetupDB()
	router := router.SetupRouter(db)
	router.Run(":3030")
}

/*
export MYSQL_USER="root"
export MYSQL_PASSWORD="root"
export MYSQL_HOST="127.0.0.1"
export MYSQL_PORT="3306"
export MYSQL_DBNAME="todo4"
*/
