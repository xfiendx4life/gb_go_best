// Данная программа предназначена для удаления дубликатов файлов
package files_deleter

import (
	"fmt"
	"hash"
	"hash/crc32"
	"io"
	"io/fs"
	"log"
	"os"
	"sync"
)

// структура для хранения хешей и имен файлов. Ключ - хеш, значение - имя
// мьютекс для асинхронной записи и чтения
type FileSet struct {
	data map[uint32]string
	sync.Mutex
}

// инициализация структуры с хешами
func newFileSet() *FileSet {
	return &FileSet{
		data: make(map[uint32]string),
	}
}

// Записываем в мапу данные ключ - хеш, значение имя файла
func (f *FileSet) WriteData(filename string, h uint32) {
	f.Lock()
	defer f.Unlock()
	f.data[h] = filename
}

func (f *FileSet) String() string {
	f.Lock()
	defer f.Unlock()
	res := ""
	for k, v := range f.data {
		res += fmt.Sprintf("%d - %s\n", k, v)
	}
	return res
}

// читаем данные с мьютексом из мапы с хешами
func (f *FileSet) ReadData(h uint32) (filename string, ok bool) {
	f.Lock()
	defer f.Unlock()
	d, ok := f.data[h]
	return d, ok
}

// Создаем хеш, получаем на вход имя файла, читаем его и хешируем
func CreateHash(fileName string) (hash.Hash32, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("can't open file %s", err)
	}
	defer file.Close()
	h := crc32.NewIEEE()
	if _, err := io.Copy(h, file); err != nil {
		return nil, fmt.Errorf("cant get hash of the file %s", err)
	}
	return h, nil
}

// Функция делает запрос к мапе с хэшами файлов, если хэш есть, значит это дубликат
// и его нужно удалиить. Если нет, то добавляем в мапу хеш
// возвращаем ошибку. В канал отправляем имя файла если обнаружен дубликат
func FindDopplerOrWrite(filename string, fs *FileSet, res chan [2]string, errChan chan error) {
	newH, err := CreateHash(filename)
	if err != nil {
		errChan <- fmt.Errorf("check hash %s", err)
		return
	}
	// fmt.Println(fs)
	if name, ok := fs.ReadData(newH.Sum32()); ok {
		res <- [2]string{name, filename}
		return
	}

	fs.WriteData(filename, newH.Sum32())
}

func walkTheDir(dirList *[]fs.DirEntry, dirPath string, resChan chan [2]string, errChan chan error, fset *FileSet) error {
	wg := sync.WaitGroup{}
	defer wg.Wait()

	for _, file := range *dirList {
		if !file.Type().IsDir() {
			filename := fmt.Sprintf("%s/%s", dirPath, file.Name())
			wg.Add(1)
			go func(filename string, fs *FileSet, wg *sync.WaitGroup) {
				defer wg.Done()
				FindDopplerOrWrite(filename, fs, resChan, errChan)

			}(filename, fset, &wg)
		} else {
			intDirPath := fmt.Sprintf("%s/%s", dirPath, file.Name())
			internalDir, err := os.ReadDir(intDirPath)
			if err != nil {
				log.Printf("can't read directory %s", err)
				return err
			}
			err = walkTheDir(&internalDir, intDirPath, resChan, errChan, fset)
			if err != nil {
				log.Printf("error while reading dir %s", err)
				return err
			}
		}
	}
	return nil
}

func Delete(filedir string, isDel bool) (n int, err error) {
	dirList, err := os.ReadDir(filedir)
	if err != nil {
		log.Printf("can't read directory %s", err)
	}
	fset := newFileSet()
	chanBuf := len(dirList)
	resChan := make(chan [2]string, chanBuf)
	errChan := make(chan error, chanBuf)
	walkTheDir(&dirList, filedir, resChan, errChan, fset)
	for {
		select {
		case err = <-errChan:
		case names := <-resChan:
			n += 1
			fmt.Printf("Found duplicates %s, %s\n", names[0], names[1])
			if isDel {
				err = os.Remove(names[1])
				if err != nil {
					return 0, err
				}
				log.Printf("Duplicate %s deleted", names[1])
			}
		default:
			return
		}
	}
}
