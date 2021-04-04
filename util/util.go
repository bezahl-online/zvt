package util

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/bezahl-online/zvt/instr"
)

// Save saves data to persistence
func Save(data *[]byte, i *instr.CtrlField, sender string) string {
	if len(os.Getenv("ZVT_DUMPFILEPATH")) == 0 {
		return ""
	}
	ctrlField := []byte{(*i).Class, (*i).Instr}
	ms := time.Now().UnixNano() / int64(time.Millisecond)
	dumpfilePath := ENVFilePath("ZVT_DUMPFILEPATH", fmt.Sprintf("%d%s.hex", ms, sender))
	dump := ctrlField
	dump = append(dump, *data...)
	err := ioutil.WriteFile(dumpfilePath, dump, 0644)
	if err != nil {
		log.Fatal(err)
	}
	return dumpfilePath
}

// Load returns data from persitence
func Load(fileName string) ([]byte, error) {
	return ioutil.ReadFile(fileName)
}

func ENVFilePath(envName string, fileName string) string {
	var filePath string
	if v := os.Getenv(envName); len(v) > 0 {
		filePath = v
	} else {
		log.Fatalf("please set environment varibalbe '%s'", envName)
	}
	filePath = strings.TrimRight(filePath, "/")
	filePath += "/" + fileName
	return filePath
}
