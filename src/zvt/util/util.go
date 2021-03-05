package util

import (
	"fmt"
	"io/ioutil"
	"time"
)

var fileNr int = 0

// Save saves data to persistence
func Save(data *[]byte, length int) (string, error) {
	_, _, d := time.Now().Date()
	m, s, ms := time.Now().Clock()
	fileName := fmt.Sprintf("data%02d%02d%02d%03d.bin", d, m, s, ms)
	err := ioutil.WriteFile(fileName, (*data)[:length], 0644)
	if err != nil {
		return "", err
	}
	return fileName, nil
}

// Load returns data from persitence
func Load(fileName string) ([]byte, error) {
	return ioutil.ReadFile(fileName)
}
