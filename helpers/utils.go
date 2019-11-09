package helpers

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/fatih/color"
)

// SearchAndReplaceFiles searches all files in fullPath and replace
func SearchAndReplaceFiles(fullPath string, replacers map[string]string) error {
	fileOrDirList := []string{}
	err := filepath.Walk(fullPath,
		func(path string, f os.FileInfo, err error) error {
			fileOrDirList = append(fileOrDirList, path)
			return nil
		})

	if err != nil {
		return err
	}

	for _, fileOrDir := range fileOrDirList {
		fileInfo, _ := os.Stat(fileOrDir)
		if !fileInfo.IsDir() {
			for oldString, newString := range replacers {
				contentBytes, _ := ioutil.ReadFile(fileOrDir)
				newContentBytes :=
					bytes.Replace(contentBytes, []byte(oldString), []byte(newString), -1)

				err := ioutil.WriteFile(fileOrDir, newContentBytes, fileInfo.Mode())
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// ExitOnError will print out the error and exit
func ExitOnError(err error) {
	if err != nil {
		c := color.New(color.FgRed)
		c.Printf("%s\n", err.Error())
		os.Exit(1)
	}
}
