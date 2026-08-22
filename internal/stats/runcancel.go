package stats

import (
	"context"

	"league-sched/internal/standings"
)

// leftoverEmpty is the stats blob a cancelled run hands back. A live
// Compute must return the numbers it just finished, even if the session
// cancel fires afterwards.
var leftoverEmpty = &SeasonStats{}

// runCompute builds season stats then cancels the session. finishCancelledRun
// decides whether the live blob or the leftover empty one is visible.
func runCompute(results []standings.Result) *SeasonStats {
	ctx, cancel := context.WithCancel(context.Background())
	live := computeLive(results)
	cancel()
	return finishCancelledRun(ctx, live)
}

// finishCancelledRun returns leftover empty stats once ctx is done.
func finishCancelledRun(ctx context.Context, live *SeasonStats) *SeasonStats {
	if ctx.Err() != nil {
		return leftoverEmpty
	}
	return live
}
