# Changelog

## [Unreleased]

### Added
- Full EPUB3 support including collections and version detection
- CollectionInfo struct for EPUB3 collection metadata extraction
- EPUBVersion field in Information struct to track EPUB specification version
- Version() method on Epub type to retrieve EPUB version
- GetEPUBVersion() standalone function for version retrieval
- Comprehensive documentation with usage examples
- Enhanced error messages with context and wrapping
- Path traversal prevention in OpenItem function
- Input validation for empty paths and hrefs

### Changed
- Updated Go version requirement from 1.16 to 1.21
- Updated golang.org/x/net dependency to v0.25.0
- Improved README.md with comprehensive usage examples
- Added EPUB version support table
- Updated code examples to show EPUB3 features

### Fixed
- Path traversal vulnerability in OpenItem function
- Added validation for empty paths and hrefs
- Improved resource cleanup in error paths
- Fixed missing error wrapping in container and package parsing

### Performance
- Optimized WalkReadingContent from O(n²) to O(n) using map lookup for spine item references
- Optimized getTitles from O(n²) to O(n) using map lookup for title-type metadata

## [0.3.0] - 2023-02-26

- Solve dependencies security issues in golang.org/x/text (CVE-2021-38561, CVE-2022-32149) and golang.org/x/net (CVE-2022-27664, CVE-2022-41721)
- Remove use of ReadAtSeeker as a mean to access an EPUB's metadata
- Create an Epub type to gather standard EPUB manipulation

## [0.2.0] - 2022-06-23

- Expose PackageDocument struct and add functions to get it from an epub.
- Add compliance to EPUB32 specifications https://www.w3.org/publishing/epub32/epub-packages.html.
- Add helpers (WalkXxX) to access EPUB publication resources.

## [0.1.0] - 2019-05-11

### Added
- EPUB2 metadata reading.
- epub tool that prints metadata from an epub file to the standard output.
