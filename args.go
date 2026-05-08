package main

import "strconv"

const codexProfileProvider = "cl"
const codexDefaultBaseURL = "https://api.openai.com/v1"

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

	if profile.APIKey != "" || profile.BaseURL != "" {
		baseURL := profile.BaseURL
		if baseURL == "" {
			baseURL = codexDefaultBaseURL
		}
		addConfig("model_provider", codexProfileProvider)
		addConfig("model_providers."+codexProfileProvider+".name", codexProfileProvider)
		addConfig("model_providers."+codexProfileProvider+".base_url", baseURL)
		addConfig("model_providers."+codexProfileProvider+".wire_api", "responses")
		if profile.APIKey != "" {
			addConfig("model_providers."+codexProfileProvider+".env_key", "OPENAI_API_KEY")
		}
	}
	addConfig("model", profile.Model)
	addConfig("model_reasoning_effort", profile.ModelReasoningEffort)

	return args
}
