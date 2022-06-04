package epub

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/pirmd/verify"
)

const (
	testdataPath = "./testdata"
)

func TestGetMetadataFromFile(t *testing.T) {
	testCases, err := filepath.Glob(filepath.Join(testdataPath, "*.epub"))
	if err != nil {
		t.Fatalf("cannot read test data in %s:%v", testdataPath, err)
	}

	out := []*MetaInformation{}
	for _, tc := range testCases {
		m, err := GetMetadataFromFile(tc)
		if err != nil {
			t.Errorf("Fail to get metadata for %s: %v", tc, err)
		}
		out = append(out, m)
	}

	got, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		t.Fatalf("Fail to marshal test output to json: %v", err)
	}

	if failure := verify.MatchGolden(t.Name(), string(got)); failure != nil {
		t.Fatalf("Metadata is not as expected:\n%v", failure)
	}
}

func TestGetPackageFromFile(t *testing.T) {
	testCases, err := filepath.Glob(filepath.Join(testdataPath, "*.epub"))
	if err != nil {
		t.Fatalf("cannot read test data in %s:%v", testdataPath, err)
	}

	out := []*PackageDocument{}
	for _, tc := range testCases {
		opf, err := GetPackageFromFile(tc)
		if err != nil {
			t.Errorf("Fail to get package for %s: %v", tc, err)
		}
		out = append(out, opf)
	}

	got, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		t.Fatalf("Fail to marshal test output to json: %v", err)
	}

	if failure := verify.MatchGolden(t.Name(), string(got)); failure != nil {
		t.Fatalf("Package Document is not as expected:\n%v", failure)
	}
}
