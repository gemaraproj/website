---
layout: page
title: Write and validate your first control catalog
description: "A first Gemara tutorial: write a two-control catalog in YAML, validate it with CUE, fix a validation error and export the catalog to OSCAL."
---

In this tutorial you write a small Gemara control catalog for a web service, check it against the Gemara schema, break it on purpose to see what a validation error looks like, fix it, and convert it to OSCAL.

It takes about 15 minutes. You don't need to know Gemara or CUE first.

## Before you start

You need [Go](https://go.dev/doc/install) 1.25 or later and a terminal. Install the CUE command-line tool:

```bash
go install cuelang.org/go/cmd/cue@latest
```

Check that it runs:

```bash
cue version
```

If your shell can't find `cue`, add Go's install folder to your `PATH`: `export PATH="$PATH:$(go env GOPATH)/bin"`.

Then make an empty folder to work in:

```bash
mkdir gemara-tutorial && cd gemara-tutorial
```

## Step 1: Write the catalog

A control catalog is a list of [controls](../../model/02-definitions.html#control). Each control has an objective and one or more assessment requirements: conditions someone can check to see whether the objective is met.

Create a file called `my-catalog.yaml` and paste this into it:

```yaml
title: Example Web Service Controls
metadata:
  id: EXAMPLE.WEB
  type: ControlCatalog
  gemara-version: "1.6.0"
  description: Two controls for a small web service, written for the Gemara tutorial.
  version: 0.1.0
  author:
    id: your-github-handle
    name: Your Name
    type: Human
  applicability-groups:
    - id: production
      title: Production
      description: Live deployments that serve users.

groups:
  - id: EXAMPLE.WEB.ACCESS
    title: Access control
    description: Controls that decide who can reach the service and what they can do.

controls:
  - id: EXAMPLE.WEB.CTL01
    title: Require multi-factor authentication for administrators
    objective: |
      Make a stolen password alone insufficient to take over an administrator account.
    group: EXAMPLE.WEB.ACCESS
    assessment-requirements:
      - id: EXAMPLE.WEB.CTL01.AR01
        text: Every administrator account MUST require a second factor at sign-in.
        applicability: ["production"]
  - id: EXAMPLE.WEB.CTL02
    title: Serve all traffic over TLS
    objective: |
      Keep credentials and session data from being read or changed in transit.
    group: EXAMPLE.WEB.ACCESS
    assessment-requirements:
      - id: EXAMPLE.WEB.CTL02.AR01
        text: The service MUST reject plain-HTTP requests or redirect them to HTTPS.
        applicability: ["production"]
```

Notice three things about this file:

- `metadata.type` says what kind of Gemara artifact this is, and `metadata.gemara-version` says which version of the specification it follows.
- Every control names a group in `group`. The group `EXAMPLE.WEB.ACCESS` is declared once, under `groups`.
- Every assessment requirement lists where it applies in `applicability`. The value `production` is declared once, under `metadata.applicability-groups`.

## Step 2: Validate it

Ask CUE to check the file against the `#ControlCatalog` definition in Gemara v1.6.0:

```bash
cue vet -c -d '#ControlCatalog' github.com/gemaraproj/gemara@v1.6.0 my-catalog.yaml
```

The first run downloads the Gemara module from the CUE Central Registry, so it may take a few seconds.

If the catalog is valid, the command prints nothing and exits. No news is good news.

## Step 3: Break it and read the error

Now make a mistake on purpose. In the first control, change `group: EXAMPLE.WEB.ACCESS` to `group: EXAMPLE.WEB.ACESS`, dropping one `C`. Save the file and run the same command again:

```bash
cue vet -c -d '#ControlCatalog' github.com/gemaraproj/gemara@v1.6.0 my-catalog.yaml
```

This time CUE reports an error (the cache paths in the middle lines will differ on your machine):

```text
_groupValidation."0": invalid value ["EXAMPLE.WEB.ACCESS"] (does not satisfy list.Contains("EXAMPLE.WEB.ACESS")):
    .../github.com/gemaraproj/gemara@v1.6.0/controlcatalog.cue:28:47
    .../github.com/gemaraproj/gemara@v1.6.0/controlcatalog.cue:28:30
    ./my-catalog.yaml:27:12
```

Read it from the outside in:

- `_groupValidation."0"` is the check that failed: the group of control number 0, the first control.
- `["EXAMPLE.WEB.ACCESS"]` is the list of groups the catalog declares, and `"EXAMPLE.WEB.ACESS"` is the value that isn't in it.
- The last line, `./my-catalog.yaml:27:12`, points at the line and column in your file.

The YAML itself was fine. The schema caught a reference to a group that doesn't exist, which a plain YAML linter would miss.

## Step 4: Fix it

Change `EXAMPLE.WEB.ACESS` back to `EXAMPLE.WEB.ACCESS` and run the command once more:

```bash
cue vet -c -d '#ControlCatalog' github.com/gemaraproj/gemara@v1.6.0 my-catalog.yaml
```

It prints nothing again. Your catalog is valid.

## Step 5: Export it to OSCAL

A Gemara catalog is meant to be used by tools. To see one at work, convert yours to [OSCAL](https://pages.nist.gov/OSCAL/), the compliance format many assessors and government programs accept. Install the `oscalexport` command from the [Go SDK](../../sdk/go-sdk.html):

```bash
go install github.com/gemaraproj/go-gemara/cmd/oscalexport@v0.12.0
```

Then run it on your catalog:

```bash
oscalexport catalog my-catalog.yaml --output my-catalog.oscal.json
```

It prints:

```text
Successfully wrote OSCAL content to my-catalog.oscal.json
```

Open `my-catalog.oscal.json`. Your group became an OSCAL group holding both controls, and each assessment requirement became an OSCAL sub-control whose statement is your requirement text:

```json
"prose": "Every administrator account MUST require a second factor at sign-in."
```

## What you did

You wrote a Gemara control catalog, validated it against the published schema, read and fixed a validation error, and turned the catalog into OSCAL.

## Next steps

- To write a catalog for a real system, linked to the threats it mitigates, follow [Write a control catalog](../../how-to/control-catalog/).
- To see every field a control catalog can have, read the [Control Catalog schema](../../schema/controlcatalog.html).
- To see where control catalogs sit among the other artifacts, read [The Gemara model](../../model/).
