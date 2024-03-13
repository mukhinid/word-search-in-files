package searcher

import (
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

// Производит поиск слова word по файлам в папке файловой системы.
func (s *Searcher) Search(word string) ([]string, error) {
	if s.wordFilesMap == nil {
		err := s.init()
		if err != nil {
			return nil, err
		}
	}

	files := s.wordFilesMap[word]

	return files, nil
}

// Инициализирует словарь для поиска.
func (s *Searcher) init() error {
	s.wordFilesMap = make(map[string][]string)

	files, err := dir.FilesFS(s.FS, ".")

	if err != nil {
		return err
	}

	ch := make(chan error)
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}

	wg.Add(len(files))
	for _, file := range files {
		go s.processFile(file, mu, wg, ch)
	}

	wg.Wait()
	close(ch)

	for err = range ch {
		return err
	}

	return nil
}

func (s *Searcher) processFile(filename string, mu *sync.Mutex, wg *sync.WaitGroup, ch chan error) {
	defer wg.Done()

	fileWihtoutExt := getFileNameWithoutExtension(filename)
	content, err := fs.ReadFile(s.FS, filename)

	if err != nil {
		ch <- err
		return
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
}

// Возвращает имя файла filename без расширения.
func getFileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}
