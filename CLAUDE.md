# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

`cl` 是一个 CLI 工具，用于在终端内从任意项目目录启动 Claude Code 或 Codex，并自动注入对应的中转站/官方账号配置（Provider 认证）。它通过环境变量注入 API Key、Base URL、Model 等配置，不修改 `~/.claude/settings.json`，与 `cc switch`（管理 Skills 和 MCP）互不干扰。

## 核心架构

- **配置目录**：`~/.cl/`
- **profiles.yaml**：存储所有 Provider 配置（api_key、base_url、model、cli、自定义 env），包含 `default` 字段和 `defaults.env` 通用环境变量块
- 每个 profile 可通过 `cli` 字段指定启动的 CLI 工具（`claude` 或 `codex`，默认 `claude`）
- 砍掉了 bind/unbind 和 projects.yaml，裸 `cl` 走全局默认 profile

### 启动流程

1. 解析 CLI 参数：第一个位置参数为 profile 名或子命令（list/edit/default/help 为保留字，不可用作 profile 名），`cl` 自身无任何 flag，其余参数全部透传给目标 CLI
2. 清理冲突环境变量：`ANTHROPIC_*`（AUTH_TOKEN、API_KEY、BASE_URL、MODEL、SMALL_FAST_MODEL）和 `OPENAI_*`（API_KEY、BASE_URL、MODEL）
3. 按 profile 的 `cli` 字段决定环境变量映射：claude → `ANTHROPIC_*`，codex → `OPENAI_*`
4. 注入配置，优先级：顶层字段 > profile.env > defaults.env
5. 启动对应 CLI 子进程（`claude` 或 `codex`），附带透传参数

### CLI 命令

| 命令 | 说明 |
|------|------|
| `cl` | 使用全局默认 profile 启动 |
| `cl <profile>` | 指定 profile 启动 |
| `cl <profile> -r` | 指定 profile，`-r` 透传给目标 CLI |
| `cl list` | 列出所有 profile |
| `cl edit` | 用 `$EDITOR` 打开 profiles.yaml |
| `cl default <profile>` | 设置全局默认 profile |

## 多 CLI 支持

- Profile 的 `cli` 字段指定启动哪个 CLI 工具，白名单：`claude`（默认）、`codex`
- `cli: claude` → 顶层字段映射到 `ANTHROPIC_*` 环境变量，启动 `claude` 二进制
- `cli: codex` → 顶层字段映射到 `OPENAI_*` 环境变量（`api_key` → `OPENAI_API_KEY`，`base_url` → `OPENAI_BASE_URL`，`model` → `OPENAI_MODEL`），启动 `codex` 二进制；`auth_token` 和 `small_fast_model` 对 codex 无意义，设置时被忽略
- Codex 特有配置（如 `CODEX_CONFIG_DIR`）通过 `env` 字段注入

## 关键约束

- 敏感信息（api_key）明文存本地，依赖文件权限 `600` 保护
- 不修改 `~/.claude/settings.json`，只通过环境变量覆盖认证
- 未设置默认 profile 时裸 `cl` 报错，不做猜测
- 启动前必须 unset 所有冲突环境变量（`ANTHROPIC_*` 和 `OPENAI_*`），防止 Auth conflict
- `cli` 字段白名单校验：仅允许 `claude`、`codex`，非法值在加载配置时报错
- 保留字校验：list、edit、default、help 不可用作 profile 名
