---
layout: page
title: Distribution Formats Ship with the Spec
---

- **ADR:** 0026
- **Proposal Author(s):** @jpower432
- **Status:** Accepted

## Context

ADR 0023 puts the bundle contract and the resolve, lock, and analysis formats in the spec repository, but leaves open whether they get their own module and release cycle. The first of these, the OCI bundle format in [gemaraproj/gemara#648](https://github.com/gemaraproj/gemara/pull/648), is in review. The others are not defined yet.

These formats are consumed by tools that move faster than the spec: go-gemara, grcli, grc-store-clientkit, gemara-publish-action, and registry backends. The spec releases on its own schedule as one CUE module, `github.com/gemaraproj/gemara@v1`.

A distribution format can carry its own compatibility version. The OCI bundle format puts it in its media types (`application/vnd.gemara.bundle.v1` and related types), and any change a version 1 consumer cannot handle requires new media types. Spec releases do not change that version.

CUE allows several modules in one repository, but it has no documented tag convention for modules in subdirectories, and it recommends one module per repository for most projects. The bundle schema imports nothing from the main schemas, so moving it later is mechanical.

## Action

Distribution formats live under `distribution/` in the spec repository, one directory per format:

- `distribution/bundle/`: the OCI bundle format
- `distribution/resolve/`: dependency resolution
- `distribution/lock/`: lock files
- `distribution/analysis/`: analysis output

They ship in the main spec module and release with it. Each format matures through the existing schema status lifecycle (`experimental`, `stable`, `deprecated`), independent of the schemas around it.

Distribution formats move to their own module and release cycle when either of these happens:

- a distribution fix or conformance case is ready but is blocked waiting for a spec release, or
- a tool has to adopt a newer spec release only to pick up a distribution format change.

If that happens, all of `distribution/` moves as one module, not one module per format.

## Consequences

Positive: one module, one tag scheme, and one release pipeline; related formats share one home with room to grow; the split criteria are stated up front instead of argued case by case.

Negative: a distribution format fix waits for the next spec release; tools see spec version bumps that do not change any distribution format.

Neutral: where a format defines its own compatibility version, as the bundle's media types do, that version governs compatibility, not the spec version.

## Alternatives Considered

- **Own module and release cycle now** — rejected: adds a second pipeline, an invented tag scheme, and another version number, when the bundle's media types already separate compatibility from spec releases.
- **`packaging/` as the umbrella** — rejected: resolve, lock, and analysis formats are not packaging.
- **`package/` as a directory name** — rejected: `package` is a keyword in CUE and Go, so neither a CUE package nor a mirrored Go package could use the name.
