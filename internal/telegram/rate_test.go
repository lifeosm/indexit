package telegram

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gotd/td/tgerr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRateGuard_SuccessNoRetry(t *testing.T) {
	var calls int
	var slept bool
	guard := RateGuard{
		Sleeper: SleeperFunc(func(context.Context, time.Duration) error {
			slept = true
			return nil
		}),
	}
	err := guard.Do(t.Context(), func(context.Context) error {
		calls++
		return nil
	})
	require.NoError(t, err)
	assert.Equal(t, 1, calls)
	assert.False(t, slept)
}

func TestRateGuard_RetriesFloodWait(t *testing.T) {
	var calls int
	var slept time.Duration
	guard := RateGuard{
		MaxFloodWait: time.Second,
		Sleeper: SleeperFunc(func(ctx context.Context, d time.Duration) error {
			slept = d
			return nil
		}),
	}

	err := guard.Do(t.Context(), func(ctx context.Context) error {
		calls++
		if calls == 1 {
			return tgerr.New(420, "FLOOD_WAIT_1")
		}
		return nil
	})

	require.NoError(t, err)
	assert.Equal(t, 2, calls)
	assert.Equal(t, time.Second, slept)
}

func TestRateGuard_RejectsLongFloodWait(t *testing.T) {
	guard := RateGuard{MaxFloodWait: time.Second}

	err := guard.Do(t.Context(), func(context.Context) error {
		return tgerr.New(420, "FLOOD_WAIT_2")
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "exceeds max")
}

func TestRateGuard_PeerFloodFailsFast(t *testing.T) {
	var calls int
	var slept bool
	guard := RateGuard{
		Sleeper: SleeperFunc(func(context.Context, time.Duration) error {
			slept = true
			return nil
		}),
	}
	err := guard.Do(t.Context(), func(context.Context) error {
		calls++
		return tgerr.New(420, "PEER_FLOOD")
	})
	require.Error(t, err)
	assert.True(t, tgerr.Is(err, "PEER_FLOOD"))
	assert.Equal(t, 1, calls, "PEER_FLOOD must not retry")
	assert.False(t, slept, "PEER_FLOOD must not sleep")
}

func TestRateGuard_NonTgErrorPassThrough(t *testing.T) {
	guard := RateGuard{}
	want := errors.New("network down")
	err := guard.Do(t.Context(), func(context.Context) error { return want })
	require.ErrorIs(t, err, want)
}

func TestRateGuard_RejectsFloodWaitExceedingDeadline(t *testing.T) {
	guard := RateGuard{}
	ctx, cancel := context.WithDeadline(t.Context(), time.Now().Add(100*time.Millisecond))
	defer cancel()

	err := guard.Do(ctx, func(context.Context) error {
		return tgerr.New(420, "FLOOD_WAIT_5")
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "exceeds command deadline")
}

func TestRateGuard_SleeperErrorPropagates(t *testing.T) {
	want := errors.New("cancelled while sleeping")
	guard := RateGuard{
		Sleeper: SleeperFunc(func(context.Context, time.Duration) error { return want }),
	}
	err := guard.Do(t.Context(), func(context.Context) error {
		return tgerr.New(420, "FLOOD_WAIT_1")
	})
	require.ErrorIs(t, err, want)
}

func TestRateGuard_ContextCancelDuringSleep(t *testing.T) {
	// Default sleeper honours ctx.Done(). Cancel right away → returns ctx.Err.
	guard := RateGuard{}
	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	err := guard.Do(ctx, func(context.Context) error {
		return tgerr.New(420, "FLOOD_WAIT_1")
	})
	require.ErrorIs(t, err, context.Canceled)
}

func TestRateGuard_MaxFloodWaitZeroIsUnbounded(t *testing.T) {
	var slept time.Duration
	guard := RateGuard{
		Sleeper: SleeperFunc(func(_ context.Context, d time.Duration) error {
			slept = d
			return nil
		}),
	}
	var calls int
	err := guard.Do(t.Context(), func(context.Context) error {
		calls++
		if calls == 1 {
			return tgerr.New(420, "FLOOD_WAIT_3600")
		}
		return nil
	})
	require.NoError(t, err)
	assert.Equal(t, time.Hour, slept, "MaxFloodWait=0 should allow any duration")
}
