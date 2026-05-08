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
	Name                 string            `yaml:"name"`
	CLI                  string            `yaml:"cli"`
	BaseURL              string            `yaml:"base_url"`
	APIKey               string            `yaml:"api_key"`
	AuthToken            string            `yaml:"auth_token"`
	Model                string            `yaml:"model"`
	ModelReasoningEffort string            `yaml:"model_reasoning_effort"`
	Env                  map[string]string `yaml:"env"`
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
		return nil, fmt.Errorf(tr(msgYAMLParseFailed), err)
	}
	if cfg.Profiles == nil {
		cfg.Profiles = make(map[string]Profile)
	}
	for name, p := range cfg.Profiles {
		if p.CLI != "" && !validCLIs[p.CLI] {
			return nil, fmt.Errorf(tr(msgInvalidCLI), name, p.CLI)
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
		return nil, fmt.Errorf(tr(msgProfileMissing), name)
	}
	return &p, nil
}

func defaultConfigTemplate() string {
	if appLanguage() == "zh" {
		return `# cl 配置文件
# 步骤：
# 1. 把 default 改成下方某个 profile 名，例如 claude 或 codex。
# 2. 填写 base_url、api_key 和需要的 model。
# 3. 保存后运行 cl 或 cl <profile> 启动。
default: ""

defaults:
  env:
    CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC: "1"

profiles:
  claude:
    name: Claude Code 中转配置
    cli: claude
    base_url: https://api.example.com
    api_key: sk-your-key-here
    model: ""
    env: {}

  codex:
    name: Codex 中转配置
    cli: codex
    base_url: https://api.example.com/v1
    api_key: sk-your-key-here
    model: ""
    model_reasoning_effort: ""
    env: {}
`
	}
	return `# cl configuration
# Steps:
# 1. Set default to one of the profile keys below, for example claude or codex.
# 2. Fill in base_url, api_key, and model if your provider requires one.
# 3. Save the file, then run cl or cl <profile>.
default: ""

defaults:
  env:
    CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC: "1"

profiles:
  claude:
    name: Claude Code via your provider
    cli: claude
    base_url: https://api.example.com
    api_key: sk-your-key-here
    model: ""
    env: {}

  codex:
    name: Codex via your provider
    cli: codex
    base_url: https://api.example.com/v1
    api_key: sk-your-key-here
    model: ""
    model_reasoning_effort: ""
    env: {}
`
}
