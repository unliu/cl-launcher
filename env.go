package main

import (
	"net"
	"net/url"
	"os"
	"strings"
)

var conflictVars = []string{
	"ANTHROPIC_AUTH_TOKEN",
	"ANTHROPIC_API_KEY",
	"ANTHROPIC_BASE_URL",
	"ANTHROPIC_MODEL",
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
			"OPENAI_API_KEY": profile.APIKey,
		}
	default:
		topLevel = map[string]string{
			"ANTHROPIC_API_KEY":    profile.APIKey,
			"ANTHROPIC_AUTH_TOKEN": profile.AuthToken,
			"ANTHROPIC_BASE_URL":   profile.BaseURL,
			"ANTHROPIC_MODEL":      profile.Model,
		}
	}
	for k, v := range topLevel {
		if v != "" {
			env[k] = v
		}
	}

	if isLoopback(profile.BaseURL) {
		addNoProxy(env, profile.BaseURL)
	}

	result := make([]string, 0, len(env))
	for k, v := range env {
		result = append(result, k+"="+v)
	}
	return result
}

func isLoopback(baseURL string) bool {
	u, err := url.Parse(baseURL)
	if err != nil {
		return false
	}
	host := u.Hostname()
	if host == "localhost" {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

func addNoProxy(env map[string]string, baseURL string) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return
	}
	host := u.Hostname()

	targets := []string{host}
	if host == "localhost" {
		targets = append(targets, "127.0.0.1", "::1")
	} else if host == "127.0.0.1" || host == "::1" {
		targets = append(targets, "localhost")
	}

	for _, key := range []string{"no_proxy", "NO_PROXY"} {
		existing := env[key]
		var missing []string
		for _, t := range targets {
			if !containsNoProxyEntry(existing, t) {
				missing = append(missing, t)
			}
		}
		if len(missing) > 0 {
			if existing != "" {
				env[key] = existing + "," + strings.Join(missing, ",")
			} else {
				env[key] = strings.Join(missing, ",")
			}
		}
	}
}

func containsNoProxyEntry(noProxy, entry string) bool {
	for _, e := range strings.Split(noProxy, ",") {
		if strings.TrimSpace(e) == entry {
			return true
		}
	}
	return false
}
