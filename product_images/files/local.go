package files

import (
	"golang.org/x/xerrors"
	"io"
	"os"
	"path/filepath"
)

// Local is an implementation of the Storage interface which works with the
// local disk on the current machine
// Local - это рефлизация интерфейса Storage, который работает с
// локальным диском на текущей машине
type Local struct {
	maxFileSize int // maximum number of bytes for files.
	// Максимальное кол-во байтов для файлов
	basePath string
}

// NewLocal creates a new Local filesystem with the given base path
// basePath is the base directory to save files to
// maxSize is the max number of bytes that a file can be
// NewLocal создает новую локальную файловую систему с заданым базовым путем
// basePath - это базовый каталог для сохранения файлов
// maxSize - это максимальное количество байтов которое может быть в файле
func NewLocal(basePath string, maxSize int) (*Local, error) {
	p, err := filepath.Abs(basePath)
	if err != nil {
		return nil, err
	}

	return &Local{basePath: p}, nil
}

// Save the contents of the Writer to the given path
// path is a relative path, basePath will be appended
// Save - содержит Writer по заданному пути
// path - относительный путь, к нему будет добавлен basePath
func (l *Local) Save(path string, contents io.Reader) error {
	// get the full path for the file
	// получить полный путь к файлу
	fp := l.fullPath(path)

	// get the directory and make sure it exists
	// получить каталог и убедиться что он существует
	d := filepath.Dir(fp)
	err := os.MkdirAll(d, os.ModePerm)
	if err != nil {
		return xerrors.Errorf("Unable to create directory: %w", err)
	}

	// if the file exists delete it
	// если файл существует, удалить его
	_, err = os.Stat(fp)
	if err == nil {
		err = os.Remove(fp)
		if err != nil {
			return xerrors.Errorf("Unable to delete file: %w", err)
		}
	} else if !os.IsNotExist(err) {
		// if this is anything other than a not exists error
		// если это что-то другое, кроме несуществующей ошибки
		return xerrors.Errorf("Unable to get file: %w", err)
	}

	// create a new file at the path
	// создать новый файл согласно пути
	f, err := os.Create(fp)
	if err != nil {
		return xerrors.Errorf("Unable to create file: %w", err)
	}
	defer f.Close()

	// write the contents to the new file
	// ensure that we are not writing greater than max bytes
	// записать содержимое в новый файл
	// гарантия того, что не пишется больше, чем максимальное кол-во байтов
	_, err = io.Copy(f, contents)
	if err != nil {
		return xerrors.Errorf("Unable to write to file: %w", err)
	}

	return nil
}

// Get the file at the given path and return a Reader
// the calling function is responsible for closing the reader
// Получить файл по заданому пути и вернуть Reader
// вызывающая функция отвечает за закрытие ридера
func (l *Local) Get(path string) (*os.File, error) {
	// get the full path for the file
	// получить полный путь к файлу
	fp := l.fullPath(path)

	// open the file
	// открытие файла
	f, err := os.Open(fp)
	if err != nil {
		return nil, xerrors.Errorf("Unable to open file: %w", err)
	}

	return f, nil
}

// returns the absolute path
// возвращает абсолютный путь
func (l *Local) fullPath(path string) string {
	// append the given path to the base path
	return filepath.Join(l.basePath, path)
}
