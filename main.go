package main

import (
	"fmt"
	"os"
	"strings"
	"syscall"
)

var reservedWords = map[string]bool{
	"list":    true,
	"edit":    true,
	"default": true,
	"help":    true,
	"version": true,
}

var version = "dev"

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		execLaunchDefault(nil)
		return
	}

	switch args[0] {
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
			fatal("配置文件不存在，请运行 cl edit 创建")
		}
		fatal("加载配置失败: %v", err)
	}
	if cfg.Default == "" {
		fatal("未设置默认 profile，请使用 cl <profile> 或 cl default <profile>")
	}
	profile, err := cfg.GetProfile(cfg.Default)
	if err != nil {
		fatal("默认 profile %q 不存在，请检查配置", cfg.Default)
	}
	env := BuildEnv(cfg, profile)
	args := BuildArgs(profile, cliArgs)
	if err := Launch(profile.GetCLI(), env, args); err != nil {
		fatal("%v", err)
	}
}

func execLaunchProfile(name string, cliArgs []string) {
	cfg, err := LoadConfig()
	if err != nil {
		if os.IsNotExist(err) {
			fatal("配置文件不存在，请运行 cl edit 创建")
		}
		fatal("加载配置失败: %v", err)
	}
	profile, err := cfg.GetProfile(name)
	if err != nil {
		available := make([]string, 0, len(cfg.Profiles))
		for k := range cfg.Profiles {
			available = append(available, k)
		}
		fatal("profile %q 不存在，可用: %s", name, strings.Join(available, ", "))
	}
	env := BuildEnv(cfg, profile)
	args := BuildArgs(profile, cliArgs)
	if err := Launch(profile.GetCLI(), env, args); err != nil {
		fatal("%v", err)
	}
}

func execList() {
	cfg, err := LoadConfig()
	if err != nil {
		if os.IsNotExist(err) {
			fatal("配置文件不存在，请运行 cl edit 创建")
		}
		fatal("加载配置失败: %v", err)
	}
	if len(cfg.Profiles) == 0 {
		fmt.Println("暂无 profile，请运行 cl edit 添加")
		return
	}
	for k, p := range cfg.Profiles {
		marker := ""
		if k == cfg.Default {
			marker = " (default)"
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
		os.MkdirAll(dir, 0700)
		os.WriteFile(path, []byte(configTemplate), 0600)
		fmt.Fprintf(os.Stderr, "已创建模板配置: %s\n", path)
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
		fatal("未找到编辑器，请设置 $EDITOR 环境变量")
	}

	editorPath, err := findExecutable(editor)
	if err != nil {
		fatal("编辑器 %q 未找到: %v", editor, err)
	}

	syscall.Exec(editorPath, []string{editor, path}, os.Environ())
}

func execDefault(args []string) {
	if len(args) == 0 {
		cfg, err := LoadConfig()
		if err == nil && cfg.Default != "" {
			fmt.Printf("当前默认 profile: %s\n", cfg.Default)
		} else {
			fmt.Println("未设置默认 profile")
		}
		return
	}

	name := args[0]
	if reservedWords[name] {
		fatal("%q 是保留字，不可用作 profile 名", name)
	}

	cfg, err := LoadConfig()
	if err != nil {
		if os.IsNotExist(err) {
			fatal("配置文件不存在，请运行 cl edit 创建")
		}
		fatal("加载配置失败: %v", err)
	}

	if _, err := cfg.GetProfile(name); err != nil {
		fatal("%v", err)
	}

	cfg.Default = name
	if err := SaveConfig(cfg); err != nil {
		fatal("保存配置失败: %v", err)
	}
	fmt.Printf("默认 profile 已设置为: %s\n", name)
}

func execHelp() {
	fmt.Print(`cl — Claude Code / Codex 多环境启动器

用法:
  cl                    使用默认 profile 启动
  cl <profile>          指定 profile 启动
  cl <profile> [args]   指定 profile，透传参数给 CLI
  cl list               列出所有 profile
  cl edit               编辑 profiles.yaml
  cl default [profile]  查看或设置默认 profile
  cl help               显示帮助
  cl version            显示版本

配置文件: ~/.cl/profiles.yaml
支持的 CLI: claude (默认), codex
`)
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
	return "", fmt.Errorf("%s not found in PATH", name)
}

func fatal(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "cl: "+format+"\n", args...)
	os.Exit(1)
}
