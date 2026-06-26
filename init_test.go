package main

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setupTestHome(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	t.Setenv("CL_LANG", "en")
	return dir
}

func TestRunInit_NewConfig(t *testing.T) {
	setupTestHome(t)

	input := strings.Join([]string{
		"myprofile",      // profile name
		"claude",         // cli
		"sk-test-key",    // api key
		"https://my.api", // base url
		"claude-opus-4",  // model
	}, "\n") + "\n"

	var out bytes.Buffer
	runInit(strings.NewReader(input), &out)

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if cfg.Default != "myprofile" {
		t.Errorf("expected default=myprofile, got %q", cfg.Default)
	}
	p, err := cfg.GetProfile("myprofile")
	if err != nil {
		t.Fatalf("GetProfile failed: %v", err)
	}
	if p.CLI != "claude" {
		t.Errorf("expected cli=claude, got %q", p.CLI)
	}
	if p.APIKey != "sk-test-key" {
		t.Errorf("expected api_key=sk-test-key, got %q", p.APIKey)
	}
	if p.BaseURL != "https://my.api" {
		t.Errorf("expected base_url=https://my.api, got %q", p.BaseURL)
	}
	if p.Model != "claude-opus-4" {
		t.Errorf("expected model=claude-opus-4, got %q", p.Model)
	}
	if cfg.Defaults.Env["CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC"] != "1" {
		t.Error("expected CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC=1 in defaults.env")
	}
}

func TestRunInit_DefaultValues(t *testing.T) {
	setupTestHome(t)

	input := strings.Join([]string{
		"",            // profile name → "my"
		"",            // cli → "claude"
		"sk-test-key", // api key (required)
		"",            // base url → default (not saved)
		"",            // model → empty
	}, "\n") + "\n"

	var out bytes.Buffer
	runInit(strings.NewReader(input), &out)

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if cfg.Default != "my" {
		t.Errorf("expected default=my, got %q", cfg.Default)
	}
	p, err := cfg.GetProfile("my")
	if err != nil {
		t.Fatalf("GetProfile failed: %v", err)
	}
	if p.CLI != "claude" {
		t.Errorf("expected cli=claude, got %q", p.CLI)
	}
	if p.BaseURL != "" {
		t.Errorf("expected empty base_url when using default, got %q", p.BaseURL)
	}
}

func TestRunInit_CodexDefaultBaseURL(t *testing.T) {
	setupTestHome(t)

	input := strings.Join([]string{
		"codex-profile", // profile name
		"codex",         // cli
		"sk-test-key",   // api key
		"",              // base url → codex default (not saved)
		"",              // model
	}, "\n") + "\n"

	var out bytes.Buffer
	runInit(strings.NewReader(input), &out)

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	p, err := cfg.GetProfile("codex-profile")
	if err != nil {
		t.Fatalf("GetProfile failed: %v", err)
	}
	if p.CLI != "codex" {
		t.Errorf("expected cli=codex, got %q", p.CLI)
	}
	if p.BaseURL != "" {
		t.Errorf("expected empty base_url when using codex default, got %q", p.BaseURL)
	}

	output := out.String()
	if !strings.Contains(output, codexDefaultBaseURL) {
		t.Errorf("expected prompt to show codex default URL %q, output: %s", codexDefaultBaseURL, output)
	}
}

func TestRunInit_AddToExisting(t *testing.T) {
	dir := setupTestHome(t)

	os.MkdirAll(filepath.Join(dir, ".cl"), 0700)
	os.WriteFile(filepath.Join(dir, ".cl", "profiles.yaml"), []byte(`default: existing
profiles:
  existing:
    cli: claude
    api_key: sk-old
`), 0600)

	input := strings.Join([]string{
		"newprofile",  // profile name
		"codex",       // cli
		"sk-new-key",  // api key
		"https://new", // base url
		"gpt-5",       // model
		"n",           // don't set as default
	}, "\n") + "\n"

	var out bytes.Buffer
	runInit(strings.NewReader(input), &out)

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if cfg.Default != "existing" {
		t.Errorf("expected default unchanged, got %q", cfg.Default)
	}
	if len(cfg.Profiles) != 2 {
		t.Errorf("expected 2 profiles, got %d", len(cfg.Profiles))
	}
	p, err := cfg.GetProfile("newprofile")
	if err != nil {
		t.Fatalf("GetProfile failed: %v", err)
	}
	if p.APIKey != "sk-new-key" {
		t.Errorf("expected api_key=sk-new-key, got %q", p.APIKey)
	}
}

func TestRunInit_AddToExisting_SetDefault(t *testing.T) {
	dir := setupTestHome(t)

	os.MkdirAll(filepath.Join(dir, ".cl"), 0700)
	os.WriteFile(filepath.Join(dir, ".cl", "profiles.yaml"), []byte(`default: old
profiles:
  old:
    cli: claude
    api_key: sk-old
`), 0600)

	input := strings.Join([]string{
		"newdefault",  // profile name
		"claude",      // cli
		"sk-new-key",  // api key
		"https://new", // base url
		"",            // model
		"y",           // set as default
	}, "\n") + "\n"

	var out bytes.Buffer
	runInit(strings.NewReader(input), &out)

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if cfg.Default != "newdefault" {
		t.Errorf("expected default=newdefault, got %q", cfg.Default)
	}
}

func TestRunInit_DuplicateProfileRetry(t *testing.T) {
	dir := setupTestHome(t)

	os.MkdirAll(filepath.Join(dir, ".cl"), 0700)
	os.WriteFile(filepath.Join(dir, ".cl", "profiles.yaml"), []byte(`default: taken
profiles:
  taken:
    cli: claude
    api_key: sk-old
`), 0600)

	input := strings.Join([]string{
		"taken",       // duplicate → rejected
		"fresh",       // valid name
		"claude",      // cli
		"sk-test-key", // api key
		"",            // base url
		"",            // model
		"n",           // don't set as default
	}, "\n") + "\n"

	var out bytes.Buffer
	runInit(strings.NewReader(input), &out)

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if _, err := cfg.GetProfile("fresh"); err != nil {
		t.Error("expected profile 'fresh' to be created")
	}
	output := out.String()
	if !strings.Contains(output, "already exists") {
		t.Error("expected duplicate warning in output")
	}
}

func TestRunInit_ReservedWordRetry(t *testing.T) {
	setupTestHome(t)

	input := strings.Join([]string{
		"list",        // reserved → rejected
		"myprofile",   // valid name
		"claude",      // cli
		"sk-test-key", // api key
		"",            // base url
		"",            // model
	}, "\n") + "\n"

	var out bytes.Buffer
	runInit(strings.NewReader(input), &out)

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if _, err := cfg.GetProfile("myprofile"); err != nil {
		t.Error("expected profile 'myprofile' to be created")
	}
	output := out.String()
	if !strings.Contains(output, "reserved") {
		t.Error("expected reserved word warning in output")
	}
}

func TestRunInit_InvalidCLIRetry(t *testing.T) {
	setupTestHome(t)

	input := strings.Join([]string{
		"myprofile",   // profile name
		"vim",         // invalid cli → rejected
		"claude",      // valid cli
		"sk-test-key", // api key
		"",            // base url
		"",            // model
	}, "\n") + "\n"

	var out bytes.Buffer
	runInit(strings.NewReader(input), &out)

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	p, err := cfg.GetProfile("myprofile")
	if err != nil {
		t.Fatalf("GetProfile failed: %v", err)
	}
	if p.CLI != "claude" {
		t.Errorf("expected cli=claude, got %q", p.CLI)
	}
}

func TestRunInit_EmptyAPIKeyRetry(t *testing.T) {
	setupTestHome(t)

	input := strings.Join([]string{
		"myprofile",   // profile name
		"claude",      // cli
		"",            // empty api key → rejected
		"sk-test-key", // valid api key
		"",            // base url
		"",            // model
	}, "\n") + "\n"

	var out bytes.Buffer
	runInit(strings.NewReader(input), &out)

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	p, err := cfg.GetProfile("myprofile")
	if err != nil {
		t.Fatalf("GetProfile failed: %v", err)
	}
	if p.APIKey != "sk-test-key" {
		t.Errorf("expected api_key=sk-test-key, got %q", p.APIKey)
	}
}

func TestPromptWithDefault(t *testing.T) {
	var out bytes.Buffer
	scanner := bufio.NewScanner(strings.NewReader("custom\n"))
	val := promptWithDefault(scanner, &out, "Field", "fallback")
	if val != "custom" {
		t.Errorf("expected custom, got %q", val)
	}

	out.Reset()
	scanner = bufio.NewScanner(strings.NewReader("\n"))
	val = promptWithDefault(scanner, &out, "Field", "fallback")
	if val != "fallback" {
		t.Errorf("expected fallback, got %q", val)
	}

	output := out.String()
	if !strings.Contains(output, "[fallback]") {
		t.Errorf("expected default shown in brackets, got %q", output)
	}

	out.Reset()
	scanner = bufio.NewScanner(strings.NewReader("\n"))
	val = promptWithDefault(scanner, &out, "Field", "")
	if val != "" {
		t.Errorf("expected empty string, got %q", val)
	}
	output = out.String()
	if strings.Contains(output, "[") {
		t.Errorf("expected no brackets for empty default, got %q", output)
	}
}

func TestPromptRequired(t *testing.T) {
	var out bytes.Buffer
	scanner := bufio.NewScanner(strings.NewReader("\nvalue\n"))
	val := promptRequired(scanner, &out, "Required")
	if val != "value" {
		t.Errorf("expected value, got %q", val)
	}
	output := out.String()
	if !strings.Contains(output, "required") {
		t.Error("expected required warning on empty input")
	}
}

func TestPromptChoice(t *testing.T) {
	var out bytes.Buffer
	scanner := bufio.NewScanner(strings.NewReader("bad\nclaude\n"))
	val := promptChoice(scanner, &out, "CLI", []string{"claude", "codex"}, "claude")
	if val != "claude" {
		t.Errorf("expected claude, got %q", val)
	}
	output := out.String()
	if !strings.Contains(output, "claude, codex") {
		t.Error("expected valid choices listed on invalid input")
	}

	out.Reset()
	scanner = bufio.NewScanner(strings.NewReader("\n"))
	val = promptChoice(scanner, &out, "CLI", []string{"claude", "codex"}, "codex")
	if val != "codex" {
		t.Errorf("expected default codex, got %q", val)
	}
}
