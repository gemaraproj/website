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

| Type | Target Requirement | Resolver Text Rule | Description / Effect |
| :--- | :--- | :--- | :--- |
| **`Add`** | MUST NOT exist | Required | Appends a brand-new assessment requirement to the target control. |
| **`Modify`** | MUST exist | Optional | Partially updates specified fields (`text`, `applicability`, `recommendation`); untouched fields retain baseline values. |
| **`Remove`** | MUST exist | N/A | Excises or deactivates the requirement from the effective set. |
| **`Replace`** | MUST exist | Required | Completely substitutes the requirement definition; unmentioned optional fields (`recommendation`) are cleared. |
| **`Override`** | MUST exist | Required | Enforces an authoritative organizational policy mandate superseding the baseline requirement. |

{% include themed-image.html
   light="/assets/diagrams/requirements-modification-decision-tree-light.svg"
   dark="/assets/diagrams/requirements-modification-decision-tree-dark.svg"
   alt="Modification Type Decision Tree"
   caption="Figure 2: Assessment Requirement Modification Processing Order"
   class="spec-diagram-decision"
%}

### Resolver Validation & Schema Rules

1. **Modification Rationale:** The schema requires `modification-rationale` on every modifier (`#AssessmentRequirementModifier`). Resolvers MUST reject any modification entry missing a rationale.
2. **Text Field Validation:** In the CUE schema, `text` is optional (`text?: string`), but conforming resolvers enforce operational validation rules:
   * `Add`, `Replace`, and `Override` MUST provide `text`.
   * `Modify` MAY omit `text` if only patching `applicability` or `recommendation`.
   * `Remove` ignores `text`.
3. **Recommendation Handling:**
   * `Modify` updates `recommendation` if specified, or preserves the baseline value if omitted.
   * `Replace` and `Override` clear any existing `recommendation` unless explicitly provided in the modifier.
4. **`Override` vs. `Replace`:** While both substitute requirement text, `Replace` represents local adaptation within the existing governance framework, whereas `Override` represents an explicit policy mandate that supersedes the source standard (flagged for audit review).
5. **Parent Control Association for `Add` (Open Question):**
   * `#AssessmentRequirementModifier` defines `target-id` as the ID of the requirement. The schema does not currently define a dedicated `control-id` field on the modifier (tracked in `gemaraproj/gemara#492`).
   * By convention, resolvers associate an `Add` requirement with its parent control based on hierarchical ID prefix matching (e.g., `SEC.SLAM.CM.CTL02.AR04` matches control `SEC.SLAM.CM.CTL02`).

### Applicability Rule
* When a modification supplies `applicability`, it **replaces** the target requirement's applicability list.
* When the modification omits `applicability`, the target requirement's original applicability list is **preserved**.

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

controls:
  - id: "SEC.SLAM.CM.CTL02"
    title: "Require TLS/SSL with Certificate Pinning"
    objective: "Protect registry traffic and verify artifact integrity."
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
  mapping-references:
    - id: "SEC.SLAM.CM"
      title: "Container Management Security Controls"
      version: "1.0.0"

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

controls:
  - id: "SEC.SLAM.CM.CTL02"
    title: "Require TLS/SSL with Certificate Pinning"
    objective: "Protect registry traffic and verify artifact integrity."
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
