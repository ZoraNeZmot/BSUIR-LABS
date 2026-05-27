package htst

import (
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	FileName          = "HTST-File-Name"
	FilePath          = "HTST-File-Path"
	FileSize          = "HTST-File-Size"
	FileIsDirectory   = "HTST-File-Is-Directory"
	FileLastModified  = "HTST-File-Last-Modified"
	FileContentLength = "HTST-File-Content-Length"

	IsPutDirectory = "HTST-Is-Put-Directory"
	IsMove         = "HTST-Is-Move"
)

func SetHeaderFields(w http.ResponseWriter, file os.FileInfo, path string) {
	w.Header().Set(FileName, file.Name())
	w.Header().Set(FilePath, path)
	w.Header().Set(FileSize, strconv.FormatInt(file.Size(), 10))
	w.Header().Set(FileIsDirectory, strconv.FormatBool(file.IsDir()))
	w.Header().Set(FileLastModified, file.ModTime().Format(time.RFC3339))
}
