// Package oswrapper is a wrapper around the standard os package that allows more secure interactions with the operating system.
// It should be used as a replacement in production code of the standard os package.
package oswrapper

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"runtime"
)

const (
	maxFileSizeInMb = 10
)

// IsRunningOnWindows returns true if the code is running on Windows.
func IsRunningOnWindows() bool {
	return runtime.GOOS == "windows"
}

// Stat is an os.Stat wrapper.
func Stat(path string) (os.FileInfo, error) {
	log.Printf("[DEBUG] Reading the %s file info", path)
	return os.Stat(path)
}

// Getenv is an os.Getenv wrapper.
func Getenv(name string) string {
	log.Printf("[DEBUG] Reading the %s environmental variable", name)
	return os.Getenv(name)
}

// LookupEnv is an os.LookupEnv wrapper.
func LookupEnv(name string) (string, bool) {
	log.Printf("[DEBUG] Reading the %s environmental variable", name)
	return os.LookupEnv(name)
}

// ReadFileSafe checks if a file is safe to read, and then reads it.
func ReadFileSafe(path string) ([]byte, error) {
	if err := fileIsSafeToRead(path); err != nil {
		return nil, err
	}
	return readFile(path)
}

// ReadFileSkipPermissionsCheck reads a file without checking its permissions.
func ReadFileSkipPermissionsCheck(path string) ([]byte, error) {
	if err := fileIsTooBig(path); err != nil {
		return nil, err
	}
	return readFile(path)
}

func readFile(path string) ([]byte, error) {
	log.Printf("[DEBUG] Reading the %s file", path)
	return os.ReadFile(path)
}

func fileIsSafeToRead(path string) error {
	fileInfo, err := Stat(path)
	if err != nil {
		return fmt.Errorf("reading information about the config file: %w", err)
	}
	if fileInfo.Size() > maxFileSizeInMb*1024*1024 {
		return fmt.Errorf("config file %s is too big - maximum allowed size is %dMB", path, maxFileSizeInMb)
	}
	if !IsRunningOnWindows() {
		if !unixFilePermissionsAreStrict(fileInfo.Mode().Perm()) {
			return fmt.Errorf("config file %s has unsafe permissions - %#o", path, fileInfo.Mode().Perm())
		}
	} else {
		log.Println("[DEBUG] Skipped checking file permissions on a Windows system")
	}
	return nil
}

func fileIsTooBig(path string) error {
	fileInfo, err := Stat(path)
	if err != nil {
		return fmt.Errorf("reading information about the config file: %w", err)
	}
	if fileInfo.Size() > maxFileSizeInMb*1024*1024 {
		return fmt.Errorf("config file %s is too big - maximum allowed size is %dMB", path, maxFileSizeInMb)
	}
	return nil
}

func unixFilePermissionsAreStrict(perm fs.FileMode) bool {
	log.Println("[DEBUG] Checking file permissions on a Unix system...")
	// group or others have any access
	unsafeBits := os.FileMode(0o077)
	return perm&unsafeBits == 0
}

// UserHomeDir is an os.UserHomeDir wrapper.
func UserHomeDir() (string, error) {
	log.Printf("[DEBUG] Reading the user home directory location from the operating system")
	return os.UserHomeDir()
}
