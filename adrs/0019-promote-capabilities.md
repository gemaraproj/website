---
layout: page
title: Promote Capabilities to a First-class Catalog
---

- **ADR:** 0019
- **Proposal Author(s):** @eddie-knight, @jpower432
- **Status:** Accepted

## Context

Threats and other artifacts often need to reference system “capabilities” (features, components, or objects) as domain content.

Historically, capabilities were treated as nested/ad hoc lists inside other documents, which made them hard to:

- Reference consistently across artifacts
- Extend/import and version as shared domain content
- Maintain as their own area of ownership (distinct from threats/controls)

During review of [ADR 0018](./0018-promote-nested-concepts-to-catalogs.html), we agreed that **Capabilities should stand alone** as domain content, while organizational groupings (e.g., families) should generally live in the documents that use them.

## Decision

Define Capabilities as a **first-class catalog artifact** with a dedicated catalog type, `#CapabilityCatalog`.

### Schema shape

`#CapabilityCatalog` embeds the shared `#Catalog` base and defines a typed list of capability entries.

In CUE, the shape is:

```cue
#CapabilityCatalog: {
    #Catalog
    capabilities?: [...#Capability]
}

#Capability: {
    id:          string
    title:       string
    description: string
}
```

Threats (and other artifacts) reference capabilities using existing mapping mechanisms (e.g., `#MultiEntryMapping`) rather than duplicating capability definitions inline.

## Consequences

- Capabilities become reusable domain content that can be authored, versioned, imported, and extended independently of threat catalogs.
- Tooling can validate capability references uniformly (by id) and avoid special-casing nested capability lists.
- This introduces (or formalizes) an additional artifact type in the ecosystem (`CapabilityCatalog`) that producers/consumers may need to support.

## Migration plan

- Prefer capability references from threats/other catalogs via mappings to a `CapabilityCatalog`.
- Update examples, fixtures, and documentation that currently define capabilities inline to instead:
  - Define capabilities in a `CapabilityCatalog`
  - Reference them from `Threat.capabilities` using mappings

## Alternatives Considered

### Keep capabilities nested within threat catalogs

Continue defining capabilities inline in the same document as threats. Simpler for single-document use, but prevents independent lifecycle management and reuse of capability definitions across documents and teams.

