package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	isCountLines := flag.Bool("l", false, "Count lines")
	flag.Parse()

	fmt.Println(count(os.Stdin, *isCountLines))
}

func count(r io.Reader, isCountLines bool) int {
	scanner := bufio.NewScanner(r)

	if !isCountLines {
		scanner.Split(bufio.ScanWords)
	}

	var wc int
	for scanner.Scan() {
		wc++
	}

	return wc
}
