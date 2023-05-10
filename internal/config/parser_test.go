package config

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want map[string]string
	}{
		{
			name: "basic",
			in:   "FOO=bar\nBAZ=qux",
			want: map[string]string{"FOO": "bar", "BAZ": "qux"},
		},
		{
			name: "comments and blanks",
			in:   "# header\n\nFOO=bar\n   # indented comment\nBAZ=qux\n",
			want: map[string]string{"FOO": "bar", "BAZ": "qux"},
		},
		{
			name: "export prefix",
			in:   "export FOO=bar",
			want: map[string]string{"FOO": "bar"},
		},
		{
			name: "inline comment unquoted",
			in:   "FOO=bar # trailing",
			want: map[string]string{"FOO": "bar"},
		},
		{
			name: "double quoted with escapes",
			in:   `FOO="line1\nline2"`,
			want: map[string]string{"FOO": "line1\nline2"},
		},
		{
			name: "single quoted literal",
			in:   `FOO='no\nescape'`,
			want: map[string]string{"FOO": `no\nescape`},
		},
		{
			name: "value with equals",
			in:   "FOO=a=b=c",
			want: map[string]string{"FOO": "a=b=c"},
		},
		{
			name: "empty value",
			in:   "FOO=",
			want: map[string]string{"FOO": ""},
		},
		{
			name: "url-like value",
			in:   "URL=https://example.com/path?x=1",
			want: map[string]string{"URL": "https://example.com/path?x=1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parse(strings.NewReader(tt.in))
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseErrors(t *testing.T) {
	tests := []struct {
		name string
		in   string
	}{
		{"no equals", "FOO\n"},
		{"empty key", "=value\n"},
		{"invalid key", "FOO-BAR=baz\n"},
		{"leading digit", "1FOO=bar\n"},
		{"unterminated double quote", `FOO="bar`},
		{"unterminated single quote", `FOO='bar`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parse(strings.NewReader(tt.in))
			assert.Error(t, err)
		})
	}
}
