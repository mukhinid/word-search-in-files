package main

import (
	"fmt"
	"net/http"
	"os"

	"word-search-in-files/pkg/searcher"

	"github.com/gin-gonic/gin"
)

var searchModule *searcher.Searcher

func searchInFiles(c *gin.Context) {
	word := c.Param("word")

	files, _ := searchModule.Search(word)

	c.IndentedJSON(http.StatusOK, files)
}

func main() {
	var err error
	searchModule, err = searcher.NewSearcher(os.DirFS("examples"))

	if err != nil {
		fmt.Println("Произошла ошибка при старте приложения", err)
	}

	router := gin.Default()
	// Предпочёл сделать word path-параметром, так как в моём понимании path-параметры обязательны, а queryString-параметры опциональны. Параметр word обязателен для работы поиска.
	// Правилом это не является, просто мои предпочтения в дизайне API.
	router.GET("/files/search/:word", searchInFiles)

	router.Run("localhost:8080")
}
