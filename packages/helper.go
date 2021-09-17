package packages

import (
	// Add imports here
	"errors"
	"math"
	"os"
	"runtime"
	"strings"
)

const (
	B             = "b"
	KB            = "kb"
	MB            = "mb"
	GB            = "gb"
	BYTES         = 1
	KILOBYTES     = 1024
	MEGABYTES     = 1048576
	GIGABYTES     = 1073741824
)

// List type's here

var (
	unitOfMeasureDenominator uint64
)

/*
	getFileSize will return the size of the file the user enter for processing
*/
func getFileSize(fileName string) (fileSize uint64) {

	if file, err := os.Stat(fileName); err != nil {
		panic(err.Error())
	} else {
		fileSize = uint64(file.Size())
	}

	return
}

/*
	getSystemMemory will return the current heap idle value, which is the memory available for other purposes.
*/
func getSystemMemory() (memoryAvailable uint64) {

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	memoryAvailable = m.HeapIdle

	return
}

/*
	validateMemoryUnitOfMeasure checks that the Unit of measure provide matches one of the expected values.
*/
func validateMemoryUnitOfMeasure(unitOfMeasure string) (err error) {

	switch strings.ToLower(unitOfMeasure) {
	case B:
		unitOfMeasureDenominator = BYTES
	case KB:
		unitOfMeasureDenominator = KILOBYTES
	case MB:
		unitOfMeasureDenominator = MEGABYTES
	case GB:
		unitOfMeasureDenominator = GIGABYTES
	default:
		err = errors.New("only b, kb, mb, or gb are acceptable values (case insensitive)")
	}

	return
}

/*
	convertBytesTo takes the size in bytes and converts it to the unit of measure provided rounded up
*/
func ConvertBytesTo(sizeIn uint64, unitOfMeasure string) (sizeOut float64) {

	if err := validateMemoryUnitOfMeasure(unitOfMeasure); err != nil {
		panic(err.Error())
	}

	sizeOut = math.Ceil(float64(sizeIn) / float64(unitOfMeasureDenominator))

	return
}

/*
	sizeWithinLimit will return an error if the size is greater than the limitSize
*/
func sizeWithinLimit(size, limitSize uint64) (result bool) {

	if size > limitSize {
		result = false
	} else {
		result = true
	}

	return
}
