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

// Searcher позволяет искать ключевое слово в файлах.
type Searcher struct {
	FS           fs.FS               // Файловая система, по которой осуществляется поиск.
	wordFilesMap map[string][]string // Словарь соответствия слов и файлов. Данные хранятся в формате {"key1": ["file1","file2"]}
}

// Создаёт экземпляр объекта поисковика и инициализирует словарь для поиска.
func NewSearcher(f fs.FS) (*Searcher, error) {
	s := &Searcher{
		FS: f,
	}
	err := s.init()
	return s, err
}

// Производит поиск слова word по файлам в файловой системе.
func (s *Searcher) Search(word string) ([]string, error) {
	// Если wordFilesMap не создан, значит объект Searcher был создан неправильно
	if s.wordFilesMap == nil {
		return nil, errors.New("ошибка при инициализации словаря")
	}

	files := s.wordFilesMap[word]

	return files, nil
}

// Инициализирует словарь для поиска.
// Словарь послужит индексом, чтобы сам поиск слова выполнялся за O(1).
func (s *Searcher) init() error {
	s.wordFilesMap = make(map[string][]string)

	files, err := dir.FilesFS(s.FS, ".")

	if err != nil {
		s.wordFilesMap = nil
		return err
	}

	ch := make(chan error) // Канал, куда метод processFile будет писать ошибки обработки, если они будут
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}

	wg.Add(len(files))
	for _, file := range files {
		go s.processFile(file, mu, wg, ch)
	}

	// Дождёмся обработки всех файлов и закроем канал для ошибок обработки.
	wg.Wait()
	close(ch)

	// Вернём первую ошибку обработки. Если их не было - в цикл не зайдём.
	for err = range ch {
		s.wordFilesMap = nil
		return err
	}

	return nil
}

// Обрабатывает слова из файла filename. Использует [sync.Mutex] mu и [sync.WaitGroup] для синхронизации с другими обработчиками.
// В случае возникновения ошибок - отправляет ошибку в канал ch и завершает работу.
func (s *Searcher) processFile(filename string, mu *sync.Mutex, wg *sync.WaitGroup, ch chan error) {
	defer wg.Done()

	fileWihtoutExt := getFileNameWithoutExtension(filename)
	content, err := fs.ReadFile(s.FS, filename)

	if err != nil {
		ch <- err
		return
	}

	// Нарезаем весь текст в файле на отдельные слова по пробелам.
	words := strings.Fields(string(content))
	for _, word := range words {
		// Работа со словарём должна быть атомарной - после чтения из словаря нельзя в него писать, пока мы тут слово не обработаем.
		// Так что лочимся перед чтением.
		mu.Lock()

		filesList, wordContains := s.wordFilesMap[word]
		// Если слово уже есть в словаре, проверим, есть ли запись о том, что оно встречается в этом файле.
		if wordContains {
			// Если нет - добавим.
			if !slices.Contains(filesList, fileWihtoutExt) {
				filesList = append(filesList, fileWihtoutExt)
				// Как таковой необходимости сортировать нет, решил добавить, чтобы получать всегда идентичные резульаты.
				slices.Sort(filesList)
				s.wordFilesMap[word] = filesList
			}
		} else {
			// Если слова нет в словаре, добавим запись о том, что оно встречается в этом файле.
			s.wordFilesMap[word] = []string{fileWihtoutExt}
		}

		// Разлочимся после записи.
		mu.Unlock()
	}
}

// Возвращает имя файла filename без расширения.
func getFileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}
