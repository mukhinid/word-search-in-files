package main

import (
	"net/http"
	"os"
	"strconv"

	"word-search-in-files/pkg/searcher"

	"github.com/gin-gonic/gin"
)

var searchModule searcher.Searcher

func searchInFiles(c *gin.Context) {
	ignoreCase := false
	var err error

	word := c.Param("word")
	ignoreCaseStr := c.Query("ignore_case")

	if ignoreCaseStr != "" {
		ignoreCase, err = strconv.ParseBool(ignoreCaseStr)
		if err != nil {
			ignoreCase = false
		}
	}

	files, _ := searchModule.Search(word, ignoreCase)

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
