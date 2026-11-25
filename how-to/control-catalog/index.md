---
layout: page
title: Control Catalog Guide
description: Step-by-step guide to creating Gemara-compatible control catalogs
redirect_from:
  - /tutorials/controls/control-catalog-guide.html
---

## What This Is

This guide walks through creating a **control catalog** using the [Gemara](https://gemara.openssf.org/) project, building on the threats and scope you identified in the [Threat Assessment Guide](../threat-assessment/).

In technical terms:
* **[Controls](../../model/02-definitions.html#control)** are safeguards with a stated **objective** and a list of **[assessment requirements](../../model/02-definitions.html#assessment-requirement)**.
* **Groups** group related controls by domain (e.g., supply chain, access control).
* **[Threats](../../model/02-definitions.html#threat)** link each control to the threat(s) it mitigates, connecting your control catalog to your layer 2 threat catalog.

This exercise produces a structured way to develop control objectives and corresponding testable conditions to determine if the objective is met.

## Optional Workflow

Complete the [Threat Assessment Guide](../threat-assessment/) for your scope (e.g., the container management tool **SEC.SLAM.CM**). You can reference the linked `threats` and `guidelines` to support the intent of Controls in the catalog. 


## Walkthrough

Each step explains only the fields it needs. For every field and its type, see the [Control Catalog schema](../../schema/controlcatalog.html).

### Step 0: Define Scope and References

Reuse the same component or technology as your threat assessment. You will reference that threat catalog so controls can map to the threats you identified.

**Leverage existing resources**: As with threat catalogs, you can import controls from external catalogs (e.g., FINOS Common Cloud Controls, OSPS Baseline) via `mapping-references` and list imported controls in your catalog. 

We continue with the container management tool example (SEC.SLAM.CM) and assume a threat catalog for it already exists.

### Step 1: Setting Up Metadata

Declare your control catalog and the catalogs it refers to. The fields that need a decision:

| Field | What to put there |
|-------|-------------------|
| `metadata.type` | `ControlCatalog` |
| `metadata.gemara-version` | The spec version you validate against, here `"1.6.0"` |
| `metadata.mapping-references` | One entry per catalog you import from or link to, such as your threat catalog and the CCC Core catalogs. The `id` is what `reference-id` uses later. |
| `metadata.applicability-groups` | The contexts your requirements apply in, such as production or CI/CD. Required once the catalog defines controls; every `applicability` value must be one of these ids. |

**Example (YAML):**

```yaml
title: Container Management Tool Security Control Catalog
metadata:
  id: SEC.SLAM.CM
  type: ControlCatalog
  gemara-version: "1.6.0"
  description: Control catalog for container management tool security; mitigates threats from the SEC.SLAM.CM.THR threat catalog.
  version: 1.0.0
  author:
    id: example
    name: Example
    type: Human
  mapping-references:
    - id: SEC.SLAM.CM.THR
      title: Container Management Tool Security Threat Catalog
      version: "1.0.0"
      url: file://threats.yaml
      description: |
        Threat catalog for the same scope; provides threat IDs referenced from each
        control's threats.
    - id: CCC.Core.Controls
      title: CCC Core Controls
      version: v2026.11
      url: https://hub.grc.store/v1/catalogs/finos-ccc/ccc.core.controls/versions/v2026.11
      description: Common cloud controls from the FINOS Common Cloud Controls Core catalog.
    - id: CCC.Core.Threats
      title: CCC Core Threats
      version: v2026.11
      url: https://hub.grc.store/v1/catalogs/finos-ccc/ccc.core.threats/versions/v2026.11
      description: Common cloud threats from the FINOS Common Cloud Controls Core catalog.
  applicability-groups:
    - id: production
      title: Production
      description: |
        Production container workloads and clusters; controls apply to
        live environments where security posture is enforced.
    - id: all_deployments
      title: All Deployments
      description: |
        Requirements that apply whenever container images are built, pulled, or run—
        regardless of environment (dev, staging, production) or pipeline stage.
    - id: untrusted_networks
      title: Untrusted Networks
      description: |
        Applies when registry or image traffic traverses untrusted networks.
    - id: ci_cd
      title: CI/CD
      description: |
        Applies in continuous integration and deployment pipelines.
```

### Step 2: Define Control Groups

**Groups** group controls by theme. The schema requires at least one group when the catalog defines its own `controls`. Each group needs an `id`, a `title` and a `description`, and each control's `group` field must match the `id` of one of them.

**Example (YAML):**

```yaml
groups:
  - id: SEC.SLAM.CM.FAM01
    title: Image Integrity and Supply Chain
    description: |
      Controls that ensure container images are authentic, unmodified,
      and from trusted sources throughout retrieval and use.
```

### Step 3: Import or Define Controls

> **Note:** For how controls are automatically pulled into policy to mitigate high-severity risks, see the FAQ.


**Option A — Import Controls:** If an external catalog (e.g., CCC) defines controls that address your threats, reference it in `mapping-references` and list imported controls under top-level `imports`.

**Example (YAML):**

```yaml
imports:
  - reference-id: CCC.Core.Controls
    entries:
      - reference-id: CCC.Core.CN05
        remarks: Prevent access from untrusted entities
```

**Option B — Define your own controls:** For controls specific to your scope, define them under `controls`. Each control needs an `id` (this guide uses `ORG.PROJ.COMPONENT.CTL##`), a `title`, an `objective` stating its intent, a `group` from Step 2, and at least one assessment requirement. The fields that link a control to the rest of your work:

| Field | What to put there |
|-------|-------------------|
| `assessment-requirements` | Verifiable conditions, each with an `id`, a `text` (usually a MUST statement) and an `applicability` list of ids from `metadata.applicability-groups` |
| `threats` | Optional. The threats this control mitigates, as a `reference-id` from `mapping-references` plus the threat ids under `entries` |
| `state` | Optional. `Active` (the default), `Draft`, `Deprecated` or `Retired` |

**Example (YAML):**

```yaml
controls:
  - id: SEC.SLAM.CM.CTL01
    title: Use Immutable Image References by Digest
    objective: |
      Require signature validation so that only legitimate, trusted images are
      accepted; then pin each image to an immutable digest (e.g., sha256)
      after the check so that what is used matches what was verified and
      TOCTOU (time-of-check to time-of-use) attacks are prevented.
    group: SEC.SLAM.CM.FAM01
    assessment-requirements:
      - id: SEC.SLAM.CM.CTL01.AR01
        text: |
          The system MUST verify image signature before pull or run, then pin
          the image to a digest (e.g., sha256:...) after the check and use that
          digest for all subsequent use.
        applicability: ["all_deployments"]
      - id: SEC.SLAM.CM.CTL01.AR02
        text: |
          Configuration and policies MUST disallow or override use of tag-only
          references for production or sensitive workloads where supported.
        applicability: ["production"]
    threats:
      - reference-id: SEC.SLAM.CM.THR
        entries:
          - reference-id: SEC.SLAM.CM.THR01
          - reference-id: SEC.SLAM.CM.THR03
          - reference-id: SEC.SLAM.CM.THR04
      - reference-id: CCC.Core.Threats
        entries:
          - reference-id: CCC.Core.TH14
  - id: SEC.SLAM.CM.CTL02
    title: Require TLS/SSL with Certificate Pinning
    objective: |
      Mitigate MITM Container Image Interception by protecting registry traffic and verifying artifact integrity: use
      TLS/SSL with certificate pinning for all registry communication, use VPNs on untrusted networks to reduce interception risk, and verify artifact signatures or hashes so that tampered or redirected content is detected even if the channel is compromised.
    group: SEC.SLAM.CM.FAM01
    state: Active
    assessment-requirements:
      - id: SEC.SLAM.CM.CTL02.AR01
        text: |
          The system MUST use TLS/SSL for all registry communication and MUST pin to the expected server certificate or public key (or certificate chain) for the registry.
        applicability: ["all_deployments"]
        state: Active
      - id: SEC.SLAM.CM.CTL02.AR02
        text: |
          On untrusted networks, the system or deployment pipeline MUST use a VPN or other trusted path for registry traffic, or MUST restrict
          image pulls to environments where the network is trusted.
        applicability: ["untrusted_networks", "ci_cd"]
        state: Active
      - id: SEC.SLAM.CM.CTL02.AR03
        text: |
          The system MUST verify artifact signatures or hashes (e.g. signature verification, digest check) before use so that tampered or redirected artifacts are rejected.
        applicability: ["all_deployments"]
        state: Active
    threats:
      - reference-id: SEC.SLAM.CM.THR
        entries:
          - reference-id: SEC.SLAM.CM.THR02
      - reference-id: CCC.Core.Threats
        entries:
          - reference-id: CCC.Core.TH02
```

### Step 4: Validate Against the Gemara Schema

Validate the final catalog with CUE:

**Validation commands:**

```bash
go install cuelang.org/go/cmd/cue@latest
cue vet -c -d '#ControlCatalog' github.com/gemaraproj/gemara@v1.6.0 your-control-catalog.yaml
```

Fix any reported errors (e.g., missing required fields, invalid `group` reference, or malformed `threats`) so the catalog is schema-consistent.

### Step 5: Assemble the Full Catalog and Validate

Combine metadata, mapping-references, groups, `imports` (if any), and controls into a single YAML document. A complete, schema-valid catalog for the SEC.SLAM.CM scenario is in [control-catalog.yaml](control-catalog.yaml).

**Complete example (SEC.SLAM.CM):**

```yaml
title: Container Management Tool Security Control Catalog

metadata:
  id: SEC.SLAM.CM
  type: ControlCatalog
  gemara-version: "1.6.0"
  description: |
    Control catalog for container management tool security; mitigates threats
    from the SEC.SLAM.CM.THR threat catalog.
  version: 1.0.0
  author:
    id: example
    name: Example
    type: Human
  mapping-references:
    - id: SEC.SLAM.CM.THR
      title: Container Management Tool Security Threat Catalog
      version: "1.0.0"
      url: https://gemara.openssf.org/how-to/threat-assessment/threat-catalog.yaml
      description: |
        Threat catalog for the same scope; provides threat IDs referenced from each
        control's threats.
    - id: CCC.Core.Controls
      title: CCC Core Controls
      version: v2026.11
      url: https://hub.grc.store/v1/catalogs/finos-ccc/ccc.core.controls/versions/v2026.11
      description: Common cloud controls from the FINOS Common Cloud Controls Core catalog.
    - id: CCC.Core.Threats
      title: CCC Core Threats
      version: v2026.11
      url: https://hub.grc.store/v1/catalogs/finos-ccc/ccc.core.threats/versions/v2026.11
      description: Common cloud threats from the FINOS Common Cloud Controls Core catalog.
  applicability-groups:
    - id: production
      title: Production
      description: |
        Production container workloads and clusters; controls apply to
        live environments where security posture is enforced.
    - id: all_deployments
      title: All Deployments
      description: |
        Requirements that apply whenever container images are built, pulled, or run—
        regardless of environment (dev, staging, production) or pipeline stage.
    - id: untrusted_networks
      title: Untrusted Networks
      description: |
        Applies when registry or image traffic traverses untrusted networks.
    - id: ci_cd
      title: CI/CD
      description: |
        Applies in continuous integration and deployment pipelines.

groups:
  - id: SEC.SLAM.CM.FAM01
    title: Image Integrity and Supply Chain
    description: |
      Controls that ensure container images are authentic, unmodified,
      and from trusted sources throughout retrieval and use.

imports:
  - reference-id: CCC.Core.Controls
    entries:
      - reference-id: CCC.Core.CN05
        remarks: Prevent access from untrusted entities

controls:
  - id: SEC.SLAM.CM.CTL01
    title: Use Immutable Image References by Digest
    objective: |
      Require signature validation so that only legitimate, trusted images are
      accepted; then pin each image to an immutable digest (e.g., sha256)
      after the check so that what is used matches what was verified and
      TOCTOU (time-of-check to time-of-use) attacks are prevented.
    group: SEC.SLAM.CM.FAM01
    assessment-requirements:
      - id: SEC.SLAM.CM.CTL01.AR01
        text: |
          The system MUST verify image signature before pull or run, then pin
          the image to a digest (e.g., sha256:...) after the check and use that
          digest for all subsequent use.
        applicability: ["all_deployments"]
      - id: SEC.SLAM.CM.CTL01.AR02
        text: |
          Configuration and policies MUST disallow or override use of tag-only
          references for production or sensitive workloads where supported.
        applicability: ["production"]
    threats:
      - reference-id: SEC.SLAM.CM.THR
        entries:
          - reference-id: SEC.SLAM.CM.THR01
          - reference-id: SEC.SLAM.CM.THR03
          - reference-id: SEC.SLAM.CM.THR04
      - reference-id: CCC.Core.Threats
        entries:
          - reference-id: CCC.Core.TH14
  - id: SEC.SLAM.CM.CTL02
    title: Require TLS/SSL with Certificate Pinning
    objective: |
      Mitigate MITM Container Image Interception by protecting registry
      traffic and verifying artifact integrity: use TLS/SSL with certificate
      pinning for all registry communication, use VPNs on untrusted networks
      to reduce interception risk, and verify artifact signatures or hashes so
      that tampered or redirected content is detected even if the channel is
      compromised.
    group: SEC.SLAM.CM.FAM01
    state: Active
    assessment-requirements:
      - id: SEC.SLAM.CM.CTL02.AR01
        text: |
          The system MUST use TLS/SSL for all registry communication and MUST
          pin to the expected server certificate or public key (or certificate
          chain) for the registry.
        applicability: ["all_deployments"]
        state: Active
      - id: SEC.SLAM.CM.CTL02.AR02
        text: |
          On untrusted networks, the system or deployment pipeline MUST use a
          VPN or other trusted path for registry traffic, or MUST restrict
          image pulls to environments where the network is trusted.
        applicability: ["untrusted_networks", "ci_cd"]
        state: Active
      - id: SEC.SLAM.CM.CTL02.AR03
        text: |
          The system MUST verify artifact signatures or hashes (e.g. signature
          verification, digest check) before use so that tampered or redirected
          artifacts are rejected.
        applicability: ["all_deployments"]
        state: Active
    threats:
      - reference-id: SEC.SLAM.CM.THR
        entries:
          - reference-id: SEC.SLAM.CM.THR02
      - reference-id: CCC.Core.Threats
        entries:
          - reference-id: CCC.Core.TH02
```
**Validation commands:**

Using the **published** module:

```bash
go install cuelang.org/go/cmd/cue@latest
cue vet -c -d '#ControlCatalog' github.com/gemaraproj/gemara@v1.6.0 your-control-catalog.yaml
```

Fix any reported errors (e.g., missing required fields, invalid `group` reference, or malformed `threats`) so the catalog is schema-consistent.

## Summary: From Threat Assessment to Control Catalog

| From threat assessment        | Use in control catalog                                      |
|-------------------------------|-------------------------------------------------------------|
| Threat catalog id (e.g. SEC.SLAM.CM.THR) | A `mapping-references` entry, and the `reference-id` in each control's `threats` |
| Threat IDs (e.g. SEC.SLAM.CM.THR01, CCC.Core.TH14) | `threats` on each control that mitigates those threats |
| Understanding of what goes wrong | Control objectives and assessment requirements that prevent or detect it |

## What's Next

Use the control catalog in downstream Gemara layers: for example, map controls to Layer 1 guidelines via `guidelines`, or use the catalog in Layer 5 evaluations to record which assessment requirements have been verified.

**Layer 2 schema documentation:** 
- [Threat Catalog](https://gemara.openssf.org/schema/threatcatalog.html)
- [Capability Catalog](https://gemara.openssf.org/schema/capabilitycatalog.html)
- [Control Catalog](https://gemara.openssf.org/schema/controlcatalog.html)
