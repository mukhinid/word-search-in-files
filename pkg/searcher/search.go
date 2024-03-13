package searcher

import (
	"io/fs"
	"path/filepath"
	"strings"
	"sync"
	"word-search-in-files/pkg/internal/dir"
)

// Searcher позволяет искать ключевое слово в файлах указанной папки.
type Searcher struct {
	FS fs.FS // Файловая система, предоставляющая папку, в файлах которой осуществляется поиск.
}

// Производит поиск слова word по файлам в папке файловой системы.
func (s *Searcher) Search(word string) ([]string, error) {
	files, err := dir.FilesFS(s.FS, ".")

	if err != nil {
		return files, err
	}

	result := make([]string, 0, len(files))

	var wg sync.WaitGroup
	ch := make(chan string)

	for _, file := range files {
		wg.Add(1)

		go func(file string) {
			searchInFile(s.FS, file, word, ch)
			wg.Done()
		}(file)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for file := range ch {
		result = append(result, getFileNameWithoutExtension(file))
	}

	return result, err
}

// Производит поиск слова word по файлу filename в файловой системе filesystem. Если слово содержится в файле - отправляет filename в канал ch.
func searchInFile(filesystem fs.FS, filename string, word string, ch chan string) {
	content, err := fs.ReadFile(filesystem, filename)

	if err != nil {
		return
	}

	words := strings.Fields(string(content))

	for _, w := range words {
		if w == word {
			ch <- filename
			return
		}
	}
}

// Возвращает имя файла filename без расширения.
func getFileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}
