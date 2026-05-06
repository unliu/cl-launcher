package main

import (
	"os"
	"strings"
)

var conflictVars = []string{
	"ANTHROPIC_AUTH_TOKEN",
	"ANTHROPIC_API_KEY",
	"ANTHROPIC_BASE_URL",
	"ANTHROPIC_MODEL",
	"ANTHROPIC_SMALL_FAST_MODEL",
	"OPENAI_API_KEY",
	"OPENAI_BASE_URL",
	"OPENAI_MODEL",
}

func BuildEnv(cfg *Config, profile *Profile) []string {
	env := make(map[string]string)
	for _, e := range os.Environ() {
		k, v, _ := strings.Cut(e, "=")
		env[k] = v
	}

	for _, k := range conflictVars {
		delete(env, k)
	}

	if cfg.Defaults.Env != nil {
		for k, v := range cfg.Defaults.Env {
			env[k] = v
		}
	}

	if profile.Env != nil {
		for k, v := range profile.Env {
			env[k] = v
		}
	}

	var topLevel map[string]string
	switch profile.GetCLI() {
	case "codex":
		topLevel = map[string]string{
			"OPENAI_API_KEY":  profile.APIKey,
			"OPENAI_BASE_URL": profile.BaseURL,
			"OPENAI_MODEL":    profile.Model,
		}
	default:
		topLevel = map[string]string{
			"ANTHROPIC_API_KEY":          profile.APIKey,
			"ANTHROPIC_AUTH_TOKEN":       profile.AuthToken,
			"ANTHROPIC_BASE_URL":         profile.BaseURL,
			"ANTHROPIC_MODEL":            profile.Model,
			"ANTHROPIC_SMALL_FAST_MODEL": profile.SmallFastModel,
		}
	}
	for k, v := range topLevel {
		if v != "" {
			env[k] = v
		}
	}

	result := make([]string, 0, len(env))
	for k, v := range env {
		result = append(result, k+"="+v)
	}
	return result
}
