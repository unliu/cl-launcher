<div align="center">

# cl — CLI launcher for Claude Code and Codex

[中文](README.md) | English

> Switch between API providers for Claude Code and Codex in one command

![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go)

</div>

CLI launcher for [Claude Code](https://docs.anthropic.com/en/docs/claude-code) and [Codex](https://github.com/openai/codex) with provider profiles.

Pairs well with [CC Switch](https://github.com/farion1231/cc-switch): CC Switch manages Skills and MCP, `cl` manages Provider auth — independent, non-conflicting, best together.

## Features

- Switch between multiple provider profiles, launch with one command
- Support both Claude Code and Codex CLI tools
- Inject config via environment variables and Codex CLI config overrides, no modification to `~/.claude/settings.json`
- Complements CC Switch: Skills/MCP management and Provider auth are independent

## Install

**Prerequisites**: [Claude Code](https://docs.anthropic.com/en/docs/claude-code) or [Codex](https://github.com/openai/codex) must be installed.

```bash
brew install unliu/tap/cl
```

Or with Go:

```bash
go install github.com/unliu/cl-launcher@latest
```

## Quick Start

```bash
# Create config (opens editor)
cl edit

# Set default profile
cl default myrelay

# Launch Claude Code with default profile
cl

# Launch with a specific profile
cl myrelay

# Pass flags through to the underlying CLI
cl myrelay -r
```

## Local Test Build

For development or pre-release checks, build a local test binary first:

```bash
./scripts/build-local.sh
```

The default output is `dist/local/cl-dev`. The script runs `go test ./...` first, then builds a test binary with version `local`. Use it in place of the installed `cl` for smoke tests:

```bash
./dist/local/cl-dev <profile> debug models --bundled
./dist/local/cl-dev <profile>
```

You can also choose a custom output path:

```bash
./scripts/build-local.sh /tmp/cl-dev
```

## Configuration

Config lives at `~/.cl/profiles.yaml`:

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
    name: Kimi k2.6 | https://www.kimi.com/membership/subscription
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

### Profile Fields

| Field | Description |
|---|---|
| `name` | Display name (shown in `cl list`) |
| `cli` | Target CLI: `claude` (default) or `codex` |
| `base_url` | API base URL |
| `api_key` | API key |
| `auth_token` | Auth token (Claude only, alternative to api_key) |
| `model` | Model override |
| `model_reasoning_effort` | Codex reasoning effort override (Codex only) |
| `env` | Extra environment variables |

### Environment Variable Mapping

- `cli: claude` — top-level fields map to `ANTHROPIC_*` env vars
- `cli: codex` — `api_key` maps to `OPENAI_API_KEY`; `model` and `model_reasoning_effort` are passed as `codex -c ...` overrides. When `api_key` or `base_url` is set, `cl` injects a temporary custom `model_provider` named `cl`, uses `base_url` (or `https://api.openai.com/v1` when omitted), and reads the key from `OPENAI_API_KEY` without modifying or depending on `~/.codex/auth.json`.

### Priority

Top-level fields > `profile.env` > `defaults.env`

All conflicting `ANTHROPIC_*` and `OPENAI_*` env vars are cleared before injection to prevent auth conflicts.

## Commands

| Command | Description |
|---|---|
| `cl` | Launch with default profile |
| `cl <profile>` | Launch with specified profile |
| `cl <profile> [args]` | Launch with profile, pass args to CLI |
| `cl list` | List all profiles |
| `cl edit` | Edit profiles.yaml in `$EDITOR` |
| `cl default [profile]` | Show or set default profile |
| `cl help` | Show help |
| `cl version` / `cl --version` | Show version |

## Security

`profiles.yaml` contains API keys in plaintext, protected by `0600` file permissions. The config directory `~/.cl/` is created with `0700` permissions.

## License

MIT
