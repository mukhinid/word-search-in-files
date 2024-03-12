package main

import (
	"net/http"
	"os"

	"word-search-in-files/pkg/searcher"

	"github.com/gin-gonic/gin"
)

var searchModule searcher.Searcher

func searchInFiles(c *gin.Context) {
	word := c.Param("word")

	files, _ := searchModule.Search(word)

	c.IndentedJSON(http.StatusOK, files)
}

func main() {
	searchModule = searcher.Searcher{
		FS: os.DirFS("examples"),
	}

	router := gin.Default()
	router.GET("/files/search/:word", searchInFiles)

	router.Run("localhost:8080")
}
