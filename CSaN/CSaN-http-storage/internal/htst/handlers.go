package htst

import (
	"encoding/json"
	"errors"
	"htst/pkg/config"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

func HTSTHandler(w http.ResponseWriter, r *http.Request) {
	if !AdminCheck(r.URL.Path) {
		w.WriteHeader(http.StatusForbidden)
		logrus.Error("Admin check failed: ", r.URL.Path)
		return
	}
	status, err := PrepareRequest(r)
	if err != nil {
		w.WriteHeader(status)
		logrus.WithError(err).Info("Failed to prepare request: ", r.URL.Path)
		return
	}

	switch r.Method {
	case http.MethodGet:
		HTSTGetHandler(w, r)
	case http.MethodHead:
		HTSTHeadHandler(w, r)
	case http.MethodDelete:
		HTSTDeleteHandler(w, r)
	case http.MethodPut:
		HTSTPutHandler(w, r)
	case http.MethodOptions:
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		logrus.Error("Method not allowed: ", r.Method)
		return
	}
}

func HTSTGetHandler(w http.ResponseWriter, r *http.Request) {
	config := config.GetConfig()
	path := filepath.Join(config.App.HTST.Root, r.URL.Path)

	stat, err := os.Stat(path)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		logrus.WithError(err).Info("Failed to stat file: ", r.URL.Path)
		return
	}
	if stat.IsDir() {
		HTSTGetDirectoryHandler(w, r, path)
	} else {
		HTSTGetFileHandler(w, r, path)
	}
}

func HTSTHeadHandler(w http.ResponseWriter, r *http.Request) {
	config := config.GetConfig()
	path := filepath.Join(config.App.HTST.Root, r.URL.Path)
	stat, err := os.Stat(path)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		logrus.WithError(err).Info("Failed to stat file: ", r.URL.Path)
		return
	}
	SetHeaderFields(w, stat, path)
	w.WriteHeader(http.StatusOK)
	logrus.Info("HEAD request for file: ", r.URL.Path)
}

func HTSTGetDirectoryHandler(w http.ResponseWriter, r *http.Request, path string) {
	// config := config.GetConfig()
	logrus.Info("GET request for directory: ", r.URL.Path)
	files, err := os.ReadDir(path)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.WithError(err).Info("Failed to read directory: ", r.URL.Path)
		return
	}
	fileList := make(FileJSONList, len(files))
	for i, file := range files {
		fileStat, err := file.Info()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logrus.WithError(err).Info("Failed to get file info: ", file.Name())
			return
		}
		fileList[i] = FileJSON{
			Name:         fileStat.Name(),
			FullPath:     filepath.Join(path, fileStat.Name()),
			IsDirectory:  fileStat.IsDir(),
			Size:         fileStat.Size(),
			LastModified: fileStat.ModTime(),
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	bytes, err := json.Marshal(fileList)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.WithError(err).Error("Failed to encode file list: ", r.URL.Path)
		return
	}
	_, err = w.Write(bytes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.WithError(err).Error("Failed to write file list: ", r.URL.Path)
		return
	}
}

func HTSTGetFileHandler(w http.ResponseWriter, r *http.Request, path string) {
	logrus.Info("GET request for file: ", r.URL.Path)
	file, err := os.Open(path)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		logrus.WithError(err).Info("Failed to open file: ", r.URL.Path)
		return
	}
	defer file.Close()

	w.WriteHeader(http.StatusOK)
	_, err = io.Copy(w, file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.WithError(err).Error("Failed to copy file: ", r.URL.Path)
		return
	}
}

func HTSTDeleteHandler(w http.ResponseWriter, r *http.Request) {
	config := config.GetConfig()
	path := filepath.Join(config.App.HTST.Root, r.URL.Path)
	stat, err := os.Stat(path)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		logrus.WithError(err).Info("Failed to stat file: ", r.URL.Path)
		return
	}
	if stat.IsDir() {
		err := os.RemoveAll(path)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logrus.WithError(err).Error("Failed to delete directory: ", r.URL.Path)
			return
		}
		w.WriteHeader(http.StatusOK)
		logrus.Info("Deleted directory: ", r.URL.Path)
	} else {
		err := os.Remove(path)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logrus.WithError(err).Error("Failed to delete file: ", r.URL.Path)
			return
		}
		w.WriteHeader(http.StatusOK)
		logrus.Info("Deleted file: ", r.URL.Path)
	}
}

func HTSTPutHandler(w http.ResponseWriter, r *http.Request) {
	config := config.GetConfig()
	path := filepath.Join(config.App.HTST.Root, r.URL.Path)

	if r.Header.Get(IsPutDirectory) != "" {
		logrus.Info("PUT request for directory: ", r.URL.Path)
		w.WriteHeader(http.StatusCreated)
		err := os.MkdirAll(path, 0777)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logrus.WithError(err).Error("Failed to create directory: ", r.URL.Path)
			return
		}
	} else {
		if r.Header.Get(IsMove) != "" {
			destPath := filepath.Join(config.App.HTST.Root, r.Header.Get(IsMove))
			logrus.Info("MOVE request for : ", r.URL.Path)
			err := os.Rename(path, destPath)
			if errors.Is(err, os.ErrNotExist) {
				w.WriteHeader(http.StatusNotFound)
				logrus.WithError(err).Error("Failed to move file: ", path, " to ", destPath)
				return
			} else if errors.Is(err, os.ErrPermission) {
				w.WriteHeader(http.StatusForbidden)
				logrus.WithError(err).Error("Failed to move file: ", path, " to ", destPath)
				return
			} else if errors.Is(err, os.ErrExist) {
				w.WriteHeader(http.StatusConflict)
				logrus.WithError(err).Error("Failed to move file: ", path, " to ", destPath)
				return
			} else if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				logrus.WithError(err).Error("Failed to move file: ", path, " to ", destPath)
				return
			}
			w.WriteHeader(http.StatusOK)
			logrus.Infof("Moved file: %s to %s", r.URL.Path, destPath)
		} else {
			logrus.Info("PUT request for file: ", r.URL.Path)
			file, err := os.Create(path)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				logrus.WithError(err).Error("Failed to create file: ", r.URL.Path)
				return
			}
			defer file.Close()
			w.WriteHeader(http.StatusCreated)
			_, err = io.Copy(file, r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				logrus.WithError(err).Error("Failed to copy file: ", r.URL.Path)
				return
			}
		}

	}
}
