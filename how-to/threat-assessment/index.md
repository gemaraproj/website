---
layout: page
title: Threat Assessment Guide
description: Step-by-step guide to performing Gemara-compatible threat assessments
redirect_from:
  - /tutorials/controls/threat-assessment-guide.html
---

## What This Is

This guide walks through a threat assessment using the [Gemara](https://gemara.openssf.org/) project.

**The basic idea:** Think of a project like a house. First, you identify what the house can do: its **[capabilities](../../model/02-definitions.html#capability)** (e.g., "allow entry/exit", "store belongings"). Then, you identify **[threats](../../model/02-definitions.html#threat)**, what could go wrong with those capabilities (e.g., "unauthorized entry through unlocked door", "theft of stored belongings").

In technical terms:
 * **[Capabilities](../../model/02-definitions.html#capability)** define what the technology can do. These form a primary component of the **attack surface** because every intended function represents a potential path for unintended use.
 * **[Threats](../../model/02-definitions.html#threat)** define specific ways those capabilities could be misused or exploited.

This exercise helps you systematically identify what could go wrong so you can build appropriate defenses.

Gemara splits these into two artifact kinds: a `CapabilityCatalog` for capability definitions, and a `ThreatCatalog` for threats. The threat catalog points back at the capability catalog through a mapping reference, so an assessment produces two files that reference each other. The examples below target Gemara **v1.6.0**.

## Walkthrough

### Step 0: Define Scope

Select a component or technology to assess (service, API, infrastructure component, or technology stack).

**Leverage existing resources**: Gemara supports importing entries from external catalogs so you don't have to start from scratch. The FINOS Common Cloud Controls (CCC) Core [catalogs](https://github.com/finos/common-cloud-controls/releases/tag/core%2Fccc%2Fv2026.06-rc1) define well-vetted capabilities and threats that apply broadly across cloud services. These pre-built items can help accelerate your assessment.

We will explore how this is leveraged below as we dive into our container management tool example (i.e., `SEC.SLAM.CM.CAP` for capabilities and `SEC.SLAM.CM` for threats).

### Step 1: Setting Up Metadata

Each artifact declares its own scope and mapping references. Key fields:

| Field | What It Is | Why |
|-------|------------|-----|
| `title` | Display name for the catalog (top-level field) | Human-readable label used in reports and tooling output |
| `metadata.id` | Identifier for this catalog | Other artifacts point at this catalog by this value |
| `metadata.type` | `CapabilityCatalog` or `ThreatCatalog` | Identifies the artifact for `#CapabilityCatalog` or `#ThreatCatalog` validation (required) |
| `metadata.gemara-version` | Spec version the artifact conforms to, here `"1.6.0"` | Lets tooling pick the right schema and migration path (required) |
| `mapping-references` with `id: CCC` | Pointer to the CCC Core catalog release | Resolves the imported CCC capability and threat IDs used later |
| Top-level `imports` (Steps 2 and 3) | List of `#MultiEntryMapping` rows | Pulls in common capabilities and threats without redefining them; requires `mapping-references` |
| `groups` | Categories that entries in this catalog belong to | Required once a catalog defines capabilities or threats; keeps large catalogs readable |

The threat catalog adds one more mapping reference: a pointer to the capability catalog, so its threats can link to the capabilities defined there (see Step 3).

**Example (YAML)**: capability catalog metadata:

```yaml
title: Container Management Tool Capability Catalog
metadata:
  id: SEC.SLAM.CM.CAP
  type: CapabilityCatalog
  gemara-version: "1.6.0"
  version: "1.0.0"
  description: Capabilities of the container management tool under assessment
  author:
    id: example
    name: Example
    type: Human
  mapping-references:
    - id: CCC
      title: Common Cloud Controls Core
      version: v2026.06-rc1
      url: https://github.com/finos/common-cloud-controls/releases/download/core%2Fccc%2Fv2026.06-rc1/capabilities.yaml
      description: |
        Foundational repository of reusable security controls, capabilities,
        and threat models maintained by FINOS.
```

### Step 2: Identify Capabilities (Capability catalog)

Capabilities are the core functions or features within the scope. In Gemara they live in a `CapabilityCatalog`.

**Start with the imported capabilities** you can leverage from FINOS CCC. Ask: "Which common cloud capabilities does this technology have?"

A container management tool actively reaches out to registries to pull images and configuration. Since CCC Core already defines that as **CP29** (Active Ingestion), we import it rather than redefining it. Image tags also function as version identifiers: the tool resolves a tag like `latest` or `v1.0` to a specific image. CCC Core defines that behavior as **CP18** (Resource Versioning), so we import that as well.

**Example (YAML)**: `imports` on the capability catalog:

```yaml
imports:
  - reference-id: CCC
    entries:
      - reference-id: CCC.Core.CP29
        remarks: Active Ingestion
      - reference-id: CCC.Core.CP18
        remarks: Resource Versioning
```

**Then, define specific capabilities** unique to your target. Required fields:

| Field | Required | Description |
|-------|----------|-------------|
| `id` | Yes | Unique identifier. The schema enforces no pattern; this guide uses `ORG.PROJ.COMPONENT.CAP##` as a convention |
| `title` | Yes | Clear name for the capability |
| `description` | Yes | What this capability does |
| `group` | Yes | `id` of a **group** declared in this catalog's `groups` list |

**Example (YAML)**:

```yaml
groups:
  - id: image-management
    title: Image Management
    description: |
      Capabilities for retrieving, resolving, and running container images.
capabilities:
  - id: SEC.SLAM.CM.CAP01
    title: Image Retrieval by Tag
    description: |
      Ability to retrieve container images from registries using mutable tag names
      (e.g., 'latest', 'v1.0').
    group: image-management
```

### Step 3: Identify Threats (Threat catalog)

 **[Threats](../../model/02-definitions.html#threat)** are specific ways **[capabilities](../../model/02-definitions.html#capability)** can be misused, exploited, or cause problems. For each **capability**, identify potential **threats**. These go in the `ThreatCatalog`, which points back at the capability catalog through a mapping reference.

**Example (YAML)**: threat catalog mapping references:

```yaml
mapping-references:
  - id: CCC
    title: Common Cloud Controls Core
    version: v2026.06-rc1
    url: https://github.com/finos/common-cloud-controls/releases
    description: |
      Foundational repository of reusable security controls, capabilities,
      and threat models maintained by FINOS.
  - id: SEC.SLAM.CM.CAP
    title: Container Management Tool Capability Catalog
    version: "1.0.0"
    description: |
      Capabilities defined for the container management tool, exploited by
      the threats in this catalog.
```

> **Note:** When a mapping reference has no `url`, its `id` is the referenced artifact's `metadata.id`. That is why the second reference uses `SEC.SLAM.CM.CAP`.

**Check for imported threats first.** As with **[capabilities](../../model/02-definitions.html#capability)**, review the CCC Core catalog for threats linked to the capabilities you imported. If a threat fits your scope, import it. In this example, CCC Core defines **TH14** ("Older Resource Versions are Used") which is linked to **CP18**. It applies because mutable image tags let the tool resolve to a stale or compromised version.

Imported threats and the capabilities you reference from the capability catalog share the single top-level `imports` list.

**Example (YAML)**: `imports` on the threat catalog:

```yaml
imports:
  - reference-id: CCC
    entries:
      - reference-id: CCC.Core.TH14
        remarks: Older Resource Versions are Used
  - reference-id: SEC.SLAM.CM.CAP
    entries:
      - reference-id: SEC.SLAM.CM.CAP01
        remarks: Image Retrieval by Tag
```

**Then, define specific threats** unique to your target. Required fields:

| Field | Required | Description |
|-------|----------|-------------|
| `id` | Yes | Unique identifier. The schema enforces no pattern; this guide uses `ORG.PROJ.COMPONENT.THR##` as a convention |
| `title` | Yes | Short name for the threat |
| `description` | Yes | What goes wrong and why it matters |
| `group` | Yes | `id` of a **group** declared in this catalog's `groups` list |
| `capabilities` | Yes | Links this threat to the capabilities it exploits, by mapping reference |
| `vectors` | No | Optional link to vector catalog entries |
| `actors` | No | Optional threat actors (`#Actor`) |

**Example (YAML)**: a custom threat (*Container Image Tampering or Poisoning*) linked to the capabilities it exploits: **CCC.Core.CP29** (Active Ingestion), **CCC.Core.CP18** (Resource Versioning), and **SEC.SLAM.CM.CAP01** (Image Retrieval by Tag) via the capability catalog reference (`SEC.SLAM.CM.CAP`).

```yaml
groups:
  - id: image-integrity
    title: Image Integrity
    description: |
      Threats to the authenticity and provenance of container images the tool
      retrieves and runs.
threats:
  - id: SEC.SLAM.CM.THR01
    title: Container Image Tampering or Poisoning
    description: |
      Attackers may replace a legitimately published image tag with a malicious image
      by exploiting tag mutability in image registries, especially when the container
      management tool retrieves images by tag name rather than digest. This enables
      unauthorized access, data exfiltration, and system compromise.
    group: image-integrity
    capabilities:
      - reference-id: CCC
        entries:
          - reference-id: CCC.Core.CP29
          - reference-id: CCC.Core.CP18
      - reference-id: SEC.SLAM.CM.CAP
        entries:
          - reference-id: SEC.SLAM.CM.CAP01
```

### Step 4: Validate

The final capability catalog should look something like this:

```yaml
title: Container Management Tool Capability Catalog
metadata:
  id: SEC.SLAM.CM.CAP
  type: CapabilityCatalog
  gemara-version: "1.6.0"
  version: "1.0.0"
  description: Capabilities of the container management tool under assessment
  author:
    id: example
    name: Example
    type: Human
  mapping-references:
    - id: CCC
      title: Common Cloud Controls Core
      version: v2026.06-rc1
      url: https://github.com/finos/common-cloud-controls/releases/download/core%2Fccc%2Fv2026.06-rc1/capabilities.yaml
      description: |
        Foundational repository of reusable security controls, capabilities,
        and threat models maintained by FINOS.
imports:
  - reference-id: CCC
    entries:
      - reference-id: CCC.Core.CP29
        remarks: Active Ingestion
      - reference-id: CCC.Core.CP18
        remarks: Resource Versioning
groups:
  - id: image-management
    title: Image Management
    description: |
      Capabilities for retrieving, resolving, and running container images.
capabilities:
  - id: SEC.SLAM.CM.CAP01
    title: Image Retrieval by Tag
    description: |
      Ability to retrieve container images from registries using mutable tag names
      (e.g., 'latest', 'v1.0').
    group: image-management
```

And the final threat catalog like this:

```yaml
title: Container Management Tool Threat Catalog
metadata:
  id: SEC.SLAM.CM
  type: ThreatCatalog
  gemara-version: "1.6.0"
  version: "1.0.0"
  description: Threat catalog for container management tool security assessment
  author:
    id: example
    name: Example
    type: Human
  mapping-references:
    - id: CCC
      title: Common Cloud Controls Core
      version: v2026.06-rc1
      url: https://github.com/finos/common-cloud-controls/releases
      description: |
        Foundational repository of reusable security controls, capabilities,
        and threat models maintained by FINOS.
    - id: SEC.SLAM.CM.CAP
      title: Container Management Tool Capability Catalog
      version: "1.0.0"
      description: |
        Capabilities defined for the container management tool, exploited by
        the threats in this catalog.
imports:
  - reference-id: CCC
    entries:
      - reference-id: CCC.Core.TH14
        remarks: Older Resource Versions are Used
  - reference-id: SEC.SLAM.CM.CAP
    entries:
      - reference-id: SEC.SLAM.CM.CAP01
        remarks: Image Retrieval by Tag
groups:
  - id: image-integrity
    title: Image Integrity
    description: |
      Threats to the authenticity and provenance of container images the tool
      retrieves and runs.
threats:
  - id: SEC.SLAM.CM.THR01
    title: Container Image Tampering or Poisoning
    description: |
      Attackers may replace a legitimately published image tag with a malicious image
      by exploiting tag mutability in image registries, especially when the container
      management tool retrieves images by tag name rather than digest. This enables
      unauthorized access, data exfiltration, and system compromise.
    group: image-integrity
    capabilities:
      - reference-id: CCC
        entries:
          - reference-id: CCC.Core.CP29
          - reference-id: CCC.Core.CP18
      - reference-id: SEC.SLAM.CM.CAP
        entries:
          - reference-id: SEC.SLAM.CM.CAP01
```

For a longer example with more CCC imports and additional threats, see [capabilities.yaml](capabilities.yaml) and [threat-catalog.yaml](threat-catalog.yaml).

**Validation commands:**

```bash
go install cuelang.org/go/cmd/cue@latest
cue vet -c -d '#CapabilityCatalog' github.com/gemaraproj/gemara@v1.6.0 your-capabilities.yaml
cue vet -c -d '#ThreatCatalog' github.com/gemaraproj/gemara@v1.6.0 your-threat-catalog.yaml
```

No output means the artifact is valid. Swap `@v1.6.0` for `@latest` to validate against the newest published schema.

## What's Next

Create a Gemara control catalog that maps security controls to these threats using the [Control Catalog Guide](../control-catalog/).

**Layer 2 schema documentation:** 
- [Capability Catalog](https://gemara.openssf.org/schema/capabilitycatalog.html)
- [Threat Catalog](https://gemara.openssf.org/schema/threatcatalog.html)
- [Control Catalog](https://gemara.openssf.org/schema/controlcatalog.html)
