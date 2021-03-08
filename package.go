package epub

import (
	"encoding/xml"
)

type packageXML struct {
	XMLName  xml.Name  `xml:"http://www.idpf.org/2007/opf package"`
	Version  string    `xml:"version,attr"`
	Metadata *Metadata `xml:"metadata"`
	Manifest *Manifest `xml:"manifest"`
}

type Metadata struct {
	Title       []string     `xml:"title"`
	Language    []string     `xml:"language"`
	Identifier  []Identifier `xml:"identifier"`
	Creator     []Author     `xml:"creator"`
	Subject     []string     `xml:"subject"`
	Description []string     `xml:"description"`
	Publisher   []string     `xml:"publisher"`
	Contributor []Author     `xml:"contributor"`
	Date        []Date       `xml:"date"`
	Type        []string     `xml:"type"`
	Format      []string     `xml:"format"`
	Source      []string     `xml:"source"`
	Relation    []string     `xml:"relation"`
	Coverage    []string     `xml:"coverage"`
	Rights      []string     `xml:"rights"`
	Meta        []*Meta      `xml:"meta"`
}

type Identifier struct {
	Value  string `xml:",chardata"`
	ID     string `xml:"id,attr"`
	Scheme string `xml:"scheme,attr"`
}

type Author struct {
	FullName string `xml:",chardata"`
	FileAs   string `xml:"file-as,attr"`
	Role     string `xml:"role,attr"`
}

type Date struct {
	Stamp string `xml:",chardata"`
	Event string `xml:"event,attr"`
}

type Meta struct {
	Name    string `xml:"name,attr"`
	Content string `xml:"content,attr"`
}

type Manifest struct {
	Item []*Item `xml:"item"`
}

type Item struct {
	ID        string `xml:"id,attr"`
	Href      string `xml:"href,attr"`
	MediaType string `xml:"media-type,attr"`
}
