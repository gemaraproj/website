---
layout: page
title: Ecosystem
nav-title: Ecosystem
---

Gemara is a small ecosystem of repositories, each with a deliberate scope. This page records what belongs where, so contributors know where a change should go and reviews stay focused.

| Repository | Scope |
| --- | --- |
| **[gemara](https://github.com/gemaraproj/gemara)** (spec) | Declarative definitions: the model, CUE schemas, artifact shapes, the bundle contract, and the resolve/lock/analysis formats. Deliberately slow-moving; changes are controlled and reviewed carefully. The bundle contract lives in its own module with its own release pipeline. |
| **SDKs** (Go, Python, React) | Deliberately kept small. The exact scope of the SDKs will be defined by a future spec. |
| **hub-clientkit** | Publishing, consuming, and interacting with registries: bundle management, trust and signing wiring shared by tools that talk to a hub. Types are generated from the spec, not hand-written. |
| **[grcli](https://github.com/gemaraproj/grcli)** | The workbench CLI. Executes and interprets what the spec declares. |
| **[gemara-ai](https://github.com/gemaraproj/gemara-ai)** | AI skills, OpenPackage bundles, and MCP configurations for agents working with Gemara artifacts. |

A few boundaries worth calling out:

- If it's a shape or a format that clients and the hub must agree on, it goes in the spec — not in an SDK, where it would become "one SDK's opinion."
- If it's shared behavior for publishing or consumption, it goes in hub-clientkit, not the SDKs.
- Skills and agent tooling install and use grcli as their deterministic tool; grcli does not manage skills.
