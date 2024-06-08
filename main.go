package main

import (
	"github.com/caesar003/day3-golang-praisindo-advanced-unit-test/router"
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	router.SetupRouter(r)

	r.Run(":2345")
}
