package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

const claudeDefaultBaseURL = "https://api.anthropic.com"

func execInit() {
	if !isTerminal(os.Stdin) {
		fatal(tr(msgInitNotInteractive))
	}
	runInit(os.Stdin, os.Stderr)
}

func runInit(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	var cfg *Config
	var isNew bool
	loaded, err := LoadConfig()
	if err != nil {
		if !os.IsNotExist(err) {
			fatal(tr(msgLoadConfigFailed), err)
		}
		isNew = true
		cfg = &Config{
			Defaults: Defaults{
				Env: map[string]string{
					"CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC": "1",
				},
			},
			Profiles: make(map[string]Profile),
		}
	} else {
		cfg = loaded
	}

	fmt.Fprintln(out, tr(msgInitHeader))

	profileName := promptProfileName(scanner, out, cfg, "my")
	cli := promptChoice(scanner, out, tr(msgInitPromptCLI), []string{"claude", "codex"}, "claude")

	apiKey := promptRequired(scanner, out, tr(msgInitPromptAPIKey))

	defaultBaseURL := claudeDefaultBaseURL
	if cli == "codex" {
		defaultBaseURL = codexDefaultBaseURL
	}
	baseURL := promptWithDefault(scanner, out, tr(msgInitPromptBaseURL), defaultBaseURL)
	if baseURL == defaultBaseURL {
		baseURL = ""
	}

	model := promptWithDefault(scanner, out, tr(msgInitPromptModel), "")

	setDefault := isNew || len(cfg.Profiles) == 0
	if !setDefault {
		answer := promptChoice(scanner, out, tr(msgInitPromptSetDefault), []string{"y", "n"}, "y")
		setDefault = answer == "y"
	}

	profile := Profile{
		CLI:     cli,
		APIKey:  apiKey,
		BaseURL: baseURL,
		Model:   model,
	}

	cfg.Profiles[profileName] = profile
	if setDefault {
		cfg.Default = profileName
	}

	if err := SaveConfig(cfg); err != nil {
		fatal(tr(msgSaveConfigFailed), err)
	}

	if setDefault {
		fmt.Fprintf(out, "\n"+tr(msgInitSuccessDefault)+"\n", profileName, ConfigPath())
	} else {
		fmt.Fprintf(out, "\n"+tr(msgInitSuccess)+"\n", profileName, ConfigPath())
	}
}

func promptProfileName(scanner *bufio.Scanner, out io.Writer, cfg *Config, defaultName string) string {
	for {
		name := promptWithDefault(scanner, out, tr(msgInitPromptProfile), defaultName)
		if reservedWords[name] {
			fmt.Fprintf(out, "  "+tr(msgReservedProfileName)+"\n", name)
			continue
		}
		if _, exists := cfg.Profiles[name]; exists {
			fmt.Fprintf(out, "  "+tr(msgInitProfileExists)+"\n", name)
			continue
		}
		return name
	}
}

func promptWithDefault(scanner *bufio.Scanner, out io.Writer, label string, defaultVal string) string {
	if defaultVal != "" {
		fmt.Fprintf(out, "  %s [%s]: ", label, defaultVal)
	} else {
		fmt.Fprintf(out, "  %s: ", label)
	}
	if !scanner.Scan() {
		return defaultVal
	}
	val := strings.TrimSpace(scanner.Text())
	if val != "" {
		return val
	}
	return defaultVal
}

func promptRequired(scanner *bufio.Scanner, out io.Writer, label string) string {
	for {
		fmt.Fprintf(out, "  %s: ", label)
		if !scanner.Scan() {
			fmt.Fprintln(out)
			os.Exit(1)
		}
		val := strings.TrimSpace(scanner.Text())
		if val != "" {
			return val
		}
		fmt.Fprintf(out, "  %s\n", tr(msgInitRequired))
	}
}

func promptChoice(scanner *bufio.Scanner, out io.Writer, label string, choices []string, defaultVal string) string {
	for {
		val := promptWithDefault(scanner, out, label, defaultVal)
		for _, c := range choices {
			if val == c {
				return val
			}
		}
		fmt.Fprintf(out, "  "+tr(msgInitInvalidChoice)+"\n", strings.Join(choices, ", "))
	}
}
