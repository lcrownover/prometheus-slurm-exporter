package util

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"runtime"
)

// getTestDataDir returns the path to the `testdata` directory in the project
func getTestDataDir() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf(dir + "/testdata/")
}

// GetTestDataFilePath returns the full filepath to the specified filename
// of test data.
//
// Example: GetTestDataFilePath("SomeTestData.json") ->
//
//	/home/me/prometheus-slurm-exporter/testdata/SomeTestData.json
func GetTestDataFilePath(filename string) string {
	testDataDir := getTestDataDir()
	return fmt.Sprintf("%s/%s", testDataDir, filename)
}

// ReadTestDataBytes takes the short filename of the desired test data file
// and returns that data as bytes.
func ReadTestDataBytes(filename string) []byte {
	filepath := GetTestDataFilePath(filename)
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("failed to open file: %v\n", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("failed to read file: %v\n", err)
	}

	return data
}
