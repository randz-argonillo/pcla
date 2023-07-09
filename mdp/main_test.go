package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

const (
	inputFile  = "./testdata/test1.md"
	goldenFile = "./testdata/test1.md.html"
)

func TestConvertToHtml(t *testing.T) {
	markdown, err := os.ReadFile(inputFile)
	if err != nil {
		t.Fatal(err)
	}

	result, err := convertToHtml(markdown, "")
	if err != nil {
		t.Fatal(err)
	}

	expected, err := os.ReadFile(goldenFile)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(expected, result) {
		t.Logf("golden: \n%s\n", expected)
		t.Logf("result: \n%s\n", result)
		t.Error("Result is not the same with the golden file")
	}
}

func TestRun(t *testing.T) {
	var outputWriter bytes.Buffer

	if err := run(inputFile, "", &outputWriter, true); err != nil {
		t.Fatal(err)
	}

	golden, err := os.ReadFile(goldenFile)
	if err != nil {
		t.Fatal(err)
	}

	resultFile := strings.TrimSpace(outputWriter.String())

	result, err := os.ReadFile(resultFile)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(golden, result) {
		t.Logf("golden: \n%s\n", golden)
		t.Logf("result: \n%s\n", result)
		t.Error("Result is not the same with golden file")
	}

	os.Remove(resultFile)
}
