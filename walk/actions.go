package main

import (
	"compress/gzip"
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

func archiveFile(archiveDir string, root string, file string) error {
	if err := ensureDir(archiveDir); err != nil {
		return err
	}

	relativeDir, err := filepath.Rel(root, filepath.Dir(file))
	if err != nil {
		return err
	}

	// create target dir
	destFile := fmt.Sprintf("%s.gz", filepath.Base(file))
	targetFilePath := filepath.Join(archiveDir, relativeDir, destFile)
	if err := ensureDir(filepath.Dir(targetFilePath)); err != nil {
		return err
	}

	outputFile, err := os.OpenFile(targetFilePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	inputFile, err := os.Open(file)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	writer := gzip.NewWriter(inputFile)
	writer.Name = filepath.Base(file)

	if _, err := io.Copy(writer, inputFile); err != nil {
		return err
	}

	if err := writer.Close(); err != nil {
		return err
	}

	return outputFile.Close()

}

func ensureDir(dirPath string) error {
	fi, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		return os.MkdirAll(dirPath, 0755)
	}

	if !fi.IsDir() {
		return fmt.Errorf("'%s' is not a directory", dirPath)
	}

	return nil
}
