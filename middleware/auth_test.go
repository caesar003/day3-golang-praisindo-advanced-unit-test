package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/caesar003/day3-unit-test/handler"
	"github.com/caesar003/day3-unit-test/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleWare(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()

	// Public routes
	router.GET("/", handler.RootHandler)
	router.GET("/api/user/", handler.GetUsers)
	router.GET("/api/user/:id", handler.GetUser)

	// Private routes
	private := router.Group("/api/user")
	private.Use(middleware.AuthMiddleWare())
	private.POST("/", handler.AddUser)
	private.PUT("/:id", handler.UpdateUser)
	private.DELETE("/:id", handler.DeleteUser)

	t.Run("No Authorization", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/user/", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Authorization basic token required")
	})

	t.Run("Invalid Authorization", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/user/", nil)
		req.SetBasicAuth("wronguser", "wrongpassword")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid authorization token")
	})

	t.Run("Valid Authorization", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/user/", nil)
		req.SetBasicAuth("superadmin", "supersecretpassword")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Check if we get a 400 because the body is missing (or adapt based on your actual validation logic)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
