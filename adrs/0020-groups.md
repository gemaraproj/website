---
layout: page
title: Merge Family and Category into Group
---

- **ADR:** 0020
- **Proposal Author(s):** @eddie-knight, @jpower432
- **Status:** Proposed

## Context

Across the schema and documentation we currently use multiple terms to describe the same underlying concept: **a named classification used to organize entries within an artifact**.

Examples:

- Control catalogs use **families** (`families` list; each control references `family`).
- Risk catalogs use **categories** (`categories` list; each risk references `category`).
- Metadata already defines a generic `#Group` type (used by `metadata.applicability-categories`), but the terminology still varies in YAML and in docs.

This terminology split increases cognitive load and invites inconsistent modeling:

- Readers must learn whether a given artifact uses “family” vs “category” despite the same structural role (an in-document grouping keyed by `id`).
- Tooling must special-case field names for what is effectively the same relationship (“entry belongs to a group defined in the same artifact”).
- ADR-0018’s rejection feedback reinforced that these organizational groupings should generally live **within the documents that use them**, rather than as standalone catalogs (see [ADR 0018](./0018-promote-nested-concepts-to-catalogs.html)).

## Decision

Adopt **Group** as the single, consistent term for organizational groupings across the schema and documentation.

### Schema and field name changes

Wherever “family” or “category” is used to mean “an organizational grouping with an `id`, `title`, and `description`”, rename the fields to **group**:

- **Control and guidance grouping**
  - Rename top-level `families` → `groups`
  - Rename per-entry `family` → `group`
- **Risk grouping**
  - Rename top-level `categories` (risk catalog groupings) → `groups`
  - Rename per-entry `category` → `group`
- **Metadata applicability grouping**
  - Rename `metadata.applicability-categories` → `metadata.applicability-groups`

The grouping entry shape remains `#Group` (already defined in `base.cue`), and the semantics remain unchanged.

## Consequences

- Schema becomes easier to understand and document: one concept, one name.
- Tooling can treat “belongs-to-group” uniformly across artifact types.
- This is a breaking change for serialized YAML/JSON field names and requires a migration step for existing examples and test data.

### Keep “family” and “category” per artifact

Retain existing per-domain vocabulary (controls have families; risks have categories). This preserves current YAML shapes, but keeps conceptual duplication and complicates docs/tooling.

### Introduce “group” only as a shared type, not a field name

Standardize on `#Group` for the entry shape but keep `families`/`categories` as field names. This reduces some implementation duplication but does not address user-facing inconsistency in schema usage and documentation.
