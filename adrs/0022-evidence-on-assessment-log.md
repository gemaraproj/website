---
layout: page
title: Add Evidence to AssessmentLog
description: "ADR-0022: reshape the #Evidence type so evaluation and audit logs both record the evidence behind their opinions."
---

- **ADR:** 0022
- **Proposal Author(s):** @jpower432, @eddie-knight
- **Status:** Accepted

## Context

Gemara's measurement layers produce opinions, structured conclusions about the state of a target. An `EvaluationLog` is a tool's opinion: "based on what I observed, this control passes or fails." An `AuditLog` is an auditor's opinion: "based on the evidence I reviewed, this control is effective or not." Both opinions are rooted in evidence, the data that was actually observed that helped inform the actual conclusion. The tool sees it firsthand (an API response, a file, a config dump). The auditor sees it through the tool's record. Different proximity to the raw data, same root.

Currently, the evidence root is invisible. An `AssessmentLog` records a pass/fail result with a message, but not "I consulted this specific resource at this time." The tool's opinion exists without a record of what it was rooted in, limiting downstream consumers of the log.

## Action

Reshape the existing `#Evidence` type in the experimental `auditlog.cue` to serve both the evaluation layer (`AssessmentLog`) and the audit layer (`AuditLog`). 
The `#Evidence` type captures what was cited to support an opinion for a specific activity, raw data for the tool, evaluation and enforcement artifacts for the auditor.

```cue
#Evidence: {
    // id uniquely identifies this evidence
    id?: string

    // type categorizes the kind of evidence (open enum)
    type: #EvidenceType

    // collected is the timestamp when the evidence was gathered
    "collected-at": #Datetime

    // address identifies where the evidence can be found
    address?: string

    // payload is the raw evidence data collected
    payload?: _

    // digest is a hash of the evidence content at collection time for integrity verification
    // Enables verification that mutable-address evidence (S3 objects, HTTP URLs, API responses) has not changed since the tool observed it.
    // Not needed for content-addressable systems (OCI, git) or inline payloads.
    digest?: string

    // description explains what this evidence represents
    description?: string
}
```
`#EvidenceType` remains an open enum. Recommended values include artifact types already known to Gemara (e.g., `EvaluationLog`, `EnforcementLog`) plus categories for common evidence forms like provenance.

**Changes to the Assessment Log**:
```cue
#AssessmentLog: {
    // ... existing fields ...

    // evidence records the raw data cited to support this assessment's opinion
    evidence?: [#Evidence, ...#Evidence]
}
```

The field is optional (one or more when present) to maintain backward compatibility with existing EvaluationLog documents.

## Consequences

### Positive

- Assessment logs can record exactly what the tool consulted, when, and what was observed, making the root of the tool's opinion inspectable.
- The `address` field enables SARIF converters to populate `PhysicalLocation` and `Region` with file paths and line numbers from evidence data.
- `address` as a plain string works for file paths, API endpoints, ARNs, container image references, physical locations, and any other evidence source.

### Negative

Embedding raw evidence via the `payload` field can make assessment logs large. 
It is the evaluator's responsibility to manage payload size, trim sensitive data, and ensure evidence does not cross trust boundaries inappropriately.

## Alternatives Considered

### Separate `#Provenance` type

A standalone `#Provenance` type on `#AssessmentLog` with fields for `address`, `collected`, and `observed`. This was the original proposal in [issue #417](https://github.com/gemaraproj/gemara/issues/417).

**Rejected because:** Both layers produce opinions rooted in evidence. The tool cites raw data to support its opinion. The auditor cites evaluation artifacts to support theirs. Introducing a separate type creates an artificial distinction when the relationship (opinion rooted in evidence) is the same at both layers.


