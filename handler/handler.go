package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/caesar003/day3-unit-test/entity"
	"github.com/gin-gonic/gin"
)

var (
	Users  []entity.User
	NextID int
)

func Init() {
	Users = []entity.User{}
	NextID = 1
}

func RootHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "this is running",
	})
}

func GetUsers(c *gin.Context) {
	c.JSON(200, Users)
}

func GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	var user *entity.User

	for _, u := range Users {
		if u.ID == id {
			user = &u
			break
		}
	}

	if user == nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	c.JSON(200, user)
}

func AddUser(c *gin.Context) {
	var user entity.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.ID = NextID
	NextID++
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	Users = append(Users, user)
	c.JSON(http.StatusCreated, user)
}

func UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	var updatedUser entity.User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user *entity.User
	for i, u := range Users {
		if u.ID == id {
			Users[i].Name = updatedUser.Name
			Users[i].Password = updatedUser.Password
			Users[i].Email = updatedUser.Email
			Users[i].UpdatedAt = time.Now()
			user = &Users[i]
			break
		}
	}

	if user == nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	c.JSON(200, user)
}

func DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	var index int
	var user *entity.User
	for i, u := range Users {
		if u.ID == id {
			user = &u
			index = i
			break
		}
	}

	if user == nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	Users = append(Users[:index], Users[index+1:]...)

	c.JSON(200, gin.H{"message": "User deleted"})
}
