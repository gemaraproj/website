---
layout: page
title: SDKs
description: "Gemara's Go, Python and React SDKs, the spec version each one supports, and how SDK releases relate to spec releases."
---

Language-specific SDKs provide programmatic access to Gemara documents. SDK types are generated from CUE schemas, ensuring consistency between validation and programmatic access.

## Available SDKs

- [Go SDK](./go-sdk.html) — Type-safe Go APIs for reading, writing, and manipulating Gemara documents
- [Python SDK](./python-sdk.html) — Typed Pydantic v2 models and loader utilities for Gemara documents
- [React SDK](./react-sdk.html) — Headless React renderers and UI components for Gemara documents

## Supported specification versions

| SDK | Latest release | Generated from spec |
|:--|:--|:--|
| Go (`go-gemara`) | v0.12.0 | v1.5.0 |
| Python (`gemara-python`) | 0.1.0 (pre-release) | v1.5.0 |
| React (`@gemara/react`) | 0.2.2 | v1.3.0 |

The latest specification is v1.6.0. Check each SDK's release notes for newer versions.

## Versioning and Maintenance

The specification and each SDK keep separate [SemVer](https://semver.org/) lifecycles. The policy for both is on the [Schemas page](../schema/#versioning-and-maintenance).

### SDK Release Versioning

SDK releases are versioned independently to allow for rapid iteration on tooling, bug fixes, and utilities without requiring core specification version bumps.

* The CUE schemas in the core specification release are the source of truth for validation.
* SDKs must explicitly document which core specification release version they support.
* When a new core specification release is published, SDKs regenerate their types from the updated schemas and release new versions that support the updated specification.

## Relationship to Other Components

### [The Model](../model)
Provides the conceptual foundation. SDK types correspond to elements in the model.

### [The Schemas](../schema/)
SDK types are generated from the CUE schemas. The schemas are the source of truth for data structures.
