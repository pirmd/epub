# EPUB

[![GoDoc](https://pkg.go.dev/badge/github.com/pirmd/epub)](https://pkg.go.dev/github.com/pirmd/epub)
[![Go Report Card](https://goreportcard.com/badge/github.com/pirmd/epub)](https://goreportcard.com/report/github.com/pirmd/epub)
[![License](https://img.shields.io/badge/license-BSD-blue.svg)](LICENSE)

`epub` is a Go package for reading metadata from EPUB files. It provides comprehensive support for both EPUB 2 and EPUB 3.2 specifications.

## Features

- **Metadata Extraction**: Read titles, authors, identifiers (ISBN, DOI, UUID), languages, dates, publishers, and more
- **Package Document Parsing**: Full support for EPUB OPF (Open Packaging Format) files
- **Resource Access**: Walk through EPUB files, publication resources, and reading content
- **EPUB 2 & 3.2 Compatible**: Works with both legacy and modern EPUB formats
- **Series Support**: Extract series information from Calibre and EPUB3 metadata

## Installation

```bash
go get github.com/pirmd/epub
```

## Usage

### Basic Metadata Extraction

```go
package main

import (
    "fmt"
    "github.com/pirmd/epub"
)

func main() {
    // Read metadata from an EPUB file
    metadata, err := epub.GetMetadataFromFile("book.epub")
    if err != nil {
        panic(err)
    }

    // Access metadata
    fmt.Printf("Title: %v\n", metadata.Title)
    fmt.Printf("Authors: %v\n", metadata.Creator)
    fmt.Printf("Language: %v\n", metadata.Language)
    fmt.Printf("Publisher: %v\n", metadata.Publisher)
    
    // Access identifiers (ISBN, DOI, etc.)
    for _, id := range metadata.Identifier {
        fmt.Printf("Identifier (%s): %s\n", id.Scheme, id.Value)
    }
}
```

### Access Package Document

```go
package main

import (
    "fmt"
    "github.com/pirmd/epub"
)

func main() {
    // Get the full PackageDocument
    pkg, err := epub.GetPackageFromFile("book.epub")
    if err != nil {
        panic(err)
    }

    fmt.Printf("EPUB Version: %s\n", pkg.Version)
    fmt.Printf("Unique Identifier: %s\n", pkg.UniqueIdentifier)
    
    // Access manifest items
    for _, item := range pkg.Manifest.Items {
        fmt.Printf("Item: %s (type: %s)\n", item.ID, item.MediaType)
    }
}
```

### Walking EPUB Resources

```go
package main

import (
    "fmt"
    "io"
    "github.com/pirmd/epub"
)

func main() {
    // Walk all files in the EPUB
    err := epub.WalkFiles("book.epub", func(r io.Reader, info interface{}) error {
        fmt.Printf("File: %s\n", info.Name())
        return nil
    })
    if err != nil {
        panic(err)
    }

    // Walk only publication resources (from manifest)
    err = epub.WalkPublicationResources("book.epub", func(r io.Reader, info interface{}) error {
        fmt.Printf("Publication Resource: %s\n", info.Name())
        return nil
    })

    // Walk reading content (from spine, in order)
    err = epub.WalkReadingContent("book.epub", func(r io.Reader, info interface{}) error {
        fmt.Printf("Reading Content: %s\n", info.Name())
        return nil
    })
}
```

### Using the Epub Type

```go
package main

import (
    "fmt"
    "github.com/pirmd/epub"
)

func main() {
    // Open an EPUB file
    e, err := epub.Open("book.epub")
    if err != nil {
        panic(err)
    }
    defer e.Close()

    // Get package document
    pkg, err := e.Package()
    if err != nil {
        panic(err)
    }

    // Get simplified metadata
    info, err := e.Information()
    if err != nil {
        panic(err)
    }

    fmt.Printf("Title: %v\n", info.Title)

    // Open a specific item by href
    file, err := e.OpenItem("chapter1.xhtml")
    if err != nil {
        panic(err)
    }
    defer file.Close()
    
    // Read the file content
    // ...
}
```

## Command Line Tool

A simple command-line tool is available to print EPUB metadata in JSON format.

### Installation

```bash
go install github.com/pirmd/epub/cmd/epub@latest
```

### Usage

```bash
# Print all metadata
epub mybook.epub

# Pretty-print JSON output (default)
epub --pretty mybook.epub

# Compact JSON output
epub --compact mybook.epub
```

### Example Output

```json
{
  "Title": ["The Great Novel"],
  "Creator": [
    {
      "FullName": "Jane Doe",
      "FileAs": "Doe, Jane",
      "Role": "author"
    }
  ],
  "Identifier": [
    {
      "Scheme": "ISBN",
      "Value": "978-1-23456-789-0"
    }
  ],
  "Language": ["en"],
  "Publisher": ["Example Publisher"],
  "Series": "My Series",
  "SeriesIndex": "1"
}
```

## EPUB Specifications

This library aims to support:
- [EPUB 3.2 Specification](https://www.w3.org/publishing/epub32/epub-packages.html)
- EPUB 2 (legacy support)

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the BSD License - see the [LICENSE](LICENSE) file for details.

[modeline]: # ( vim: set fenc=utf-8 spell spl=en: )
