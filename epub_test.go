package epub

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"
)

const (
	testdataPath = "./testdata"
)

func TestGetMetadata(t *testing.T) {
	testCases, err := filepath.Glob(filepath.Join(testdataPath, "*.epub"))
	if err != nil {
		t.Fatalf("cannot read test data in %s:%v", testdataPath, err)
	}

	out := []*Metadata{}
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

	if err := matchGolden("TestGetMetadata", string(got)); err != nil {
		t.Error()
	}
}

// MatchGolden compares a test result to the content of a 'golden' file
// If 'update' command flag is used, update the 'golden' file
func matchGolden(name string, got string) error {
	goldenPath := filepath.Join(testdataPath, name+".golden")

	want, err := readGolden(goldenPath)
	if err != nil {
		return fmt.Errorf("cannot read golden file %s: %s.\nTest output is:\n%s", goldenPath, err, got)
	}

	if len(want) == 0 {
		return fmt.Errorf("no existing or empty golden file.\nTest output is:\n%s", got)
	}

	if !reflect.DeepEqual(got, string(want)) {
		return errors.New("Output different than expected")
	}
	return nil
}

func readGolden(path string) ([]byte, error) {
	want, err := ioutil.ReadFile(path)
	if err != nil {
		return []byte{}, fmt.Errorf("cannot read golden file %s: %v", path, err)
	}
	return want, nil
}
