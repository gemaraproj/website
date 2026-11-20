---
layout: page
title: Versioning and lifecycle
description: "How the Gemara specification and SDKs are versioned, and how schemas move from experimental to stable to deprecated."
---

Project release deliverables are divided into the **Core Specification** and language-specific **SDKs**.

* **Core Specification Release:** The CUE schemas are versioned and released as the core specification. These schemas implement the Model, which is published separately.
* **SDK Releases:** Language-specific implementations that provide tooling, types, and helpers to work with Gemara documents. SDK types are generated from the CUE schemas.

Each **maintains its own** independent [SemVer](https://semver.org/) lifecycle.

## Specification Release Versioning

The core specification release versions the CUE schemas as a single unit.

| Change Type | Version Bump | Examples                                                     |
|:------------|:-------------|:-------------------------------------------------------------|
| Major       | v2.0.0       | Breaking changes to Stable schemas.                         |
| Minor       | v1.(x+1).0   | Additive changes, schema promotions, or new optional fields. |
| Patch       | v1.x.(y+1)   | Bug fixes in schema logic or documentation.                  |

## Schema Lifecycle and Major Version Example

Possible schema states include: **Experimental** → **Stable** → **Deprecated**. Each schema file declares its state with a `@gemara(status="...")` attribute, for example `@gemara(status="stable")`.

The following table illustrates how schemas progress through their lifecycle and how major version changes are handled:

| Version | Status              | Change Type | Example Scenario                        |
|:--------|:--------------------|:------------|:----------------------------------------|
| v1.0.0  | Experimental        | Initial     | Schema first published                  |
| v1.1.0  | Experimental        | Minor       | Optional fields added                   |
| v1.2.0  | Stable              | Minor       | Promoted to Stable                      |
| v1.3.0  | Stable              | Minor       | Additive changes                        |
| v1.3.1  | Stable              | Patch       | Bug fix                                 |
| v1.4.0  | Stable              | Minor       | Field deprecated, replacement added     |
| v1.5.0  | Stable → Deprecated | Minor       | Original deprecated, replacement Stable |
| v2.0.0  | Stable              | Major       | Deprecated schema removed               |
| v2.1.0  | Stable              | Minor       | Additive changes                        |

## Experimental Status
* All new schemas start as Experimental.
* Adding Experimental schemas or making breaking changes triggers minor version increments.
* Breaking changes and performance issues may occur during this phase.
* These schemas may not be feature-complete.

## Stable Status

* Schemas promote independently and at different times (e.g., ControlCatalog can be Stable while Policy remains Experimental).
* When a schema is promoted to Stable, all types it references must also be Stable (e.g., promoting `#ControlCatalog` requires `#Catalog`, `#Metadata`, `#Group`, and any other referenced types to already be Stable).
* Stable schemas may only reference other Stable types.
* Stable schemas maintain backward and forward compatibility within major versions, allowing **additive optional changes** only.
* Breaking changes require major version increments and should be avoided in all normal circumstances.

## Deprecated Status

* Schemas or fields within schemas may be deprecated when replaced.
* Replacement schema or field must be added in Experimental status and promoted to Stable before deprecation.
* Deprecated schemas and fields maintain the same support guarantees as Stable schemas and remain functional.
* Deprecated schemas are removed in the next major version release.

## SDK Release Versioning

SDK releases are versioned independently to allow for rapid iteration on tooling, bug fixes, and utilities without requiring core specification version bumps.

* The CUE schemas in the core specification release are the source of truth for validation.
* SDKs must explicitly document which core specification release version they support.
* When a new core specification release is published, SDKs regenerate their types from the updated schemas and release new versions that support the updated specification.

The spec version each SDK supports is listed on the [SDKs page](../../sdk/#supported-specification-versions).
