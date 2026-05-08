package main

import (
	"reflect"
	"testing"
)

func TestBuildArgs_CodexConfigOverrides(t *testing.T) {
	profile := &Profile{
		CLI:                  "codex",
		BaseURL:              "http://127.0.0.1:8317/v1",
		APIKey:               "sk-test",
		Model:                "gpt-5.5",
		ModelReasoningEffort: "xhigh",
	}

	got := BuildArgs(profile, []string{"--sandbox", "workspace-write"})
	want := []string{
		"-c", `model_provider="cl"`,
		"-c", `model_providers.cl.name="cl"`,
		"-c", `model_providers.cl.base_url="http://127.0.0.1:8317/v1"`,
		"-c", `model_providers.cl.wire_api="responses"`,
		"-c", `model_providers.cl.env_key="OPENAI_API_KEY"`,
		"-c", `model="gpt-5.5"`,
		"-c", `model_reasoning_effort="xhigh"`,
		"--sandbox", "workspace-write",
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected args:\n got: %#v\nwant: %#v", got, want)
	}
}

func TestBuildArgs_CodexAPIKeyUsesDefaultOpenAIBaseURL(t *testing.T) {
	profile := &Profile{
		CLI:    "codex",
		APIKey: "sk-test",
		Model:  "gpt-5.5",
	}

	got := BuildArgs(profile, nil)
	want := []string{
		"-c", `model_provider="cl"`,
		"-c", `model_providers.cl.name="cl"`,
		"-c", `model_providers.cl.base_url="https://api.openai.com/v1"`,
		"-c", `model_providers.cl.wire_api="responses"`,
		"-c", `model_providers.cl.env_key="OPENAI_API_KEY"`,
		"-c", `model="gpt-5.5"`,
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected args:\n got: %#v\nwant: %#v", got, want)
	}
}

func TestBuildArgs_CodexBaseURLWithoutAPIKeySkipsEnvKey(t *testing.T) {
	profile := &Profile{
		CLI:     "codex",
		BaseURL: "http://127.0.0.1:8317/v1",
	}

	got := BuildArgs(profile, nil)
	want := []string{
		"-c", `model_provider="cl"`,
		"-c", `model_providers.cl.name="cl"`,
		"-c", `model_providers.cl.base_url="http://127.0.0.1:8317/v1"`,
		"-c", `model_providers.cl.wire_api="responses"`,
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected args:\n got: %#v\nwant: %#v", got, want)
	}
}

func TestBuildArgs_CodexSkipsEmptyConfig(t *testing.T) {
	profile := &Profile{
		CLI:   "codex",
		Model: "gpt-5.5",
	}

	got := BuildArgs(profile, nil)
	want := []string{"-c", `model="gpt-5.5"`}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected args:\n got: %#v\nwant: %#v", got, want)
	}
}

func TestBuildArgs_ClaudePassThroughOnly(t *testing.T) {
	profile := &Profile{
		BaseURL: "https://api.example.com",
		Model:   "claude-sonnet-4-20250514",
	}

	got := BuildArgs(profile, []string{"-r"})
	want := []string{"-r"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected args:\n got: %#v\nwant: %#v", got, want)
	}
}
