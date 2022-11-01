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

	repositoryActivityGroup := repository.NewRepositoryActivityGroup(db)
	serviceActivityGroup := service.NewServiceActivityGroup(repositoryActivityGroup)
	handlerActivityGroup := handler.NewActivityGroupHandler(serviceActivityGroup)

	// Route activity groups
	activityGroup := router.Group("/activity-groups")
	activityGroup.GET("", handlerActivityGroup.GetAll)
	activityGroup.GET("/:id", handlerActivityGroup.GetOne)
	activityGroup.POST("", handlerActivityGroup.Create)
	activityGroup.PATCH("/:id", handlerActivityGroup.Update)
	activityGroup.DELETE("/:id", handlerActivityGroup.Delete)

	return router
}
