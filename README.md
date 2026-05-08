<div align="center">

# cl — CLI launcher for Claude Code and Codex

中文 | [English](README_en.md)

> 一条命令切换 Claude Code / Codex 的 API 供应商配置

![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go)

</div>

[Claude Code](https://docs.anthropic.com/en/docs/claude-code) 和 [Codex](https://github.com/openai/codex) 的 CLI 启动器，通过 Provider 配置文件管理认证信息。

与 [CC Switch](https://github.com/farion1231/cc-switch) 各司其职：CC Switch 管理 Skills 和 MCP，`cl` 管理 Provider 认证，互不干扰，搭配使用。

## 功能特性

- 多 Provider 配置文件切换，一条命令启动不同 API 供应商
- 支持 Claude Code 和 Codex 两种 CLI 工具
- 通过环境变量和 Codex CLI 配置覆盖注入配置，不修改 `~/.claude/settings.json`
- 与 CC Switch 各司其职，管理 Skills/MCP 与 Provider 认证互不干扰

## 安装

**前置条件**：已安装 [Claude Code](https://docs.anthropic.com/en/docs/claude-code) 或 [Codex](https://github.com/openai/codex)。

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

## 本地测试版本

开发或验证未发布改动时，先生成本地测试二进制：

```bash
./scripts/build-local.sh
```

默认产物为 `dist/local/cl-dev`。它会先运行 `go test ./...`，再构建版本号为 `local` 的测试二进制。可以直接用它替代系统 `cl` 做 smoke test：

```bash
./dist/local/cl-dev <profile> debug models --bundled
./dist/local/cl-dev <profile>
```

也可以指定输出路径：

```bash
./scripts/build-local.sh /tmp/cl-dev
```

## 配置

配置文件位于 `~/.cl/profiles.yaml`：

```yaml
default: myrelay

defaults:
  env:
    CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC: "1"

profiles:
  cc:
    name: Anthropic Direct
    api_key: sk-ant-xxx
    env: {}

  kimi:
    name: Kimi k2.6 | https://www.kimi.com/code/console
    base_url: https://api.kimi.com/coding
    api_key: sk-xxx
    model: kimi-k2.6
    env: {}

  glm:
    name: SiliconFlow GLM-5.1 | https://cloud.siliconflow.cn/me/expensebill
    base_url: https://api.siliconflow.cn/
    api_key: sk-xxx
    model: Pro/zai-org/GLM-5.1
    env: {}

  ds:
    name: DeepSeek v4-pro | https://platform.deepseek.com
    base_url: https://api.deepseek.com/anthropic
    api_key: sk-xxx
    model: deepseek-v4-pro
    env: {}

  myrelay:
    name: SomeRelay opus-4.7 | https://somerelay.example.com/bill_address
    base_url: https://somerelay.example.com/anthropic
    api_key: sk-xxx
    model: claude-opus-4-7
    env: {}

  codex-provider:
    name: Codex via Relay
    cli: codex
    base_url: https://relay.example.com
    api_key: sk-xxx
    model: gpt-5.5
    model_reasoning_effort: xhigh
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
| `model_reasoning_effort` | Codex 推理强度覆盖（仅 Codex） |
| `env` | 额外的环境变量 |

### 环境变量映射

- `cli: claude` — 顶层字段映射到 `ANTHROPIC_*` 环境变量
- `cli: codex` — `api_key` 映射到 `OPENAI_API_KEY`；`model`、`model_reasoning_effort` 通过 `codex -c ...` 覆盖 Codex 配置；设置了 `api_key` 或 `base_url` 时会临时注入一个名为 `cl` 的自定义 `model_provider`，使用 `base_url`（未设置时为 `https://api.openai.com/v1`）并从 `OPENAI_API_KEY` 读取 key，不修改也不依赖 `~/.codex/auth.json`

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
| `cl version` / `cl --version` | 显示版本 |

## 安全性

`profiles.yaml` 明文存储 API Key，通过 `0600` 文件权限保护。配置目录 `~/.cl/` 以 `0700` 权限创建。

## 许可证

MIT
