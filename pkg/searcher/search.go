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

func (s *Searcher) Search(word string) ([]string, error) {
	files, err := dir.FilesFS(s.FS, ".")

	if err != nil {
		return files, err
	}

	result := make([]string, 0, len(files))

	for _, file := range files {
		isFound, e := searchInFile(s.FS, file, word)

		if e != nil {
			return nil, e
		}

		if isFound {
			result = append(result, getFileNameWithoutExtension(file))
		}
	}

	return result, err
}

func searchInFile(filesystem fs.FS, filename string, word string) (bool, error) {
	content, err := fs.ReadFile(filesystem, filename)

	if err != nil {
		return false, err
	}

	words := strings.Fields(string(content))

	for _, w := range words {
		if w == word {
			return true, nil
		}
	}
	return false, nil
}

func getFileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}
