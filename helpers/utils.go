package helpers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/kardianos/osext"
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
				var newContentBytes []byte
				if strings.HasPrefix(oldString, "REGEXP-") {
					regStr := strings.Split(oldString, "REGEXP-")[1]
					r, err := regexp.Compile(regStr)
					if err != nil {
						panic(err)
					}
					newContentBytes =
						r.ReplaceAll(contentBytes, []byte(newString))
				} else {
					newContentBytes =
						bytes.Replace(contentBytes, []byte(oldString), []byte(newString), -1)
				}

				if strings.HasPrefix(fileOrDir, ".git") == false {
					err := ioutil.WriteFile(fileOrDir, newContentBytes, fileInfo.Mode())
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func goPaths() []string {
	return strings.Split(os.Getenv("GOPATH"), ":")
}

// GetLatestBaseFolder helps to identify what is the current base folder
func GetLatestBaseFolder() (string, error) {
	executableDir, err := osext.ExecutableFolder()
	ExitOnError(err)

	ret := filepath.Join(executableDir, "templates")
	if _, err = os.Stat(ret); err == nil {
		return executableDir, nil
	}

	currDir, err := os.Getwd()
	ret = filepath.Join(currDir, "templates")
	if _, err = os.Stat(ret); err == nil {
		return currDir, nil
	}

	base := filepath.Join("pkg", "mod", "github.com", "jameshuynh")
	srcDir := filepath.Join(filepath.Dir(executableDir), base)
	files, _ := ioutil.ReadDir(srcDir)

	latestFolder := ""
	for _, f := range files {
		if strings.Contains(f.Name(), "corn@v") &&
			!strings.Contains(f.Name(), "templates") {
			latestFolder = f.Name()
		}
	}

	ret = filepath.Join(srcDir, latestFolder)

	if _, err = os.Stat(ret); err == nil && latestFolder != "" {
		return ret, nil
	}

	base = filepath.Join("src", "github.com", "jameshuynh")
	srcDir = filepath.Join(filepath.Dir(executableDir), base)
	files, _ = ioutil.ReadDir(srcDir)

	latestFolder = ""
	for _, f := range files {
		if strings.Contains(f.Name(), "corn") &&
			!strings.Contains(f.Name(), "templates") {
			latestFolder = f.Name()
		}
	}

	ret = filepath.Join(srcDir, latestFolder)

	if _, err = os.Stat(ret); err == nil && latestFolder != "" {
		return ret, nil
	}

	return "", fmt.Errorf("Unable to find corn's directory")
}

// ExitOnError will print out the error and exit
func ExitOnError(err error) {
	if err != nil {
		c := color.New(color.FgRed)
		c.Printf("%s\n", err.Error())
		os.Exit(1)
	}
}
