package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"syscall"
	"unicode/utf8"
)

var reservedWords = map[string]bool{
	"list":    true,
	"edit":    true,
	"default": true,
	"help":    true,
	"version": true,
	"init":    true,
}

var version = "dev"

var urlPattern = regexp.MustCompile(`https?://[^\s<>"'` + "`" + `]+`)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		execLaunchDefault(nil)
		return
	}

	switch args[0] {
	case "init":
		execInit()
	case "list":
		execList()
	case "edit":
		execEdit()
	case "default":
		execDefault(args[1:])
	case "help":
		execHelp()
	case "version", "--version", "-v":
		execVersion()
	default:
		execLaunchProfile(args[0], args[1:])
	}
}

func execLaunchDefault(cliArgs []string) {
	cfg, err := LoadConfig()
	if err != nil {
		if os.IsNotExist(err) {
			fatal(tr(msgConfigMissingInit))
		}
		fatal(tr(msgLoadConfigFailed), err)
	}
	if cfg.Default == "" {
		fatal(tr(msgDefaultNotSet))
	}
	profile, err := cfg.GetProfile(cfg.Default)
	if err != nil {
		fatal(tr(msgDefaultProfileMissing), cfg.Default)
	}
	env := BuildEnv(cfg, profile)
	args := BuildArgs(profile, cliArgs)
	printProfileName(profile)
	if err := Launch(profile.GetCLI(), env, args); err != nil {
		fatal("%v", err)
	}
}

func execLaunchProfile(name string, cliArgs []string) {
	cfg, err := LoadConfig()
	if err != nil {
		if os.IsNotExist(err) {
			fatal(tr(msgConfigMissingInit))
		}
		fatal(tr(msgLoadConfigFailed), err)
	}
	profile, err := cfg.GetProfile(name)
	if err != nil {
		available := make([]string, 0, len(cfg.Profiles))
		for k := range cfg.Profiles {
			available = append(available, k)
		}
		fatal(tr(msgProfileMissingAvailable), name, strings.Join(available, ", "))
	}
	env := BuildEnv(cfg, profile)
	args := BuildArgs(profile, cliArgs)
	printProfileName(profile)
	if err := Launch(profile.GetCLI(), env, args); err != nil {
		fatal("%v", err)
	}
}

func printProfileName(profile *Profile) {
	if profile.Name != "" {
		name := profile.Name
		if isTerminal(os.Stderr) {
			name = linkURLs(name)
		}
		fmt.Fprintln(os.Stderr, name)
	}
}

func isTerminal(file *os.File) bool {
	info, err := file.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}

func linkURLs(text string) string {
	urls := profileNameURLs(text)
	if len(urls) == 0 {
		return text
	}

	var b strings.Builder
	last := 0
	for _, match := range urls {
		start, end := match.start, match.end
		b.WriteString(text[last:start])
		url := text[start:end]
		b.WriteString(terminalLink(url, url))
		last = end
	}
	b.WriteString(text[last:])
	return b.String()
}

type urlMatch struct {
	start int
	end   int
}

func profileNameURLs(text string) []urlMatch {
	matches := urlPattern.FindAllStringIndex(text, -1)
	urls := make([]urlMatch, 0, len(matches))
	for _, match := range matches {
		start, end := match[0], trimURLMatchEnd(text, match[0], match[1])
		if start == end {
			continue
		}
		urls = append(urls, urlMatch{start: start, end: end})
	}
	return urls
}

func trimURLMatchEnd(text string, start, end int) int {
	for end > start {
		r, size := rune(text[end-1]), 1
		if r >= 0x80 {
			r, size = utf8.DecodeLastRuneInString(text[start:end])
		}
		if !strings.ContainsRune(".,;:!?)]}", r) {
			break
		}
		end -= size
	}
	return end
}

func terminalLink(url, text string) string {
	return "\x1b]8;;" + sanitizeTerminalLinkURL(url) + "\a" + text + "\x1b]8;;\a"
}

func sanitizeTerminalLinkURL(url string) string {
	return strings.Map(func(r rune) rune {
		if r < 0x20 || r == 0x7f {
			return -1
		}
		return r
	}, url)
}

func execList() {
	cfg, err := LoadConfig()
	if err != nil {
		if os.IsNotExist(err) {
			fatal(tr(msgConfigMissingInit))
		}
		fatal(tr(msgLoadConfigFailed), err)
	}
	if len(cfg.Profiles) == 0 {
		fmt.Println(tr(msgNoProfiles))
		return
	}
	for k, p := range cfg.Profiles {
		marker := ""
		if k == cfg.Default {
			marker = " " + tr(msgDefaultMarker)
		}
		displayName := k
		if p.Name != "" {
			displayName = fmt.Sprintf("%s (%s)", k, p.Name)
		}
		fmt.Printf("  %s%s\n", displayName, marker)
	}
}

func execEdit() {
	path := ConfigPath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		dir := ConfigDir()
		if err := os.MkdirAll(dir, 0700); err != nil {
			fatal(tr(msgCreateConfigDirFailed), err)
		}
		if err := os.WriteFile(path, []byte(defaultConfigTemplate()), 0600); err != nil {
			fatal(tr(msgWriteConfigFailed), err)
		}
		fmt.Fprintf(os.Stderr, tr(msgCreatedTemplate)+"\n", path)
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		for _, e := range []string{"vim", "vi", "nano"} {
			if p, err := findExecutable(e); err == nil {
				editor = p
				break
			}
		}
	}
	if editor == "" {
		fatal(tr(msgEditorMissing))
	}

	editorPath, err := findExecutable(editor)
	if err != nil {
		fatal(tr(msgEditorNotFound), editor, err)
	}

	syscall.Exec(editorPath, []string{editor, path}, os.Environ())
}

func execDefault(args []string) {
	if len(args) == 0 {
		cfg, err := LoadConfig()
		if err == nil && cfg.Default != "" {
			fmt.Printf(tr(msgCurrentDefault)+"\n", cfg.Default)
		} else {
			fmt.Println(tr(msgDefaultNotSetShort))
		}
		return
	}

	name := args[0]
	if reservedWords[name] {
		fatal(tr(msgReservedProfileName), name)
	}

	cfg, err := LoadConfig()
	if err != nil {
		if os.IsNotExist(err) {
			fatal(tr(msgConfigMissingInit))
		}
		fatal(tr(msgLoadConfigFailed), err)
	}

	if _, err := cfg.GetProfile(name); err != nil {
		fatal("%v", err)
	}

	cfg.Default = name
	if err := SaveConfig(cfg); err != nil {
		fatal(tr(msgSaveConfigFailed), err)
	}
	fmt.Printf(tr(msgDefaultSet)+"\n", name)
}

func execHelp() {
	fmt.Print(tr(msgHelp))
}

func execVersion() {
	fmt.Printf("cl %s\n", version)
}

func findExecutable(name string) (string, error) {
	// Handle absolute paths
	if strings.HasPrefix(name, "/") {
		return name, nil
	}
	return findInPath(name)
}

func findInPath(name string) (string, error) {
	path, err := lookPath(name)
	if err != nil {
		return "", err
	}
	return path, nil
}

func lookPath(name string) (string, error) {
	pathEnv := os.Getenv("PATH")
	for _, dir := range strings.Split(pathEnv, ":") {
		full := dir + "/" + name
		if info, err := os.Stat(full); err == nil && !info.IsDir() {
			return full, nil
		}
	}
	return "", fmt.Errorf(tr(msgExecutableNotFound), name)
}

func fatal(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "cl: "+format+"\n", args...)
	os.Exit(1)
}
