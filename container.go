package epub

import (
	"encoding/xml"
	"fmt"
	"io"
)

const (
	containerPath = "META-INF/container.xml"
)

//TODO: add support for multiple Rootfiles

type container struct {
	XMLName   xml.Name `xml:"urn:oasis:names:tc:opendocument:xmlns:container container"`
	Rootfiles rootfile `xml:"rootfiles>rootfile"`
}

type rootfile struct {
	XMLName  xml.Name `xml:"rootfile"`
	FullPath string   `xml:"full-path,attr"`
}

func newContainer(r io.Reader) (*container, error) {
	c := &container{}
	if err := decodeXML(r, &c); err != nil {
		return nil, fmt.Errorf("failed to decode container XML: %w", err)
	}

	return c, nil
}
