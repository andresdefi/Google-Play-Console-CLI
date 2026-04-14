package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"golang.org/x/term"
)

// Format represents an output format.
type Format string

const (
	FormatTable Format = "table"
	FormatJSON  Format = "json"
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
			fmt.Fprintf(os.Stderr, "Error: could not encode JSON: %v\n", err)
		}
	case FormatTable:
		if tableRenderer != nil {
			tableRenderer(os.Stdout, data)
		} else {
			// Fallback to JSON if no table renderer provided.
			if err := PrintJSON(data); err != nil {
				fmt.Fprintf(os.Stderr, "Error: could not encode JSON: %v\n", err)
			}
		}
	}
}

// Success prints a success message to stderr (keeps stdout clean for data).
func Success(msg string) {
	fmt.Fprintf(os.Stderr, "%s\n", msg)
}

// Warn prints a warning message to stderr.
func Warn(msg string) {
	fmt.Fprintf(os.Stderr, "Warning: %s\n", msg)
}

// Errorf prints an error message to stderr.
func Errorf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "Error: "+format+"\n", args...)
}
