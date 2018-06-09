package tagmap

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type JSONFileStore struct {
	path string
}

func NewJSONFileStore(path string) *JSONFileStore {
	if path == "" {
		panic("path required")
	}

	return &JSONFileStore{
		path: path,
	}
}

func (tmf *JSONFileStore) Load() (*TM, error) {

	if _, err := os.Stat(tmf.path); err == nil {
		f, err := os.Open(tmf.path)
		defer f.Close()

		if err != nil {
			return nil, err
		}

		var tm TM
		dec := json.NewDecoder(f)
		err = dec.Decode(&tm)
		if err != nil {
			return nil, err
		}

		return &tm, nil
	}

	return New(), nil
}

func (tmf *JSONFileStore) Put(tm *TM) error {

	tm = tm.Normalize()

	jsonBytes, err := json.Marshal(tm)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(tmf.path, jsonBytes, 0755)
	if err != nil {
		return err
	}

	return nil
}
