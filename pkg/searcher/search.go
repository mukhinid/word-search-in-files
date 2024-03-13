package searcher

import (
	"errors"
	"io/fs"
	"path/filepath"
	"slices"
	"strings"
	"word-search-in-files/pkg/internal/dir"
)

// Searcher позволяет искать ключевое слово в файлах указанной папки.
type Searcher struct {
	FS           fs.FS               // Файловая система, предоставляющая папку, в файлах которой осуществляется поиск.
	wordFilesMap map[string][]string // Словарь соответствия слов и файлов.
}

// Инициализирует словарь для поиска.
func (s *Searcher) Init() error {
	s.wordFilesMap = make(map[string][]string)

	files, err := dir.FilesFS(s.FS, ".")

	if err != nil {
		return err
	}

	for _, file := range files {
		fileWihtoutExt := getFileNameWithoutExtension(file)
		content, fileErr := fs.ReadFile(s.FS, file)

		if fileErr != nil {
			continue
		}

		words := strings.Fields(string(content))
		for _, word := range words {
			entry, wordContains := s.wordFilesMap[word]
			if wordContains {
				if !slices.Contains(entry, fileWihtoutExt) {
					s.wordFilesMap[word] = append(entry, fileWihtoutExt)
				}
			} else {
				s.wordFilesMap[word] = []string{fileWihtoutExt}
			}
		}
	}

	return nil
}

// Производит поиск слова word по файлам в папке файловой системы.
func (s *Searcher) Search(word string) ([]string, error) {
	if s.wordFilesMap == nil {
		return nil, errors.New("непроинциализирован словарь поиска, вызовите функцию Searcher.Init(), прежде чем искать")
	}

	files := s.wordFilesMap[word]

	return files, nil
}

// Возвращает имя файла filename без расширения.
func getFileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}
