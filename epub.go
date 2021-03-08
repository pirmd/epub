// epub package provides a way to retrieve stored metadata from epub files.

package epub

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"

	"golang.org/x/net/html/charset"
)

//GetMetadata reads metadata from the given epub opened as a readatSeeker
func GetMetadata(r readatSeeker) (*Metadata, error) {
	opf, err := GetOPFData(r)
	if err != nil {
		return nil, fmt.Errorf("not a valid Epub: %v", err)
	}

	return opf.Metadata, nil
}

//GetMetadataFromFile reads metadata from the given epub file
func GetMetadataFromFile(path string) (*Metadata, error) {
	r, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	return GetMetadata(r)
}

func getContainerData(r readatSeeker) (*containerXML, error) {
	f, err := openFromZip(r, containerPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	c := &containerXML{}
	err = decodeXML(f, &c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// GetOPFData reads the whole OPF from the given epub file
func GetOPFData(r readatSeeker) (*packageXML, error) {
	c, err := getContainerData(r)
	if err != nil {
		return nil, err
	}

	f, err := openFromZip(r, c.Rootfiles.FullPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	opf := &packageXML{}
	err = decodeXML(f, &opf)
	if err != nil {
		return nil, err
	}

	return opf, nil
}

func decodeXML(f io.Reader, v interface{}) error {
	decoder := xml.NewDecoder(f)
	decoder.Entity = xml.HTMLEntity
	decoder.CharsetReader = charset.NewReaderLabel
	return decoder.Decode(v)
}
