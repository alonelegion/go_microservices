package handlers

import (
	"github.com/alonelegion/go_microservices/product_images/files"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
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

// UploadREST implements the http.Handler interface
// UploadREST реализует интерфейс http.Handler
func (f *Files) UploadREST(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]
	fn := vars["filename"]

	f.log.Info("Handler POST", "id", id, "filename", fn)

	// no need to check for invalid id or filename as the mux router
	// will not send requests here unless they have the correct parameters
	// Нет необходимости проверять недопустимый идентификатор или имя файла,
	// так как mux router не будет отправлять запросы здесь, если у них нет
	// правильных параметров

	f.saveFile(id, fn, w, req.Body)
}

// UploadMultipart something
func (f *Files) UploadMultipart(w http.ResponseWriter, req *http.Request) {
	err := req.ParseMultipartForm(128 * 1024)
	if err != nil {
		f.log.Error("Bad request", "error", err)
		http.Error(w, "Expected multipart from data", http.StatusBadRequest)
		return
	}

	id, idErr := strconv.Atoi(req.FormValue("id"))
	f.log.Info("Process form for id", "id", id)

	if idErr != nil {
		f.log.Error("Bad request", "error", err)
		http.Error(w, "Expected integer id", http.StatusBadRequest)
		return
	}

	ff, mh, err := req.FormFile("file")
	if err != nil {
		f.log.Error("Bad request", "error", err)
		http.Error(w, "Expected file", http.StatusBadRequest)
		return
	}

	f.saveFile(req.FormValue("id"), mh.Filename, w, ff)
}

// saveFile saves the contents of the request to a file
// saveFile сохраняет содержимое запроса в файл
func (f *Files) saveFile(id, path string, w http.ResponseWriter, r io.ReadCloser) {
	f.log.Info("Save file for product", "id", id, "path", path)

	fp := filepath.Join(id, path)
	err := f.store.Save(fp, r)
	if err != nil {
		f.log.Error("Unable to save file", "error", err)
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
	}
}

func (f *Files) invalidURI(uri string, w http.ResponseWriter) {
	f.log.Error("Invalid path", "path", uri)
	http.Error(w, "Invalid file path should be in the format: /[id]/[filepath]", http.StatusBadRequest)
}
