package router

import (
	"github.com/caesar003/day3-golang-praisindo-advanced-unit-test/handler"
	"github.com/caesar003/day3-golang-praisindo-advanced-unit-test/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	r.GET("/", handler.RootHandler)

	userPublic := r.Group("/api/user")
	userPublic.GET("/", handler.GetUsers)
	userPublic.GET("/:id", handler.GetUser)

	userPrivate := r.Group("/api/user")
	userPrivate.Use(middleware.AuthMiddleWare())
	userPrivate.POST("/", handler.AddUser)
	userPrivate.PUT("/:id", handler.UpdateUser)
	userPrivate.DELETE("/:id", handler.DeleteUser)
}
