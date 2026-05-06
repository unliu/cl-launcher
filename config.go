package main

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Default  string             `yaml:"default"`
	Defaults Defaults           `yaml:"defaults"`
	Profiles map[string]Profile `yaml:"profiles"`
}

type Defaults struct {
	Env map[string]string `yaml:"env"`
}

var validCLIs = map[string]bool{
	"claude": true,
	"codex":  true,
}

type Profile struct {
	Name           string            `yaml:"name"`
	CLI            string            `yaml:"cli"`
	BaseURL        string            `yaml:"base_url"`
	APIKey         string            `yaml:"api_key"`
	AuthToken      string            `yaml:"auth_token"`
	Model          string            `yaml:"model"`
	SmallFastModel string            `yaml:"small_fast_model"`
	Env            map[string]string `yaml:"env"`
}

func (p *Profile) GetCLI() string {
	if p.CLI == "" {
		return "claude"
	}
	return p.CLI
}

func ConfigDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".cl")
}

func ConfigPath() string {
	return filepath.Join(ConfigDir(), "profiles.yaml")
}

func LoadConfig() (*Config, error) {
	data, err := os.ReadFile(ConfigPath())
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("profiles.yaml 解析失败: %w", err)
	}
	if cfg.Profiles == nil {
		cfg.Profiles = make(map[string]Profile)
	}
	for name, p := range cfg.Profiles {
		if p.CLI != "" && !validCLIs[p.CLI] {
			return nil, fmt.Errorf("profile %q: cli %q 不合法，仅支持 claude、codex", name, p.CLI)
		}
	}
	return &cfg, nil
}

func SaveConfig(cfg *Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	dir := ConfigDir()
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	return os.WriteFile(ConfigPath(), data, 0600)
}

func (c *Config) GetProfile(name string) (*Profile, error) {
	p, ok := c.Profiles[name]
	if !ok {
		return nil, fmt.Errorf("profile %q 不存在", name)
	}
	return &p, nil
}

const configTemplate = `default: ""

defaults:
  env:
    CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC: "1"

profiles:
  example:
    name: 示例 Provider
    cli: claude  # claude 或 codex
    base_url: https://api.example.com
    api_key: sk-your-key-here
    model: ""
    small_fast_model: ""
    env: {}
`
