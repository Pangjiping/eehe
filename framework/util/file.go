package util

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// fileExist returns true if given path exists.
func fileExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return false
}

// IsHiddenDir returns true if given path is hidden path.
func IsHiddenDir(path string) bool {
	return len(path) > 1 && strings.HasPrefix(filepath.Base(path), ".")
}

// SubDir returns all subDirectory name of folder.
func SubDir(folder string) ([]string, error) {
	subs, err := ioutil.ReadDir(folder)
	if err != nil {
		return nil, err
	}

	subDirs := []string{}
	for _, sub := range subs {
		if sub.IsDir() {
			subDirs = append(subDirs, sub.Name())
		}
	}
	return subDirs, nil
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

// CopyFolder copies a folder from src to dest.
func CopyFolder(src, dest string) error {
	var err error = filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		var relPath string = strings.Replace(path, src, "", 1)
		if relPath == "" {
			return nil
		}
		if info.IsDir() {
			return os.Mkdir(filepath.Join(dest, relPath), 0755)
		} else {
			var data, err1 = ioutil.ReadFile(filepath.Join(src, relPath))
			if err1 != nil {
				return err1
			}
			return ioutil.WriteFile(filepath.Join(dest, relPath), data, 0777)
		}
	})
	return err
}

// CopyFile copies a file from src to dest.
func CopyFile(src, dest string) error {
	var data, err1 = ioutil.ReadFile(src)
	if err1 != nil {
		return err1
	}
	return ioutil.WriteFile(dest, data, 0777)
}
