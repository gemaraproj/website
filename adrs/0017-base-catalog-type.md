---
layout: page
title: Base Catalog Type
---

- **ADR:** 0017
- **Proposal Author(s):** @eddie-knight, @jpower432
- **Status:** Accepted

## Context

[ADR-0012](./0012-schema-types) established "Catalog" as the term for artifacts containing multiple definitions. Historically, each catalog-like artifact in the schema was represented with its **own standalone struct**, re-declaring the same core fields across multiple types and files.

[ADR-0016](./0016-base-log-type.html) establishes a common CUE type for all "Log" type artifacts.

## Decision

Introduce a `#Catalog` base definition in `catalog.cue` containing the fields common to all catalog artifacts:

| **Field** | **Type** | **Purpose** |
|:---|:---|:---|
| `title` | `string` | Short human-readable purpose/name |
| `metadata` | `#Metadata` | Standard artifact metadata |
| `extends?` | `...#ArtifactMapping` | References to catalogs this catalog builds upon |
| `imports?` | `[#MultiEntryMapping, ...#MultiEntryMapping]` | Pull entries from other artifacts |

Concrete catalog types embed `#Catalog` and add their specific entries and other essential fields. For example:

```cue
#ThreatCatalog: {
	#Catalog
	threats?: [...#Threat]
}
```

## Consequences

- Catalog types no longer redeclare core catalog fields; changes to shared catalog fields propagate from one definition
- New catalog types require only embedding `#Catalog` plus their own typed `entries`
- No breaking change to the serialized YAML/JSON shape for existing catalogs — embedded fields are structurally identical to inline fields in CUE

## Alternatives Considered

### Keep fields inline on each catalog type

Each catalog type declares its own `title`, `metadata`, `extends`, and `imports`. Simpler per-type definitions in isolation, but violates single-source-of-truth as the catalog family grows per ADR-0012 and makes consistent evolution harder.

