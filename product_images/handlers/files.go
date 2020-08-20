package handlers

import (
	"github.com/alonelegion/go_microservices/product_images/files"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"net/http"
	"path/filepath"
)

// Files is a handler for reading and writing files
// Files - обработчик для чтения и записи файлов
type Files struct {
	log   hclog.Logger
	store files.Storage
}

// NewFiles creates a new File handler
// NewFiles создает новый обработчик файлов
func NewFiles(s files.Storage, l hclog.Logger) *Files {
	return &Files{store: s, log: l}
}

// ServeHTTP implements the http.Handler interface
// ServeHTTP реализует интерфейс http.Handler
func (f *Files) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]
	fn := vars["filename"]

	f.log.Info("Handler POST", "id", id, "filename", fn)

	// no need to check for invalid id or filename as the mux router
	// will not send requests here unless they have the correct parameters
	// Нет необходимости проверять недопустимый идентификатор или имя файла,
	// так как mux router не будет отправлять запросы здесь, если у них нет
	// правильных параметров

	f.saveFile(id, fn, w, req)
}

// saveFile saves the contents of the request to a file
// saveFile сохраняет содержимое запроса в файл
func (f *Files) saveFile(id, path string, w http.ResponseWriter, req *http.Request) {
	f.log.Info("Save file for product", "id", id, "path", path)

	fp := filepath.Join(id, path)
	err := f.store.Save(fp, req.Body)
	if err != nil {
		f.log.Error("Unable to save file", "error", err)
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
	}
}

func (f *Files) invalidURI(uri string, w http.ResponseWriter) {
	f.log.Error("Invalid path", "path", uri)
	http.Error(w, "Invalid file path should be in the format: /[id]/[filepath]", http.StatusBadRequest)
}
