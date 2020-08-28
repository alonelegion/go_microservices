package files

import "io"

// Storage defines the behavior for file operations
// Implementations may be of the time local disk, or cloud storage, etc
// Хранилище определяет поведение файловых операций
// Реализации могут быть как локальный диск, или облачное хранилищие и т.д.
type Storage interface {
	Save(path string, file io.Reader) error
}
