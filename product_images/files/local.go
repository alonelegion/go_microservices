package files

import "path/filepath"

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

//
