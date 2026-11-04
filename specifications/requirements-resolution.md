---
layout: page
title: Requirements Resolution Specification
nav-title: Requirements Resolution
---

This specification defines how a Gemara **[Policy](../model/05.3-Layer-3.html)** (Layer 3) resolves a single **[Control Catalog](../model/05.2-Layer-2.html)** import into an **Effective Requirements Set** (represented as a derived `ControlCatalog`).

The behavior defined here reflects the Gemara schemas (`policy.cue`, `controlcatalog.cue`).

> **Scope & Deferred Items:**
> * **Single Catalog Imports:** This document covers the resolution of a single `imports.catalogs[]` entry.
> * **Guidance Imports (`imports.guidance[]`):** Guidance imports accept exclusions and constraints, but do not support requirement modifications. Their resolution will be detailed in a follow-up specification.
> * **Policy Imports (`imports.policies[]`):** A policy can import another policy as an inheritance model. An importing policy brings in the referenced policy's resolved requirements wholesale and cannot alter them. A modifier of any type targeting an inherited policy requirement is a resolution error.
> * **Static Resolution and Policy Scope:** Static resolution does not filter requirements based on the policy's `scope`. Policy `scope` defines macro organizational boundaries, while requirement `applicability` governs evaluation-time check activation (tracked in `gemaraproj/gemara#548`).
> * **Applicability Namespacing & Collisions:** Prefixing applicability groups with the source catalog's `metadata.id` (e.g., `SEC.SLAM.CM:all_deployments`) and handling collisions across multiple imports are deferred to the multi-catalog resolution specification. In this single-catalog baseline specification, applicability identifiers remain in their local catalog-scoped form.

---

## 1. Resolution Inputs

A conforming resolver accepts:
1. **The Policy Document:** The primary Layer 3 artifact containing `imports.catalogs[]` and `metadata.mapping-references`.
2. **The Source Control Catalog:** The Layer 2 artifact identified by the import's `reference-id`.

### Schema References
* `policy.cue`: `#Imports`, `#CatalogImport`, `#Constraint`, `#AssessmentRequirementModifier`, `#ModType`.
* `controlcatalog.cue`: `#ControlCatalog`, `#Control`, `#AssessmentRequirement`.

---

## 2. Resolution Workflow (Single Catalog Import)

Resolution of an imported catalog proceeds in four deterministic steps:

{% include themed-image.html
   light="/assets/diagrams/requirements-resolution-pipeline-light.svg"
   dark="/assets/diagrams/requirements-resolution-pipeline-dark.svg"
   alt="Requirements Resolution Pipeline"
   caption="Figure 1: Single Catalog Import Resolution Pipeline"
   class="spec-diagram-pipeline"
%}

### Step 1: Artifact Acquisition & Identification
The resolver locates the catalog using `policy.metadata.mapping-references` matching `policy.imports.catalogs[].reference-id`.
* The mapping reference `id` MUST match the referenced artifact's `metadata.id`.
* When `version` is specified in the mapping reference, it MUST match the artifact's `metadata.version`.
* The transport mechanism used to retrieve the artifact (local filesystem, HTTP, OCI registry) is left to the tool implementation.

### Step 2: Apply Exclusions
If `exclusions` are defined on the import:
* Excluded control IDs or assessment requirement IDs are excised from the candidate catalog before constraints or modifications are processed.

### Step 3: Apply Constraints
If `constraints` are declared on the import:
* In `policy.cue`, `#Constraint` is defined as `{id, target-id, text}`: prescriptive requirement text attached to a specific control or guidance item (for example, "Review access at least annually").
* The resolver attaches these constraint statements to the matching target control in the candidate set.

### Step 4: Apply Assessment-Requirement Modifications
Modifications declared under `assessment-requirement-modifications` are evaluated against the candidate catalog.
* **Cardinality Rule:** A policy MUST have at most one modification per target requirement. Specifying multiple modifications targeting the same requirement ID is a resolution error.

---

## 3. Modification Types (`#ModType`)

Gemara defines five modification types for assessment requirements:

| Type | Target Requirement | Description / Effect |
| :--- | :--- | :--- |
| **`Add`** | MUST NOT exist | Appends a brand-new assessment requirement to the target control. |
| **`Modify`** | MUST exist | Partially updates specified fields (`text`, `applicability`, `recommendation`); untouched fields retain baseline values. |
| **`Remove`** | MUST exist | Excises or deactivates the requirement from the effective set. |
| **`Replace`** | MUST exist | Completely substitutes the requirement definition with the fields provided in the modifier. |
| **`Override`** | MUST exist | Enforces an authoritative organizational policy mandate superseding the baseline requirement. |

{% include themed-image.html
   light="/assets/diagrams/requirements-modification-decision-tree-light.svg"
   dark="/assets/diagrams/requirements-modification-decision-tree-dark.svg"
   alt="Modification Type Decision Tree"
   caption="Figure 2: Assessment Requirement Modification Processing Order"
   class="spec-diagram-decision"
%}

### Schema Rules & Open Design Questions

1. **Modification Rationale:** The schema requires `modification-rationale` on every modifier (`#AssessmentRequirementModifier`). Conforming policies MUST provide a non-empty rationale for each modification.
2. **Cardinality:** A policy MUST declare at most one modifier per target requirement. Multiple modifiers targeting the same requirement ID are a resolution error.
3. **Applicability Replacement Rule:** When a modifier supplies `applicability`, it replaces the target requirement's applicability list. When omitted, the baseline applicability is preserved.
4. **Single-Catalog Baseline Identifiers:** In single catalog imports, applicability identifiers remain in their local form (e.g., `all_deployments`). Cross-catalog namespacing (`<catalog-id>:<group-name>`) is deferred to multi-catalog resolution.

#### Open Questions (for Schema Stabilization)

As the experimental Policy schema progresses toward stabilization, the following operational questions are under active maintainer review:

* **Text and Field Requirements:** In `policy.cue`, `text` and `recommendation` are optional on `#AssessmentRequirementModifier`. Whether resolvers should enforce that `Add`, `Replace`, and `Override` provide `text`, or whether `Replace`/`Override` should implicitly clear unmentioned optional fields like `recommendation`, remains an open design decision.
* **`Override` vs. `Replace` Output Semantics:** Currently, `Replace` and `Override` can produce identical resolved document outputs. Whether the schema should distinguish them structurally (such as emitting an audit or policy-mandate flag) or consolidate them into fewer modifier options is under evaluation.
* **Parent Control Association for `Add` ([gemaraproj/gemara#492](https://github.com/gemaraproj/gemara/issues/492)):** `#AssessmentRequirementModifier` specifies only `target-id`. Resolvers currently associate added requirements with parent controls using ID prefix matching (e.g., `SEC.SLAM.CM.CTL02.AR04` $\rightarrow$ `SEC.SLAM.CM.CTL02`). Whether to add an explicit `control-id` to the schema is tracked in #492.
* **Applicability Provenance & Undeclared Groups ([gemaraproj/gemara#548](https://github.com/gemaraproj/gemara/issues/548)):** `ControlCatalog` requires all requirement `applicability` values to be defined in `metadata.applicability-groups`. When a policy modifier introduces an applicability tag (e.g. `edge_gateways`), whether those groups must be pre-declared in the source catalog, inherited from policy-level `applicability-groups`, or rejected as a validation error is an open question.

---

## 4. Worked Example & Test Vectors

The following example demonstrates resolving a single catalog import. All fixtures conform to the Gemara schemas.

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
    - id: "edge_gateways"
      title: "Edge Gateways"
      description: "Applies to perimeter and edge gateway deployments."

groups:
  - id: "SEC.SLAM.CM.FAM01"
    title: "Network and Image Transport Security"
    description: "Controls governing network transport and image integrity."

controls:
  - id: "SEC.SLAM.CM.CTL02"
    title: "Require TLS/SSL with Certificate Pinning"
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
```

### 4.2 Input: Policy with Import & Modifications (`org-policy-001`)

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
      assessment-requirement-modifications:
        # Override AR01: Strict mandate with updated text and replaced applicability
        - id: "MOD-AR01"
          target-id: "SEC.SLAM.CM.CTL02.AR01"
          modification-type: "Override"
          modification-rationale: "Mandate certificate pinning across all environments."
          text: "The system MUST use TLS/SSL and MUST pin to the trusted registry certificate authority."
          applicability: ["all_deployments", "edge_gateways"]

        # Modify AR02: Update text only (applicability is preserved!)
        - id: "MOD-AR02"
          target-id: "SEC.SLAM.CM.CTL02.AR02"
          modification-type: "Modify"
          modification-rationale: "Permit WireGuard tunnels in addition to corporate VPN."
          text: "On untrusted networks, the system MUST use an approved WireGuard or IPSec VPN."

        # Add AR04: Inject a new requirement under CTL02
        - id: "MOD-AR04"
          target-id: "SEC.SLAM.CM.CTL02.AR04"
          modification-type: "Add"
          modification-rationale: "Enforce Sigstore cosign cryptographic verification."
          text: "The system MUST verify signatures using the organizational Sigstore keyless cosign root."
          applicability: ["production"]

adherence:
  evaluation-methods:
    - id: "EV-AUTO-01"
      type: "Behavioral"
      mode: "Automated"
      required: true
      description: "Automated CI/CD policy verification checks."
```

### 4.3 Output: Resolved Effective Control Catalog

The resolved effective requirements set forms a schema-valid `ControlCatalog`:

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
    - id: "untrusted_networks"
      title: "Untrusted Networks"
      description: "Applies when communication traverses untrusted networks."
    - id: "production"
      title: "Production"
      description: "Applies to live production workloads."
    - id: "edge_gateways"
      title: "Edge Gateways"
      description: "Applies to perimeter and edge gateway deployments."

groups:
  - id: "SEC.SLAM.CM.FAM01"
    title: "Network and Image Transport Security"
    description: "Controls governing network transport and image integrity."

controls:
  - id: "SEC.SLAM.CM.CTL02"
    title: "Require TLS/SSL with Certificate Pinning"
    objective: "Protect registry traffic and verify artifact integrity."
    group: "SEC.SLAM.CM.FAM01"
    assessment-requirements:
      # AR01: Overridden (text updated, applicability replaced)
      - id: "SEC.SLAM.CM.CTL02.AR01"
        text: "The system MUST use TLS/SSL and MUST pin to the trusted registry certificate authority."
        applicability: ["all_deployments", "edge_gateways"]

      # AR02: Modified text (original applicability preserved)
      - id: "SEC.SLAM.CM.CTL02.AR02"
        text: "On untrusted networks, the system MUST use an approved WireGuard or IPSec VPN."
        applicability: ["untrusted_networks"]

      # AR03: Unmodified baseline requirement
      - id: "SEC.SLAM.CM.CTL02.AR03"
        text: "The system MUST verify image digests before execution."
        applicability: ["all_deployments"]

      # AR04: Newly added requirement appended to control
      - id: "SEC.SLAM.CM.CTL02.AR04"
        text: "The system MUST verify signatures using the organizational Sigstore keyless cosign root."
        applicability: ["production"]
```
