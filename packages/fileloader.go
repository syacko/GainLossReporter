package packages

import (
	"bufio"
	"os"
)

// List type's here

func init() {
}

func LoadFileAsString(fileName string) (file string, err error) {

	var (
		bytesFile []byte
	)
	bytesFile, err = LoadFileAsBytes(fileName)
	file = string(bytesFile)

	return
}

func LoadFileAsBytes(fileName string) (file []byte, err error) {

	if file, err = os.ReadFile(fileName); err != nil {
		panic(err.Error())
	}

	return
}

func LoadFileAsStrings(fileName string) (file []string, err error) {

	var (
		fileHandle *os.File
	)

	if fileHandle, err = os.Open(fileName); err != nil {
		panic(err.Error())
	}
	defer fileHandle.Close()

	scanner := bufio.NewScanner(fileHandle)
	for scanner.Scan() {
		file = append(file, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		panic(err.Error())
	}

	return
}