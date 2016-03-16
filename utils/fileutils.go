package fileutils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Exists returns whether or not the given file exists
func Exists(p string) bool {
	if _, err := os.Stat(p); err == nil {
		return true
	}
	return false
}

// WriteData writes data to the given path, creating the file if necessary.
func WriteData(path string, data []byte) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		fmt.Println(err)
	}
	// TODO: check for Close() errors
	// https://github.com/ncw/swift/blob/master/swift.go#L170
	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		fmt.Println(err)
	}
}

// ReadData returns file data as an array of lines from the file.
func ReadData(p string) *[]string {
	f, err := os.Open(p)
	if err != nil {
		return nil
	}
	defer f.Close()

	lines := []string{}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil
	}

	return &lines
}

// RemoveLine searches for the line that starts with the given value and,
// if found, removes it and saves the updated file.
func RemoveLine(p, startsWith string) {
	f, err := os.Open(p)
	if err != nil {
		return
	}
	defer f.Close()

	var outText string
	found := false

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), startsWith) {
			found = true
		} else {
			outText += scanner.Text() + string('\n')
		}
	}

	if err := scanner.Err(); err != nil {
		return
	}

	if found {
		WriteData(p, []byte(outText))
	}
}

// FindLine searches the given file for a line that begins with the given
// string.
func FindLine(p, startsWith string) string {
	f, err := os.Open(p)
	if err != nil {
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), startsWith) {
			return scanner.Text()
		}
	}

	if err := scanner.Err(); err != nil {
		return ""
	}

	return ""
}
