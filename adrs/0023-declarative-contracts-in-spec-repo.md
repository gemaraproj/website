---
layout: page
title: Declarative Contracts Live in the Spec Repository
---

- **ADR:** 0023
- **Proposal Author(s):** @jpower432
- **Status:** Accepted

## Context

Gemara has grown from a single spec repository into an ecosystem with language SDKs, a client kit, and the grcli CLI. New work — a bundle contract, dependency resolution formats, analysis formats — forced the question of where declarative definitions belong. When shapes live in an implementation, they become that implementation's opinion: the Python SDK cannot consume a Go SDK's types, and vendored copies drift from the source.

## Action

All declarative contracts live in the spec repository: the model, CUE schemas, artifact shapes, the bundle contract, and the resolve/lock/analysis formats. The bundle contract lives in its own module with its own release pipeline. Anything that clients and a hub must agree on is defined here before it is implemented anywhere.

## Consequences

Positive: one authoritative source for every shape; implementations generate or consume rather than invent; the spec remains the slow, controlled repository.

Negative: defining a format requires spec changes before implementation can start.

## Alternatives Considered

- **Define shapes in the Go SDK** — rejected: would not serve other SDKs; no authoritative source, only one SDK's vendored opinion.
