---
layout: page
title: Requirements Resolution Specification
nav-title: Requirements Resolution
---

This specification defines how a Gemara **[Policy](../model/05.3-Layer-3.html)** (Layer 3) resolves a single **[Control Catalog](../model/05.2-Layer-2.html)** import into an **Effective Requirements Set** (represented as a derived `ControlCatalog`).

The normative behavior defined here reflects the Gemara schemas (`policy.cue`, `controlcatalog.cue`) and the profile model established in [gemaraproj/gemara#560](https://github.com/gemaraproj/gemara/issues/560), [gemaraproj/gemara#548](https://github.com/gemaraproj/gemara/issues/548), and [gemaraproj/gemara#554](https://github.com/gemaraproj/gemara/pull/554).

> **Scope & Deferred Items:**
> * **Single Catalog Imports:** This document specifies resolution for a single `imports.catalogs[]` entry.
> * **Guidance Imports (`imports.guidance[]`):** Guidance imports accept exclusions and constraints following the same pipeline; their multi-artifact resolution will be detailed in a follow-up specification.
> * **Policy Imports (`imports.policies[]`):** A policy can import another policy as an inheritance model. An importing policy brings in the referenced policy's resolved requirements wholesale and cannot alter them.
> * **Static Resolution and Policy Scope:** Policy `scope` defines organizational, architectural, and geopolitical boundaries. Requirements resolution is governed statically and exclusively by declared `exclusions`, imported catalog `applicability` selection, and `constraints`.
> * **Cross-Catalog Namespacing:** Prefixing identifiers across multiple catalog imports is deferred to the multi-catalog resolution specification. In this single-catalog baseline specification, identifiers remain in their local catalog-scoped form.

---

## 1. Resolution Inputs

A conforming resolver accepts:
1. **The Policy Document:** The Layer 3 artifact containing `imports.catalogs[]` and `metadata.mapping-references`.
2. **The Source Control Catalog:** The Layer 2 artifact identified by the import's `reference-id`.

### Schema References
* `policy.cue`: `#Imports`, `#CatalogImport`, `#Constraint`.
* `controlcatalog.cue`: `#ControlCatalog`, `#Control`, `#AssessmentRequirement`.

### Schema Validation vs. Tool Resolution Rules

Implementations must recognize the distinction between rules enforced syntactically by CUE schemas versus semantic rules enforced by tool resolvers:
* **Enforced by CUE:** Schema structure, required types, and field boundaries (e.g., `#CatalogImport` declaring `exclusions?: [...string]`, `#Constraint` requiring `id`, `target-id`, and `text`, `#MappingReference` requiring `id` and `version`).
* **Enforced by Tool Resolvers:** Operational semantics that CUE schemas do not evaluate across document boundaries. For example:
  * Resolving `policy.imports.catalogs[].reference-id` against `policy.metadata.mapping-references`.
  * Verifying that each ID in `exclusions` actually exists in the acquired catalog.
  * Verifying that each ID in `#CatalogImport.applicability` exists in `metadata."applicability-groups"`.
  * Verifying that `#Constraint.target-id` references a control that remains in play.
  * Pruning empty controls after requirement selection.

---

## 2. Resolution Workflow (Single Catalog Import)

A Policy is a **profile**: it selects, constrains, and parametrizes; it does not rewrite catalog content ([gemaraproj/gemara#560](https://github.com/gemaraproj/gemara/issues/560)). Resolution follows a deterministic 4-step pipeline: **Acquire ➔ Select ➔ Exclude ➔ Constrain**.

{% include themed-image.html
   light="/assets/diagrams/requirements-resolution-pipeline-light.svg"
   dark="/assets/diagrams/requirements-resolution-pipeline-dark.svg"
   alt="Requirements Resolution Pipeline"
   caption="Figure 1: Deterministic Requirements Resolution Pipeline"
   class="spec-diagram-pipeline"
%}

### Step 1: Artifact Acquisition & Identification
The resolver locates the catalog using `policy.metadata.mapping-references` matching `policy.imports.catalogs[].reference-id`.
* `#MappingReference.version` is required by schema and MUST match the acquired catalog's `metadata.version`.
* When `#MappingReference.url` is omitted, `#MappingReference.id` MUST equal the referenced catalog's `metadata.id`. When `url` is present, `id` acts as a local alias pointing to the artifact located at `url`.
* The transport mechanism used to retrieve the artifact (local filesystem, HTTP, OCI registry) is implementation-specific.

### Step 2: Select and Exclude
Selection and exclusion run against the acquired source catalog to establish which requirements are in play:

1. **Applicability Selection ([gemaraproj/gemara#554](https://github.com/gemaraproj/gemara/pull/554)):**
   * Selection comes strictly from `applicability` declared on `#CatalogImport`, not from the policy's own metadata.
   * Each ID in `#CatalogImport.applicability` MUST match an entry in the imported catalog's `metadata."applicability-groups"`. If the import specifies an applicability group ID that the catalog never defined, the resolver MUST reject the policy with a resolution error.
   * When `applicability` is specified on the import, a requirement is **in play** if any of its declared `applicability` tags matches a selected group. Requirements with no matching applicability tag are dropped.
   * If `applicability` is omitted on the import, all requirements are candidate in-play (unfiltered by applicability groups).
2. **Exclusions:**
   * An entry in `exclusions` may name a control ID or an individual assessment requirement ID.
   * Every ID in `exclusions` MUST exist in the acquired catalog; naming an ID the catalog never defined is a resolution error.
   * An exclusion naming an ID that applicability selection already dropped is harmless and accepted.
   * Any control or assessment requirement matching an entry in `exclusions` is excised (dropped completely).
3. **Pruning Empty Controls:**
   * A control that has no assessment requirements remaining in play after selection and exclusion is dropped from the resolved catalog.

### Step 3: Apply Constraints
If `constraints` are declared on the import:
* In `policy.cue`, `#Constraint` is defined as `{id: string, target-id: string, text: string}`: prescriptive requirement text attached to a specific control (for catalog imports) or guidance item (for guidance imports).
* Every constraint's `target-id` MUST target a control that remains **in play** after Step 2. A constraint that targets an ID that was excluded, pruned (because all its assessment requirements were dropped), or never defined in the catalog is a **resolution error**.
* Constraints do not mutate `#Control` in the catalog schema; rather, resolvers track validated constraints alongside the in-play set for downstream evaluation and enforcement.

### Step 4: Resolved Output & In-Play Set
The resulting set represents the effective requirements of the policy import:
* The resolved set is a strict subset of the source catalog.
* Every control and assessment requirement retains its source ID, objective, and normative text.
* When policy `assessment-plans` are declared, they MUST only target requirements that remain in play; referencing an excluded or dropped requirement is a tool validation error.

---

## 3. Applicability Semantics & Evaluation

Gemara clearly separates static policy commitment from runtime evaluation:

* **Static Policy Commitment (In-Play):**
  * Which requirements is the policy committed to? This is static and decided by the policy author.
  * Applicability selection (`#CatalogImport.applicability`) and `exclusions` answer this question.
  * When declared, policy `assessment-plans` MUST only target in-play requirements.
* **Dynamic Target Evaluation (Not Applicable):**
  * Does an in-play requirement apply to a specific target system or environment? This is dynamic and decided by the evaluator at assessment time.
  * When an in-play requirement does not apply to a target, the evaluation log records it with disposition `Not Applicable` (along with justifying evidence or rationale).
* **Role of `metadata.applicability-groups`:**
  * `applicability-groups` declared in metadata is meaningful on a **Control Catalog** (Layer 2) to categorize requirements.
  * It is inert on a **Policy** (Layer 3) or an **EvaluationLog** (Layer 4), because policies do not define requirements to tag. A policy only selects groups from the catalogs it imports.
  * In an `EvaluationLog`, an entry's applicability reflects the requirement's own catalog tags (qualified by the control's `reference-id`), not a policy-level selection.

---

## 4. Identity, Provenance, and Audit

Because the resolved set is a strict subset of the source catalog:
1. **Source Identifiers Preserved:** Every requirement retains its upstream source identifier. Resolvers do not mint synthetic IDs.
2. **Deterministic Re-derivation:** An auditor can re-derive the exact in-play requirements set at any time from the policy and the pinned catalog version declared in `metadata.mapping-references`.
3. **Cross-Document Semantic Matching:** A policy selection and a downstream log entry match when each document's `mapping-references` resolve them to the same catalog and group, not merely when arbitrary local `reference-id` strings match.
4. **Validation Rules (Normative):**
   * **Plan Coverage:** If an assessment plan references a requirement that is not in play, tools MUST flag a validation error.
   * **Log Completeness:** An evaluation log MUST NOT record results for requirements that were dropped during policy resolution.

---

## 5. Worked Example & Test Vectors

The following example demonstrates resolving a single catalog import.

### 5.1 Input: Base Control Catalog (`SEC.SLAM.CM`)

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

### 5.2 Input: Policy with Import, Selection, Exclusions, and Constraints (`org-policy-001`)

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
      # Select applicability group declared in the catalog
      applicability:
        - "all_deployments"
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
  assessment-plans:
    - id: "plan-ctl02-ar01"
      requirement-id: "SEC.SLAM.CM.CTL02.AR01"
      frequency: "every push"
      evaluation-methods:
        - id: "EV-AUTO-01"
          type: "Behavioral"
          mode: "Automated"
          required: true
    - id: "plan-ctl02-ar05"
      requirement-id: "SEC.SLAM.CM.CTL02.AR05"
      frequency: "every push"
      evaluation-methods:
        - id: "EV-AUTO-01"
          type: "Behavioral"
          mode: "Automated"
          required: true
```

### 5.3 Output: Resolved Effective Control Catalog

The resolved effective requirements set is a strict subset of the source catalog:
* `AR01` and `AR05` are in play (tagged with selected applicability group `all_deployments`).
* `AR02` is dropped (its tag `untrusted_networks` was not selected in `imports.catalogs[].applicability`).
* `AR03` is excised by explicit entry in `exclusions`.
* Policy constraint `CST-CTL02-01` applies to `CTL02` as companion policy metadata (constraints govern evaluation alongside the catalog and do not mutate `#Control` in `controlcatalog.cue`).

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
    assessment-requirements:
      - id: "SEC.SLAM.CM.CTL02.AR01"
        text: "The system MUST use TLS/SSL for all registry communication."
        applicability: ["all_deployments"]
      - id: "SEC.SLAM.CM.CTL02.AR05"
        text: "The system MUST log all container image retrieval events."
        applicability: ["all_deployments"]
```

---

## 6. Open Questions

1. **Exclusion Rationales:** Should `exclusions` support optional rationales (e.g. why a specific requirement was excluded), analogous to how adopted applicability groups carry descriptions? An excluded ID currently carries no explanation in the schema. This would require a separate schema enhancement and depends on concrete industry use cases.
2. **Constraint Representation in Derived Artifacts:** In `controlcatalog.cue`, `#Control` has no `constraints` field, so derived `ControlCatalog` outputs cannot embed constraints directly while remaining schema-valid. Today, constraints travel alongside the catalog at the Policy layer (`#CatalogImport.constraints`) or within composite packaging (e.g. an OCI bundle or runner manifest). Whether `#Control` should eventually gain an optional `constraints` field in the catalog schema or remain strictly external companion metadata is under maintainer discussion.
