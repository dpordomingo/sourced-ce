// Package dir provides functions to manage the $HOME/.sourced and /tmp/srcd
// directories
package dir

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

var BaseDir = ".sourced"

// Path returns the absolute path for $HOME/.sourced
func Path() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrap(err, "could not detect home directory")
	}

	return filepath.Join(homedir, BaseDir), nil
}

type ComposeGetMethodType int

const (
	ComposeGetRemote ComposeGetMethodType = iota
	ComposeGetLocal
)

var ComposeGetMethod = ComposeGetRemote

// DownloadURL downloads the given url to a file to the
// dst path, creating the directory if it's needed
func DownloadURL(url, dst string) error {

	if err := os.MkdirAll(filepath.Dir(dst), os.ModePerm); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	var source io.ReadCloser
	if ComposeGetMethod == ComposeGetRemote {
		resp, err := http.Get(url)
		if err != nil {
			return err
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("HTTP status %v", resp.Status)
		}

		source = resp.Body
	} else {
		url = os.Getenv("DOCKER_COMPOSE_DEFAULT_PATH")
		if url == "" {
			return fmt.Errorf("DOCKER_COMPOSE_DEFAULT_PATH was not set")
		}

		sourceFileStat, err := os.Stat(url)
		if err != nil {
			return err
		}

		if !sourceFileStat.Mode().IsRegular() {
			return fmt.Errorf("%s is not a regular file", url)
		}

		source, err = os.Open(url)
		if err != nil {
			return err
		}

		defer source.Close()
	}

	_, err = io.Copy(out, source)
	return err
}

// TmpPath returns the absolute path for /tmp/srcd
func TmpPath() string {
	return filepath.Join(os.TempDir(), "srcd")
}
