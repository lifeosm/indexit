package log

import (
	"bytes"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetup_DefaultLevel(t *testing.T) {
	var buf bytes.Buffer
	s := Setup(Options{Out: &buf})

	s.Logger.Debug("debug-line")
	s.Logger.Info("info-line")
	s.Logger.Warn("warn-line")
	s.Logger.Error("err-line")

	out := buf.String()
	assert.NotContains(t, out, "debug-line", "debug should be hidden at default level")
	assert.Contains(t, out, "info-line")
	assert.Contains(t, out, "warn: warn-line")
	assert.Contains(t, out, "error: err-line")
}

func TestSetup_Verbose(t *testing.T) {
	var buf bytes.Buffer
	s := Setup(Options{Out: &buf, Verbose: 1})
	s.Logger.Debug("debug-line")
	s.Logger.Info("info-line")
	assert.Contains(t, buf.String(), "debug: debug-line")
	assert.Contains(t, buf.String(), "info-line")
}

func TestSetup_Quiet(t *testing.T) {
	var buf bytes.Buffer
	s := Setup(Options{Out: &buf, Quiet: true})
	s.Logger.Info("info-line")
	s.Logger.Warn("warn-line")
	s.Logger.Error("err-line")
	out := buf.String()
	assert.NotContains(t, out, "info-line")
	assert.NotContains(t, out, "warn-line")
	assert.Contains(t, out, "error: err-line")
}

func TestSetup_AttrsFormatting(t *testing.T) {
	var buf bytes.Buffer
	s := Setup(Options{Out: &buf})
	s.Logger.Info("connecting", "dc", 2, "host", "10.0.0.1")
	out := buf.String()
	assert.Contains(t, out, "connecting")
	assert.Contains(t, out, "dc=2")
	assert.Contains(t, out, "host=10.0.0.1")
}

func TestSetup_AttrQuoting(t *testing.T) {
	var buf bytes.Buffer
	s := Setup(Options{Out: &buf})
	s.Logger.Info("note", "path", "/tmp/has spaces/.env")
	assert.Contains(t, buf.String(), `path="/tmp/has spaces/.env"`)
}

func TestContext(t *testing.T) {
	ctx := WithSettings(t.Context(), &Settings{Logger: slog.Default()})
	got := FromContext(ctx)
	assert.NotNil(t, got)
	assert.Same(t, slog.Default(), got.Logger)
}

func TestFromContext_DefaultIsSafe(t *testing.T) {
	got := FromContext(t.Context())
	assert.NotNil(t, got)
	assert.NotNil(t, got.Logger)
	assert.NotNil(t, got.GotdLog)
}
