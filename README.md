# cl — CLI launcher for Claude Code and Codex

[Claude Code](https://docs.anthropic.com/en/docs/claude-code) 和 [Codex](https://github.com/openai/codex) 的 CLI 启动器，支持 Provider 配置文件。

一条命令即可在不同 API 供应商/中转站之间切换 — 无需手动切换环境变量或修改 `~/.claude/settings.json`。

与 [CC Switch](https://github.com/farion1231/cc-switch) 各司其职：CC Switch 管理 Skills 和 MCP，`cl` 管理 Provider 认证，互不干扰，搭配使用。

[English](README_en.md)

## 安装

```bash
brew install unliu/tap/cl
```

或使用 Go 安装：

```bash
go install github.com/unliu/cl-launcher@latest
```

## 快速开始

```bash
# 创建配置（打开编辑器）
cl edit

# 设置默认 profile
cl default myrelay

# 使用默认 profile 启动 Claude Code
cl

# 使用指定 profile 启动
cl myrelay

# 透传参数到底层 CLI
cl myrelay -r
```

## 配置

配置文件位于 `~/.cl/profiles.yaml`：

```yaml
default: myrelay

defaults:
  env:
    CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC: "1"

profiles:
  myrelay:
    name: My Relay
    base_url: https://relay.example.com
    api_key: sk-xxx
    model: claude-sonnet-4-20250514
    env: {}

  official:
    name: Anthropic Direct
    api_key: sk-ant-xxx
    env: {}

  codex-provider:
    name: Codex via Relay
    cli: codex
    base_url: https://relay.example.com
    api_key: sk-xxx
    model: o3
    env:
      CODEX_CONFIG_DIR: ~/.codex-envs/relay
```

### Profile 字段

| 字段 | 说明 |
|---|---|
| `name` | 显示名称（在 `cl list` 中展示） |
| `cli` | 目标 CLI：`claude`（默认）或 `codex` |
| `base_url` | API Base URL |
| `api_key` | API Key |
| `auth_token` | Auth Token（仅 Claude，替代 api_key） |
| `model` | 模型覆盖 |
| `env` | 额外的环境变量 |

### 环境变量映射

- `cli: claude` — 顶层字段映射到 `ANTHROPIC_*` 环境变量
- `cli: codex` — 顶层字段映射到 `OPENAI_*` 环境变量

### 优先级

顶层字段 > `profile.env` > `defaults.env`

启动前会清除所有冲突的 `ANTHROPIC_*` 和 `OPENAI_*` 环境变量，防止认证冲突。

## 命令

| 命令 | 说明 |
|---|---|
| `cl` | 使用默认 profile 启动 |
| `cl <profile>` | 使用指定 profile 启动 |
| `cl <profile> [args]` | 使用指定 profile 启动，透传参数 |
| `cl list` | 列出所有 profile |
| `cl edit` | 用 `$EDITOR` 打开 profiles.yaml |
| `cl default [profile]` | 查看或设置默认 profile |
| `cl help` | 显示帮助 |

## 安全性

`profiles.yaml` 明文存储 API Key，通过 `0600` 文件权限保护。配置目录 `~/.cl/` 以 `0700` 权限创建。

## 许可证

MIT
