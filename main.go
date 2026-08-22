package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"league-sched/internal/fixtures"
	"league-sched/internal/standings"
)

func fail(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func reorderFlags(args []string) []string {
	var flags, pos []string
	for i := 0; i < len(args); i++ {
		a := args[i]
		if strings.HasPrefix(a, "-") && a != "-" && a != "--" {
			flags = append(flags, a)
			if !strings.Contains(a, "=") && i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				flags = append(flags, args[i+1])
				i++
			}
			continue
		}
		pos = append(pos, a)
	}
	return append(flags, pos...)
}

func readLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var out []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		s := strings.TrimSpace(sc.Text())
		if s == "" || strings.HasPrefix(s, "#") {
			continue
		}
		out = append(out, s)
	}
	return out, sc.Err()
}

func readTeams(path string) ([]string, error) {
	lines, err := readLines(path)
	if err != nil {
		return nil, err
	}
	if len(lines) == 0 {
		return nil, fmt.Errorf("%s: no teams found", path)
	}
	return lines, nil
}

func readResults(path string) ([]standings.Result, error) {
	lines, err := readLines(path)
	if err != nil {
		return nil, err
	}
	var out []standings.Result
	for i, line := range lines {
		fields := strings.Fields(line)
		if len(fields) != 4 {
			return nil, fmt.Errorf("%s:%d: want \"HOME AWAY HOME_GOALS AWAY_GOALS\", got %q", path, i+1, line)
		}
		hg, err := strconv.Atoi(fields[2])
		if err != nil || hg < 0 {
			return nil, fmt.Errorf("%s:%d: bad home goals %q", path, i+1, fields[2])
		}
		ag, err := strconv.Atoi(fields[3])
		if err != nil || ag < 0 {
			return nil, fmt.Errorf("%s:%d: bad away goals %q", path, i+1, fields[3])
		}
		out = append(out, standings.Result{Home: fields[0], Away: fields[1], HomeGoals: hg, AwayGoals: ag})
	}
	return out, nil
}

func cmdFixtures(args []string) {
	fs := flag.NewFlagSet("fixtures", flag.ExitOnError)
	double := fs.Bool("double", false, "schedule home and away legs")
	fail(fs.Parse(args))
	if fs.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "usage: league-sched fixtures [-double] <teams-file>")
		os.Exit(2)
	}
	teams, err := readTeams(fs.Arg(0))
	fail(err)
	rounds, err := fixtures.Generate(teams, *double)
	fail(err)
	for _, round := range rounds {
		for _, m := range round {
			fmt.Printf("R%-2d  %-12s vs %-12s\n", m.Round, m.Home, m.Away)
		}
	}
}

func cmdTable(args []string) {
	fs := flag.NewFlagSet("table", flag.ExitOnError)
	fail(fs.Parse(args))
	if fs.NArg() != 2 {
		fmt.Fprintln(os.Stderr, "usage: league-sched table <teams-file> <results-file>")
		os.Exit(2)
	}
	teams, err := readTeams(fs.Arg(0))
	fail(err)
	results, err := readResults(fs.Arg(1))
	fail(err)
	rows, err := standings.Table(teams, results)
	fail(err)
	fmt.Printf("%-14s %2s %2s %2s %2s %3s %3s %3s %3s\n",
		"Team", "P", "W", "D", "L", "GF", "GA", "GD", "Pts")
	for _, r := range rows {
		fmt.Printf("%-14s %2d %2d %2d %2d %3d %3d %3d %3d\n",
			r.Team, r.Played, r.Won, r.Drawn, r.Lost, r.GF, r.GA, r.GD, r.Points)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, `usage: league-sched <command> [flags] <args>

commands:
  fixtures [-double] <teams-file>       print the round-robin schedule
  table <teams-file> <results-file>     print the standings table`)
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	cmd := os.Args[1]
	args := reorderFlags(os.Args[2:])
	switch cmd {
	case "fixtures":
		cmdFixtures(args)
	case "table":
		cmdTable(args)
	case "help", "-h":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n\n", cmd)
		usage()
		os.Exit(2)
	}
}
