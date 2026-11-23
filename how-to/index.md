---
layout: page
title: How-to guides
description: "Pick the Gemara how-to guide for your task: threat assessments, catalogs, mappings, policies, OSCAL export and OCI bundles."
---

## Start here

New to Gemara? Work through the tutorial [Write and validate your first control catalog](../tutorials/first-control-catalog/) first, and read [The Gemara Model](../model/) for the seven layers these guides build on.

---

## Find your guide

Pick your goal — each path leads to the right guide.

### Performing a threat assessment

For a system or component → [Threat Assessment Guide](threat-assessment/) — identify capabilities and threats, map them to attack surfaces (Layer 2).

### Defining security controls

That mitigate those threats → [Control Catalog Guide](control-catalog/) — create a control catalog with assessment requirements and threat-mappings (Layer 2).

### Understanding what threats and controls exist

Before writing policy → [Threat Assessment Guide](threat-assessment/) for threats and capabilities, and [Control Catalog Guide](control-catalog/) for threat-informed controls and assessment requirements (Layer 2).

### Reviewing the controls to reference in a policy

→ [Control Catalog Guide](control-catalog/) — control structure, assessment requirements, and threat links (Layer 2).

### Understanding the security posture of consumed software

→ [Threat Assessment Guide](threat-assessment/) — review threat catalogs for your dependencies (Layer 2).

→ **COMING SOON:** Use control catalogs (e.g. OSPS, CCC) as hardening guides (Layer 2).

### Creating a guidance catalog from best practices

From a spreadsheet or checklist — create a guidance catalog (guidelines, groups, mapping-references) that threat-informed controls can reference; express relationships to other frameworks in a [Mapping Document Guide](mapping-document/) (see also the [Mapping Document schema](https://gemara.openssf.org/schema/mappingdocument.html)). → [Guidance Catalog Guide](guidance-catalog/).

### Creating a mapping document between artifacts

Express how entries in a source artifact (e.g., guidance, principle, or control catalog) relate to a target artifact (e.g., regulation or framework) — [Mapping Document Guide](mapping-document/) — typed source/target references (`#TypedMapping`), `source` plus `targets` (`#MappingTarget`), relationship types, and applicability (cross-artifact).

### Creating organizational policy

Create a policy document that translates risk appetite into mandatory rules — [Policy Guide](policy/) — scope, imports, adherence, and risks (Layer 3).

### Creating a risk catalog

When you need a structured inventory of organizational or system risks—**risk categories** (appetite, optional max-severity), per-risk **severity**, optional **`rank`** for ordering within the catalog, optional RACI **owner** and **impact**, and optional **threats** links backed by `metadata.mapping-references`—so policies can reference mitigated or accepted risks → [Risk Catalog Guide](risk-catalog/) (Layer 3).

### Exporting to OSCAL

Convert Control Catalogs, Guidance Catalogs and Evaluation Logs to OSCAL JSON with the Go SDK → [Export to OSCAL](export-oscal/).

### Publishing artifacts

Package an artifact and everything it imports as an OCI bundle and push it to a registry → [Publish artifacts as an OCI bundle](publish-bundle/).

## What You'll Build

| Layer | Artifact | Guide |
|-------|----------|-------|
| **Layer 1** — Guidance | Guidance Catalog (guidelines, groups, mapping-references); [Principle Catalog](https://gemara.openssf.org/schema/principlecatalog.html) (principles, groups) | [Guidance Catalog Guide](guidance-catalog/) |
| **Layer 2** — Controls | Threat Catalog + Control Catalog (assessment requirements, threats) | [Threat Assessment](threat-assessment/), [Control Catalog](control-catalog/) |
| **Layer 3** — Policy   | Policy Document (scope, imports, adherence)                  | [Policy Guide](policy/) |
| **Layer 3** — Risks    | Risk Catalog (risk categories, appetite, risks, optional rank and threat mappings) | [Risk Catalog Guide](risk-catalog/) |
| **Cross-artifact**     | Mapping Document (typed `source`/`target` references, `targets` per mapping, relationship types; entry types per schema include guidelines, controls, Principle, threats, risks, and others) | [Mapping Document Guide](mapping-document/) |

## What You'll Need

- `go` installed
- `cue` installed for validation, for example `cue vet -c -d '#MappingDocument' github.com/gemaraproj/gemara@v1.6.0 your-mapping-document.yaml`

## Have Ideas?

- Reach out via Slack in `#gemara`
- Discuss in one of our bi-weekly meetings on the [OpenSSF calendar](https://calendar.google.com/calendar/u/0?cid=czYzdm9lZmhwNWk5cGZsdGI1cTY3bmdwZXNAZ3JvdXAuY2FsZW5kYXIuZ29vZ2xlLmNvbQ)
- Open a [GitHub Issue](https://github.com/gemaraproj/gemara/issues)
