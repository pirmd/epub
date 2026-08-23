package epub

// Information gathers meta information about an epub as a simpler version
// of Metadata to offer a more direct access to an Epub's metadata for simple
// use cases.
type Information struct {
	// Identifier contains an identifier associated with the given
	// Rendition, such as a UUID, DOI or ISBN.
	Identifier []Identifier
	// Title represents the EPUB titles.
	Title []string
	// SubTitle represents the EPUB sub-titles.
	SubTitle []string `json:",omitempty"`

	// Language element specifies the language of the content of the
	// given Rendition.
	Language []string

	// Contributor represents the name of a person, organization, etc.
	// that played a secondary role in the creation of the content of an
	// EPUB Publication.
	Contributor []Author `json:",omitempty"`
	// Coverage gives the extent or scope of the publication's content.
	Coverage []string `json:",omitempty"`
	// Creator represents the name of a person, organization, etc.
	// responsible for the creation of the content of the Rendition.
	Creator []Author
	// Date lists events associated to the EPUB like publication, creation...
	Date []Date `json:",omitempty"`
	// Description provides a description of the publication's content.
	Description []string `json:",omitempty"`
	// Format identifies the media type or dimensions of the resource.
	Format []string `json:",omitempty"`
	// Publisher identifies the publication's publisher.
	Publisher []string `json:",omitempty"`
	// Relation is an identifier of an auxiliary resource and its
	// relationship to the publication.
	Relation []string `json:",omitempty"`
	// Rights provides a statement about rights, or a reference to one.
	Rights []string `json:",omitempty"`
	// Sources provides information regarding a prior resource from which
	// the publication was derived.
	Source []string `json:",omitempty"`
	// Subject identifies the subject of the EPUB Publication.
	Subject []string `json:",omitempty"`
	// Type is used to indicate that the given EPUB Publication is of a
	// specialized type.
	Type []string `json:",omitempty"`

	// Meta element provides a generic means of including package
	// metadata.
	Meta []GenericMetadata `json:",omitempty"`

	// Series is the series to which this book belongs to.
	Series string `json:",omitempty"`
	// SeriesIndex is the position in the series to which the book belongs to.
	SeriesIndex string `json:",omitempty"`
}

// Identifier represents an identifier.
type Identifier struct {
	Scheme string
	Value  string
}

// Author represents an author.
type Author struct {
	FullName string
	FileAs   string
	Role     string
}

// Date represents an event.
type Date struct {
	Stamp string
	Event string
}

// GenericMetadata represents a generic metadata.
type GenericMetadata struct {
	Name    string
	Content string
}

// GetMetadataFromFile reads metadata from an epub file.
func GetMetadataFromFile(path string) (*Information, error) {
	e, err := Open(path)
	if err != nil {
		return nil, err
	}
	defer e.Close()

	return e.Information()
}

// Title returns the first title if available, otherwise an empty string.
func (i *Information) Title() string {
	if len(i.Title) > 0 {
		return i.Title[0]
	}
	return ""
}

// Author returns the first author's full name if available, otherwise an empty string.
func (i *Information) Author() string {
	if len(i.Creator) > 0 {
		return i.Creator[0].FullName
	}
	return ""
}

// Language returns the first language if available, otherwise an empty string.
func (i *Information) Language() string {
	if len(i.Language) > 0 {
		return i.Language[0]
	}
	return ""
}

// Publisher returns the first publisher if available, otherwise an empty string.
func (i *Information) Publisher() string {
	if len(i.Publisher) > 0 {
		return i.Publisher[0]
	}
	return ""
}

// Description returns the first description if available, otherwise an empty string.
func (i *Information) Description() string {
	if len(i.Description) > 0 {
		return i.Description[0]
	}
	return ""
}

// HasLanguage checks if the specified language is in the Language list.
func (i *Information) HasLanguage(lang string) bool {
	for _, l := range i.Language {
		if l == lang {
			return true
		}
	}
	return false
}

// SeriesInfo returns a formatted string with series and index.
// Returns empty string if no series information is available.
func (i *Information) SeriesInfo() string {
	if i.Series == "" {
		return ""
	}
	if i.SeriesIndex == "" {
		return i.Series
	}
	return i.Series + " #" + i.SeriesIndex
}

func getMeta(mdata *Metadata) *Information {
	m := &Information{
		Language:    elt2str(mdata.Language),
		Subject:     elt2str(mdata.Subject),
		Description: elt2str(mdata.Description),
		Publisher:   elt2str(mdata.Publisher),
		Type:        elt2str(mdata.Type),
		Format:      elt2str(mdata.Format),
		Source:      elt2str(mdata.Source),
		Relation:    elt2str(mdata.Relation),
		Coverage:    elt2str(mdata.Coverage),
		Rights:      elt2str(mdata.Rights),
	}

	m.Title, m.SubTitle = getTitles(mdata.Title, mdata.Meta)

	for _, id := range mdata.Identifier {
		m.Identifier = append(m.Identifier, Identifier{
			Value:  id.Value,
			Scheme: id.Scheme,
		})
	}

	for _, auth := range mdata.Creator {
		m.Creator = append(m.Creator, getAuth(auth, mdata.Meta))
	}

	for _, auth := range mdata.Contributor {
		m.Contributor = append(m.Contributor, getAuth(auth, mdata.Meta))
	}

	for _, date := range mdata.Date {
		m.Date = append(m.Date, Date{
			Stamp: date.Value,
			Event: date.Event,
		})
	}

	m.Series, m.SeriesIndex = getSeries(mdata.Meta)

	for _, meta := range mdata.Meta {
		if meta.Name != "" && meta.Content != "" {
			m.Meta = append(m.Meta, GenericMetadata{
				Name:    meta.Name,
				Content: meta.Content,
			})
		}
	}

	return m
}

// elt2str converts a slice of Element to a slice of strings.
func elt2str(elt []Element) []string {
	s := make([]string, 0, len(elt))
	for _, e := range elt {
		s = append(s, e.Value)
	}
	return s
}

// getAuth extracts author information from AuthorElt and enhances it with metadata.
func getAuth(auth AuthorElt, meta []MetaLegacy) Author {
	a := Author{
		FullName: auth.Value,
		Role:     auth.Role,
		FileAs:   auth.FileAs,
	}

	for _, m := range meta {
		if m.Refines != "#"+auth.ID {
			continue
		}

		switch m.Property {
		case "role":
			a.Role = m.Value
		case "file-as":
			a.FileAs = m.Value
		}
	}

	return a
}

// getTitles extracts titles and subtitles from elements and metadata.
// It prioritizes elements with title-type metadata to identify subtitles.
func getTitles(elt []Element, meta []MetaLegacy) (title []string, subtitle []string) {
	// Build a map of element ID to its metadata for O(1) lookup
	metaMap := make(map[string][]MetaLegacy)
	for _, m := range meta {
		if m.Refines != "" {
			metaMap[m.Refines] = append(metaMap[m.Refines], m)
		}
	}

	for _, e := range elt {
		// Check if this element has title-type metadata
		if metas, exists := metaMap["#"+e.ID]; exists {
			for _, m := range metas {
				if m.Property == "title-type" {
					switch m.Value {
					case "subtitle":
						subtitle = append(subtitle, e.Value)
						goto nextElt
					}
				}
			}
		}
		title = append(title, e.Value)
	nextElt:
	}

	return title, subtitle
}

// getSeries extracts series information from Meta.
// It supports both Calibre-like EPUB 2 series coding and EPUB3 collection metadata.
// If both are available, EPUB3 is preferred.
func getSeries(meta []MetaLegacy) (series string, seriesIndex string) {
	// First, look for EPUB3 collection metadata (priority)
	for _, m := range meta {
		if m.Property == "belongs-to-collection" {
			series = m.Value

			// Look for group-position that refines this collection
			for _, mm := range meta {
				if mm.Refines == "#"+m.ID && mm.Property == "group-position" {
					seriesIndex = mm.Value
					break
				}
			}
			// If we found series info via EPUB3, return immediately
			if series != "" {
				return series, seriesIndex
			}
		}
	}

	// Fall back to Calibre-style metadata
	for _, m := range meta {
		switch m.Name {
		case "calibre:series":
			series = m.Content
		case "calibre:series_index":
			seriesIndex = m.Content
		}
	}

	return series, seriesIndex
}
