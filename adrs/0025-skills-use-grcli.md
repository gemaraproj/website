---
layout: page
title: Skills and Agent Tooling Use grcli, Not the Reverse
---

- **ADR:** 0025
- **Proposal Author(s):** @jpower432
- **Status:** Accepted

## Context

AI skills and MCP-based agents need deterministic tools to invoke. An earlier idea had grcli fetching and managing skills and MCP updates itself, so a CLI version could pull newer skill content without a release.

## Action

Skills and agent tooling install and invoke grcli as their deterministic tool. grcli does not manage, fetch, or version skills.

Today gemara-ai's plugin runs gemara-mcp, not grcli. ADR 0027 proposes folding gemara-mcp into grcli as `grcli mcp serve`.

## Consequences

Positive: grcli stays a workbench CLI rather than a skills package manager; skill content evolves independently of CLI releases.

Negative: agents must install grcli separately rather than getting it via another tool they already have.

## Alternatives Considered

- **grcli pulls the latest skills/MCP content** — rejected: turns the CLI into a package manager for content it does not own, and version-management responsibility lands in the wrong place.
