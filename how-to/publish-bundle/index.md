---
layout: page
title: Publish artifacts as an OCI bundle
description: "Package a Gemara artifact and everything it imports as an OCI bundle with the Go SDK, and push it to a container registry."
---

This guide packages a Gemara artifact, such as a Policy, together with every artifact it imports, and pushes the result to an OCI registry. Consumers then pull one versioned, self-contained bundle instead of chasing each import by URL.

The simplest way to publish is still to attach your YAML files to a release in your repository. Use a bundle when consumers need the full dependency tree pinned together.

A bundle is an OCI artifact with type `application/vnd.gemara.bundle.v1`. The [Go SDK](../../sdk/go-sdk.html)'s `bundle` package builds, packs and reads it, and [ORAS](https://oras.land/) moves it to and from registries.

## Before you start

- Go installed, and a Go module with the SDK:

  ```bash
  go get github.com/gemaraproj/go-gemara@v0.12.0
  ```

- An artifact that passes `cue vet`. The example YAML files in these guides won't assemble as published: some mapping-references have no `url`, and others point at web pages rather than YAML files.
- Every artifact it extends or imports listed in `metadata.mapping-references` with a `url` the SDK can fetch (`file://`, `http://`, `https://` or a plain path). The assembler follows those URLs, then the imports of each fetched artifact, until the tree is complete.
- Push access to a registry.

## Assemble, pack and push

```go
package main

import (
	"context"
	"log"
	"os"

	"github.com/gemaraproj/go-gemara/bundle"
	"github.com/gemaraproj/go-gemara/fetcher"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/oci"
	"oras.land/oras-go/v2/registry/remote"
)

func main() {
	ctx := context.Background()

	// 1. Assemble the policy and everything it imports.
	data, err := os.ReadFile("policy.yaml")
	if err != nil {
		log.Fatal(err)
	}
	m := bundle.Manifest{BundleVersion: "1", GemaraVersion: "v1.6.0"}
	b, err := bundle.NewAssembler(&fetcher.URI{}).Assemble(ctx, m, bundle.File{Name: "policy.yaml", Data: data}, bundle.WithContinueOnError())
	if err != nil {
		log.Fatal(err)
	}

	// 2. Pack it into a local OCI layout and tag it.
	store, err := oci.New("./bundle-output")
	if err != nil {
		log.Fatal(err)
	}
	desc, err := bundle.Pack(ctx, store, b)
	if err != nil {
		log.Fatal(err)
	}
	if err := store.Tag(ctx, desc, "v1.0.0"); err != nil {
		log.Fatal(err)
	}

	// 3. Push the tagged bundle to a registry.
	repo, err := remote.NewRepository("registry.example.com/org/policy-bundle")
	if err != nil {
		log.Fatal(err)
	}
	if _, err := oras.Copy(ctx, store, "v1.0.0", repo, "v1.0.0", oras.DefaultCopyOptions); err != nil {
		log.Fatal(err)
	}
}
```

`GemaraVersion` records the spec version the artifacts were written against. `bundle.WithContinueOnError()` makes `Assemble` record a mapping-reference ID that matches no artifact in `b.Warnings` instead of returning an error; check `b.Warnings` before you push. Leave the option out to make any such mismatch fail the build. To authenticate to the registry, set `repo.Client` to an ORAS `auth.Client`.

## Pull a bundle

```go
b, err := bundle.Unpack(ctx, repo, "v1.0.0")
if err != nil {
	log.Fatal(err)
}
// b.Source is the policy; b.Imports holds every artifact it depends on.
```

## Related

- [Go SDK reference](../../sdk/go-sdk.html)
- [Write an organizational policy](../policy/)
