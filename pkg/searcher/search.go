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

	result := make([]string, 0, len(files))

	if err != nil {
		return
	}

	for _, file := range files {
		content, err := fs.ReadFile(s.FS, file)

		if err != nil {
			continue
		}

		words := strings.Fields(string(content))

		for _, w := range words {
			if w == word {
				result = append(result, fileNameWithoutExtension(file))
				break
			}
		}
	}

	return result, err
}

func fileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}
