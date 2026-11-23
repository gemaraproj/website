---
layout: page
title: Publishing and Consumption Live in the Client Kit
---

- **ADR:** 0024
- **Proposal Author(s):** @jpower432
- **Status:** Accepted

## Context

Some consumers need to publish bundles and interact with registries from Go code, as a library, without shelling out to a CLI. That behavior cannot go in the SDKs without growing them, and duplicating it per consumer guarantees drift.

## Action

Publishing, consumption, and hub/registry interaction — including bundle management and trust/signing wiring — live in the client kit. Its types are generated from the spec, not hand-written.

## Consequences

Positive: one Go implementation of registry interaction, shared by every consumer that needs it; SDKs stay small.

Negative: only serves Go consumers directly; other languages need their own path or the CLI.

## Alternatives Considered

- **Put it in the SDKs** — rejected: grows every SDK and multiplies maintenance across languages.
- **Shell out to grcli** — rejected: unsuitable as a library interface for embedded use.
