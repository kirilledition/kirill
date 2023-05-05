package cmd

import (
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strings"
	"testing"
)

func Test_validatePDBId(t *testing.T) {
	testCases := []struct {
		input    []string
		expected []string
		err      bool
	}{
		{
			input:    []string{"1abc", "2DEF", "3GhI"},
			expected: []string{"1abc", "2def", "3ghi"},
			err:      false,
		},
		{
			input:    []string{"1abc", "2DEF", "3Gh"},
			expected: nil,
			err:      true,
		},
		{
			input:    []string{"1abc", "2DEF", "3GhJk"},
			expected: nil,
			err:      true,
		},
		{
			input:    []string{"1abc", "2DEF", ""},
			expected: nil,
			err:      true,
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("Test case %d", i+1), func(t *testing.T) {
			result, err := validatePDBId(testCase.input)

			if (err != nil) != testCase.err {
				t.Errorf("Expected error: %v, got: %v", testCase.err, err)
			}

			if !equalStringSlices(result, testCase.expected) {
				t.Errorf("Expected result: %v, got: %v", testCase.expected, result)
			}
		})
	}
}

func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func Test_readPDBIdList(t *testing.T) {

	logger = log.New(ioutil.Discard, "", 0)

	tmpFile, err := ioutil.TempFile("", "pdbids_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	content := []byte("1abc\n2def\n3ghi\n")
	if _, err := tmpFile.Write(content); err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temporary file: %v", err)
	}

	testCases := []struct {
		input    []string
		expected []string
		err      bool
	}{
		{
			input:    []string{"1abc", "2DEF", "3GhI"},
			expected: []string{"1abc", "2def", "3ghi"},
			err:      false,
		},
		{
			input:    []string{tmpFile.Name()},
			expected: []string{"1abc", "2def", "3ghi"},
			err:      false,
		},
		{
			input:    []string{"nonexistent.txt"},
			expected: nil,
			err:      true,
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("Test case %d", i+1), func(t *testing.T) {
			result, err := readPDBIdList(testCase.input)

			if (err != nil) != testCase.err {
				t.Errorf("Expected error: %v, got: %v", testCase.err, err)
			}

			if !equalStringSlices(result, testCase.expected) {
				t.Errorf("Expected result: %v, got: %v", testCase.expected, result)
			}
		})
	}
}

func Test_PDBClient_fetch(t *testing.T) {
	testPDBID := "1abc"
	testPDBData := "dummy pdb data"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/octet-stream")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gz.Write([]byte(testPDBData))
	}))
	defer ts.Close()

	workspace, ok := os.LookupEnv("GITHUB_WORKSPACE")
	if !ok {
		workspace = ""
	}
	outputPath, err := ioutil.TempDir(workspace, "pdbtest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(outputPath)

	client := &PDBClient{
		scheme: ts.URL,
		client: &http.Client{},
	}

	logger = log.New(ioutil.Discard, "", 0)

	err = client.fetch(testPDBID, outputPath)
	if err != nil {
		t.Fatalf("fetch() returned error: %v", err)
	}

	filename := path.Join(outputPath, testPDBID+".pdb")
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatalf("could not read fetched file: %v", err)
	}

	if string(content) != testPDBData {
		t.Errorf("expected content: %s, got: %s", testPDBData, string(content))
	}
}

func Test_fetchPDB(t *testing.T) {
	testPDBID := "1abc"
	testPDBData := "dummy pdb data"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/octet-stream")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gz.Write([]byte(testPDBData))
	}))
	defer ts.Close()

	workspace, ok := os.LookupEnv("GITHUB_WORKSPACE")
	if !ok {
		workspace = ""
	}
	outputPath, err := ioutil.TempDir(workspace, "pdbtest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(outputPath)

	input := []string{testPDBID}

	client := &PDBClient{
		scheme: ts.URL,
		client: &http.Client{},
	}

	logger = log.New(ioutil.Discard, "", 0)

	fetchPDB(input, outputPath, client)

	filename := path.Join(outputPath, strings.ToUpper(testPDBID)+".pdb")
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatalf("could not read fetched file: %v", err)
	}

	if string(content) != testPDBData {
		t.Errorf("expected content: %s, got: %s", testPDBData, string(content))
	}
}

func Test_fetchPDB_API(t *testing.T) {
	testPDBID := "3NIR"

	workspace, ok := os.LookupEnv("GITHUB_WORKSPACE")
	if !ok {
		workspace = ""
	}
	outputPath, err := ioutil.TempDir(workspace, "pdbtest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(outputPath)

	input := []string{testPDBID}

	client := &PDBClient{
		scheme: "https",
		host:   "files.rcsb.org",
		path:   "download",
		client: &http.Client{},
	}

	logger = log.New(ioutil.Discard, "", 0)

	fetchPDB(input, outputPath, client)

	filename := path.Join(outputPath, strings.ToUpper(testPDBID)+".pdb")
	_, err = os.Stat(filename)
	if err != nil {
		t.Fatalf("could not find fetched file: %v", err)
	}
}
