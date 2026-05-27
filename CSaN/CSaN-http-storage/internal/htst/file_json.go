package htst

import (
	"encoding/json"
	"time"
)

type FileJSONList []FileJSON

func (list FileJSONList) MarshalJSON() ([]byte, error) {
	res := []byte("[")
	for _, file := range list {
		json, err := json.Marshal(file)
		if err != nil {
			return nil, err
		}
		res = append(res, json...)
		res = append(res, ',')
	}
	if len(res) > 1 {
		res = res[:len(res)-1]
	}
	res = append(res, ']')
	return res, nil
}

type FileJSON struct {
	Name         string    `json:"name"`
	FullPath     string    `json:"full_path"`
	IsDirectory  bool      `json:"is_directory"`
	Size         int64     `json:"size"`
	LastModified time.Time `json:"last_modified"`
}
