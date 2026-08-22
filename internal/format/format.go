// Package format provides multi-format output for schedules and standings:
// JSON, CSV, and Markdown table rendering.
package format

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"league-sched/internal/fixtures"
	"league-sched/internal/standings"
)

// WriteTableJSON writes standings as JSON.
func WriteTableJSON(w io.Writer, rows []standings.Row) error {
	data, err := json.MarshalIndent(rows, "", "  ")
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

// WriteTableCSV writes standings as CSV.
func WriteTableCSV(w io.Writer, rows []standings.Row) error {
	cw := csv.NewWriter(w)
	cw.Write([]string{"Team", "P", "W", "D", "L", "GF", "GA", "GD", "Pts"})
	for _, r := range rows {
		cw.Write([]string{
			r.Team,
			itoa(r.Played), itoa(r.Won), itoa(r.Drawn), itoa(r.Lost),
			itoa(r.GF), itoa(r.GA), itoa(r.GD), itoa(r.Points),
		})
	}
	cw.Flush()
	return cw.Error()
}

// WriteTableMarkdown writes standings as a Markdown table.
func WriteTableMarkdown(w io.Writer, rows []standings.Row) error {
	fmt.Fprintln(w, "| Team | P | W | D | L | GF | GA | GD | Pts |")
	fmt.Fprintln(w, "|------|---|---|---|---|----|----|----|----|")
	for _, r := range rows {
		fmt.Fprintf(w, "| %s | %d | %d | %d | %d | %d | %d | %d | %d |\n",
			r.Team, r.Played, r.Won, r.Drawn, r.Lost, r.GF, r.GA, r.GD, r.Points)
	}
	return nil
}

// WriteFixturesJSON writes a schedule as JSON.
func WriteFixturesJSON(w io.Writer, rounds [][]fixtures.Match) error {
	data, err := json.MarshalIndent(rounds, "", "  ")
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

// WriteFixturesCSV writes a schedule as CSV.
func WriteFixturesCSV(w io.Writer, rounds [][]fixtures.Match) error {
	cw := csv.NewWriter(w)
	cw.Write([]string{"Round", "Home", "Away"})
	for _, round := range rounds {
		for _, m := range round {
			cw.Write([]string{itoa(m.Round), m.Home, m.Away})
		}
	}
	cw.Flush()
	return cw.Error()
}

// WriteFixturesMarkdown writes a schedule as a Markdown table.
func WriteFixturesMarkdown(w io.Writer, rounds [][]fixtures.Match) error {
	fmt.Fprintln(w, "| Round | Home | Away |")
	fmt.Fprintln(w, "|-------|------|------|")
	for _, round := range rounds {
		for _, m := range round {
			fmt.Fprintf(w, "| %d | %s | %s |\n", m.Round, m.Home, m.Away)
		}
	}
	return nil
}

func itoa(n int) string {
	return fmt.Sprintf("%d", n)
}

// FormatRecord returns a W-D-L string for a team's record.
func FormatRecord(r standings.Row) string {
	return fmt.Sprintf("%dW-%dD-%dL", r.Won, r.Drawn, r.Lost)
}

// FormatFormString converts a form slice ["W","D","L",...] to a compact string "WDL...".
func FormatFormString(form []string) string {
	return strings.Join(form, "")
}
