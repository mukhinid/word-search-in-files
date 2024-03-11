package searcher

import (
	"io/fs"
	"path/filepath"
	"strings"
	"word-search-in-files/pkg/internal/dir"
)

type Searcher struct {
	FS fs.FS
}

func (s *Searcher) Search(word string) (files []string, err error) {
	files, err = dir.FilesFS(s.FS, ".")

	if err != nil {
		return
	}

	result := make([]string, 0, len(files))

	for _, file := range files {
		content, e := fs.ReadFile(s.FS, file)

		if e != nil {
			return nil, e
		}

		words := strings.Fields(string(content))

		for _, w := range words {
			if w == word {
				result = append(result, getFileNameWithoutExtension(file))
				break
			}
		}
	}

	return result, err
}

func getFileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}
