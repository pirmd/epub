package epub

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"testing"

	"github.com/pirmd/verify"
)

func TestWalkFiles(t *testing.T) {
	testCases, err := filepath.Glob(filepath.Join(testdataPath, "*.epub"))
	if err != nil {
		t.Fatalf("cannot read test data in %s:%v", testdataPath, err)
	}

	got := new(bytes.Buffer)
	for _, tc := range testCases {
		fmt.Fprintf(got, "Files from %s\n", tc)
		if err := WalkFiles(tc, func(r io.Reader, fi fs.FileInfo) error {
			fmt.Fprintf(got, "- %s\n", fi.Name())
			return nil
		}); err != nil {
			t.Errorf("Fail to get walk %s files: %v", tc, err)
		}
	}

	if failure := verify.MatchGolden(t.Name(), got.String()); failure != nil {
		t.Fatalf("Metadata is not as expected:\n%v", failure)
	}
}

func TestWalkPublicationResources(t *testing.T) {
	testCases, err := filepath.Glob(filepath.Join(testdataPath, "*.epub"))
	if err != nil {
		t.Fatalf("cannot read test data in %s:%v", testdataPath, err)
	}

	got := new(bytes.Buffer)
	for _, tc := range testCases {
		fmt.Fprintf(got, "Files from %s\n", tc)
		if err := WalkPublicationResources(tc, func(r io.Reader, fi fs.FileInfo) error {
			fmt.Fprintf(got, "- %s\n", fi.Name())
			return nil
		}); err != nil {
			t.Errorf("Fail to get walk %s files: %v", tc, err)
		}
	}

	if failure := verify.MatchGolden(t.Name(), got.String()); failure != nil {
		t.Fatalf("Metadata is not as expected:\n%v", failure)
	}
}

func TestWalkReadingContent(t *testing.T) {
	testCases, err := filepath.Glob(filepath.Join(testdataPath, "*.epub"))
	if err != nil {
		t.Fatalf("cannot read test data in %s:%v", testdataPath, err)
	}

	got := new(bytes.Buffer)
	for _, tc := range testCases {
		fmt.Fprintf(got, "Files from %s\n", tc)
		if err := WalkReadingContent(tc, func(r io.Reader, fi fs.FileInfo) error {
			fmt.Fprintf(got, "- %s\n", fi.Name())
			return nil
		}); err != nil {
			t.Errorf("Fail to get walk %s files: %v", tc, err)
		}
	}

	if failure := verify.MatchGolden(t.Name(), got.String()); failure != nil {
		t.Fatalf("Metadata is not as expected:\n%v", failure)
	}
}
