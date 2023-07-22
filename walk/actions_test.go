package main

import (
	"os"
	"testing"
)

func TestFilterOut(t *testing.T) {
	testCases := []struct {
		name     string
		file     string
		ext      string
		minSize  int64
		expected bool
	}{
		{"FilterOutNoExtension", "testdata/dir.log", "", 0, false},
		{"FilterOutExtensionMatch", "testdata/dir.log", ".log", 0, false},
		{"FilterOutExtensionNoMatch", "testdata/dir.log", ".sh", 0, true},
		{"FilterOutExtensionSizeMatch", "testdata/dir.log", ".log", 10, false},
		{"FilterOutExtensionSizeNoMatch", "testdata/dir.log", ".log", 20, true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(subTest *testing.T) {
			info, err := os.Stat(testCase.file)
			if err != nil {
				subTest.Fatal(err)
			}

			expected := filterOut(testCase.file, testCase.ext, testCase.minSize, info)
			if expected != testCase.expected {
				t.Fatalf("Expected %t but got %t instead", testCase.expected, expected)
			}
		})
	}
}
