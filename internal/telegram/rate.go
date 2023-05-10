package telegram

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/gotd/td/tgerr"
)

type Sleeper interface {
	Sleep(context.Context, time.Duration) error
}

type SleeperFunc func(context.Context, time.Duration) error

func (f SleeperFunc) Sleep(ctx context.Context, d time.Duration) error {
	return f(ctx, d)
}

type RateGuard struct {
	MaxFloodWait time.Duration
	Sleeper      Sleeper
}

func (g RateGuard) Do(ctx context.Context, call func(context.Context) error) error {
	err := call(ctx)
	if err == nil {
		return nil
	}
	if tgerr.Is(err, "PEER_FLOOD") {
		return err
	}
	wait, ok := tgerr.AsFloodWait(err)
	if !ok {
		return err
	}
	if wait < 0 {
		wait = 0
	}
	if g.MaxFloodWait > 0 && wait > g.MaxFloodWait {
		return fmt.Errorf("telegram flood wait %s exceeds max %s: %w", wait, g.MaxFloodWait, err)
	}
	if deadline, ok := ctx.Deadline(); ok && time.Until(deadline) < wait {
		return fmt.Errorf("telegram flood wait %s exceeds command deadline: %w", wait, err)
	}
	sleeper := g.Sleeper
	if sleeper == nil {
		sleeper = SleeperFunc(func(ctx context.Context, d time.Duration) error {
			timer := time.NewTimer(d)
			defer timer.Stop()
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-timer.C:
				return nil
			}
		})
	}
	slog.Default().Info("telegram: flood wait", "duration", wait.Round(time.Second))
	if err := sleeper.Sleep(ctx, wait); err != nil {
		return err
	}
	return call(ctx)
}
