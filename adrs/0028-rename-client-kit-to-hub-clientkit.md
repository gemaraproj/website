---
layout: page
title: Rename the Client Kit to hub-clientkit
---

- **ADR:** 0028
- **Proposal Author(s):** @jpower432
- **Status:** Proposed

## Context

ADR 0024 makes `grc-store-clientkit` the project's shared Go library for publishing, consumption, and hub interaction. Its name comes from grc-store, one hub implementation. A project library named after one hub suggests that hub is the default, which works against a vendor-neutral ecosystem.

The name does match the code today. The client kit's wire shapes come from `grc-store-protocol`, which is maintained outside the project, and some of its behavior follows that protocol directly. Whether those shapes move under the spec is undecided (ADR 0024).

The Go module path, `github.com/gemaraproj/grc-store-clientkit`, is imported by grcli and by repositories outside gemaraproj, including privateer-sdk and several Privateer plugins. Go rejects a module whose declared path does not match its import path, so a rename requires every importer to update.

## Action

Rename `grc-store-clientkit` to `hub-clientkit` once it speaks a hub contract the project owns. Until then it keeps its current name, so the name does not promise a neutrality the code does not have yet.

## Consequences

Positive: once renamed, the project's shared client does not favor one hub implementation.

Negative: until the rename, the shared client carries one hub implementation's name. The rename changes the Go module path, so every importing repository, including ones outside gemaraproj, has to migrate together.

## Alternatives Considered

- **Rename now** — rejected: the name would describe a neutral hub client while the code speaks only grc-store's protocol, and importers would migrate for no change in behavior.
- **Keep the grc-store name permanently** — rejected: it ties the project's one shared client to a single hub, and another hub would need its own client kit, which works against ADR 0024's single implementation.
