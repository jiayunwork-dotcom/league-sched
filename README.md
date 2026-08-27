# league-sched

Round-robin fixture generation and league standings, standard library only.

## Build & Test

```sh
go build ./...
go test ./...
```

## Usage

```sh
league-sched fixtures example/teams.txt            # single round robin
league-sched fixtures -double example/teams.txt    # home and away legs
league-sched table example/teams.txt example/results.txt
```

`teams.txt`: one team name per line. `results.txt`: one match per line,
`HOME AWAY HOME_GOALS AWAY_GOALS`. Standings sort by points, then goal
difference, then goals for, then name.

## Packages

- `internal/fixtures` — round-robin scheduling (circle method, bye handling)
- `internal/standings` — result validation and league table computation

## Example

`example/teams.txt` and `example/results.txt` hold a six-team league.
