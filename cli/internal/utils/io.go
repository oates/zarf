package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"

	"github.com/defenseunicorns/zarf/cli/internal/message"
	"github.com/otiai10/copy"
	"github.com/pterm/pterm"
)

var TempPathPrefix = "zarf-"

func MakeTempDir() (string, error) {
	tmp, err := ioutil.TempDir("", TempPathPrefix)
	message.Debugf("Creating temp path %s", tmp)
	return tmp, err
}

// VerifyBinary returns true if binary is available
func VerifyBinary(binary string) bool {
	_, err := exec.LookPath(binary)
	return err == nil
}

// CreateDirectory creates a directory for the given path and file mode
func CreateDirectory(path string, mode os.FileMode) error {
	if InvalidPath(path) {
		return os.MkdirAll(path, mode)
	}
	return nil
}

// InvalidPath checks if the given path exists
func InvalidPath(path string) bool {
	_, err := os.Stat(path)
	return os.IsNotExist(err)
}

func ListDirectories(directory string) ([]string, error) {
	var directories []string
	paths, err := os.ReadDir(directory)
	if err != nil {
		return directories, fmt.Errorf("unable to load the directory %s: %w", directory, err)
	}

	for _, entry := range paths {
		if entry.IsDir() {
			directories = append(directories, filepath.Join(directory, entry.Name()))
		}
	}

	return directories, nil
}

func WriteFile(path string, data []byte) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("unable to create the file at %s to write the contents: %w", path, err)
	}

	_, err = f.Write(data)
	if err != nil {
		_ = f.Close()
		return fmt.Errorf("unable to write the file at %s contents:%w", path, err)
	}

	err = f.Close()
	if err != nil {
		return fmt.Errorf("error saving file %s: %w", path, err)
	}

	return nil
}

func ReplaceText(path string, old string, new string) {
	input, err := ioutil.ReadFile(path)
	if err != nil {
		message.Fatalf(err, "Unable to load %s", path)
	}

	output := bytes.Replace(input, []byte(old), []byte(new), -1)

	if err = ioutil.WriteFile(path, output, 0600); err != nil {
		message.Fatalf(err, "Unable to update %s", path)
	}
}

// RecursiveFileList walks a path with an optional regex pattern and returns a slice of file paths
func RecursiveFileList(root string, pattern *regexp.Regexp) []string {
	var files []string

	err := filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				if pattern != nil {
					if len(pattern.FindStringIndex(path)) > 0 {
						files = append(files, path)
					}
				} else {
					files = append(files, path)
				}
			}
			return nil
		})

	if err != nil {
		message.Fatalf(err, "Unable to walk the directory %s", root)
	}

	return files
}

func CreateFilePath(destination string) error {
	parentDest := path.Dir(destination)
	return CreateDirectory(parentDest, 0700)
}

func CreatePathAndCopy(source string, destination string) {
	if err := CreateFilePath(destination); err != nil {
		message.Fatalf(err, "unable to copy the file %s", source)
	}

	// Copy the asset
	if err := copy.Copy(source, destination); err != nil {
		message.Fatalf(err, "unable to copy the file %s", source)
	}
	pterm.Success.Printfln("Copying %s", source)
}
