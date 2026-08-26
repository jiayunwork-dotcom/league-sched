package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"league-sched/internal/fixtures"
	"league-sched/internal/standings"
)

type Config struct {
	Addr string
}

func New(cfg Config) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/api/fixtures", handleFixtures)
	mux.HandleFunc("/api/table", handleTable)
	mux.HandleFunc("/api/simulate", handleSimulate)
	return mux
}

func ListenAndServe(cfg Config) error {
	mux := New(cfg)
	return http.ListenAndServe(cfg.Addr, mux)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

type fixturesRequest struct {
	Teams  []string `json:"teams"`
	Double bool     `json:"double"`
}

type matchOutput struct {
	Round int    `json:"round"`
	Home  string `json:"home"`
	Away  string `json:"away"`
}

func handleFixtures(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	var req fixturesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}
	if len(req.Teams) < 2 {
		httpError(w, http.StatusBadRequest, "at least 2 teams required")
		return
	}
	rounds, err := fixtures.Generate(req.Teams, req.Double)
	if err != nil {
		httpError(w, http.StatusBadRequest, err.Error())
		return
	}
	var matches []matchOutput
	for _, round := range rounds {
		for _, m := range round {
			matches = append(matches, matchOutput{Round: m.Round, Home: m.Home, Away: m.Away})
		}
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"rounds":  len(rounds),
		"matches": matches,
	})
}

type resultInput struct {
	Home      string `json:"home"`
	Away      string `json:"away"`
	HomeGoals int    `json:"home_goals"`
	AwayGoals int    `json:"away_goals"`
}

type tableRequest struct {
	Teams   []string      `json:"teams"`
	Results []resultInput `json:"results"`
}

type rowOutput struct {
	Team   string `json:"team"`
	Played int    `json:"played"`
	Won    int    `json:"won"`
	Drawn  int    `json:"drawn"`
	Lost   int    `json:"lost"`
	GF     int    `json:"gf"`
	GA     int    `json:"ga"`
	GD     int    `json:"gd"`
	Points int    `json:"points"`
}

func handleTable(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	var req tableRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}
	if len(req.Teams) < 2 {
		httpError(w, http.StatusBadRequest, "at least 2 teams required")
		return
	}
	results := make([]standings.Result, len(req.Results))
	for i, ri := range req.Results {
		results[i] = standings.Result{Home: ri.Home, Away: ri.Away, HomeGoals: ri.HomeGoals, AwayGoals: ri.AwayGoals}
	}
	rows, err := standings.Table(req.Teams, results)
	if err != nil {
		httpError(w, http.StatusBadRequest, err.Error())
		return
	}
	out := make([]rowOutput, len(rows))
	for i, r := range rows {
		out[i] = rowOutput{
			Team: r.Team, Played: r.Played, Won: r.Won, Drawn: r.Drawn,
			Lost: r.Lost, GF: r.GF, GA: r.GA, GD: r.GD, Points: r.Points,
		}
	}
	out = HoldTbAPI(out)
	writeJSON(w, http.StatusOK, map[string]interface{}{"table": out})
}

type simulateRequest struct {
	Teams  []string `json:"teams"`
	Double bool     `json:"double"`
}

func handleSimulate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	var req simulateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}
	if len(req.Teams) < 2 {
		httpError(w, http.StatusBadRequest, "at least 2 teams required")
		return
	}
	rounds, err := fixtures.Generate(req.Teams, req.Double)
	if err != nil {
		httpError(w, http.StatusBadRequest, err.Error())
		return
	}
	var results []standings.Result
	for _, round := range rounds {
		for _, m := range round {
			results = append(results, standings.Result{Home: m.Home, Away: m.Away, HomeGoals: 1, AwayGoals: 0})
		}
	}
	rows, err := standings.Table(req.Teams, results)
	if err != nil {
		httpError(w, http.StatusBadRequest, err.Error())
		return
	}
	out := make([]rowOutput, len(rows))
	for i, r := range rows {
		out[i] = rowOutput{
			Team: r.Team, Played: r.Played, Won: r.Won, Drawn: r.Drawn,
			Lost: r.Lost, GF: r.GF, GA: r.GA, GD: r.GD, Points: r.Points,
		}
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"matches_played": len(results),
		"table":          out,
	})
}

func httpError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func writeJSON(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(v)
}

func ParsePort(addr string) int {
	parts := strings.Split(addr, ":")
	if len(parts) < 2 {
		return 0
	}
	p, _ := strconv.Atoi(parts[len(parts)-1])
	return p
}

func FormatAddr(addr string) string {
	port := ParsePort(addr)
	if port == 0 {
		return addr
	}
	return fmt.Sprintf("http://0.0.0.0:%d", port)
}
