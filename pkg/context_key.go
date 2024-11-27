package pkg

import "context"

type contextKey string

const checkpointKey contextKey = "checkpoint"

func SetCheckpoint(ctx context.Context, point string, success bool) context.Context {
	checkpoints, ok := ctx.Value(checkpointKey).(map[string]bool)
	if !ok {
		checkpoints = make(map[string]bool)
	}
	checkpoints[point] = success
	return context.WithValue(ctx, checkpointKey, checkpoints)
}

func GetCheckpoint(ctx context.Context, point string) (bool, bool) {
	checkpoints, ok := ctx.Value(checkpointKey).(map[string]bool)
	if !ok {
		return false, false
	}
	success, exists := checkpoints[point]
	return success, exists
}

func GetAllCheckpoints(ctx context.Context) map[string]bool {
	checkpoints, ok := ctx.Value(checkpointKey).(map[string]bool)
	if !ok {
		return make(map[string]bool)
	}
	return checkpoints
}
