package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
)

func generateFile(filename string, data []byte) error {
	data = append(bytes.TrimSpace(data), '\n')
	if err := os.MkdirAll(filepath.Dir(filename), os.ModePerm); err != nil {
		return err
	}
	return ioutil.WriteFile(filename, data, os.ModePerm)
}
