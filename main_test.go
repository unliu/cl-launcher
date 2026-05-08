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
