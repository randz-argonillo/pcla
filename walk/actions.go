package main

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

func filterOut(path string, ext string, size int64, info fs.FileInfo) bool {
	if info.IsDir() {
		return true
	}

	fileSize := info.Size()
	if fileSize < size {
		return true
	}

	if ext != "" && filepath.Ext(path) != ext {
		return true
	}

	return false
}

func listFile(path string, output io.Writer) error {
	_, err := fmt.Fprintln(output, path)
	return err
}

func deleteFile(path string, logger *log.Logger) error {
	if err := os.Remove(path); err != nil {
		return err
	}

	logger.Printf("deleted file %s", path)
	return nil
}
