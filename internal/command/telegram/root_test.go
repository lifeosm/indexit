package telegram

import (
	"bytes"
	"errors"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.octolab.org/toolset/indexit/internal/exitcode"
)

func TestDebugUID(t *testing.T) {
	var out bytes.Buffer
	command := New()
	command.SetOut(&out)
	command.SetErr(&bytes.Buffer{})
	command.SetArgs([]string{"debug", "uid", "@telegram"})

	require.NoError(t, command.Execute())
	assert.Contains(t, out.String(), `"Username": "telegram"`)
}

func TestAuthStatusWithoutSession(t *testing.T) {
	var out bytes.Buffer
	command := New()
	command.SetOut(&out)
	command.SetErr(&bytes.Buffer{})
	command.SetArgs([]string{
		"--session", filepath.Join(t.TempDir(), "session.json"),
		"auth", "status",
	})

	require.NoError(t, command.Execute())
	assert.Contains(t, out.String(), "authorized: false")
}

func TestFetchMessagesBadUIDIsUsageError(t *testing.T) {
	command := New()
	command.SetOut(&bytes.Buffer{})
	command.SetErr(&bytes.Buffer{})
	command.SetArgs([]string{"fetch", "messages", "--dialog", "123"})

	err := command.Execute()
	require.Error(t, err)

	var usage *exitcode.Error
	require.True(t, errors.As(err, &usage))
	assert.Equal(t, exitcode.Usage, usage.Code)
}
