package os

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

const (
	IsRunningOnWindows = runtime.GOOS == "windows"
	maxFileSizeInMb    = 10
)

func Stat(path string) (os.FileInfo, error) {
	log.Printf("[DEBUG] reading the %s file info", path)
	return os.Stat(path)
}

func Getenv(name string) string {
	log.Printf("[DEBUG] reading the %s environmental variable", name)
	return os.Getenv(name)
}

func LookupEnv(name string) (string, bool) {
	log.Printf("[DEBUG] reading the %s environmental variable", name)
	return os.LookupEnv(name)
}

// ReadFileSafe checks if a file is safe to read, and then reads it.
func ReadFileSafe(path string) ([]byte, error) {
	if err := validateFile(path); err != nil {
		return nil, err
	}
	return readFile(path)
}

func readFile(path string) ([]byte, error) {
	log.Printf("[DEBUG] reading the %s file", path)
	return os.ReadFile(path)
}

func validateFile(path string) error {
	fileinfo, err := Stat(path)
	if err != nil {
		return fmt.Errorf("reading information about the config file: %w", err)
	}
	if fileinfo.Size() > maxFileSizeInMb*1024*1024 {
		return fmt.Errorf("config file %s is too big - maximum allowed size is %dMB", path, maxFileSizeInMb)
	}
	return nil
}

func UserHomeDir() (string, error) {
	log.Printf("[DEBUG] reading the user home directory location from the operating system")
	return os.UserHomeDir()
}
