package main

import (
	"os"
	"strings"
)

type messageID string

const (
	msgConfigMissing           messageID = "config_missing"
	msgLoadConfigFailed        messageID = "load_config_failed"
	msgDefaultNotSet           messageID = "default_not_set"
	msgDefaultProfileMissing   messageID = "default_profile_missing"
	msgProfileMissingAvailable messageID = "profile_missing_available"
	msgNoProfiles              messageID = "no_profiles"
	msgDefaultMarker           messageID = "default_marker"
	msgCreatedTemplate         messageID = "created_template"
	msgCreateConfigDirFailed   messageID = "create_config_dir_failed"
	msgWriteConfigFailed       messageID = "write_config_failed"
	msgEditorMissing           messageID = "editor_missing"
	msgEditorNotFound          messageID = "editor_not_found"
	msgCurrentDefault          messageID = "current_default"
	msgDefaultNotSetShort      messageID = "default_not_set_short"
	msgReservedProfileName     messageID = "reserved_profile_name"
	msgSaveConfigFailed        messageID = "save_config_failed"
	msgDefaultSet              messageID = "default_set"
	msgHelp                    messageID = "help"
	msgYAMLParseFailed         messageID = "yaml_parse_failed"
	msgInvalidCLI              messageID = "invalid_cli"
	msgProfileMissing          messageID = "profile_missing"
	msgCLINotFound             messageID = "cli_not_found"
	msgExecutableNotFound      messageID = "executable_not_found"
)

var localizedMessages = map[string]map[messageID]string{
	"en": {
		msgConfigMissing:           "config file not found, run cl edit to create it",
		msgLoadConfigFailed:        "failed to load config: %v",
		msgDefaultNotSet:           "default profile is not set, use cl <profile> or cl default <profile>",
		msgDefaultProfileMissing:   "default profile %q does not exist, please check your config",
		msgProfileMissingAvailable: "profile %q does not exist, available: %s",
		msgNoProfiles:              "no profiles yet, run cl edit to add one",
		msgDefaultMarker:           "(default)",
		msgCreatedTemplate:         "created template config: %s",
		msgCreateConfigDirFailed:   "failed to create config directory: %v",
		msgWriteConfigFailed:       "failed to write template config: %v",
		msgEditorMissing:           "editor not found, set the $EDITOR environment variable",
		msgEditorNotFound:          "editor %q not found: %v",
		msgCurrentDefault:          "current default profile: %s",
		msgDefaultNotSetShort:      "default profile is not set",
		msgReservedProfileName:     "%q is a reserved word and cannot be used as a profile name",
		msgSaveConfigFailed:        "failed to save config: %v",
		msgDefaultSet:              "default profile set to: %s",
		msgHelp: `cl - Claude Code / Codex multi-environment launcher

Usage:
  cl                    Launch with the default profile
  cl <profile>          Launch with a specific profile
  cl <profile> [args]   Launch with a specific profile and pass args to the CLI
  cl list               List all profiles
  cl edit               Edit profiles.yaml
  cl default [profile]  Show or set the default profile
  cl help               Show help
  cl version            Show version

Config file: ~/.cl/profiles.yaml
Supported CLIs: claude (default), codex
`,
		msgYAMLParseFailed:    "failed to parse profiles.yaml: %w",
		msgInvalidCLI:         "profile %q: invalid cli %q, only claude and codex are supported",
		msgProfileMissing:     "profile %q does not exist",
		msgCLINotFound:        "%s not found, please install it and make sure it is in PATH",
		msgExecutableNotFound: "%s not found in PATH",
	},
	"zh": {
		msgConfigMissing:           "配置文件不存在，请运行 cl edit 创建",
		msgLoadConfigFailed:        "加载配置失败: %v",
		msgDefaultNotSet:           "未设置默认 profile，请使用 cl <profile> 或 cl default <profile>",
		msgDefaultProfileMissing:   "默认 profile %q 不存在，请检查配置",
		msgProfileMissingAvailable: "profile %q 不存在，可用: %s",
		msgNoProfiles:              "暂无 profile，请运行 cl edit 添加",
		msgDefaultMarker:           "(默认)",
		msgCreatedTemplate:         "已创建模板配置: %s",
		msgCreateConfigDirFailed:   "创建配置目录失败: %v",
		msgWriteConfigFailed:       "写入模板配置失败: %v",
		msgEditorMissing:           "未找到编辑器，请设置 $EDITOR 环境变量",
		msgEditorNotFound:          "编辑器 %q 未找到: %v",
		msgCurrentDefault:          "当前默认 profile: %s",
		msgDefaultNotSetShort:      "未设置默认 profile",
		msgReservedProfileName:     "%q 是保留字，不可用作 profile 名",
		msgSaveConfigFailed:        "保存配置失败: %v",
		msgDefaultSet:              "默认 profile 已设置为: %s",
		msgHelp: `cl - Claude Code / Codex 多环境启动器

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
`,
		msgYAMLParseFailed:    "profiles.yaml 解析失败: %w",
		msgInvalidCLI:         "profile %q: cli %q 不合法，仅支持 claude、codex",
		msgProfileMissing:     "profile %q 不存在",
		msgCLINotFound:        "%s 未找到，请确认已安装并在 PATH 中",
		msgExecutableNotFound: "%s 不在 PATH 中",
	},
}

func tr(id messageID) string {
	lang := appLanguage()
	if messages, ok := localizedMessages[lang]; ok {
		if message, ok := messages[id]; ok {
			return message
		}
	}
	return localizedMessages["en"][id]
}

func appLanguage() string {
	if lang := normalizeLanguage(os.Getenv("CL_LANG")); lang != "" {
		return lang
	}
	for _, key := range []string{"LC_ALL", "LC_MESSAGES", "LANG", "LANGUAGE"} {
		if lang := normalizeLanguage(os.Getenv(key)); lang != "" {
			return lang
		}
	}
	return "en"
}

func normalizeLanguage(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	if value == "" || value == "c" || value == "posix" {
		return ""
	}
	value = strings.Split(value, ":")[0]
	value = strings.ReplaceAll(value, "-", "_")
	if strings.HasPrefix(value, "zh") {
		return "zh"
	}
	return "en"
}
