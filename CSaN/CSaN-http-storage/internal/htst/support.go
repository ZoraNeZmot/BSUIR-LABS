package htst

import (
	"errors"
	"htst/pkg/config"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

/*
returns true if the URL is valid
*/
func AdminCheck(url string) bool {
	return config.GetConfig().IS_DEBUG || checkPath(url)
}

func checkPath(url string) bool {
	dotCount := 0
	for _, el := range url {
		if el == '/' {
			if dotCount == 2 {
				return false
			}
			dotCount = 0
		} else {
			if el == '.' {
				dotCount++
			} else {
				dotCount = 0
			}
		}
	}
	// DEBUG: Check if URL has / in the end
	return dotCount != 2
}

func PrepareRequest(r *http.Request) (int, error) {
	var err error
	r.URL.Path, err = url.QueryUnescape(r.URL.Path)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	path := filepath.Join(config.GetConfig().App.HTST.Root, r.URL.Path)

	// check for HEAD directory request
	if r.Method == http.MethodHead {
		stat, err := os.Stat(path)
		if err != nil {
			return http.StatusNotFound, err
		}
		if stat.IsDir() {
			return http.StatusMethodNotAllowed, errors.New("method not allowed")
		}
	}
	return http.StatusOK, nil
}
