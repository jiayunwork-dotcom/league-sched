package format

import (
	"bytes"
	"strings"
	"testing"

	"league-sched/internal/fixtures"
	"league-sched/internal/standings"
)

func sampleRows() []standings.Row {
	return []standings.Row{
		{Team: "Alpha", Played: 3, Won: 2, Drawn: 1, Lost: 0, GF: 7, GA: 2, GD: 5, Points: 7},
		{Team: "Beta", Played: 3, Won: 1, Drawn: 1, Lost: 1, GF: 4, GA: 4, GD: 0, Points: 4},
	}
}

func sampleRounds() [][]fixtures.Match {
	return [][]fixtures.Match{
		{
			{Round: 1, Home: "Alpha", Away: "Beta"},
			{Round: 1, Home: "Gamma", Away: "Delta"},
		},
		{
			{Round: 2, Home: "Beta", Away: "Gamma"},
			{Round: 2, Home: "Delta", Away: "Alpha"},
		},
	}
}

func TestWriteTableJSON(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteTableJSON(&buf, sampleRows()); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "Alpha") || !strings.Contains(out, "Beta") {
		t.Fatalf("JSON missing teams: %s", out)
	}
}

func TestWriteTableCSV(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteTableCSV(&buf, sampleRows()); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "Alpha") {
		t.Fatalf("CSV missing Alpha: %s", out)
	}
}

func TestWriteTableMarkdown(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteTableMarkdown(&buf, sampleRows()); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "|") {
		t.Fatalf("Markdown missing table syntax: %s", out)
	}
}

func TestWriteFixturesJSON(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteFixturesJSON(&buf, sampleRounds()); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "Alpha") {
		t.Fatal("fixtures JSON missing Alpha")
	}
}

func TestWriteFixturesCSV(t *testing.T) {
	var buf bytes.Buffer
	if err := WriteFixturesCSV(&buf, sampleRounds()); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "Gamma") {
		t.Fatal("fixtures CSV missing Gamma")
	}
}

func TestFormatRecord(t *testing.T) {
	r := standings.Row{Team: "Test", Played: 2, Won: 1, Drawn: 0, Lost: 1, GF: 3, GA: 2, GD: 1, Points: 3}
	s := FormatRecord(r)
	if !strings.Contains(s, "1W") {
		t.Fatalf("FormatRecord missing win count: %s", s)
	}
}

func TestFormatFormString(t *testing.T) {
	form := []string{"W", "L", "D"}
	s := FormatFormString(form)
	if s == "" {
		t.Fatal("expected non-empty form string")
	}
}
