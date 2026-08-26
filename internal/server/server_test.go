package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthEndpoint(t *testing.T) {
	mux := New(Config{Addr: ":8080"})
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestFixturesEndpoint(t *testing.T) {
	mux := New(Config{Addr: ":8080"})
	payload := fixturesRequest{Teams: []string{"A", "B", "C", "D"}, Double: false}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/fixtures", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var resp struct {
		Rounds  int           `json:"rounds"`
		Matches []matchOutput `json:"matches"`
	}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	if resp.Rounds != 3 {
		t.Errorf("expected 3 rounds for 4 teams, got %d", resp.Rounds)
	}
	if len(resp.Matches) != 6 {
		t.Errorf("expected 6 matches, got %d", len(resp.Matches))
	}
}

func TestFixturesEndpoint_TooFew(t *testing.T) {
	mux := New(Config{Addr: ":8080"})
	body := []byte(`{"teams":["A"]}`)
	req := httptest.NewRequest(http.MethodPost, "/api/fixtures", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestTableEndpoint(t *testing.T) {
	mux := New(Config{Addr: ":8080"})
	payload := tableRequest{
		Teams: []string{"X", "Y", "Z"},
		Results: []resultInput{
			{Home: "X", Away: "Y", HomeGoals: 2, AwayGoals: 1},
			{Home: "Y", Away: "Z", HomeGoals: 0, AwayGoals: 0},
			{Home: "Z", Away: "X", HomeGoals: 1, AwayGoals: 3},
		},
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/table", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var resp struct {
		Table []rowOutput `json:"table"`
	}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	if len(resp.Table) != 3 {
		t.Fatalf("expected 3 rows, got %d", len(resp.Table))
	}
	if resp.Table[0].Team != "X" || resp.Table[0].Points != 6 {
		t.Errorf("expected X with 6 pts at top, got %s with %d", resp.Table[0].Team, resp.Table[0].Points)
	}
}

func TestSimulateEndpoint(t *testing.T) {
	mux := New(Config{Addr: ":8080"})
	payload := simulateRequest{Teams: []string{"A", "B", "C"}, Double: false}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/simulate", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var resp struct {
		MatchesPlayed int         `json:"matches_played"`
		Table         []rowOutput `json:"table"`
	}
	json.Unmarshal(rec.Body.Bytes(), &resp)
	if resp.MatchesPlayed == 0 {
		t.Error("expected matches played > 0")
	}
	if len(resp.Table) != 3 {
		t.Errorf("expected 3 rows, got %d", len(resp.Table))
	}
}

func TestMethodNotAllowed(t *testing.T) {
	mux := New(Config{Addr: ":8080"})
	endpoints := []string{"/api/fixtures", "/api/table", "/api/simulate"}
	for _, ep := range endpoints {
		req := httptest.NewRequest(http.MethodGet, ep, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("%s: expected 405, got %d", ep, rec.Code)
		}
	}
}
