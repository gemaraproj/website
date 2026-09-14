---
layout: page
title: Capability Extensions Guide
description: Step-by-step guide to extending Gemara with CUE
---

## What This Is

This guide shows how to write a **Capability Extension**: a domain-specific extension of Gemara's `#Capability` as its own CUE type.

**The basic idea:** `#Capability` is deliberately neutral — `id`, `title`, `description`, `group`. It stays that way. If your domain needs more, you define those fields in your own CUE module and unify them with `gemara.#Capability` rather than asking Gemara to add them to core.

This is the **extend** direction. If instead you need to *restrict* the base type, a profile or baseline that tightens what values are allowed without adding fields. See [Extending vs Constraining Gemara](extending-vs-constraining), which explains when to reach for each and how they differ.

## Walkthrough

### Step 1: Create Your Extension Module

Extensions live in their own CUE module. Create a directory and initialize it:

```sh
mkdir acme-k8s && cd acme-k8s
cue mod init acme.example/k8s@v0
```

Add Gemara as a dependency in `cue.mod/module.cue`, then let CUE resolve it:

```cue
module: "acme.example/k8s@v0"
language: version: "v0.17.0"
deps: {
	"github.com/gemaraproj/gemara@v1": v: "v1.2.0"
}
```

```sh
cue mod tidy
```

### Step 2: Define Your Extension Type

Define a new type embedding `#Capability`:

```cue
package kubernetes

import gemara "github.com/gemaraproj/gemara@v1"

#KubernetesCapability: {
    gemara.#Capability
    "api-group"?:   string
    "api-resource": string
    verb:           "get" | "list" | "watch" | "create" | "update" | "patch" | "delete" | "deletecollection"
    namespaced:     bool
}

#KubernetesCapabilityCatalog: {
    gemara.#CapabilityCatalog
    capabilities: [#KubernetesCapability, ...#KubernetesCapability]
}
```

`gemara.#Capability` is embedded in `#KubernetesCapability` and unified with that struct so it has every `Capability` base field plus the additions. Because it keeps every base field, an extended catalog is a **superset** of a base catalog (see [the concept guide](extending-vs-constraining) for what that buys you).

### Step 3: Author Your Catalog

Write a `CapabilityCatalog` document as you normally would, adding your extension's fields alongside the base ones:

```yaml
title: Acme Kubernetes Capability Catalog

metadata:
  id: ACME-K8S
  type: CapabilityCatalog
  gemara-version: "1.2.0"
  version: "0.1.0"
  description: Kubernetes workload capabilities for the Acme platform.
  author:
    id: acme-platform
    name: Acme Platform Team
    type: Human

groups:
  - id: workloads
    title: Workloads
    description: Capabilities related to pod and deployment management.

capabilities:
  - id: CAP-K8S-001
    title: Create deployments
    description: The cluster can create deployment resources in application namespaces.
    group: workloads
    api-group: apps
    api-resource: deployments
    verb: create
    namespaced: true
```

See [`capability-extension-example.yaml`](capability-extension-example.yaml) for the full file.

### Step 4: Validate Against Your Extension

From inside your module directory, vet the document against your extension definition:

```sh
cue vet -c -d '#KubernetesCapabilityCatalog' . capability-extension-example.yaml
```

Validate against the **extension** definition you authored, not base `#CapabilityCatalog`. Vetting an extended document against base fails — base `#Capability` is a closed definition, so the extra fields are rejected (`field not allowed`). That is expected: the contract your document conforms to is your extension type. Base consumers still read the document fine (see below).

## Consuming an Extended Catalog in Go

Gemara generates Go types from its CUE schema, and your extension module does the same. This is what makes an extension a superset rather than a fork: the same document is readable by both extension-aware and generic base tooling.

**Extension-aware consumer** — embeds the generated base struct and adds the typed fields:

```go
type KubernetesCapability struct {
    gemara.Capability          // embedded base type
    APIGroup    *string `json:"api-group,omitempty"`
    APIResource string  `json:"api-resource"`
    Verb        string  `json:"verb"`
    Namespaced  bool    `json:"namespaced"`
}
```

Unmarshalling populates both the base fields (via the embed) and your extensions.

**Generic base consumer** — a tool that only knows `gemara.Capability` unmarshals the same document and reads the base fields. Go's `encoding/json` silently ignores unknown fields by default, so the extension attributes are dropped, no error:

```go
var c gemara.Capability
json.Unmarshal(doc, &c) // id, title, description, group populated; extras ignored
```

**One caveat:** a consumer that opts into strict decoding with `Decoder.DisallowUnknownFields()` **will** error on the extension fields (`json: unknown field "api-group"`). If you distribute a base-only tool meant to read extended catalogs, do not enable strict unknown-field rejection.

## Publishing and Discovery

Publish your extension module like any CUE module (to the public CUE registry or your own) and open a PR on [`awesome-gemara`](https://github.com/gemaraproj/awesome-gemara) to add it to the index. Today that index is how extensions are shared and discovered; see [Extending vs Constraining Gemara](extending-vs-constraining) for how sanctioned support is expected to evolve.

## Have Ideas?

- Reach out via Slack in `#gemara`
- Discuss in one of our bi-weekly meetings on the [OpenSSF calendar](https://calendar.google.com/calendar/u/0?cid=czYzdm9lZmhwNWk5cGZsdGI1cTY3bmdwZXNAZ3JvdXAuY2FsZW5kYXIuZ29vZ2xlLmNvbQ)
- Open a [GitHub Issue](https://github.com/gemaraproj/gemara/issues)
