package output

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"golang.org/x/term"
	"gopkg.in/yaml.v3"
)

// Format represents an output format.
type Format string

const (
	FormatTable Format = "table"
	FormatJSON  Format = "json"
	FormatCSV   Format = "csv"
	FormatYAML  Format = "yaml"
)

// IsTTY reports whether stdout is a terminal.
func IsTTY() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}

// DetectFormat returns the appropriate format based on TTY status and explicit override.
func DetectFormat(explicit string) Format {
	if explicit != "" {
		switch strings.ToLower(explicit) {
		case "json":
			return FormatJSON
		case "table":
			return FormatTable
		case "csv":
			return FormatCSV
		case "yaml", "yml":
			return FormatYAML
		}
	}
	if IsTTY() {
		return FormatTable
	}
	return FormatJSON
}

// PrintJSON writes data as indented JSON to stdout.
func PrintJSON(data any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}

// PrintYAML writes data as YAML to stdout.
func PrintYAML(data any) error {
	// Round-trip through JSON to normalize the data (handles json.RawMessage, etc).
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	var normalized any
	if err := json.Unmarshal(jsonBytes, &normalized); err != nil {
		return err
	}
	enc := yaml.NewEncoder(os.Stdout)
	enc.SetIndent(2)
	defer func() { _ = enc.Close() }()
	return enc.Encode(normalized)
}

// PrintCSV writes data as CSV to stdout. The csvRenderer extracts headers and rows.
func PrintCSV(headers []string, rows [][]string) error {
	w := csv.NewWriter(os.Stdout)
	if err := w.Write(headers); err != nil {
		return err
	}
	for _, row := range rows {
		if err := w.Write(row); err != nil {
			return err
		}
	}
	w.Flush()
	return w.Error()
}

// NewTable creates a new table writer configured for terminal output.
func NewTable(w io.Writer, headers ...string) table.Writer {
	t := table.NewWriter()
	t.SetOutputMirror(w)
	t.SetStyle(table.StyleLight)

	row := make(table.Row, len(headers))
	for i, h := range headers {
		row[i] = h
	}
	t.AppendHeader(row)
	return t
}

// Print outputs data in the given format. The tableRenderer is called for table format.
func Print(format Format, data any, tableRenderer func(io.Writer, any)) {
	switch format {
	case FormatJSON:
		if err := PrintJSON(data); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error: could not encode JSON: %v\n", err)
		}
	case FormatYAML:
		if err := PrintYAML(data); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error: could not encode YAML: %v\n", err)
		}
	case FormatCSV:
		// CSV requires structured data; fall back to JSON for raw data.
		if err := PrintJSON(data); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error: could not encode JSON: %v\n", err)
		}
	case FormatTable:
		if tableRenderer != nil {
			tableRenderer(os.Stdout, data)
		} else {
			if err := PrintJSON(data); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Error: could not encode JSON: %v\n", err)
			}
		}
	}
}

// PrintWithCSV outputs data with full format support including CSV.
func PrintWithCSV(format Format, data any, tableRenderer func(io.Writer, any), csvHeaders []string, csvRows [][]string) {
	if format == FormatCSV && csvHeaders != nil {
		if err := PrintCSV(csvHeaders, csvRows); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error: could not encode CSV: %v\n", err)
		}
		return
	}
	Print(format, data, tableRenderer)
}

// Success prints a success message to stderr (keeps stdout clean for data).
func Success(msg string) {
	_, _ = fmt.Fprintf(os.Stderr, "%s\n", msg)
}

// Warn prints a warning message to stderr.
func Warn(msg string) {
	_, _ = fmt.Fprintf(os.Stderr, "Warning: %s\n", msg)
}

// Errorf prints an error message to stderr.
func Errorf(format string, args ...any) {
	_, _ = fmt.Fprintf(os.Stderr, "Error: "+format+"\n", args...)
}
