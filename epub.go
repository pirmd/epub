// Package epub provides a way to retrieve stored metadata from epub files.
package epub

import (
	"archive/zip"
	"os"
)

// GetPackage reads an epub's Open Package Document from an epub opened as a ReadAtSeeker.
// Deprecated: interacting with epub through ReadAtSeeker will be suppressed in
// next version, seems not useful in practice.
func GetPackage(r ReadAtSeeker) (*PackageDocument, error) {
	zr, err := newZipReader(r)
	if err != nil {
		return nil, err
	}

	return getPackage(zr)
}

// GetPackageFromFile reads an epub's Open Package Document from an epub  file.
func GetPackageFromFile(path string) (*PackageDocument, error) {
	rzip, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer rzip.Close()

	return GetPackage(rzip)
}

func getContainer(zr *zip.Reader) (*container, error) {
	r, err := zr.Open(containerPath)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	return newContainer(r)
}

func getPackage(zr *zip.Reader) (*PackageDocument, error) {
	c, err := getContainer(zr)
	if err != nil {
		return nil, err
	}

	r, err := zr.Open(c.Rootfiles.FullPath)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	return newPackageDocument(r)
}
