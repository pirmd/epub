package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pirmd/epub"
)

func stringify(v interface{}) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintf("%+v", v)
	}

	return string(b)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("USAGE: %s <epub>\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	metadata, err := epub.GetMetadataFromFile(os.Args[1])
	if err != nil {
		fmt.Printf("Error while parsing epub %s: %v\n", os.Args[1], err)
		os.Exit(1)
	}

	fmt.Printf("%s\n", stringify(metadata))
	os.Exit(0)
}
