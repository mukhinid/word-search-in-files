package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func searchInFiles(c *gin.Context) {
	word := c.Param("word")

	c.IndentedJSON(http.StatusOK, word)
}

func main() {
	router := gin.Default()
	router.GET("/files/search/:word", searchInFiles)

	router.Run("localhost:8080")
}
