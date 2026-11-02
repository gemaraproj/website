---
layout: page
title: "Specification: Cross-Framework Mapping Resolution"
---

## Status

* **Status:** Draft / Active
* **Type:** Tool Specification
* **Domain:** Compliance Mapping & Graph Resolution
* **Applies To:** Gemara Tooling, SDKs (`go-gemara`), GRC Engines, Auditing & Reporting Platforms

---

## 1. Abstract

Gemara uses `MappingDocument` (Layer 4) and inline mapping structures to connect security controls, assessment requirements, guidelines, threats, and risks to external compliance frameworks (such as NIST SP 800-53, ISO/IEC 27001, OpenSSF Scorecard, and CIS Controls).

While `mappingdocument.cue` specifies the syntactic schema of mapping documents, compliance tools, reporting dashboards, and auditors require a normative specification defining how compliant tools deterministically ingest, traverse, combine, and resolve cross-framework relationships across complex multi-standard graphs.

This specification establishes normative algorithms for mapping traversal, relationship composition algebra, conflict handling, cycle prevention, and test vector verification.

---

## 2. Terminology & Input Artifacts

* **Source Artifact:** The origin document defining entries being mapped from (e.g., an internal ControlCatalog).
* **Target Artifact:** The destination framework or regulation being mapped to (e.g., NIST SP 800-53 Rev 5).
* **Mapping Entry:** An individual atomic statement associating a source entry (`source: string`) with one or more target entries (`targets: []#MappingTarget`) under a declared `relationship`.
* **Intermediate Artifact:** A bridging document connecting Framework A to Framework B, used when evaluating transitive chains (e.g., Framework A &rarr; Framework B &rarr; Framework C).
* **Resolved Cross-Reference Matrix:** The normalized, computed evaluation report stating coverage, strength, and confidence between entries across frameworks.

---

## 3. Relationship Types & Mathematical Semantics

Conformant resolution engines **MUST** evaluate mapping relationships according to the following semantics defined by Gemara Layer 4:

| Relationship Type | Symbol | Mathematical & Compliance Semantics |
| :--- | :---: | :--- |
| `equivalent` | $\equiv$ | Symmetric and bidirectional equivalence: fulfilling Source satisfies Target, and fulfilling Target satisfies Source. |
| `implements` | $\vdash$ | Directed satisfaction: Source fully operationalizes and satisfies Target. |
| `implemented-by` | $\dashv$ | Inverse satisfaction: Source is satisfied by Target (equivalent to Target $\vdash$ Source). |
| `subsumes` | $\supseteq$ | Scope inclusion: Source encompasses Target plus additional scope (Target is a subset of Source). |
| `supports` | $\rightarrow$ | Partial contribution: Source contributes towards fulfilling Target, but is not alone sufficient. |
| `supported-by` | $\leftarrow$ | Inverse partial contribution: Target contributes towards Source. |
| `relates-to` | $\sim$ | Contextual relevance: Semantic association without asserting coverage or fulfillment. |
| `no-match` | $\emptyset$ | Explicit non-alignment: Source has no counterpart in the target framework (`targets` omitted). |

---

## 4. Traversal & Derivation Algorithm

When a tool evaluates coverage for a source control across multi-hop mapping chains (e.g., finding the compliance posture against Framework C when mappings exist only from A &rarr; B and B &rarr; C), it **MUST** execute the deterministic traversal algorithm described below.

### 4.1 Relationship Composition Algebra

When composing two consecutive edges $R_1(A \rightarrow B)$ and $R_2(B \rightarrow C)$, the resulting relationship $R_{comp}(A \rightarrow C)$ **MUST** follow the conservative composition table:

| $R_1 \downarrow \quad R_2 \rightarrow$ | `equivalent` | `implements` | `subsumes` | `supports` | `relates-to` |
| :--- | :---: | :---: | :---: | :---: | :---: |
| **`equivalent`** | `equivalent` | `implements` | `subsumes` | `supports` | `relates-to` |
| **`implements`** | `implements` | `implements` | `supports` | `supports` | `relates-to` |
| **`subsumes`** | `subsumes` | `supports` | `subsumes` | `supports` | `relates-to` |
| **`supports`** | `supports` | `supports` | `supports` | `supports` | `relates-to` |
| **`relates-to`** | `relates-to` | `relates-to` | `relates-to` | `relates-to` | `relates-to` |

*Rule:* Any composition involving `no-match` produces `no-match`. Transitive derivation **SHALL NOT** promote relationship strength higher than the weakest link in the traversal path.

### 4.2 Attribute Attenuation Along Paths

When intermediate mappings declare numerical `strength` (range 1–10) or qualitative `confidence-level` (`Low`, `Medium`, `High`), derived multi-hop edges **MUST** attenuate:

1. **Derived Strength:** The effective strength along path $P = (e_1, e_2, \dots, e_n)$ **MUST** be the minimum strength of any edge along the path:
   $$\text{strength}(P) = \min_{e \in P}(\text{strength}(e))$$
2. **Derived Confidence:** The effective confidence **MUST** take the lowest confidence grade present on the path ($\text{Low} < \text{Medium} < \text{High}$).

---

## 5. Conflict Resolution & Cycle Prevention

### 5.1 Graph Cycle Prevention
Because cross-framework mappings are authored independently by multiple organizations, circular references (e.g., $A \rightarrow B \rightarrow C \rightarrow A$) are common in combined corpora.

1. Resolution engines **MUST** maintain a visited node set during traversal.
2. If an edge returns to an already-visited entry in the current path, traversal along that branch **MUST** terminate immediately. Circular traversals **SHALL NOT** result in infinite recursion.

### 5.2 Conflicting Mapping Assertions
When multiple mapping documents make differing assertions between the exact same pair of $(Entry_A, Entry_B)$:

1. **Precedence by Specificity:** If one document maps specific `AssessmentRequirement` units and another maps whole `Control` units, the finer-grained requirement-level mapping **MUST** take precedence.
2. **Tie-Breaking:** If two documents make incompatible claims at the same granularity (e.g., Doc 1 claims `equivalent` while Doc 2 claims `supports`):
   - Conforming tools **MUST** record a `MappingConflictWarning`.
   - In automated scoring modes, tools **MUST** adopt the conservative minimum relationship (e.g. `supports` over `equivalent`).

---

## 6. Conformance Test Vectors

Conformant resolution implementations can verify their derivation logic against these canonical test cases:

### Test Vector 1: Transitive Satisfaction
* **Inputs:**
  * Control `OSPS-AC-01` $\vdash$ `CRA-AnnexI-Part1-2.1` (`implements`, strength: 8)
  * `CRA-AnnexI-Part1-2.1` $\equiv$ `NIS2-Article21-2a` (`equivalent`, strength: 10)
* **Expected Output:**
  * `OSPS-AC-01` $\vdash$ `NIS2-Article21-2a` (`implements`, resolved strength: 8)

### Test Vector 2: Attenuation to Support
* **Inputs:**
  * Control `OSPS-QA-02` $\supseteq$ `OWASP-A06` (`subsumes`, strength: 7)
  * `OWASP-A06` $\vdash$ `NIST-800-53-SI-2` (`implements`, strength: 9)
* **Expected Output:**
  * `OSPS-QA-02` $\rightarrow$ `NIST-800-53-SI-2` (`supports`, resolved strength: 7)

---

## 7. Implementation & SDK Guidance

* Language SDKs (`go-gemara`) **SHOULD** provide an in-memory graph structure indexing mapping documents by `(source-reference, target-reference)` tuples.
* CLI tools (`grcli`) evaluating policies against multiple standard baselines **MUST** output clear provenance metadata detailing which mapping documents and intermediate traversal paths established each compliance claim.
