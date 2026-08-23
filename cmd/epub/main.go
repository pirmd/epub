package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pirmd/epub"
)

var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func stringify(v interface{}, pretty bool) string {
	var b []byte
	var err error
	
	if pretty {
		b, err = json.MarshalIndent(v, "", "  ")
	} else {
		b, err = json.Marshal(v)
	}
	
	if err != nil {
		return fmt.Sprintf("%+v", v)
	}

	return string(b)
}

func printField(metadata *epub.Information, field string) {
	switch field {
	case "title":
		fmt.Printf("%v\n", metadata.Title)
	case "subtitle":
		fmt.Printf("%v\n", metadata.SubTitle)
	case "author", "creator":
		for _, a := range metadata.Creator {
			fmt.Printf("%s\n", a.FullName)
		}
	case "contributor":
		for _, a := range metadata.Contributor {
			fmt.Printf("%s\n", a.FullName)
		}
	case "identifier":
		for _, id := range metadata.Identifier {
			if id.Scheme != "" {
				fmt.Printf("%s: %s\n", id.Scheme, id.Value)
			} else {
				fmt.Printf("%s\n", id.Value)
			}
		}
	case "language":
		fmt.Printf("%v\n", metadata.Language)
	case "publisher":
		fmt.Printf("%v\n", metadata.Publisher)
	case "date":
		for _, d := range metadata.Date {
			if d.Event != "" {
				fmt.Printf("%s: %s\n", d.Event, d.Stamp)
			} else {
				fmt.Printf("%s\n", d.Stamp)
			}
		}
	case "description":
		fmt.Printf("%v\n", metadata.Description)
	case "subject":
		fmt.Printf("%v\n", metadata.Subject)
	case "series":
		if metadata.Series != "" {
			fmt.Printf("%s", metadata.Series)
			if metadata.SeriesIndex != "" {
				fmt.Printf(" #%s", metadata.SeriesIndex)
			}
			fmt.Println()
		}
	case "rights":
		fmt.Printf("%v\n", metadata.Rights)
	case "type":
		fmt.Printf("%v\n", metadata.Type)
	case "format":
		fmt.Printf("%v\n", metadata.Format)
	case "source":
		fmt.Printf("%v\n", metadata.Source)
	case "coverage":
		fmt.Printf("%v\n", metadata.Coverage)
	case "relation":
		fmt.Printf("%v\n", metadata.Relation)
	default:
		// Try to access via reflection for custom fields
		fmt.Printf("Unknown field: %s\n", field)
		os.Exit(1)
	}
}

func printVersion(w io.Writer) {
	fmt.Fprintf(w, "epub version %s\n", version)
	fmt.Fprintf(w, "commit: %s\n", commit)
	fmt.Fprintf(w, "built: %s\n", date)
}

func main() {
	// Define flags
	var (
		pretty   bool
		compact  bool
		versionFlag bool
		fields   string
		help     bool
	)

	flag.BoolVar(&pretty, "pretty", true, "Pretty-print JSON output")
	flag.BoolVar(&compact, "compact", false, "Compact JSON output (no indentation)")
	flag.BoolVar(&versionFlag, "version", false, "Print version and exit")
	flag.StringVar(&fields, "field", "", "Print only specific field(s). Comma-separated list: title,author,language,identifier,publisher,date,description,series,etc.")
	flag.BoolVar(&help, "help", false, "Print help and exit")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] <epub files>\n\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(flag.CommandLine.Output(), "Command-line tool for reading EPUB metadata.\n\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), "\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Fields:\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  title, subtitle, author, creator, contributor,\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  identifier, language, publisher, date, description,\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  subject, series, rights, type, format, source,\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  coverage, relation\n")
		fmt.Fprintf(flag.CommandLine.Output(), "\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Examples:\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  %s book.epub                    # Print all metadata (pretty JSON)\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(flag.CommandLine.Output(), "  %s --compact book.epub         # Print all metadata (compact JSON)\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(flag.CommandLine.Output(), "  %s --field title book.epub     # Print only title\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(flag.CommandLine.Output(), "  %s --field author,language book.epub  # Print author and language\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(flag.CommandLine.Output(), "  %s file1.epub file2.epub       # Process multiple files\n", filepath.Base(os.Args[0]))
	}

	flag.Parse()

	// Handle help flag
	if help {
		flag.Usage()
		os.Exit(0)
	}

	// Handle version flag
	if versionFlag {
		printVersion(os.Stdout)
		os.Exit(0)
	}

	// Validate flags
	if pretty && compact {
		fmt.Fprintf(os.Stderr, "Error: cannot use both --pretty and --compact\n")
		os.Exit(1)
	}

	// Get remaining arguments (epub files)
	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Error: no EPUB file specified\n")
		flag.Usage()
		os.Exit(1)
	}

	// Process each file
	hasError := false
	for i, arg := range args {
		// Print filename header if processing multiple files
		if len(args) > 1 {
			if i > 0 {
				fmt.Println()
			}
			fmt.Printf("=== %s ===\n", arg)
		}

		metadata, err := epub.GetMetadataFromFile(arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while parsing epub %s: %v\n", arg, err)
			hasError = true
			continue
		}

		// If specific fields are requested
		if fields != "" {
			fieldList := strings.Split(fields, ",")
			for _, field := range fieldList {
				field = strings.TrimSpace(field)
				if field != "" {
					printField(metadata, field)
				}
			}
		} else {
			// Print all metadata as JSON
			fmt.Printf("%s\n", stringify(metadata, !compact))
		}
	}

	if hasError {
		os.Exit(1)
	}
	os.Exit(0)
}
