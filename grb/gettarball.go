package grb

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/pkg/errors"
)

func GetTarball(url, path string) error {
	if url == "" || path == "" {
		return fmt.Errorf("invalid url or download directory")
	}

	// Get the filename from our URL
	parts := strings.Split(url, "/")
	l := len(parts)
	file := parts[l-1]
	filepath := fmt.Sprintf("%s/%s", path, file)

	// Check if the file already exists
	if fileExists(filepath) {
		fmt.Fprintf(os.Stderr, "tarball already exists, skipping download\n")
		return nil
	}

	// Download the file
	resp, err := http.Get(url)
	if err != nil {
		return errors.Wrap(err, "failed to download tarball")
	}
	defer resp.Body.Close()

	// Create the file on disk
	tarball, err := os.Create(filepath)
	if err != nil {
		return errors.Wrap(err, "failed to create tarball file")
	}
	defer tarball.Close()

	// Fill out the file
	_, err = io.Copy(tarball, resp.Body)
	if err != nil {
		return errors.Wrap(err, "failed to save tarball file")
	}

	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}
