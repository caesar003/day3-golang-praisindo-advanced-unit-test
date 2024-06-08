package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/caesar003/day3-golang-praisindo-advanced-unit-test/entity"
	"github.com/caesar003/day3-golang-praisindo-advanced-unit-test/handler"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRootHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler.Init() // Reset the state

	router := gin.Default()
	router.GET("/", handler.RootHandler)

	req, _ := http.NewRequest("GET", "/", nil)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	expectedBody := `{"message":"this is running"}`
	assert.JSONEq(t, expectedBody, w.Body.String())
}

func TestGetUsers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler.Init() // Reset the state

	router := gin.Default()
	handler.Users = []entity.User{
		{ID: 1, Name: "John Doe", Email: "john.doe@example.com"},
		{ID: 2, Name: "Jane Doe", Email: "jane.doe@example.com"},
	}
	router.GET("/users", handler.GetUsers)

	req, _ := http.NewRequest("GET", "/users", nil)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	expectedBody, _ := json.Marshal(handler.Users)
	assert.JSONEq(t, string(expectedBody), w.Body.String())
}

func TestGetUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler.Init() // Reset the state

	router := gin.Default()
	handler.Users = []entity.User{
		{ID: 1, Name: "John Doe", Email: "john.doe@example.com"},
	}
	router.GET("/users/:id", handler.GetUser)

	t.Run("Positive Case - Existing User", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/users/1", nil)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		expectedBody, _ := json.Marshal(handler.Users[0])
		assert.JSONEq(t, string(expectedBody), w.Body.String())
	})

	t.Run("Negative Case - Non-existent User", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/users/2", nil)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "User not found")
	})

	t.Run("Negative Case - Invalid User ID", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/users/abc", nil)

		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid user ID")
	})
}

func TestAddUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler.Init() // Reset the state

	router := gin.Default()
	router.POST("/users", handler.AddUser)

	newUser := entity.User{
		Name:     "Test User",
		Password: "password",
		Email:    "test@example.com",
	}
	requestBodyBytes, _ := json.Marshal(newUser)

	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(requestBodyBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// Verify the response body
	var createdUser entity.User
	err := json.Unmarshal(w.Body.Bytes(), &createdUser)
	require.NoError(t, err)
	assert.Equal(t, newUser.Name, createdUser.Name)
	assert.Equal(t, newUser.Email, createdUser.Email)
	assert.NotZero(t, createdUser.ID)
	assert.NotZero(t, createdUser.CreatedAt)
	assert.NotZero(t, createdUser.UpdatedAt)
}

func TestUpdateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler.Init() // Reset the state

	router := gin.Default()
	router.PUT("/users/:id", handler.UpdateUser)

	// Add a test user to the users slice for updating
	handler.Users = []entity.User{
		{ID: 1, Name: "Test User", Email: "test@example.com", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	updatedUser := entity.User{
		Name:     "Updated User",
		Password: "newpassword",
		Email:    "updated@example.com",
	}
	requestBodyBytes, _ := json.Marshal(updatedUser)

	req, _ := http.NewRequest("PUT", "/users/1", bytes.NewBuffer(requestBodyBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify the response body
	var user entity.User
	err := json.Unmarshal(w.Body.Bytes(), &user)
	require.NoError(t, err)
	assert.Equal(t, updatedUser.Name, user.Name)
	assert.Equal(t, updatedUser.Email, user.Email)
	assert.Equal(t, updatedUser.Password, user.Password)
	assert.Equal(t, 1, user.ID)
}

func TestDeleteUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler.Init() // Reset the state

	router := gin.Default()
	router.DELETE("/users/:id", handler.DeleteUser)

	// Add a test user to the users slice for deletion
	handler.Users = []entity.User{
		{ID: 1, Name: "Test User", Email: "test@example.com", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	req, _ := http.NewRequest("DELETE", "/users/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	expectedBody := `{"message":"User deleted"}`
	assert.JSONEq(t, expectedBody, w.Body.String())
}
