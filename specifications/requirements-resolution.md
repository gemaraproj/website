---
layout: page
title: Requirements Resolution Specification
nav-title: Requirements Resolution
---

This specification defines how a Gemara **[Policy](../model/05.3-Layer-3.html)** (Layer 3) resolves a single **[Control Catalog](../model/05.2-Layer-2.html)** import into an **Effective Requirements Set** (represented as a derived `ControlCatalog`).

The normative behavior defined here reflects the Gemara schemas (`policy.cue`, `controlcatalog.cue`) and the profile model established in [gemaraproj/gemara#560](https://github.com/gemaraproj/gemara/issues/560) and [gemaraproj/gemara#548](https://github.com/gemaraproj/gemara/issues/548).

> **Scope & Deferred Items:**
> * **Single Catalog Imports:** This document covers the resolution of a single `imports.catalogs[]` entry.
> * **Guidance Imports (`imports.guidance[]`):** Guidance imports accept exclusions and constraints. Their multi-artifact resolution will be detailed in a follow-up specification.
> * **Policy Imports (`imports.policies[]`):** A policy can import another policy as an inheritance model. An importing policy brings in the referenced policy's resolved requirements wholesale and cannot alter them.
> * **Static Resolution and Policy Scope:** Policy `scope` defines organizational, architectural, and geopolitical boundaries. Requirements filtering is governed exclusively by declared `exclusions` and adopted `applicability-groups` (see [gemaraproj/gemara#548](https://github.com/gemaraproj/gemara/issues/548)).
> * **Cross-Catalog Namespacing:** Prefixing identifiers across multiple catalog imports is deferred to the multi-catalog resolution specification. In this single-catalog baseline specification, identifiers remain in their local catalog-scoped form.

---

## 1. Resolution Inputs

A conforming resolver accepts:
1. **The Policy Document:** The primary Layer 3 artifact containing `imports.catalogs[]`, `metadata.mapping-references`, and optional `metadata.applicability-groups`.
2. **The Source Control Catalog:** The Layer 2 artifact identified by the import's `reference-id`.

### Schema References
* `policy.cue`: `#Imports`, `#CatalogImport`, `#Constraint`.
* `controlcatalog.cue`: `#ControlCatalog`, `#Control`, `#AssessmentRequirement`.

---

## 2. Resolution Workflow (Single Catalog Import)

In accordance with [gemaraproj/gemara#560](https://github.com/gemaraproj/gemara/issues/560), a Policy is a **profile**: it selects, constrains, and parametrizes; it does not rewrite catalog content. Consequently, resolution is a deterministic 3-step pipeline: **Acquire ➔ Exclude ➔ Constrain**.

```
+-----------------------------------+
| 1. Artifact Acquisition           |
|    - Resolve reference-id         |
|    - Verify metadata and version  |
+-----------------+-----------------+
                  |
                  v
+-----------------+-----------------+
| 2. Apply Exclusions & Selection   |
|    - Filter by applicability      |
|    - Excise excluded IDs          |
+-----------------+-----------------+
                  |
                  v
+-----------------+-----------------+
| 3. Apply Constraints              |
|    - Attach prescriptive text     |
|    - Bind ad hoc requirements     |
+-----------------+-----------------+
                  |
                  v
+-----------------------------------+
| Resolved Effective ControlCatalog |
+-----------------------------------+
```

### Step 1: Artifact Acquisition & Identification
The resolver locates the catalog using `policy.metadata.mapping-references` matching `policy.imports.catalogs[].reference-id`.
* The mapping reference `id` MUST match the referenced artifact's `metadata.id`.
* When `version` is specified in the mapping reference, it MUST match the artifact's `metadata.version`.
* The transport mechanism used to retrieve the artifact (local filesystem, HTTP, OCI registry) is implementation-specific.

### Step 2: Apply Exclusions & Applicability Filtering
Candidate controls and assessment requirements from the acquired catalog are filtered to form the candidate subset:

1. **Applicability Group Selection ([gemaraproj/gemara#548](https://github.com/gemaraproj/gemara/issues/548)):**
   * If the policy defines or adopts specific applicability groups in `metadata.applicability-groups`, any assessment requirement whose declared `applicability` tags do not intersect with the adopted groups is filtered out.
   * If the policy adopts no specific applicability groups, all baseline applicability tags are retained.
2. **Explicit Exclusions:**
   * If `exclusions` are declared on `#CatalogImport`, any control ID or assessment requirement ID matching an entry in `exclusions` is excised (dropped completely) from the candidate set.
   * Dropping all assessment requirements for a control does not drop the parent control unless the control ID itself is explicitly listed under `exclusions`.

### Step 3: Apply Constraints
If `constraints` are declared on the import:
* In `policy.cue`, `#Constraint` is defined as `{id, target-id, text}`: prescriptive requirement text attached to a specific control or guidance item (e.g., "Rotate service account credentials every 90 days").
* The resolver attaches these constraint statements to the matching target control in the candidate set as supplemental mandates.

---

## 3. Catalog Tailoring & Downstream Extension

Gemara enforces a strict separation of concerns:
* **Subset & Parametrize (In-Policy):** Use `exclusions` to narrow the requirement set and `constraints` to add minimum thresholds.
* **Reword or Restate (In-Catalog):** If an organization requires requirements worded differently, with altered recommendations, or entirely new custom requirements, it publishes a **local Control Catalog** and traces relationships back to upstream catalogs using a **[Mapping Document](../model/05.2-Layer-2.html)**.

This guarantees that:
1. The resolved catalog is always a strict subset of the source catalog.
2. Requirement IDs, objectives, and normative text preserve provenance without ambiguity or synthetic ID minting.

---

## 4. Worked Example & Test Vectors

The following example demonstrates resolving a single catalog import.

### 4.1 Input: Base Control Catalog (`SEC.SLAM.CM`)

```yaml
title: "Container Management Security Controls"
metadata:
  id: "SEC.SLAM.CM"
  type: "ControlCatalog"
  gemara-version: "1.2.0"
  version: "1.0.0"
  description: "Control catalog for container management security."
  author:
    id: "example-author"
    name: "Security Engineering"
    type: "Human"
  applicability-groups:
    - id: "all_deployments"
      title: "All Deployments"
      description: "Requirements that apply to all deployments."
    - id: "untrusted_networks"
      title: "Untrusted Networks"
      description: "Applies when communication traverses untrusted networks."
    - id: "production"
      title: "Production"
      description: "Applies to live production workloads."

groups:
  - id: "SEC.SLAM.CM.FAM01"
    title: "Network and Image Transport Security"
    description: "Controls governing network transport and image integrity."

controls:
  - id: "SEC.SLAM.CM.CTL02"
    title: "Require TLS/SSL and Digest Verification"
    objective: "Protect registry traffic and verify artifact integrity."
    group: "SEC.SLAM.CM.FAM01"
    assessment-requirements:
      - id: "SEC.SLAM.CM.CTL02.AR01"
        text: "The system MUST use TLS/SSL for all registry communication."
        applicability: ["all_deployments"]
      - id: "SEC.SLAM.CM.CTL02.AR02"
        text: "On untrusted networks, the system MUST use a VPN."
        applicability: ["untrusted_networks"]
      - id: "SEC.SLAM.CM.CTL02.AR03"
        text: "The system MUST verify image digests before execution."
        applicability: ["all_deployments"]
      - id: "SEC.SLAM.CM.CTL02.AR05"
        text: "The system MUST log all container image retrieval events."
        applicability: ["all_deployments"]
```

### 4.2 Input: Policy with Import, Exclusions, and Constraints (`org-policy-001`)

```yaml
title: "Production Infrastructure Security Policy"
metadata:
  id: "org-policy-001"
  type: "Policy"
  gemara-version: "1.2.0"
  version: "1.0.0"
  description: "Organizational policy establishing container transport requirements."
  author:
    id: "security-team"
    name: "Security Team"
    type: "Human"
  mapping-references:
    - id: "SEC.SLAM.CM"
      title: "Container Management Security Controls"
      version: "1.0.0"
      description: "Base container control catalog."
  applicability-groups:
    - id: "all_deployments"
      title: "All Deployments"
      description: "Workloads deployed across standard environments."

contacts:
  responsible:
    - name: "Platform Engineering"
      affiliation: "Engineering"
      email: "platform@example.com"
  accountable:
    - name: "CISO"
      affiliation: "Security"
      email: "ciso@example.com"

scope:
  in:
    technologies:
      - "Cloud Infrastructure"
      - "Container Runtimes"
    geopolitical:
      - "Global"

imports:
  catalogs:
    - reference-id: "SEC.SLAM.CM"
      # Exclude AR03 (e.g. handled by external admission webhook)
      exclusions:
        - "SEC.SLAM.CM.CTL02.AR03"
      # Attach prescriptive organizational constraint to CTL02
      constraints:
        - id: "CST-CTL02-01"
          target-id: "SEC.SLAM.CM.CTL02"
          text: "Registry connections MUST enforce TLS 1.3 or higher with corporate root CA pinning."

adherence:
  evaluation-methods:
    - id: "EV-AUTO-01"
      type: "Behavioral"
      mode: "Automated"
      required: true
      description: "Automated CI/CD policy verification checks."
```

### 4.3 Output: Resolved Effective Control Catalog

The resolved effective requirements set is a strict subset of the source catalog.
* `AR01` and `AR05` are retained (match adopted applicability `all_deployments`).
* `AR02` is filtered out by applicability selection (`untrusted_networks` is not adopted by the policy).
* `AR03` is excised by explicit declaration in `exclusions`.
* Control `CTL02` carries the attached constraint `CST-CTL02-01`.

```yaml
title: "Container Management Security Controls (Resolved)"
metadata:
  id: "SEC.SLAM.CM"
  type: "ControlCatalog"
  gemara-version: "1.2.0"
  version: "1.0.0"
  description: "Control catalog for container management security."
  author:
    id: "example-author"
    name: "Security Engineering"
    type: "Human"
  applicability-groups:
    - id: "all_deployments"
      title: "All Deployments"
      description: "Requirements that apply to all deployments."

groups:
  - id: "SEC.SLAM.CM.FAM01"
    title: "Network and Image Transport Security"
    description: "Controls governing network transport and image integrity."

controls:
  - id: "SEC.SLAM.CM.CTL02"
    title: "Require TLS/SSL and Digest Verification"
    objective: "Protect registry traffic and verify artifact integrity."
    group: "SEC.SLAM.CM.FAM01"
    constraints:
      - id: "CST-CTL02-01"
        target-id: "SEC.SLAM.CM.CTL02"
        text: "Registry connections MUST enforce TLS 1.3 or higher with corporate root CA pinning."
    assessment-requirements:
      - id: "SEC.SLAM.CM.CTL02.AR01"
        text: "The system MUST use TLS/SSL for all registry communication."
        applicability: ["all_deployments"]
      - id: "SEC.SLAM.CM.CTL02.AR05"
        text: "The system MUST log all container image retrieval events."
        applicability: ["all_deployments"]
```
