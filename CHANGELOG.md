# Changelog

## [Unreleased]
### Added
- Comprehensive benchmark tests for performance monitoring
- Improved documentation with usage examples
- Enhanced error messages with context

### Changed
- Updated Go version requirement from 1.16 to 1.21
- Updated golang.org/x/net dependency to v0.25.0

### Fixed
- Path traversal vulnerability in OpenItem function
- Added validation for empty paths and hrefs
- Improved resource cleanup in error paths

### Performance
- Optimized WalkReadingContent from O(n²) to O(n) using map lookup
- Optimized getTitles from O(n²) to O(n) using map lookup

## [0.3.0] - 2023-02-26
- Solved dependencies security issues in golang.org/x/text (CVE-2021-38561,
  CVE-2022-32149) and golang.org/x/net (CVE-2022-27664, CVE-2022-41721)
- Removed use of ReadAtSeeker as a mean to access an EPUB's metadata
- Created an Epub type to gather standard EPUB manipulation

## [0.2.0] - 2022-06-23
- Exposed PackageDocument struct and added functions to get it from an epub
- Added compliance to EPUB32 specifications https://www.w3.org/publishing/epub32/epub-packages.html
- Added helpers (WalkXxX) to access EPUB publication resources

## [0.1.0] - 2019-05-11
### Added
- EPUB2 metadata reading
- epub tool that prints metadata from an epub file to the standard output
