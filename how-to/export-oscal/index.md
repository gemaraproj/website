---
layout: page
title: Export to OSCAL
description: "Convert a Gemara Control Catalog, Guidance Catalog or Evaluation Log to OSCAL JSON with the Go SDK."
---

This guide converts Gemara artifacts to [OSCAL](https://pages.nist.gov/OSCAL/) JSON with the [Go SDK](../../sdk/go-sdk.html). Use it when you author in Gemara but need to hand OSCAL to a regulator, an assessor or an OSCAL-based tool.

The SDK converts three artifact types:

| Gemara artifact | OSCAL output |
|:--|:--|
| Control Catalog | Catalog |
| Guidance Catalog | Catalog and Profile |
| Evaluation Log | Assessment Results |

Other artifact types have no OSCAL conversion yet; they are tracked in [go-gemara#29](https://github.com/gemaraproj/go-gemara/issues/29). The SDK takes care of OSCAL-specific details such as UUIDs, so your source YAML stays as it is.

## Before you start

- Go installed.
- A Gemara artifact that passes `cue vet`, for example the [control catalog](../control-catalog/control-catalog.yaml) from the control catalog guide.
- A Go module with the SDK and the OSCAL types it returns:

```bash
go mod init example.com/oscal-export
go get github.com/gemaraproj/go-gemara@v0.12.0 github.com/defenseunicorns/go-oscal@v0.7.1
```

## Convert a Control Catalog

Load the catalog, convert it, and write it out. An OSCAL JSON document wraps the model in a top-level key (`catalog`, `profile` or `assessment-results`), which `oscal.OscalModels` adds for you.

```go
package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	oscal "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"
	"github.com/gemaraproj/go-gemara"
	"github.com/gemaraproj/go-gemara/fetcher"
	"github.com/gemaraproj/go-gemara/gemaraconv"
)

func main() {
	ctx := context.Background()

	catalog, err := gemara.Load[gemara.ControlCatalog](ctx, &fetcher.File{}, "control-catalog.yaml")
	if err != nil {
		log.Fatal(err)
	}

	// gemara.Load returns a pointer; the converters take a value.
	oscalCatalog, err := gemaraconv.ControlCatalog(*catalog).ToOSCAL()
	if err != nil {
		log.Fatal(err)
	}

	out, err := json.MarshalIndent(oscal.OscalModels{Catalog: &oscalCatalog}, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile("catalog.json", out, 0o644); err != nil {
		log.Fatal(err)
	}
}
```

Run it with `go run .`. The result is `catalog.json`.

## Convert a Guidance Catalog

A Guidance Catalog becomes two OSCAL documents: a catalog and a profile that imports it. Pass the path the profile should use to reach the catalog file.

```go
guidance, err := gemara.Load[gemara.GuidanceCatalog](ctx, &fetcher.File{}, "guidance-example.yaml")
if err != nil {
	log.Fatal(err)
}

oscalCatalog, profile, err := gemaraconv.GuidanceCatalog(*guidance).ToOSCAL("guidance-catalog.json")
if err != nil {
	log.Fatal(err)
}
// Write oscal.OscalModels{Catalog: &oscalCatalog} to guidance-catalog.json
// and oscal.OscalModels{Profile: &profile} to guidance-profile.json.
```

## Convert an Evaluation Log

An Evaluation Log becomes OSCAL Assessment Results. Passing the Control Catalog the log was evaluated against, with `WithCatalog`, lets the converter fill in control details.

```go
f := &fetcher.File{}

catalog, err := gemara.Load[gemara.ControlCatalog](ctx, f, "control-catalog.yaml")
if err != nil {
	log.Fatal(err)
}
evalLog, err := gemara.Load[gemara.EvaluationLog](ctx, f, "evaluation-log.yaml")
if err != nil {
	log.Fatal(err)
}

results, err := gemaraconv.EvaluationLog(*evalLog).ToOSCALAssessmentResults(
	gemaraconv.WithCatalog(catalog),
)
if err != nil {
	log.Fatal(err)
}
// Write oscal.OscalModels{AssessmentResults: &results} to assessment-results.json.
```

## Related

- [Go SDK reference](../../sdk/go-sdk.html) for loading, the other converters and options.
- [Gemara and OSCAL](../../faq/#gemara-and-oscal) in the FAQ, on when to use each.
