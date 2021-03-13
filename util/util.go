package util

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/bezahl-online/zvt/instr"
)

var fileNr int = 0

// Save saves data to persistence
func Save(data *[]byte, i *instr.CtrlField, sender string) (string, error) {
	// _, _, d := time.Now().Date()
	// h, m, s := time.Now().Clock()
	ctrlField := []byte{(*i).Class, (*i).Instr}
	fmt.Printf("%s(% X): % X\n", sender, ctrlField, *data)
	ms := time.Now().UnixNano() / int64(time.Millisecond)
	fileName := fmt.Sprintf("dump/%d%s.hex", ms, sender)
	dump := ctrlField
	dump = append(dump, *data...)
	err := ioutil.WriteFile(fileName, dump, 0644)
	if err != nil {
		return "", err
	}
	return fileName, nil
}

// Load returns data from persitence
func Load(fileName string) ([]byte, error) {
	return ioutil.ReadFile(fileName)
}
