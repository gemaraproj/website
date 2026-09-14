---
layout: page
title: Extending vs Constraining Gemara
description: How to specialize Gemara types
---

## Two Ways to Specialize a Gemara Type

Gemara's core types are deliberately neutral. When a domain needs something the base type does not give you, there are two distinct moves depending on the use case.

|                          | **Extend**                                                                        | **Constrain**                                                                   |
|--------------------------|-----------------------------------------------------------------------------------|---------------------------------------------------------------------------------|
| Goal                     | Add domain-specific fields                                                        | Restrict the values the base already allows                                     |
| CUE mechanism            | [Embed](https://cuelang.org/docs/reference/spec/#embedding) the base type (widen) | [Unify](https://cuelang.org/docs/reference/spec/#unification) with `&` (narrow) |
| Fields                   | Base fields **plus** new ones                                                     | Base fields only                                                                |
| Valid against base type? | No — extra fields are rejected by the closed base definition                      | Yes — every constrained document is still a valid base document                 |
| Use for                  | New capability attributes (K8s API resources, kernel capabilities)                | Profiles and baselines that tighten an existing catalog                         |

Pick **extend** when your domain has attributes the base type simply doesn't model. Pick **constrain** when the base type already fits, and you only need to require, restrict, or shape its values.

### Extend

Embed the base type in your own definition and add fields. Full walkthrough — module setup, authoring, validation, and Go consumption — in the [Capability Extensions Guide](capability-extensions-guide):

```cue
#KubernetesCapability: {
    gemara.#Capability
    "api-resource": string
    verb:           "get" | "list" | "watch" | "create" | "update" | "patch" | "delete" | "deletecollection"
    namespaced:     bool
}
```

### Constrain

Unify the base type with tighter rules using `&`. You add no fields, so the result is still a valid instance of the base type. This is how the [OSSF Security Baseline](https://github.com/ossf/security-baseline/tree/main/schema) layers OSPS-specific requirements onto a Gemara `#ControlCatalog`:

```cue
// #OSPSBaseline layers OSPS-specific constraints on top of the Gemara #ControlCatalog.
#OSPSBaseline: gemara.#ControlCatalog & {
    let _agID = =~"^maturity-"

    metadata: "applicability-groups": [{id: _agID}, ...{id: _agID}]
}
```

Every `#OSPSBaseline` document is, by construction, a valid `#ControlCatalog` — the constraint only narrows what is allowed. Generic Gemara tooling reads it with no special knowledge.

## Superset, Not Fork

An extension embeds every base field, so an extended document is a **superset** of a base document, not a divergent one. That distinction matters in practice:

- **Extension-aware tools** read the full typed shape — base fields plus your additions.
- **Generic base tools** read the base fields and ignore the rest. In Go, `encoding/json` drops unknown fields by default, so a base-only consumer reads an extended catalog with no error. (The one exception is strict decoding via `DisallowUnknownFields` — see the [extensions guide](capability-extensions-guide#consuming-an-extended-catalog-in-go).)

So an extension is not a fork of Gemara. It is Gemara plus typed, validated, domain-specific fields — making it shareable and verifiable.

## The Intended Path

Gemara extensions do not yet have: a **dedicated extension point** the core is aware of or tooling support. Gemara does not support a property bag solution because they lack discoverability, are difficult to verify, and create subtle fragmentation.

In Gemara today, the base type does not know that a given document is a maintained extension of it. It only sees a superset it happens to be able to read. This is intentional, for now. The plan is to let the community pattern prove itself before formalizing it:

1. **Today** — write extensions with the embed pattern and publish them to [`awesome-gemara`](https://github.com/gemaraproj/awesome-gemara). That index is the de-facto registry: how extensions are shared, discovered, and versioned.
2. **Next** — as real extensions reach critical mass and their shapes stabilize, Gemara adds formal, core-aware support so an extension can declare itself a sanctioned extension of a base type rather than an undeclared superset.

Building the extension pattern this way — typed, published, and versioned from the start — means today's extensions are already shaped for that future support.

## Have Ideas?

- Reach out via Slack in `#gemara`
- Discuss in one of our bi-weekly meetings on the [OpenSSF calendar](https://calendar.google.com/calendar/u/0?cid=czYzdm9lZmhwNWk5cGZsdGI1cTY3bmdwZXNAZ3JvdXAuY2FsZW5kYXIuZ29vZ2xlLmNvbQ)
- Open a [GitHub Issue](https://github.com/gemaraproj/gemara/issues)
