package output

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"
)

// --- DetectFormat ---

func TestDetectFormat_Explicit_JSON(t *testing.T) {
	f := DetectFormat("json")
	if f != FormatJSON {
		t.Errorf("expected FormatJSON, got %v", f)
	}
}

func TestDetectFormat_Explicit_Table(t *testing.T) {
	f := DetectFormat("table")
	if f != FormatTable {
		t.Errorf("expected FormatTable, got %v", f)
	}
}

func TestDetectFormat_Explicit_CaseInsensitive(t *testing.T) {
	for _, input := range []string{"JSON", "Json", "TABLE", "Table"} {
		f := DetectFormat(input)
		lower := strings.ToLower(input)
		switch lower {
		case "json":
			if f != FormatJSON {
				t.Errorf("input %q: expected FormatJSON, got %v", input, f)
			}
		case "table":
			if f != FormatTable {
				t.Errorf("input %q: expected FormatTable, got %v", input, f)
			}
		}
	}
}

func TestDetectFormat_Empty_NonTTY(t *testing.T) {
	// In tests, stdout is typically not a TTY, so empty should return JSON.
	f := DetectFormat("")
	if f != FormatJSON {
		t.Errorf("expected FormatJSON for non-TTY with empty explicit, got %v", f)
	}
}

func TestDetectFormat_Explicit_YAML(t *testing.T) {
	f := DetectFormat("yaml")
	if f != FormatYAML {
		t.Errorf("expected FormatYAML, got %v", f)
	}
}

func TestDetectFormat_Explicit_YML(t *testing.T) {
	f := DetectFormat("yml")
	if f != FormatYAML {
		t.Errorf("expected FormatYAML for 'yml', got %v", f)
	}
}

func TestDetectFormat_Explicit_CSV(t *testing.T) {
	f := DetectFormat("csv")
	if f != FormatCSV {
		t.Errorf("expected FormatCSV, got %v", f)
	}
}

func TestDetectFormat_UnknownFormat(t *testing.T) {
	f := DetectFormat("xml")
	// In tests (non-TTY), should return JSON.
	if f != FormatJSON {
		t.Errorf("expected FormatJSON for unknown format, got %v", f)
	}
}

// --- PrintJSON ---

func TestPrintJSON_Map(t *testing.T) {
	r, w, _ := os.Pipe()
	oldStdout := os.Stdout
	os.Stdout = w

	data := map[string]string{"key": "value"}
	err := PrintJSON(data)

	_ = w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out, _ := io.ReadAll(r)
	var m map[string]string
	if err := json.Unmarshal(out, &m); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if m["key"] != "value" {
		t.Errorf("expected key=value, got %v", m)
	}
}

func TestPrintJSON_Struct(t *testing.T) {
	r, w, _ := os.Pipe()
	oldStdout := os.Stdout
	os.Stdout = w

	type Item struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	}
	err := PrintJSON(Item{Name: "test", Count: 5})

	_ = w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out, _ := io.ReadAll(r)
	var item Item
	_ = json.Unmarshal(out, &item)
	if item.Name != "test" || item.Count != 5 {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestPrintJSON_Array(t *testing.T) {
	r, w, _ := os.Pipe()
	oldStdout := os.Stdout
	os.Stdout = w

	err := PrintJSON([]int{1, 2, 3})

	_ = w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out, _ := io.ReadAll(r)
	var arr []int
	_ = json.Unmarshal(out, &arr)
	if len(arr) != 3 {
		t.Errorf("expected 3 elements, got %d", len(arr))
	}
}

func TestPrintJSON_Nil(t *testing.T) {
	r, w, _ := os.Pipe()
	oldStdout := os.Stdout
	os.Stdout = w

	err := PrintJSON(nil)

	_ = w.Close()
	os.Stdout = oldStdout

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out, _ := io.ReadAll(r)
	trimmed := strings.TrimSpace(string(out))
	if trimmed != "null" {
		t.Errorf("expected 'null', got %q", trimmed)
	}
}

// --- NewTable ---

func TestNewTable_Headers(t *testing.T) {
	var buf bytes.Buffer
	tw := NewTable(&buf, "ID", "Name", "Status")
	tw.Render()

	output := buf.String()
	upper := strings.ToUpper(output)
	for _, h := range []string{"ID", "NAME", "STATUS"} {
		if !strings.Contains(upper, h) {
			t.Errorf("expected header %q in output, got: %s", h, output)
		}
	}
}

func TestNewTable_Render(t *testing.T) {
	var buf bytes.Buffer
	tw := NewTable(&buf, "Col1", "Col2")
	tw.AppendRow([]interface{}{"val1", "val2"})
	tw.Render()

	output := buf.String()
	if !strings.Contains(output, "val1") {
		t.Errorf("expected 'val1' in output, got: %s", output)
	}
	if !strings.Contains(output, "val2") {
		t.Errorf("expected 'val2' in output, got: %s", output)
	}
}

// --- Print ---

func TestPrint_JSON_Format(t *testing.T) {
	r, w, _ := os.Pipe()
	oldStdout := os.Stdout
	os.Stdout = w

	Print(FormatJSON, map[string]string{"a": "b"}, nil)

	_ = w.Close()
	os.Stdout = oldStdout

	out, _ := io.ReadAll(r)
	if !strings.Contains(string(out), `"a"`) {
		t.Errorf("expected JSON output, got: %s", out)
	}
}

func TestPrint_Table_Format(t *testing.T) {
	r, w, _ := os.Pipe()
	oldStdout := os.Stdout
	os.Stdout = w

	Print(FormatTable, "data", func(wr io.Writer, data any) {
		_, _ = wr.Write([]byte("TABLE_OUTPUT"))
	})

	_ = w.Close()
	os.Stdout = oldStdout

	out, _ := io.ReadAll(r)
	if !strings.Contains(string(out), "TABLE_OUTPUT") {
		t.Errorf("expected 'TABLE_OUTPUT', got: %s", out)
	}
}

func TestPrint_Table_NilRenderer_FallsBackToJSON(t *testing.T) {
	r, w, _ := os.Pipe()
	oldStdout := os.Stdout
	os.Stdout = w

	Print(FormatTable, map[string]string{"x": "y"}, nil)

	_ = w.Close()
	os.Stdout = oldStdout

	out, _ := io.ReadAll(r)
	if !strings.Contains(string(out), `"x"`) {
		t.Errorf("expected JSON fallback output, got: %s", out)
	}
}

// --- Success, Warn, Errorf ---

func TestSuccess_WritesToStderr(t *testing.T) {
	r, w, _ := os.Pipe()
	oldStderr := os.Stderr
	os.Stderr = w

	Success("all good")

	_ = w.Close()
	os.Stderr = oldStderr

	out, _ := io.ReadAll(r)
	if !strings.Contains(string(out), "all good") {
		t.Errorf("expected 'all good' in stderr, got: %s", out)
	}
}

func TestWarn_WritesToStderr(t *testing.T) {
	r, w, _ := os.Pipe()
	oldStderr := os.Stderr
	os.Stderr = w

	Warn("watch out")

	_ = w.Close()
	os.Stderr = oldStderr

	out, _ := io.ReadAll(r)
	if !strings.Contains(string(out), "Warning: watch out") {
		t.Errorf("expected 'Warning: watch out' in stderr, got: %s", out)
	}
}

func TestErrorf_WritesToStderr(t *testing.T) {
	r, w, _ := os.Pipe()
	oldStderr := os.Stderr
	os.Stderr = w

	Errorf("something %s", "failed")

	_ = w.Close()
	os.Stderr = oldStderr

	out, _ := io.ReadAll(r)
	if !strings.Contains(string(out), "Error: something failed") {
		t.Errorf("expected 'Error: something failed' in stderr, got: %s", out)
	}
}

// --- YAML output ---

func TestPrintYAML_Map(t *testing.T) {
	r, w, _ := os.Pipe()
	oldStdout := os.Stdout
	os.Stdout = w

	data := map[string]string{"name": "gpc", "version": "1.0"}
	_ = PrintYAML(data)

	_ = w.Close()
	os.Stdout = oldStdout

	out, _ := io.ReadAll(r)
	s := string(out)
	if !strings.Contains(s, "name: gpc") {
		t.Errorf("expected YAML with 'name: gpc', got: %s", s)
	}
	if !strings.Contains(s, "version: \"1.0\"") && !strings.Contains(s, "version: '1.0'") {
		// YAML may or may not quote strings
		if !strings.Contains(s, "version:") {
			t.Errorf("expected 'version' key in YAML, got: %s", s)
		}
	}
}

func TestPrintYAML_Struct(t *testing.T) {
	r, w, _ := os.Pipe()
	oldStdout := os.Stdout
	os.Stdout = w

	type item struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	}
	_ = PrintYAML(item{Name: "test", Count: 42})

	_ = w.Close()
	os.Stdout = oldStdout

	out, _ := io.ReadAll(r)
	s := string(out)
	if !strings.Contains(s, "name: test") {
		t.Errorf("expected 'name: test' in YAML, got: %s", s)
	}
	if !strings.Contains(s, "count: 42") {
		t.Errorf("expected 'count: 42' in YAML, got: %s", s)
	}
}

func TestPrint_YAML_Format(t *testing.T) {
	r, w, _ := os.Pipe()
	oldStdout := os.Stdout
	os.Stdout = w

	Print(FormatYAML, map[string]string{"key": "val"}, nil)

	_ = w.Close()
	os.Stdout = oldStdout

	out, _ := io.ReadAll(r)
	if !strings.Contains(string(out), "key: val") {
		t.Errorf("expected YAML output with 'key: val', got: %s", out)
	}
}

// --- CSV output ---

func TestPrintCSV_Simple(t *testing.T) {
	r, w, _ := os.Pipe()
	oldStdout := os.Stdout
	os.Stdout = w

	headers := []string{"Name", "Value"}
	rows := [][]string{
		{"foo", "bar"},
		{"baz", "qux"},
	}
	_ = PrintCSV(headers, rows)

	_ = w.Close()
	os.Stdout = oldStdout

	out, _ := io.ReadAll(r)
	s := string(out)
	if !strings.Contains(s, "Name,Value") {
		t.Errorf("expected CSV headers, got: %s", s)
	}
	if !strings.Contains(s, "foo,bar") {
		t.Errorf("expected 'foo,bar' row, got: %s", s)
	}
}

func TestPrintWithCSV_UsesCSV(t *testing.T) {
	r, w, _ := os.Pipe()
	oldStdout := os.Stdout
	os.Stdout = w

	headers := []string{"A", "B"}
	rows := [][]string{{"1", "2"}}
	PrintWithCSV(FormatCSV, nil, nil, headers, rows)

	_ = w.Close()
	os.Stdout = oldStdout

	out, _ := io.ReadAll(r)
	s := string(out)
	if !strings.Contains(s, "A,B") {
		t.Errorf("expected CSV output, got: %s", s)
	}
}

func TestPrintWithCSV_FallsBackToJSON(t *testing.T) {
	r, w, _ := os.Pipe()
	oldStdout := os.Stdout
	os.Stdout = w

	PrintWithCSV(FormatJSON, map[string]string{"x": "y"}, nil, nil, nil)

	_ = w.Close()
	os.Stdout = oldStdout

	out, _ := io.ReadAll(r)
	if !strings.Contains(string(out), `"x"`) {
		t.Errorf("expected JSON output, got: %s", out)
	}
}

// --- NoColor ---

func TestNoColor_NotSet(t *testing.T) {
	// t.Setenv will restore the original value after the test.
	// We need to actually unset it, not set to empty.
	_ = os.Unsetenv("NO_COLOR")
	if NoColor() {
		t.Error("expected NoColor() to return false when NO_COLOR is not set")
	}
}

func TestNoColor_Set(t *testing.T) {
	t.Setenv("NO_COLOR", "")
	if !NoColor() {
		t.Error("expected NoColor() to return true when NO_COLOR is set to empty string")
	}
}

func TestNoColor_SetWithValue(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	if !NoColor() {
		t.Error("expected NoColor() to return true when NO_COLOR is set to '1'")
	}
}

func TestNewTable_NoColor(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	var buf bytes.Buffer
	tw := NewTable(&buf, "Col1", "Col2")
	tw.AppendRow([]interface{}{"val1", "val2"})
	tw.Render()
	// Just verify no panic and some output is produced.
	if buf.Len() == 0 {
		t.Error("expected non-empty table output with NO_COLOR set")
	}
}
