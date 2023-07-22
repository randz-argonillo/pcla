package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestRun(t *testing.T) {
	testcases := []struct {
		name     string
		root     string
		cfg      config
		expected string
	}{
		{
			name:     "NoFilter",
			root:     "testdata",
			cfg:      config{ext: "", size: 0, list: true},
			expected: "testdata/dir.log\ntestdata/dir2/script.sh\n"},
		{
			name:     "FilterExtensionMatch",
			root:     "testdata",
			cfg:      config{ext: ".log", size: 0, list: true},
			expected: "testdata/dir.log\n"},
		{
			name:     "FilterExtensionSizeMatch",
			root:     "testdata",
			cfg:      config{ext: ".log", size: 10, list: true},
			expected: "testdata/dir.log\n"},
		{
			name:     "FilterExtensionSizeNoMatch",
			root:     "testdata",
			cfg:      config{ext: ".log", size: 20, list: true},
			expected: ""},
		{
			name:     "FilterExtensionNoMatch",
			root:     "testdata",
			cfg:      config{ext: ".gz", size: 0, list: true},
			expected: ""},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(st *testing.T) {
			output := bytes.Buffer{}
			if err := run(tc.root, &output, tc.cfg); err != nil {
				st.Fatal(err)
			}

			res := output.String()

			if res != tc.expected {
				st.Fatalf("Expected %q, got %q instead", tc.expected, res)
			}

		})
	}
}

func TestRunDeleteFiles(t *testing.T) {
	testCases := []struct {
		name        string
		cfg         config
		extNoDelete string
		nDelete     int
		nNoDelete   int
		expected    string
	}{
		{
			name:        "DeleteExtensionNoMatch",
			cfg:         config{ext: ".log", delete: true},
			extNoDelete: ".gz",
			nNoDelete:   10,
			nDelete:     0,
			expected:    ""},
		{
			name:        "DeleteExtensionMatch",
			cfg:         config{ext: ".log", delete: true},
			extNoDelete: "",
			nNoDelete:   0,
			nDelete:     10,
			expected:    ""},
		{
			name:        "DeleteExtensionMixed",
			cfg:         config{ext: ".log", delete: true},
			extNoDelete: ".gz",
			nNoDelete:   5,
			nDelete:     5,
			expected:    ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(st *testing.T) {
			testdir, cleanup := createTestFiles(st, map[string]int{
				tc.cfg.ext:     tc.nDelete,
				tc.extNoDelete: tc.nNoDelete,
			})

			defer cleanup()

			var output bytes.Buffer
			var delLog bytes.Buffer

			tc.cfg.delLog = &delLog

			if err := run(testdir, &output, tc.cfg); err != nil {
				st.Fatal(err)
			}

			// result := output.String()
			// if tc.nDelete > 0 && !strings.Contains(result, "deleted file") {
			// 	st.Fatalf("Expected to delete files but not")
			// }

			// if tc.nDelete == 0 && strings.Contains(result, "deleted file") {
			// 	st.Fatal("Expected not to delete any file but a file was deleted")
			// }

			// test deleted line count
			expectedDelLines := tc.nDelete + 1
			actualDelLines := bytes.Split(delLog.Bytes(), []byte("\n"))
			if len(actualDelLines) != expectedDelLines {
				t.Errorf(
					"Expected %d log lines, but %d instead",
					expectedDelLines,
					len(actualDelLines))
			}

			// Test files remaining
			filesLeft, err := os.ReadDir(testdir)
			if err != nil {
				st.Fatal(err)
			}

			if len(filesLeft) != tc.nNoDelete {
				st.Fatalf("Expected %d files left, but got %d", tc.nNoDelete, len(filesLeft))
			}

		})
	}
}

func createTestFiles(t *testing.T, files map[string]int) (dir string, cleanup func()) {
	t.Helper()

	tmpDir, err := os.MkdirTemp("", "walktest")
	if err != nil {
		t.Fatal(err)
	}

	for ext, fileCount := range files {
		for i := 1; i <= fileCount; i++ {
			fname := fmt.Sprintf("file%d.%s", i, ext)
			fullname := filepath.Join(tmpDir, fname)

			if err := os.WriteFile(fullname, []byte("dummy"), 0644); err != nil {
				t.Fatal(err)
			}
		}
	}

	return tmpDir, func() { os.RemoveAll(tmpDir) }
}
