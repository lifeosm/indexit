package output

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriterStdout(t *testing.T) {
	var out bytes.Buffer
	w, err := New("-", &out)
	require.NoError(t, err)
	require.NoError(t, w.Write(map[string]any{"_kind": "test", "id": 1}))
	require.NoError(t, w.Close())
	assert.JSONEq(t, `{"_kind":"test","id":1}`, out.String())
	assert.Equal(t, "\n", out.String()[len(out.String())-1:])
}
