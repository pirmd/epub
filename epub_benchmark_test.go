package epub

import (
	"path/filepath"
	"testing"
)

// BenchmarkGetMetadataFromFile benchmarks metadata extraction from an EPUB file.
func BenchmarkGetMetadataFromFile(b *testing.B) {
	epubFile := filepath.Join(testdataPath, "pg11-images.epub")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = GetMetadataFromFile(epubFile)
	}
}

// BenchmarkGetPackageFromFile benchmarks package document extraction.
func BenchmarkGetPackageFromFile(b *testing.B) {
	epubFile := filepath.Join(testdataPath, "pg11-images.epub")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = GetPackageFromFile(epubFile)
	}
}

// BenchmarkWalkReadingContent benchmarks walking reading content.
func BenchmarkWalkReadingContent(b *testing.B) {
	epubFile := filepath.Join(testdataPath, "pg11-images.epub")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = WalkReadingContent(epubFile, func(r interface{}, info interface{}) error {
			return nil
		})
	}
}

// BenchmarkElt2str benchmarks the elt2str conversion function.
func BenchmarkElt2str(b *testing.B) {
	elts := make([]Element, 100)
	for i := range elts {
		elts[i] = Element{Value: "Test Value"}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = elt2str(elts)
	}
}

// BenchmarkGetSeries benchmarks the getSeries function.
func BenchmarkGetSeries(b *testing.B) {
	meta := []MetaLegacy{
		{Meta: &Meta{Property: "belongs-to-collection", Value: "Test Series", ID: "series1"}},
		{Meta: &Meta{Property: "group-position", Value: "1", Refines: "#series1"}},
		{MetaLegacy: &MetaLegacy{Name: "calibre:series", Content: "Calibre Series"}},
		{MetaLegacy: &MetaLegacy{Name: "calibre:series_index", Content: "2"}},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = getSeries(meta)
	}
}

// BenchmarkGetTitles benchmarks the getTitles function.
func BenchmarkGetTitles(b *testing.B) {
	elts := make([]Element, 50)
	for i := range elts {
		elts[i] = Element{Value: "Title " + string(rune(i)), ID: "title" + string(rune(i))}
	}
	meta := []MetaLegacy{
		{Meta: &Meta{Property: "title-type", Value: "subtitle", Refines: "#title0"}},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = getTitles(elts, meta)
	}
}
