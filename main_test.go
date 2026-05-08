package main

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestPrintProfileName(t *testing.T) {
	tests := []struct {
		name    string
		profile *Profile
		want    string
	}{
		{
			name:    "prints configured name",
			profile: &Profile{Name: "Codex via Relay"},
			want:    "Codex via Relay\n",
		},
		{
			name:    "skips empty name",
			profile: &Profile{},
			want:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := captureStderr(t, func() {
				printProfileName(tt.profile)
			})
			if got != tt.want {
				t.Fatalf("unexpected stderr:\n got: %q\nwant: %q", got, tt.want)
			}
		})
	}
}

func TestLinkURLs(t *testing.T) {
	tests := []struct {
		name string
		text string
		want string
	}{
		{
			name: "leaves text without URLs unchanged",
			text: "Codex via Relay",
			want: "Codex via Relay",
		},
		{
			name: "links URL inside profile name",
			text: "Relay https://relay.example.com",
			want: "Relay \x1b]8;;https://relay.example.com\ahttps://relay.example.com\x1b]8;;\a",
		},
		{
			name: "keeps trailing punctuation outside link",
			text: "Relay (https://relay.example.com).",
			want: "Relay (\x1b]8;;https://relay.example.com\ahttps://relay.example.com\x1b]8;;\a).",
		},
		{
			name: "links multiple URLs",
			text: "A https://a.example B http://b.example",
			want: "A \x1b]8;;https://a.example\ahttps://a.example\x1b]8;;\a B \x1b]8;;http://b.example\ahttp://b.example\x1b]8;;\a",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := linkURLs(tt.text)
			if got != tt.want {
				t.Fatalf("unexpected linked text:\n got: %q\nwant: %q", got, tt.want)
			}
		})
	}
}

func captureStderr(t *testing.T, fn func()) string {
	t.Helper()

	orig := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe failed: %v", err)
	}
	os.Stderr = w

	fn()

	if err := w.Close(); err != nil {
		t.Fatalf("close stderr pipe writer failed: %v", err)
	}
	os.Stderr = orig

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("read stderr pipe failed: %v", err)
	}
	return buf.String()
}
