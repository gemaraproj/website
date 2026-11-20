---
layout: page
title: Promote Nested Concepts to First-class Catalogs
description: "ADR-0018 (rejected): promote families, categories and capabilities from nested lists to first-class catalogs."
---

- **ADR:** 0018
- **Proposal Author(s):** @eddie-knight
- **Status:** Rejected

## Context

Historically, the schema treated some concepts as *sub-structures* of other catalogs rather than as standalone artifacts. In particular, Families, Categories, and Capabilities were represented only as nested lists within other catalogs.

This caused recurring problems:

- These concepts could not be referenced, extended, imported, or versioned as independent artifacts.
- Different contributors often maintained the different asset types
- As the number of catalog types grows, keeping ad hoc nested representations consistent becomes harder and encourages drift between similar concepts.

We want these entries to be composable and version controlled independently of other assets.

## Decision

A proposal was made along with proposed changes, and subsequently rejected. This is the log of that decision.

### Rejected proposal

Promote Groups, Capabilities, and Risk Categories to **first-class catalog artifacts** with dedicated top-level catalog types:

- `#GroupCatalog`
- `#CapabilityCatalog`
- `#RiskCategoryCatalog`

Each new catalog type embeds the shared `#Catalog` base and defines a typed `entries` list for its concept.

### Rejection reason

This proposal was rejected because we do not anticipate strong cross-document reuse (or a desire to maintain additional standalone artifacts) to document group descriptions.

- **Groups (e.g. Family) and Risk Categories** are primarily organizational and should live within the document that uses them.
- **Capabilities** are agreed to be domain content and should have a standalone proposal in a subsequent ADR.

## Consequences

- **No new top-level catalog types** are introduced for Groups and Risk Categories; these organizational groupings remain defined within the documents that use them.
- **Capabilities** as a new catalog will be considered in a follow-up proposal.

## Alternatives Considered

### Keep concepts nested within other catalogs

Continue representing Groups, Capabilities, and Risk Categories only as lists inside other catalogs. This avoids introducing new top-level artifacts, but prevents independent reuse and composition.