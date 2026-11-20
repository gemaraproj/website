---
layout: page
title: Go SDK
description: "Reference for go-gemara, the Go SDK for loading, validating, converting and bundling Gemara artifacts."
---

The Go SDK provides type-safe APIs for reading, writing, converting, and bundling Gemara documents. Types are generated from the CUE schemas with `cue exp gengotypes`, then tidied by the SDK's `typestagger` tool.

**[Go Package Reference →](https://pkg.go.dev/github.com/gemaraproj/go-gemara)**

## Installation

```bash
go get github.com/gemaraproj/go-gemara
```

## Loading Documents

`gemara.Load` is a generic loader — substitute the type parameter for any document kind (`gemara.GuidanceCatalog`, `gemara.ControlCatalog`, `gemara.Policy`, `gemara.EvaluationLog`, etc.). Format (YAML or JSON) is detected from the file extension.

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/gemaraproj/go-gemara"
    "github.com/gemaraproj/go-gemara/fetcher"
)

func main() {
    f := &fetcher.File{}
    ctx := context.Background()

    catalog, err := gemara.Load[gemara.ControlCatalog](ctx, f, "path/to/controls.yaml")
    if err != nil {
        log.Fatal(err)
    }

    for _, control := range catalog.Controls {
        fmt.Printf("Control: %s - %s\n", control.Id, control.Title)
    }
}
```

The `fetcher` package provides three implementations of the `gemara.Fetcher` interface:

- `fetcher.File` — reads from the local filesystem
- `fetcher.HTTP` — fetches over HTTP(S)
- `fetcher.URI` — dispatches by scheme (`file://`, `http(s)://`, or bare paths)

`ControlCatalog` and `GuidanceCatalog` also expose `LoadFiles` for merging multiple sources, and `ControlCatalog.LoadNestedCatalog` for YAML files where the catalog is wrapped in a single key (e.g. `catalog:`).

## Converting to OSCAL

The `gemaraconv` package converts three artifact types to OSCAL:

| Wrapper | Method | Returns |
|:--|:--|:--|
| `gemaraconv.ControlCatalog(c)` | `ToOSCAL(opts...)` | `oscal.Catalog` |
| `gemaraconv.GuidanceCatalog(g)` | `ToOSCAL(catalogHref, opts...)` | `oscal.Catalog`, `oscal.Profile` |
| `gemaraconv.EvaluationLog(l)` | `ToOSCALAssessmentResults(opts...)` | `oscal.AssessmentResults` |

The wrappers take values, so dereference the pointer `gemara.Load` returns. Control and Guidance Catalogs can also be rendered to Markdown with `ToMarkdown(ctx)`. For a complete program, see [Export to OSCAL](../how-to/export-oscal/).

## Converting to SARIF

`EvaluationLog` results can be emitted as SARIF for surfacing in code-scanning tools. A `ControlCatalog` is required to resolve control metadata referenced from the log.

```go
catalog, _ := gemara.Load[gemara.ControlCatalog](ctx, f, "path/to/catalog.yaml")

evaluationLog := gemara.EvaluationLog{ /* populate */ }
sarifBytes, _ := gemaraconv.EvaluationLog(evaluationLog).ToSARIF(
    gemaraconv.WithArtifactURI("path/to/artifact.md"),
    gemaraconv.WithCatalog(catalog),
)
```

## Bundling and Distributing via OCI

The `bundle` package defines Gemara's OCI artifact format (`application/vnd.gemara.bundle.v1`). `bundle.NewAssembler(fetcher).Assemble` resolves a document's full dependency tree (`extends` and `imports`), `bundle.Pack` writes the result to any ORAS target, and `bundle.Unpack` reads it back. For the steps, see [Publish artifacts as an OCI bundle](../how-to/publish-bundle/).

## Developer Tooling

The SDK repository also includes two binaries — `oscalexport` and `typestagger` — that exist as test and development infrastructure for SDK maintainers, not as supported end-user CLIs. `oscalexport` is used to regenerate OSCAL fixtures from the bundled test catalogs (`make oscal-export`), and `typestagger` post-processes generated Go types after `make generate`. See the [`go-gemara` repository](https://github.com/gemaraproj/go-gemara) for usage in a development workflow.

## Relationship to Other Components

### [The Model](../model)
Provides the conceptual foundation. Go SDK types correspond to elements in the model.

### [The Schemas](../schema/)
Go SDK types are generated from the CUE schemas, ensuring consistency between validation and programmatic access. The schema version supported by a given SDK release is exposed as `gemara.SchemaVersion`.
