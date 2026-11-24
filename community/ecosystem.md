---
layout: page
title: Ecosystem
nav-title: Ecosystem
---

Gemara is a small ecosystem of repositories, each with a deliberate scope. This page records what belongs where, so contributors know where a change should go and reviews stay focused.

The scopes below are the intended ones. Where a repository does not match its scope yet, the row says what is true today.

| Repository | Scope |
| --- | --- |
| **[gemara](https://github.com/gemaraproj/gemara)** (spec) | Declarative definitions: the model, CUE schemas, artifact shapes, and the distribution formats under `distribution/`: the bundle contract and the resolve/lock/analysis formats. Deliberately slow-moving; changes are controlled and reviewed carefully. Distribution formats release with the spec until they need their own cycle (ADR 0026). Today the spec has the model and CUE schemas in a single module. The bundle contract is in review in [gemaraproj/gemara#648](https://github.com/gemaraproj/gemara/pull/648), and the resolve, lock, and analysis formats are not defined yet. |
| **SDKs** ([Go](https://github.com/gemaraproj/go-gemara), [Python](https://github.com/gemaraproj/gemara-python), [React](https://github.com/gemaraproj/gemara-react)) | Deliberately kept small. The exact scope of the SDKs will be defined by a future spec. |
| **[grc-store-clientkit](https://github.com/gemaraproj/grc-store-clientkit)** (rename proposed in ADR 0028) | Publishing, consuming, and interacting with registries: bundle management, trust and signing wiring shared by tools that talk to a hub. It builds on the bundle packing in `go-gemara`. Its wire shapes come from `grc-store-protocol`, which is maintained outside gemaraproj. Whether those shapes move under the spec is undecided. |
| **[grcli](https://github.com/gemaraproj/grcli)** | The workbench CLI. Executes and interprets what the spec declares. |
| **[gemara-publish-action](https://github.com/gemaraproj/gemara-publish-action)** | GitHub Action for publishing Gemara artifact bundles to OCI registries. Today it embeds its own CLI. [gemaraproj/gemara-publish-action#40](https://github.com/gemaraproj/gemara-publish-action/pull/40) moves it onto grcli. |
| **[gemara-ai](https://github.com/gemaraproj/gemara-ai)** | AI skills, OpenPackage bundles, and MCP configurations for agents working with Gemara artifacts. Today its plugin runs gemara-mcp rather than grcli. |
| **[gemara-mcp](https://github.com/gemaraproj/gemara-mcp)** | The Gemara MCP server. ADR 0027 proposes folding it into grcli as `grcli mcp serve`; discussion is in [gemaraproj/grcli#17](https://github.com/gemaraproj/grcli/issues/17). |
| **[gemara-mcp-eval-infra](https://github.com/gemaraproj/gemara-mcp-eval-infra)** | Output-determinism evaluation harness for gemara-mcp. |
| **[website](https://github.com/gemaraproj/website)** | Source for this site: documentation, tutorials, specifications, and ADRs. |
| **[community](https://github.com/gemaraproj/community)** | Community content proposals. |
| **[.github](https://github.com/gemaraproj/.github)** | Organization governance: charter, maintainers, contributor ladder, code of conduct, and security policy. |
| **[awesome-gemara](https://github.com/gemaraproj/awesome-gemara)** | A curated list of resources, tools, and projects related to Gemara. |

A few boundaries worth calling out:

- If it's a shape or a format that clients and the hub must agree on, it goes in the spec — not in an SDK, where it would become "one SDK's opinion."
- If it's shared behavior for publishing or consumption, it goes in grc-store-clientkit, not the SDKs.
- Skills and agent tooling install and use grcli as their deterministic tool; grcli does not manage skills. Today gemara-ai runs gemara-mcp instead (see ADR 0027).
