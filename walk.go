package epub

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/url"
	"path/filepath"
)

var (
	// ErrStopWalk is used as a return value from WalkFunc to
	// indicate that the Walkxxx operation need to be
	// stopped. It is not return as an error by any Walkxxx
	// function.
	ErrStopWalk = errors.New("stop walk")
)

// WalkFunc is the signature of function called by Walkxxx on EPUB's resources.
// Should an error be returned by WalkFn, Walkxxx stops and returns that error.
// Only exception is returning ErrStopWalk error that only interrupts Walkxxx.
type WalkFunc func(r io.Reader, info fs.FileInfo) error

// WalkFiles walks EPUB's files, calling walkFn for each visited resource.
func WalkFiles(path string, walkFn WalkFunc) error {
	zr, err := zip.OpenReader(path)
	if err != nil {
		return err
	}
	defer zr.Close()

	for _, f := range zr.File {
		r, err := f.Open()
		if err != nil {
			return err
		}
		defer r.Close()

		if err := walkFn(r, f.FileHeader.FileInfo()); err != nil {
			if err == ErrStopWalk {
				return nil
			}
			return err
		}
	}

	return nil
}

// WalkPublicationResources walks EPUB's publication resources as listed in
// EPUB's Manifest, calling walkFn for each visited resource.
// Limitation: resources that are not belonging to the EPUB archive itself
// (like remote resources) are silently ignored.
func WalkPublicationResources(path string, walkFn WalkFunc) error {
	zr, err := zip.OpenReader(path)
	if err != nil {
		return err
	}
	defer zr.Close()

	c, err := getContainer(&zr.Reader)
	if err != nil {
		return err
	}

	r, err := zr.Open(c.Rootfiles.FullPath)
	if err != nil {
		return err
	}
	defer r.Close()

	opf, err := newPackageDocument(r)
	if err != nil {
		return err
	}

	for _, item := range opf.Manifest.Items {
		if item.Href == "" || filepath.IsAbs(item.Href) {
			continue
		}

		itempath, err := url.PathUnescape(item.Href)
		if err != nil {
			return err
		}

		itempath = filepath.Join(filepath.Dir(c.Rootfiles.FullPath), itempath)
		f, err := zr.Open(itempath)
		if err != nil {
			return err
		}
		defer f.Close()

		fi, err := f.Stat()
		if err != nil {
			return err
		}

		if err := walkFn(f, fi); err != nil {
			if err == ErrStopWalk {
				return nil
			}
			return err
		}
	}

	return nil
}

// WalkReadingContent walks EPUB's publication resources as listed in
// EPUB's Spine, calling walkFn for each visited resource.
// Limitation: resources that are not belonging to the EPUB archive itself
// (like remote resources) are silently ignored.
func WalkReadingContent(path string, walkFn WalkFunc) error {
	zr, err := zip.OpenReader(path)
	if err != nil {
		return err
	}
	defer zr.Close()

	c, err := getContainer(&zr.Reader)
	if err != nil {
		return err
	}

	r, err := zr.Open(c.Rootfiles.FullPath)
	if err != nil {
		return err
	}
	defer r.Close()

	opf, err := newPackageDocument(r)
	if err != nil {
		return err
	}

	for _, itemref := range opf.Spine.Itemrefs {
		var item *Item
		for _, it := range opf.Manifest.Items {
			if itemref.IDref != "" && itemref.IDref == it.ID {
				item = &it
				break
			}
		}

		if item == nil {
			return fmt.Errorf("found a Spine %s entry that does not exist in Manifest", itemref.IDref)
		}

		if item.Href == "" || filepath.IsAbs(item.Href) {
			continue
		}

		itempath, err := url.PathUnescape(item.Href)
		if err != nil {
			return err
		}

		itempath = filepath.Join(filepath.Dir(c.Rootfiles.FullPath), itempath)
		f, err := zr.Open(itempath)
		if err != nil {
			return err
		}
		defer f.Close()

		fi, err := f.Stat()
		if err != nil {
			return err
		}

		if err := walkFn(f, fi); err != nil {
			if err == ErrStopWalk {
				return nil
			}
			return err
		}
	}

	return nil
}
