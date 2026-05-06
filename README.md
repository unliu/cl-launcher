# cl

CLI launcher for [Claude Code](https://docs.anthropic.com/en/docs/claude-code) and [Codex](https://github.com/openai/codex) with provider profiles.

Switch between API providers / relay services in one command — no need to manually juggle environment variables or edit `~/.claude/settings.json`.

## Install

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

## Configuration

Config lives at `~/.cl/profiles.yaml`:

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
    small_fast_model: claude-haiku-4-5-20251001
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

### Profile Fields

| Field | Description |
|---|---|
| `name` | Display name (shown in `cl list`) |
| `cli` | Target CLI: `claude` (default) or `codex` |
| `base_url` | API base URL |
| `api_key` | API key |
| `auth_token` | Auth token (Claude only, alternative to api_key) |
| `model` | Model override |
| `small_fast_model` | Small/fast model override (Claude only) |
| `env` | Extra environment variables |

### Environment Variable Mapping

- `cli: claude` — top-level fields map to `ANTHROPIC_*` env vars
- `cli: codex` — top-level fields map to `OPENAI_*` env vars

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

## Security

`profiles.yaml` contains API keys in plaintext, protected by `0600` file permissions. The config directory `~/.cl/` is created with `0700` permissions.

## License

MIT
