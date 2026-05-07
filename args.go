package main

import "strconv"

func BuildArgs(profile *Profile, passthrough []string) []string {
	if profile.GetCLI() != "codex" {
		return append([]string{}, passthrough...)
	}

	args := codexConfigArgs(profile)
	args = append(args, passthrough...)
	return args
}

func codexConfigArgs(profile *Profile) []string {
	var args []string
	addConfig := func(key, value string) {
		if value == "" {
			return
		}
		args = append(args, "-c", key+"="+strconv.Quote(value))
	}

	addConfig("openai_base_url", profile.BaseURL)
	if profile.APIKey != "" || profile.BaseURL != "" {
		addConfig("forced_login_method", "api")
	}
	addConfig("model", profile.Model)
	addConfig("model_reasoning_effort", profile.ModelReasoningEffort)

	return args
}
