package searcher

import (
	"errors"
	"io/fs"
	"path/filepath"
	"slices"
	"strings"
	"sync"
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

	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}

	wg.Add(len(files))
	for _, file := range files {
		go s.processFile(file, mu, wg)
	}

	wg.Wait()

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

func (s *Searcher) processFile(filename string, mu *sync.Mutex, wg *sync.WaitGroup) error {
	defer wg.Done()

	fileWihtoutExt := getFileNameWithoutExtension(filename)
	content, err := fs.ReadFile(s.FS, filename)

	if err != nil {
		return err
	}

	words := strings.Fields(string(content))
	for _, word := range words {
		mu.Lock()

		filesList, wordContains := s.wordFilesMap[word]
		if wordContains {
			if !slices.Contains(filesList, fileWihtoutExt) {
				filesList = append(filesList, fileWihtoutExt)
				slices.Sort(filesList)
				s.wordFilesMap[word] = filesList
			}
		} else {
			s.wordFilesMap[word] = []string{fileWihtoutExt}
		}

		mu.Unlock()
	}

	return nil
}

// Возвращает имя файла filename без расширения.
func getFileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}
