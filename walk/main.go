package main

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

type config struct {
	ext    string    // The file extension to include
	size   int64     // Minimum file size to include
	list   bool      // To list the results or not
	delete bool      // Delete found files
	delLog io.Writer // Logger for delete operation
}

func main() {
	root := flag.String("root", ".", "The root directory to start")
	list := flag.Bool("list", false, "List files only")
	ext := flag.String("ext", "", "File extension for filtering results")
	size := flag.Int64("size", 0, "Minimum file size")
	del := flag.Bool("del", false, "Delete file")
	delLogFile := flag.String("delLogFile", "", "File to log delete")

	flag.Parse()

	f, cleanup, err := prepDelLog(*delLogFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if cleanup != nil {
		defer cleanup()
	}

	c := config{
		ext:    *ext,
		size:   *size,
		list:   *list,
		delete: *del,
		delLog: f,
	}

	if err := run(*root, os.Stdout, c); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(dir string, output io.Writer, conf config) error {
	delLogger := log.New(conf.delLog, "DELETED FILE:", log.LstdFlags)

	return filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if filterOut(path, conf.ext, conf.size, info) {
			return nil
		}

		if conf.list {
			return listFile(path, output)
		}

		if conf.delete {
			return deleteFile(path, delLogger)
		}

		return listFile(path, output)
	})
}

func prepDelLog(logFName string) (writer io.Writer, cleanup func(), err error) {
	if logFName == "" {
		return os.Stdout, nil, nil
	}

	file, err := os.OpenFile(logFName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, nil, err
	}

	cleanup = func() { file.Close() }

	return file, cleanup, nil
}
