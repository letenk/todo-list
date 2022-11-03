package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/letenk/todo-list/handler"
	"github.com/letenk/todo-list/repository"
	"github.com/letenk/todo-list/service"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://*", "http://*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	repositoryActivity := repository.NewRepositoryActivity(db)
	serviceActivity := service.NewServiceActivity(repositoryActivity)
	handlerActivity := handler.NewActivityHandler(serviceActivity)

	// Route activity groups
	Activity := router.Group("/activity-groups")
	Activity.GET("", handlerActivity.GetAll)
	Activity.GET("/:id", handlerActivity.GetOne)
	Activity.POST("", handlerActivity.Create)
	Activity.PATCH("/:id", handlerActivity.Update)
	Activity.DELETE("/:id", handlerActivity.Delete)

	repositoryTodo := repository.NewRepositoryTodo(db)
	serviceTodo := service.NewServiceTodo(repositoryTodo)
	handlerTodo := handler.NewTodoHandler(serviceTodo)

	// Route todo
	todo := router.Group("/todo-items")
	todo.GET("", handlerTodo.GetAll)
	todo.GET("/:id", handlerTodo.GetOne)
	todo.POST("", handlerTodo.Create)
	todo.PATCH("/:id", handlerTodo.Update)
	todo.DELETE("/:id", handlerTodo.Delete)
	return router
}
