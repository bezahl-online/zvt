package util

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
	"unicode"

	"github.com/bezahl-online/zvt/instr"
)

const DumpFilePath = "ZVT_DUMPFILEPATH"

// Save saves data to persistence
func Save(data *[]byte, i *instr.CtrlField, sender string) (string, error) {
	if len(os.Getenv(DumpFilePath)) == 0 {
		return "", fmt.Errorf("environment variable '%s' not defined", DumpFilePath)
	}
	ms := time.Now().UnixNano() / int64(time.Millisecond)
	dumpfilePath := ENVFilePath(DumpFilePath, fmt.Sprintf("%d%s.hex", ms, sender))
	ctrlField := []byte{(*i).Class, (*i).Instr}
	dump := ctrlField
	dumpLength, err := i.Length.Marshal()
	if err != nil {
		return "", err
	}
	dump = append(dump, dumpLength...)
	dump = append(dump, *data...)
	err = ioutil.WriteFile(dumpfilePath, dump, 0644)
	if err != nil {
		log.Fatal(err)
	}
	return dumpfilePath, nil
}

// Load returns data from persitence
func Load(fileName string) ([]byte, error) {
	return ioutil.ReadFile(fileName)
}

func ENVFilePath(envName string, fileName string) string {
	var filePath string
	if v := os.Getenv(envName); len(v) > 0 {
		filePath = v
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			os.Mkdir(filePath, os.ModePerm)
		}
	} else {
		log.Fatalf("please set environment varibalbe '%s'", envName)
	}
	filePath = strings.TrimRight(filePath, "/")
	filePath += "/" + fileName
	return filePath
}

func GetPureText(text string) string {
	return strings.Map(func(r rune) rune { // FIXME: not tested
		if (unicode.IsLetter(r) ||
			unicode.IsDigit(r) ||
			unicode.IsPunct(r) ||
			unicode.IsSpace(r)) &&
			r != 0x26 && r != '\r' {
			return r
		}
		return -1
	}, string(text))
}
