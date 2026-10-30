---
layout: page
title: Tool Specifications
nav-title: Specifications
---

Gemara specifications define normative operational behavior for tools, libraries, and SDKs implementing the Gemara framework. While the [Model](../model/) provides foundational conceptual definitions and the [Schemas](../schema/) provide syntactic data validation rules, **Tool Specifications** define how implementations must deterministically ingest, resolve, execute, and exchange Gemara artifacts.

## Purpose

To ensure that independent implementations, such as the [Go SDK](../sdk/go-sdk.html), [Python SDK](../sdk/python-sdk.html), policy engines, CI/CD gates, and compliance platforms, produce identical, reproducible outcomes when processing the same inputs, Gemara publishes normative specifications for tool authors.

## Planned Specifications

* **Requirements Resolution:** Defines how a Gemara Policy deterministically resolves its imported Guidance and Control Catalogs, applies assessment-requirement modifications (`Add`, `Modify`, `Remove`, `Replace`, `Override`), and generates a derived Effective Requirements set.
* **OCI Bundle Distribution & Packaging:** Standards for packaging, referencing, and fetching Gemara artifacts from OCI-compliant container registries using ORAS and standard media types.
* **Assessment Execution & Adherence Evaluation:** Normative guidelines for automated evaluation methods, evidence ingestion, and gate enforcement based on resolved policy adherence configurations.
* **Cross-Framework Mapping Resolution:** Deterministic traversal of mapping documents across standards (e.g., NIST SP 800-53, ISO/IEC 27001, OpenSSF Scorecard).

## Conformance & Implementations

Specifications in this section use RFC 2119 requirement levels (`MUST`, `MUST NOT`, `SHOULD`, `SHOULD NOT`, `MAY`). Implementations claiming conformance to a Gemara specification must satisfy all unconditional requirements and pass the published test vectors.

Tool developers and SDK contributors should reference these specifications when developing:
* **Parsers and Resolvers:** Loading and resolving composite policies.
* **Evaluation Engines:** Automated assessment of effective requirements.
* **Reporting & Auditing Tools:** Tracking requirement provenance and compliance baselines.

---

For foundational theory, see [The Model](../model/). For schema structures, see [The Schemas](../schema/). For language libraries, see [The SDKs](../sdk/).
