package util

import (
	"fmt"
	"io/ioutil"
	"time"
)

var fileNr int = 0

// Save saves data to persistence
func Save(data *[]byte, length int, sender string) (string, error) {
	_, _, d := time.Now().Date()
	h, m, s := time.Now().Clock()
	ms := time.Now().UnixNano() / int64(time.Millisecond)
	fileName := fmt.Sprintf("dump/%s%02d%02d%02d%02d%03d.hex", sender, d, h, m, s, ms)
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
