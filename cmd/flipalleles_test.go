package cmd

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
)

func Test_flipalleles_indexOf(t *testing.T) {
	testCases := []struct {
		name       string
		collection []string
		el         string
		expected   int
	}{
		{
			name:       "Element present",
			collection: []string{"a", "b", "c", "d"},
			el:         "c",
			expected:   2,
		},
		{
			name:       "Element not present",
			collection: []string{"a", "b", "c", "d"},
			el:         "x",
			expected:   -1,
		},
		{
			name:       "Empty collection",
			collection: []string{},
			el:         "a",
			expected:   -1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			index, err := indexOf(tc.collection, tc.el)
			if index != tc.expected {
				t.Errorf("Expected index %d, got %d", tc.expected, index)
			}
			if index == -1 && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if index != -1 && err != nil {
				t.Errorf("Expected no error, got: %v", err)
			}
		})
	}
}

func Test_flipalleles_flipBeta(t *testing.T) {
	testCases := []struct {
		name     string
		beta     float64
		expected float64
	}{
		{
			name:     "Positive beta",
			beta:     2.5,
			expected: -2.5,
		},
		{
			name:     "Negative beta",
			beta:     -1.5,
			expected: 1.5,
		},
		{
			name:     "Zero beta",
			beta:     0,
			expected: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			flippedBeta := flipBeta(tc.beta)
			if flippedBeta != tc.expected {
				t.Errorf("Expected flipped beta %f, got %f", tc.expected, flippedBeta)
			}
		})
	}
}

func Test_flipalleles_flipOR(t *testing.T) {
	testCases := []struct {
		name     string
		or       float64
		expected float64
	}{
		{
			name:     "Positive OR",
			or:       2.0,
			expected: 0.5,
		},
		{
			name:     "Negative OR",
			or:       -4.0,
			expected: -0.25,
		},
		{
			name:     "One OR",
			or:       1.0,
			expected: 1.0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			flippedOR := flipOR(tc.or)
			if flippedOR != tc.expected {
				t.Errorf("Expected flipped OR %f, got %f", tc.expected, flippedOR)
			}
		})
	}
}

func Test_flipalleles_parseSumstatsFileToMap(t *testing.T) {
	testData := `rsid	effect_allele	other_allele
rs123	A	C
rs456	T	G
rs123	G	T
`

	tmpfile, err := ioutil.TempFile("", "test_sumstats")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // Clean up

	if _, err := tmpfile.Write([]byte(testData)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	result, err := parseSumstatsFileToMap(tmpfile.Name(), "rsid", "effect_allele", "other_allele")
	if err != nil {
		t.Fatalf("Error calling parseSumstatsFileToMap: %v", err)
	}

	expected := map[string][]Alleles{
		"rs123": {
			{Effect: "A", Other: "C"},
			{Effect: "G", Other: "T"},
		},
		"rs456": {
			{Effect: "T", Other: "G"},
		},
	}

	for snp, alleles := range expected {
		if _, ok := result[snp]; !ok {
			t.Errorf("SNP %s not found in result", snp)
		} else {
			for i, allele := range alleles {
				if result[snp][i].Effect != allele.Effect || result[snp][i].Other != allele.Other {
					t.Errorf("SNP %s: expected alleles (%s, %s), got (%s, %s)", snp, allele.Effect, allele.Other, result[snp][i].Effect, result[snp][i].Other)
				}
			}
		}
	}
}

func Test_flipalleles_processAndWriteFlippedStats(t *testing.T) {
	testInputData := `rsid	effect_allele	other_allele	effect
rs123	A	C	1.5
rs456	T	G	0.8
`

	tmpInputFile, err := ioutil.TempFile("", "test_input")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpInputFile.Name())

	if _, err := tmpInputFile.Write([]byte(testInputData)); err != nil {
		t.Fatal(err)
	}
	if err := tmpInputFile.Close(); err != nil {
		t.Fatal(err)
	}

	tmpOutputFile, err := ioutil.TempFile("", "test_output")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpOutputFile.Name())

	referenceSNPMapping := map[string][]Alleles{
		"rs123": {
			{Effect: "C", Other: "A"},
		},
		"rs456": {
			{Effect: "G", Other: "T"},
		},
	}

	logger = log.New(ioutil.Discard, "", 0)
	err = processAndWriteFlippedStats(tmpInputFile.Name(), tmpOutputFile.Name(), "rsid", "effect_allele", "other_allele", "effect", "BETA", referenceSNPMapping)
	if err != nil {
		t.Fatalf("Error calling processAndWriteFlippedStats: %v", err)
	}

	expectedOutputData := `rsid	effect_allele	other_allele	effect
rs123	A	C	-1.5000000
rs456	T	G	-0.8000000
	`

	outputData, err := ioutil.ReadFile(tmpOutputFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	if strings.TrimSpace(string(outputData)) != strings.TrimSpace(expectedOutputData) {
		t.Errorf("Output data does not match expected data:\nExpected:\n%s\nActual:\n%s", expectedOutputData, string(outputData))
	}
}
