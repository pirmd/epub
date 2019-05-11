package epub

import (
	"encoding/xml"
)

const (
	containerPath = "META-INF/container.xml"
)

type containerXML struct {
	XmlName   xml.Name `xml:"urn:oasis:names:tc:opendocument:xmlns:container container"`
	Rootfiles rootfile `xml:"rootfiles>rootfile"`
}

type rootfile struct {
	XmlName  xml.Name `xml:"rootfile"`
	FullPath string   `xml:"full-path,attr"`
}
