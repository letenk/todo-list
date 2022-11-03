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
