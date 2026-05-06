package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadAndSaveConfig(t *testing.T) {
	dir := t.TempDir()
	origHome := os.Getenv("HOME")
	t.Setenv("HOME", dir)
	defer os.Setenv("HOME", origHome)

	os.MkdirAll(filepath.Join(dir, ".cl"), 0700)

	yamlContent := `default: pk
defaults:
  env:
    CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC: "1"
profiles:
  pk:
    name: Packy
    base_url: https://www.packyapi.com
    api_key: sk-test
    model: ""
    env: {}
  hm:
    name: 红马CC
    base_url: https://hongmacc.com/api
    auth_token: sk-hongmacc-test
    model: ""
    env:
      API_TIMEOUT_MS: "600000"
`
	os.WriteFile(filepath.Join(dir, ".cl", "profiles.yaml"), []byte(yamlContent), 0600)

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.Default != "pk" {
		t.Errorf("expected default=pk, got %q", cfg.Default)
	}
	if len(cfg.Profiles) != 2 {
		t.Errorf("expected 2 profiles, got %d", len(cfg.Profiles))
	}

	pk, err := cfg.GetProfile("pk")
	if err != nil {
		t.Fatalf("GetProfile(pk) failed: %v", err)
	}
	if pk.APIKey != "sk-test" {
		t.Errorf("expected api_key=sk-test, got %q", pk.APIKey)
	}
	if pk.BaseURL != "https://www.packyapi.com" {
		t.Errorf("expected base_url, got %q", pk.BaseURL)
	}

	hm, err := cfg.GetProfile("hm")
	if err != nil {
		t.Fatalf("GetProfile(hm) failed: %v", err)
	}
	if hm.AuthToken != "sk-hongmacc-test" {
		t.Errorf("expected auth_token, got %q", hm.AuthToken)
	}
	if hm.Env["API_TIMEOUT_MS"] != "600000" {
		t.Errorf("expected API_TIMEOUT_MS=600000, got %q", hm.Env["API_TIMEOUT_MS"])
	}

	_, err = cfg.GetProfile("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent profile")
	}

	cfg.Default = "hm"
	if err := SaveConfig(cfg); err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}
	cfg2, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig after save failed: %v", err)
	}
	if cfg2.Default != "hm" {
		t.Errorf("expected default=hm after save, got %q", cfg2.Default)
	}
}

func TestConfigFilePermissions(t *testing.T) {
	dir := t.TempDir()
	origHome := os.Getenv("HOME")
	t.Setenv("HOME", dir)
	defer os.Setenv("HOME", origHome)

	cfg := &Config{
		Default:  "test",
		Profiles: map[string]Profile{"test": {Name: "Test"}},
	}
	if err := SaveConfig(cfg); err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	info, err := os.Stat(filepath.Join(dir, ".cl", "profiles.yaml"))
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}
	perm := info.Mode().Perm()
	if perm != 0600 {
		t.Errorf("expected 0600 permissions, got %o", perm)
	}
}

func TestLoadConfig_InvalidCLI(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	os.MkdirAll(filepath.Join(dir, ".cl"), 0700)

	yamlContent := `profiles:
  bad:
    name: Bad CLI
    cli: vim
    api_key: sk-test
`
	os.WriteFile(filepath.Join(dir, ".cl", "profiles.yaml"), []byte(yamlContent), 0600)

	_, err := LoadConfig()
	if err == nil {
		t.Fatal("expected error for invalid cli value")
	}
}

func TestLoadConfig_ValidCLIs(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	os.MkdirAll(filepath.Join(dir, ".cl"), 0700)

	yamlContent := `profiles:
  cc:
    name: Claude
    cli: claude
    api_key: sk-test
  cx:
    name: Codex
    cli: codex
    api_key: sk-test
  empty:
    name: Default CLI
    api_key: sk-test
`
	os.WriteFile(filepath.Join(dir, ".cl", "profiles.yaml"), []byte(yamlContent), 0600)

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("expected valid config, got error: %v", err)
	}
	if cfg.Profiles["cc"].CLI != "claude" {
		t.Error("expected cli=claude")
	}
	if cfg.Profiles["cx"].CLI != "codex" {
		t.Error("expected cli=codex")
	}
	if cfg.Profiles["empty"].CLI != "" {
		t.Error("expected empty cli for default")
	}
}

func TestProfile_GetCLI(t *testing.T) {
	p := Profile{}
	if p.GetCLI() != "claude" {
		t.Errorf("empty cli should default to claude, got %q", p.GetCLI())
	}
	p.CLI = "codex"
	if p.GetCLI() != "codex" {
		t.Errorf("expected codex, got %q", p.GetCLI())
	}
}
