---
layout: page
title: Requirements Resolution Specification
nav-title: Requirements Resolution
---

This specification defines how a Gemara **[Policy](../model/05.3-Layer-3.html)** (Layer 3) resolves a single **[Control Catalog](../model/05.2-Layer-2.html)** import into an **Effective Requirements Set**.

The behavior defined here reflects the Gemara schemas (`policy.cue`, `controlcatalog.cue`, `guidancecatalog.cue`).

---

## 1. Resolution Inputs

A conforming resolver accepts:
1. **The Policy Document:** The primary Layer 3 artifact containing `imports.catalogs[]` and `metadata.mapping-references`.
2. **The Source Control Catalog:** The Layer 2 artifact identified by the import's `reference-id`.

### Schema References
* `policy.cue`: `#Imports`, `#CatalogImport`, `#Constraint`, `#AssessmentRequirementModifier`, `#ModType`.
* `controlcatalog.cue`: Defines the controls and assessment requirements of the derived output.

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

### Step 1: Artifact Acquisition
The resolver locates the catalog using `policy.metadata.mapping-references` matching `policy.imports.catalogs[].reference-id`. The referenced catalog is loaded as the initial candidate set.

### Step 2: Apply Exclusions
If `exclusions` are defined on the import:
* Excluded controls or requirements are removed from the candidate set before any modifications are evaluated.

### Step 3: Apply Constraints
If `constraints` are defined on the import:
* The resolver verifies that the imported catalog satisfies each constraint (such as minimum version thresholds). If a constraint is violated, resolution fails.

### Step 4: Apply Assessment-Requirement Modifications
Modifications declared under `assessment-requirement-modifications` are evaluated sequentially in declaration order against the candidate catalog.

---

## 3. Modification Types (`#ModType`)

Gemara defines five modification types for assessment requirements:

| Type | Target Requirement | Text Field | Description / Effect |
| :--- | :--- | :--- | :--- |
| **`Add`** | MUST NOT exist | Required | Appends a brand-new assessment requirement to the target control. |
| **`Modify`** | MUST exist | Optional | Partially updates specified fields (such as text or applicability); untouched fields retain original values. |
| **`Remove`** | MUST exist | N/A | Excises or deactivates the requirement from the effective set. |
| **`Replace`** | MUST exist | Required | Completely substitutes the requirement definition; unmentioned original fields are cleared. |
| **`Override`** | MUST exist | Required | Enforces a strict organizational policy mandate over the requirement (requires `modification-rationale`). |

{% include themed-image.html
   light="/assets/diagrams/requirements-modification-decision-tree-light.svg"
   dark="/assets/diagrams/requirements-modification-decision-tree-dark.svg"
   alt="Modification Type Decision Tree"
   caption="Figure 2: Assessment Requirement Modification Processing Order"
   class="spec-diagram-decision"
%}

### Applicability Rule
* When a modification supplies `applicability`, it **replaces** the target requirement's applicability list.
* When the modification omits `applicability`, the target requirement's original applicability list is **preserved**.

---

## 4. Worked Example & Test Vectors

The following YAML files demonstrate a complete single catalog import resolution.

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

        # Add AR04: Inject a new requirement
        - id: "MOD-AR04"
          target-id: "SEC.SLAM.CM.CTL02.AR04"
          modification-type: "Add"
          modification-rationale: "Enforce Sigstore cosign cryptographic verification."
          text: "The system MUST verify signatures using the organizational Sigstore keyless cosign root."
          applicability: ["production"]
```

### 4.3 Output: Resolved Effective Requirements Set

```yaml
title: "Resolved Effective Requirements for org-policy-001"
metadata:
  resolved-from-policy: "org-policy-001"
  gemara-version: "1.2.0"

controls:
  - id: "SEC.SLAM.CM.CTL02"
    title: "Require TLS/SSL with Certificate Pinning"
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
