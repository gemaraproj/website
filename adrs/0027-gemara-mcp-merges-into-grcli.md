---
layout: page
title: gemara-mcp Merges into grcli
---

- **ADR:** 0027
- **Proposal Author(s):** @jpower432
- **Status:** Proposed

## Context

gemara-mcp predates grcli. Today gemara-ai's plugin runs gemara-mcp, and ADR 0025 says agent tooling invokes grcli. Keeping both as separate tools causes three problems:

- **Two validate paths that can disagree.** grcli runs `cue vet` from `PATH`; gemara-mcp runs CUE in-process with schemas fetched from the registry. The same artifact can get different verdicts.
- **Two release pipelines,** each with its own container image and SLSA attestations.
- **gemara-ai has to reference both repositories.**

gemara-mcp's surface is also shrinking. Its prompts are moving into gemara-ai as skills ([gemaraproj/gemara-mcp#149](https://github.com/gemaraproj/gemara-mcp/issues/149), [gemaraproj/gemara-ai#24](https://github.com/gemaraproj/gemara-ai/pull/24)). What remains is two tools, `validate_gemara_artifact` and `migrate_gemara_artifact`, and its resources. The validate tool overlaps directly with `grcli validate`.

Discussion is in [gemaraproj/grcli#17](https://github.com/gemaraproj/grcli/issues/17).

## Action

gemara-mcp becomes a grcli subcommand, `grcli mcp serve`, in the same way the OpenSSF Scorecard CLI added `scorecard mcp`:

- The first change is parity only: the same tools and resources as `gemara-mcp serve`, with no new agent capabilities.
- `grcli validate` and the MCP validate tool share one implementation, using in-process CUE.
- After a final deprecation release, `gemaraproj/gemara-mcp` is archived with a pointer to grcli.

Open questions for the parity change:

- Whether the `gemara://` resource URIs and tool names stay the same, or a mapping is documented.
- How the release pipeline consolidates into a single grcli image and attestation.

## Consequences

Positive: one validate verdict for every artifact; one release pipeline; agents install one tool; ADR 0025 holds without exceptions.

Negative: grcli grows MCP code and its dependencies, so the MCP tool and resource list has to stay lean; existing `gemara-mcp serve` configurations break and need migration; gemara-mcp-eval-infra has to retarget to grcli.

## Alternatives Considered

- **Keep both repositories and share validate through a library** — rejected: removes the drift but keeps two installs and two release pipelines, and the library would have no natural home without growing the SDKs.
- **One binary calls the other** — rejected: a workaround that adds runtime coupling between two release cycles.
