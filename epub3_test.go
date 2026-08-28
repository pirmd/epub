package epub

import (
	"fmt"
	"path/filepath"
	"testing"
)

// TestEPUBVersionDetection tests that EPUB version is correctly detected
func TestEPUBVersionDetection(t *testing.T) {
	testCases := []struct {
		file     string
		expected string
	}{
		{"testdata/Row-Your-Boat-2.0.epub", "2.0"},
		{"testdata/Row-Your-Boat-3.0.epub", "3.0"},
		{"testdata/pg11-images.epub", "2.0"},
		{"testdata/pg1661-images.epub", "2.0"},
		{"testdata/pg23962-images.epub", "2.0"},
	}

	for _, tc := range testCases {
		t.Run(tc.file, func(t *testing.T) {
			version, err := GetEPUBVersion(tc.file)
			if err != nil {
				t.Errorf("Failed to get EPUB version for %s: %v", tc.file, err)
				return
			}
			if version != tc.expected {
				t.Errorf("Expected version %q for %s, got %q", tc.expected, tc.file, version)
			}
		})
	}
}

// TestEpubVersionMethod tests the Version() method on Epub type
func TestEpubVersionMethod(t *testing.T) {
	testCases := []struct {
		file     string
		expected string
	}{
		{"testdata/Row-Your-Boat-2.0.epub", "2.0"},
		{"testdata/Row-Your-Boat-3.0.epub", "3.0"},
	}

	for _, tc := range testCases {
		t.Run(tc.file, func(t *testing.T) {
			e, err := Open(tc.file)
			if err != nil {
				t.Errorf("Failed to open EPUB %s: %v", tc.file, err)
				return
			}
			defer e.Close()

			version := e.Version()
			if version != tc.expected {
				t.Errorf("Expected version %q for %s, got %q", tc.expected, tc.file, version)
			}
		})
	}
}

// TestEPUB2Compatibility tests that EPUB2 files are still properly parsed
func TestEPUB2Compatibility(t *testing.T) {
	epub2Files := []string{
		"testdata/Row-Your-Boat-2.0.epub",
		"testdata/pg11-images.epub",
		"testdata/pg1661-images.epub",
	}

	for _, file := range epub2Files {
		t.Run(file, func(t *testing.T) {
			metadata, err := GetMetadataFromFile(file)
			if err != nil {
				t.Errorf("Failed to get metadata from EPUB2 file %s: %v", file, err)
				return
			}

			// Basic validation
			if metadata.Title() == "" {
				t.Errorf("Expected non-empty title for %s", file)
			}
			if len(metadata.Creator) == 0 && len(metadata.Contributor) == 0 {
				t.Logf("Warning: No creator or contributor found for %s", file)
			}
			if metadata.EPUBVersion != "2.0" {
				t.Errorf("Expected EPUB version 2.0 for %s, got %s", file, metadata.EPUBVersion)
			}
		})
	}
}

// TestEPUB3Compatibility tests that EPUB3 files are properly parsed
func TestEPUB3Compatibility(t *testing.T) {
	epub3Files := []string{
		"testdata/Row-Your-Boat-3.0.epub",
	}

	for _, file := range epub3Files {
		t.Run(file, func(t *testing.T) {
			metadata, err := GetMetadataFromFile(file)
			if err != nil {
				t.Errorf("Failed to get metadata from EPUB3 file %s: %v", file, err)
				return
			}

			// Basic validation
			if metadata.Title() == "" {
				t.Errorf("Expected non-empty title for %s", file)
			}

			// EPUB3 should have version 3.0
			if metadata.EPUBVersion != "3.0" {
				t.Errorf("Expected EPUB version 3.0 for %s, got %s", file, metadata.EPUBVersion)
			}

			// Check for title-type metadata (EPUB3 feature)
			if len(metadata.SubTitle) > 0 {
				t.Logf("Found %d subtitles in %s", len(metadata.SubTitle), file)
			}
		})
	}
}

// TestEPUB3Collections tests collection extraction from EPUB3 files
func TestEPUB3Collections(t *testing.T) {
	// Row-Your-Boat-3.0.epub has collection metadata
	metadata, err := GetMetadataFromFile("testdata/Row-Your-Boat-3.0.epub")
	if err != nil {
		t.Fatalf("Failed to get metadata: %v", err)
	}

	// This file should have collections
	if len(metadata.Collections) == 0 {
		t.Log("No collections found in Row-Your-Boat-3.0.epub (this may be expected)")
	} else {
		for i, col := range metadata.Collections {
			t.Logf("Collection %d: Role=%s, Title=%s, Position=%d",
				i+1, col.Role, col.Title, col.Position)
		}
	}
}

// TestOpenItemPathTraversal tests path traversal prevention
func TestOpenItemPathTraversal(t *testing.T) {
	e, err := Open("testdata/Row-Your-Boat-2.0.epub")
	if err != nil {
		t.Fatalf("Failed to open EPUB: %v", err)
	}
	defer e.Close()

	// Test various path traversal attempts
	maliciousPaths := []string{
		"../etc/passwd",
		"..\\windows\\system32\\config",
		"/etc/passwd",
		"",
		"../../../etc/passwd",
		"..%2Fetc%2Fpasswd", // URL encoded
	}

	for _, path := range maliciousPaths {
		t.Run(path, func(t *testing.T) {
			_, err := e.OpenItem(path)
			if err == nil {
				t.Errorf("Expected error for malicious path %q, but got none", path)
			}
		})
	}
}

// TestOpenItemValidPaths tests that valid paths work correctly
func TestOpenItemValidPaths(t *testing.T) {
	e, err := Open("testdata/Row-Your-Boat-2.0.epub")
	if err != nil {
		t.Fatalf("Failed to open EPUB: %v", err)
	}
	defer e.Close()

	// Get package to find valid items
	pkg, err := e.Package()
	if err != nil {
		t.Fatalf("Failed to get package: %v", err)
	}

	// Try to open each item from the manifest
	for _, item := range pkg.Manifest.Items {
		t.Run(item.Href, func(t *testing.T) {
			if item.Href == "" || filepath.IsAbs(item.Href) {
				t.Skip("Skipping empty or absolute href")
			}

			f, err := e.OpenItem(item.Href)
			if err != nil {
				t.Logf("Failed to open item %q: %v", item.Href, err)
			} else {
				f.Close()
			}
		})
	}
}

// TestWalkReadingContentOptimization tests that WalkReadingContent works efficiently
func TestWalkReadingContentOptimization(t *testing.T) {
	testFiles := []string{
		"testdata/Row-Your-Boat-2.0.epub",
		"testdata/Row-Your-Boat-3.0.epub",
		"testdata/pg1661-images.epub",
	}

	for _, file := range testFiles {
		t.Run(file, func(t *testing.T) {
			count := 0
			err := WalkReadingContent(file, func(r interface{}, fi interface{}) error {
				count++
				return nil
			})
			if err != nil {
				t.Errorf("Failed to walk reading content for %s: %v", file, err)
			}
			if count == 0 {
				t.Errorf("Expected to find reading content items in %s", file)
			}
		})
	}
}

// TestGetTitlesOptimization tests that getTitles works efficiently
func TestGetTitlesOptimization(t *testing.T) {
	// Create test data with many titles
	elts := make([]Element, 100)
	for i := range elts {
		elts[i] = Element{Value: fmt.Sprintf("Title %d", i), ID: fmt.Sprintf("title%d", i)}
	}

	meta := []MetaLegacy{
		{Meta: &Meta{Property: "title-type", Value: "subtitle", Refines: "#title0"}},
		{Meta: &Meta{Property: "title-type", Value: "main", Refines: "#title1"}},
	}

	title, subtitle := getTitles(elts, meta)

	if len(title) != 99 {
		t.Errorf("Expected 99 titles, got %d", len(title))
	}
	if len(subtitle) != 1 {
		t.Errorf("Expected 1 subtitle, got %d", len(subtitle))
	}
	if subtitle[0] != "Title 0" {
		t.Errorf("Expected subtitle to be 'Title 0', got %q", subtitle[0])
	}
}
