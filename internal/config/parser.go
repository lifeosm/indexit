package config

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func parseFile(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return parse(f)
}

func parse(r io.Reader) (map[string]string, error) {
	out := map[string]string{}
	scanner := bufio.NewScanner(r)
	lineno := 0
	for scanner.Scan() {
		lineno++
		raw := scanner.Text()
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if rest, ok := strings.CutPrefix(line, "export "); ok {
			line = strings.TrimSpace(rest)
		}
		eq := strings.IndexByte(line, '=')
		if eq <= 0 {
			return nil, fmt.Errorf("line %d: expected KEY=VALUE", lineno)
		}
		key := strings.TrimSpace(line[:eq])
		val := strings.TrimSpace(line[eq+1:])
		if !validKey(key) {
			return nil, fmt.Errorf("line %d: invalid key %q", lineno, key)
		}
		val, err := unquote(val)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", lineno, err)
		}
		out[key] = val
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func validKey(s string) bool {
	if s == "" {
		return false
	}
	for i, r := range s {
		switch {
		case r >= 'A' && r <= 'Z':
		case r >= 'a' && r <= 'z':
		case r == '_':
		case i > 0 && r >= '0' && r <= '9':
		default:
			return false
		}
	}
	return true
}

func unquote(s string) (string, error) {
	if s == "" {
		return s, nil
	}
	if s[0] != '"' && s[0] != '\'' {
		if idx := strings.Index(s, " #"); idx >= 0 {
			s = strings.TrimSpace(s[:idx])
		}
		return s, nil
	}
	q := s[0]
	if len(s) < 2 || s[len(s)-1] != q {
		return "", fmt.Errorf("unterminated %c quote", q)
	}
	inner := s[1 : len(s)-1]
	if q == '\'' {
		return inner, nil
	}
	var b strings.Builder
	b.Grow(len(inner))
	for i := 0; i < len(inner); i++ {
		c := inner[i]
		if c != '\\' || i+1 >= len(inner) {
			b.WriteByte(c)
			continue
		}
		i++
		switch inner[i] {
		case 'n':
			b.WriteByte('\n')
		case 't':
			b.WriteByte('\t')
		case 'r':
			b.WriteByte('\r')
		case '\\':
			b.WriteByte('\\')
		case '"':
			b.WriteByte('"')
		default:
			b.WriteByte('\\')
			b.WriteByte(inner[i])
		}
	}
	return b.String(), nil
}
