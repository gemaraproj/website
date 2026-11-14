---
layout: page
title: Roadmap
---

This document outlines the strategic priorities ("big rocks") for Gemara in 2027 and beyond. This roadmap is a living document and is subject to change based on community feedback and evolving industry requirements.
These are high-level strategic initiatives, not actionable tasks. They represent major directions and aspirations for the project.

These priorities build on the direction agreed upon in the Gemara Roadmapping 2026 Workshop, carrying forward the work still in progress and adding new areas of investment.
We welcome input and participation in future roadmapping discussions—join our [biweekly meetings](https://github.com/gemaraproj/gemara/blob/main/README.md#contributing) to help shape the project's direction.

**Live tracking:** Day-to-day progress against these priorities — including the issues, PRs, and discussions driving each one — is tracked on the [Gemara Roadmap project board](https://github.com/orgs/gemaraproj/projects/1). Browse it to see what's in flight, what's queued up, and where new contributors can plug in.

## 2027 Goal

Make Gemara practical to author, validate, and exchange, on top of the stable foundation built in 2026: tooling that behaves the same everywhere, a reference CLI, and AI help that produces idiomatic artifacts.

### Strategic Priorities

* **Stabilize Schemas** (continued): Promote Policy ([#529](https://github.com/gemaraproj/gemara/issues/529)), MappingDocument ([#610](https://github.com/gemaraproj/gemara/issues/610)), and GuidanceCatalog ([#611](https://github.com/gemaraproj/gemara/issues/611)) to stable, and resolve the known issues blocking RiskCatalog, EnforcementLog, and AuditLog. *Why: Adopters and SDKs can only promise stability as far as the schemas do.*

* **Mature SDKs** (continued): Publish [tool specifications](./specifications/) that define how implementations resolve, validate, and exchange Gemara artifacts, and bring the Go, Python, and React SDKs into conformance. *Why: A written contract keeps SDKs in different languages consistent and lets third-party tools interoperate.*

* **Invest in grcli**: Grow [`grcli`](https://github.com/gemaraproj/grcli) into the reference command-line tool for validating, packaging, signing, publishing, and verifying Gemara artifacts. *Why: One tool for the full artifact lifecycle lowers the barrier for authors and CI pipelines, and puts the tool specifications to work.*

* **Define Idiomatic Gemara**: Document the conventions for modeling artifacts well, beyond what the schemas enforce, and create AI skills that apply them. *Why: Schemas say what is valid, not what is good. Consistent artifacts are easier to map, compare, and reuse.*

* **Grow AI Skills, Keep MCP Lean**: Move authoring prompts into Gemara's [AI skills](https://github.com/gemaraproj/gemara-ai) and grow their compatibility with `grcli` and the [MCP server](https://github.com/gemaraproj/gemara-mcp). Keep the MCP server focused on serving Gemara's overview, schemas, and lexicon as resources and on validating artifacts. *Why: Skills carry workflow guidance to any agent, while a focused MCP server gives models current, versioned reference material and stays easy to maintain.*

* **Dogfood Gemara**: Use Gemara artifacts and tooling to run the project's own governance, and publish the results as reference examples. *Why: Using our own model finds gaps before adopters do, and gives newcomers real examples to learn from.*

## Completed Milestones

### 2026

* Published the Gemara model whitepaper as [The Model](./model/)
* Contributed Gemara definitions to the OpenSSF Glossary
* Launched [`awesome-gemara`](https://github.com/gemaraproj/awesome-gemara) to centralize community resources
* Expanded the SDKs to Python and React alongside Go
* Released the Gemara MCP server
* Started the [tool specifications](./specifications/)

### 2025

* Establishment of the layered data approach (Layers 1-4)
* Adoption by OpenSSF OSPS Baseline and FINOS CCC

## Contributing

If you have suggestions for the roadmap:

* Open an [Issue](https://github.com/gemaraproj/gemara/issues/new/choose) to start a discussion
* Join our biweekly meetings to discuss roadmap changes (see the Community [page](./community) for details)

Changes to strategic priorities require community consensus through discussion in meetings or issues, not direct edits to this document.
