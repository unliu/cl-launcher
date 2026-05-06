package main

import (
	"os"
	"sort"
	"strings"
	"testing"
)

func TestBuildEnv_Priority(t *testing.T) {
	os.Unsetenv("ANTHROPIC_API_KEY")
	os.Unsetenv("ANTHROPIC_BASE_URL")
	os.Unsetenv("ANTHROPIC_MODEL")
	os.Unsetenv("ANTHROPIC_AUTH_TOKEN")
	os.Unsetenv("OPENAI_API_KEY")
	os.Unsetenv("OPENAI_BASE_URL")
	os.Unsetenv("OPENAI_MODEL")

	cfg := &Config{
		Defaults: Defaults{
			Env: map[string]string{
				"CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC": "1",
				"SHARED_VAR": "from-defaults",
			},
		},
	}

	profile := &Profile{
		BaseURL: "https://api.example.com",
		APIKey:  "sk-top-level",
		Model:   "claude-sonnet-4-20250514",
		Env: map[string]string{
			"SHARED_VAR":        "from-profile",
			"ANTHROPIC_API_KEY": "sk-from-env",
		},
	}

	env := BuildEnv(cfg, profile)
	envMap := envToMap(env)

	// Top-level api_key should override profile.env
	if envMap["ANTHROPIC_API_KEY"] != "sk-top-level" {
		t.Errorf("expected top-level api_key to win, got %q", envMap["ANTHROPIC_API_KEY"])
	}

	// Profile env should override defaults env
	if envMap["SHARED_VAR"] != "from-profile" {
		t.Errorf("expected profile env to override defaults, got %q", envMap["SHARED_VAR"])
	}

	if envMap["ANTHROPIC_BASE_URL"] != "https://api.example.com" {
		t.Errorf("expected base_url, got %q", envMap["ANTHROPIC_BASE_URL"])
	}

	if envMap["ANTHROPIC_MODEL"] != "claude-sonnet-4-20250514" {
		t.Errorf("expected model, got %q", envMap["ANTHROPIC_MODEL"])
	}

	if envMap["CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC"] != "1" {
		t.Errorf("expected defaults env inherited")
	}
}

func TestBuildEnv_EmptyValuesSkipped(t *testing.T) {
	os.Unsetenv("ANTHROPIC_API_KEY")
	os.Unsetenv("ANTHROPIC_BASE_URL")
	os.Unsetenv("ANTHROPIC_MODEL")
	os.Unsetenv("ANTHROPIC_AUTH_TOKEN")
	os.Unsetenv("OPENAI_API_KEY")
	os.Unsetenv("OPENAI_BASE_URL")
	os.Unsetenv("OPENAI_MODEL")

	cfg := &Config{}
	profile := &Profile{
		BaseURL: "https://api.example.com",
		APIKey:  "sk-key",
		Model:   "",
	}

	env := BuildEnv(cfg, profile)
	envMap := envToMap(env)

	if _, ok := envMap["ANTHROPIC_MODEL"]; ok {
		t.Error("empty model should not be injected")
	}
	if _, ok := envMap["ANTHROPIC_AUTH_TOKEN"]; ok {
		t.Error("empty auth_token should not be injected")
	}
}

func TestBuildEnv_ConflictVarsCleared(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "old-key")
	os.Setenv("ANTHROPIC_AUTH_TOKEN", "old-token")
	os.Setenv("OPENAI_API_KEY", "old-openai-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")
	defer os.Unsetenv("ANTHROPIC_AUTH_TOKEN")
	defer os.Unsetenv("OPENAI_API_KEY")

	cfg := &Config{}
	profile := &Profile{
		BaseURL: "https://api.example.com",
		APIKey:  "new-key",
	}

	env := BuildEnv(cfg, profile)
	envMap := envToMap(env)

	if envMap["ANTHROPIC_API_KEY"] != "new-key" {
		t.Errorf("expected new-key, got %q", envMap["ANTHROPIC_API_KEY"])
	}
	if _, ok := envMap["ANTHROPIC_AUTH_TOKEN"]; ok {
		t.Error("old auth_token should be cleared")
	}
}

func TestBuildEnv_AuthToken(t *testing.T) {
	os.Unsetenv("ANTHROPIC_API_KEY")
	os.Unsetenv("ANTHROPIC_AUTH_TOKEN")

	cfg := &Config{}
	profile := &Profile{
		BaseURL:   "https://hongmacc.com/api",
		AuthToken: "sk-hongmacc-test",
	}

	env := BuildEnv(cfg, profile)
	envMap := envToMap(env)

	if envMap["ANTHROPIC_AUTH_TOKEN"] != "sk-hongmacc-test" {
		t.Errorf("expected auth_token, got %q", envMap["ANTHROPIC_AUTH_TOKEN"])
	}
	if _, ok := envMap["ANTHROPIC_API_KEY"]; ok {
		t.Error("api_key should not be set when only auth_token is configured")
	}
}

func envToMap(env []string) map[string]string {
	m := make(map[string]string)
	for _, e := range env {
		k, v, _ := strings.Cut(e, "=")
		m[k] = v
	}
	return m
}

func sortedEnv(env []string) []string {
	sorted := make([]string, len(env))
	copy(sorted, env)
	sort.Strings(sorted)
	return sorted
}

func TestBuildEnv_Codex(t *testing.T) {
	for _, k := range []string{
		"ANTHROPIC_API_KEY", "ANTHROPIC_BASE_URL", "ANTHROPIC_MODEL",
		"ANTHROPIC_AUTH_TOKEN", "OPENAI_API_KEY", "OPENAI_BASE_URL", "OPENAI_MODEL",
	} {
		os.Unsetenv(k)
	}

	cfg := &Config{
		Defaults: Defaults{
			Env: map[string]string{
				"SHARED_VAR": "from-defaults",
			},
		},
	}
	profile := &Profile{
		CLI:     "codex",
		BaseURL: "https://api.aicoding.sh",
		APIKey:  "aicoding-xxx",
		Model:   "o3",
		Env: map[string]string{
			"CODEX_CONFIG_DIR": "~/.codex-envs/mirror",
		},
	}

	env := BuildEnv(cfg, profile)
	envMap := envToMap(env)

	if envMap["OPENAI_API_KEY"] != "aicoding-xxx" {
		t.Errorf("expected OPENAI_API_KEY=aicoding-xxx, got %q", envMap["OPENAI_API_KEY"])
	}
	if envMap["OPENAI_BASE_URL"] != "https://api.aicoding.sh" {
		t.Errorf("expected OPENAI_BASE_URL, got %q", envMap["OPENAI_BASE_URL"])
	}
	if envMap["OPENAI_MODEL"] != "o3" {
		t.Errorf("expected OPENAI_MODEL=o3, got %q", envMap["OPENAI_MODEL"])
	}
	if envMap["CODEX_CONFIG_DIR"] != "~/.codex-envs/mirror" {
		t.Errorf("expected CODEX_CONFIG_DIR, got %q", envMap["CODEX_CONFIG_DIR"])
	}
	if envMap["SHARED_VAR"] != "from-defaults" {
		t.Errorf("expected SHARED_VAR from defaults, got %q", envMap["SHARED_VAR"])
	}
	if _, ok := envMap["ANTHROPIC_API_KEY"]; ok {
		t.Error("codex profile should not set ANTHROPIC_API_KEY")
	}
	if _, ok := envMap["ANTHROPIC_BASE_URL"]; ok {
		t.Error("codex profile should not set ANTHROPIC_BASE_URL")
	}
}

func TestBuildEnv_CodexClearsAnthropicVars(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "old-anthropic")
	os.Setenv("OPENAI_API_KEY", "old-openai")
	defer os.Unsetenv("ANTHROPIC_API_KEY")
	defer os.Unsetenv("OPENAI_API_KEY")

	cfg := &Config{}
	profile := &Profile{
		CLI:    "codex",
		APIKey: "new-codex-key",
	}

	env := BuildEnv(cfg, profile)
	envMap := envToMap(env)

	if envMap["OPENAI_API_KEY"] != "new-codex-key" {
		t.Errorf("expected new-codex-key, got %q", envMap["OPENAI_API_KEY"])
	}
	if _, ok := envMap["ANTHROPIC_API_KEY"]; ok {
		t.Error("ANTHROPIC_API_KEY should be cleared for codex profile")
	}
}

func TestBuildEnv_DefaultCLIIsClaude(t *testing.T) {
	os.Unsetenv("ANTHROPIC_API_KEY")
	os.Unsetenv("OPENAI_API_KEY")

	cfg := &Config{}
	profile := &Profile{
		APIKey: "sk-test",
	}

	env := BuildEnv(cfg, profile)
	envMap := envToMap(env)

	if envMap["ANTHROPIC_API_KEY"] != "sk-test" {
		t.Errorf("empty cli should default to claude, expected ANTHROPIC_API_KEY=sk-test, got %q", envMap["ANTHROPIC_API_KEY"])
	}
	if _, ok := envMap["OPENAI_API_KEY"]; ok {
		t.Error("empty cli should not set OPENAI_API_KEY")
	}
}
