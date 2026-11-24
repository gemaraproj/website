---
layout: page
title: Publishing and Consumption Live in the Client Kit
---

- **ADR:** 0024
- **Proposal Author(s):** @jpower432
- **Status:** Accepted

## Context

Some consumers need to publish bundles and interact with registries from Go code, as a library, without shelling out to a CLI. That behavior cannot go in the SDKs without growing them, and duplicating it per consumer guarantees drift. It already had: grcli and privateer-sdk each grew near-identical copies of registry and trust code, and `grc-store-clientkit` was extracted from them.

## Action

Publishing, consumption, and hub/registry interaction — including bundle management and trust/signing wiring — live in [`grc-store-clientkit`](https://github.com/gemaraproj/grc-store-clientkit). The client kit builds on `go-gemara/bundle`, which implements the bundle format defined in the spec (ADR 0026). Its wire shapes come from `grc-store-protocol`, which is maintained outside the project. Whether those shapes move under the spec (ADR 0023) is undecided.

The client kit is named for grc-store, one hub implementation, because today it speaks only that hub's protocol. It will be renamed `hub-clientkit` once it speaks a hub contract the project owns. Until then it keeps its current name, so the name does not promise a neutrality the code does not have yet.

## Consequences

Positive: one Go implementation of registry interaction, shared by every consumer that needs it; SDKs stay small.

Negative: only serves Go consumers directly; other languages need their own path or the CLI. Until the rename, the project's shared client library carries one hub implementation's name. The rename changes the Go module path, so every importing repository, including ones outside gemaraproj, has to migrate together.

## Alternatives Considered

- **Put it in the SDKs** — rejected: grows every SDK and multiplies maintenance across languages.
- **Shell out to grcli** — rejected: unsuitable as a library interface for embedded use.
- **Rename to `hub-clientkit` now** — rejected: the name would describe a neutral hub client while the code speaks only grc-store's protocol, and the importers would migrate for no change in behavior.
- **Keep the grc-store name permanently** — rejected: it ties the project's one shared client to a single hub, and another hub would need its own client kit.
