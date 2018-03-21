package sync

import (
	"context"
	"testing"
)

func Test_isDone(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	done := isDone(ctx)
	if done {
		t.Error("done has bug  1  ")
	}
	cancel()
	done = isDone(ctx)
	if !done {
		t.Error("done has bug  2  ")
	}
}
