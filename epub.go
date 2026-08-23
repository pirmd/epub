// Package epub provides a way to retrieve stored metadata from epub files.
package epub

import (
	"archive/zip"
	"fmt"
	"io/fs"
	"net/url"
	"path/filepath"
	"strings"
)

// Epub represents a read-only EPUB document.
type Epub struct {
	*zip.ReadCloser

	rootfile string
}

// Open an EPUB from a file.
// Returned Epub needs to be closed when no longer needed.
func Open(path string) (*Epub, error) {
	if path == "" {
		return nil, fmt.Errorf("empty path")
	}

	e := new(Epub)

	var err error
	if e.ReadCloser, err = zip.OpenReader(path); err != nil {
		return nil, fmt.Errorf("failed to open EPUB file: %w", err)
	}

	c, err := e.container()
	if err != nil {
		e.Close()
		return nil, fmt.Errorf("failed to read container: %w", err)
	}

	e.rootfile = c.Rootfiles.FullPath

	return e, nil
}

// OpenItem opens an EPUB Publication Resource identified by its href as
// usually found in Manifest.
// OpenItem will try to unescape href first.
// Opening Items whose Href points outside of EPUB archive will fail.
func (e *Epub) OpenItem(href string) (fs.File, error) {
	if href == "" {
		return nil, fmt.Errorf("empty href")
	}

	name, err := url.PathUnescape(href)
	if err != nil {
		return nil, fmt.Errorf("failed to unescape href %q: %w", href, err)
	}

	// Prevent path traversal: ensure the resulting path stays within the EPUB archive
	baseDir := filepath.Dir(e.rootfile)
	path := filepath.Join(baseDir, name)
	
	// Clean the path to remove any ".." or "." components
	cleanPath := filepath.Clean(path)
	
	// Ensure the cleaned path is still within the base directory
	if !filepath.IsAbs(cleanPath) {
		cleanPath = filepath.Join("/", cleanPath)
	}
	
	baseClean := filepath.Clean(filepath.Join("/", baseDir))
	if !strings.HasPrefix(cleanPath, baseClean) {
		return nil, fmt.Errorf("invalid item path %q: points outside EPUB archive", href)
	}
	
	// Remove leading "/" for zip archive access
	zipPath := strings.TrimPrefix(cleanPath, "/")
	
	// Verify the file exists in the archive
	found := false
	for _, f := range e.File {
		if f.Name == zipPath {
			found = true
			break
		}
	}
	
	if !found {
		return nil, fmt.Errorf("item %q not found in EPUB archive", href)
	}
	
	return e.Open(zipPath)
}

// container returns the EPUB Container.
func (e *Epub) container() (*container, error) {
	r, err := e.Open(containerPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open container.xml: %w", err)
	}
	defer r.Close()

	c, err := newContainer(r)
	if err != nil {
		return nil, fmt.Errorf("failed to parse container.xml: %w", err)
	}
	return c, nil
}

// Package returns the EPUB PackageDocument.
func (e *Epub) Package() (*PackageDocument, error) {
	r, err := e.Open(e.rootfile)
	if err != nil {
		return nil, fmt.Errorf("failed to open package document: %w", err)
	}
	defer r.Close()

	return newPackageDocument(r)
}

// Information returns a simplified but easier to use version of
// PackageDocument.Metadata.
func (e *Epub) Information() (*Information, error) {
	opf, err := e.Package()
	if err != nil {
		return nil, fmt.Errorf("failed to get package document: %w", err)
	}

	return getMeta(opf.Metadata), nil
}

// GetPackageFromFile reads an epub's Open Package Document from an epub  file.
func GetPackageFromFile(path string) (*PackageDocument, error) {
	e, err := Open(path)
	if err != nil {
		return nil, err
	}
	defer e.Close()

	return e.Package()
}
