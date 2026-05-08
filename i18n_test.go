package main

import (
	"strings"
	"testing"
)

func TestAppLanguage(t *testing.T) {
	tests := []struct {
		name string
		env  map[string]string
		want string
	}{
		{
			name: "defaults to english",
			env:  map[string]string{},
			want: "en",
		},
		{
			name: "detects chinese lang",
			env:  map[string]string{"LANG": "zh_CN.UTF-8"},
			want: "zh",
		},
		{
			name: "cl lang overrides locale",
			env:  map[string]string{"CL_LANG": "en", "LANG": "zh_CN.UTF-8"},
			want: "en",
		},
		{
			name: "non chinese locale uses english",
			env:  map[string]string{"LANG": "fr_FR.UTF-8"},
			want: "en",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("CL_LANG", "")
			t.Setenv("LC_ALL", "")
			t.Setenv("LC_MESSAGES", "")
			t.Setenv("LANG", "")
			t.Setenv("LANGUAGE", "")
			for key, value := range tt.env {
				t.Setenv(key, value)
			}
			if got := appLanguage(); got != tt.want {
				t.Fatalf("appLanguage() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDefaultConfigTemplateLocalizes(t *testing.T) {
	t.Setenv("CL_LANG", "en")
	en := defaultConfigTemplate()
	if !strings.Contains(en, "# cl configuration") {
		t.Fatalf("expected english template header, got:\n%s", en)
	}
	if !strings.Contains(en, "claude:") || !strings.Contains(en, "codex:") {
		t.Fatalf("expected template to include claude and codex examples, got:\n%s", en)
	}

	t.Setenv("CL_LANG", "zh")
	zh := defaultConfigTemplate()
	if !strings.Contains(zh, "# cl 配置文件") {
		t.Fatalf("expected chinese template header, got:\n%s", zh)
	}
	if !strings.Contains(zh, "claude:") || !strings.Contains(zh, "codex:") {
		t.Fatalf("expected template to include claude and codex examples, got:\n%s", zh)
	}
}
