---
layout: page
title: Evaluation Log Guide
description: Step-by-step guide to generating and interpreting Gemara Layer 5 evaluation logs
---

## What This Is

This guide walks through creating, structuring, and verifying an **Evaluation Log** using the [Gemara](https://gemara.openssf.org/) specification.

**The basic idea:** Where Layer 2 defines what controls and assessment requirements *should* exist, Layer 5 records the **observed truth** of an evaluation. An `EvaluationLog` captures automated test runs, compliance scans, or manual audits against specific targets (software packages, infrastructure repositories, organizations, or hardware).

In technical terms:
* **[Layer 5 Evaluation](../../model/07.2-Layer-5.html)** measures operational implementation against Layer 2 controls via `EvaluationLog`, `ControlEvaluation`, and `AssessmentLog`.
* **Target:** Every evaluation log evaluates a single target (`Software`, `Infrastructure`, `Organization`, `Hardware`, or `Service`) identified by name, version, and optional repository or deployment metadata.
* **Control Evaluations:** Each entry references a specific Layer 2 control via `control.reference-id` and documents the overall control compliance `result`.
* **Assessment Logs:** Within each control evaluation, granular assessment logs record individual test executions against specific control requirements, capturing execution duration (`start`/`end`), step counts (`steps-executed`), outcomes, and cryptographic or payload evidence.

## Walkthrough

### Step 0: Understand the Core Schemas

Layer 5 evaluation data models follow a strict hierarchy defined in `evaluationlog.cue`:

```
#EvaluationLog
├── metadata (id, gemara-version, author, target-digest)
├── target (type, name, version, uri)
├── result (Passed | Failed | Partially Passed | Inconclusive | Not Applicable | Not Run)
└── evaluations: list of #ControlEvaluation
    ├── control (reference-id, catalog-reference)
    ├── result (Passed | Failed | Partially Passed | Inconclusive | Not Applicable | Not Run)
    └── assessment-logs: list of #AssessmentLog
        ├── requirement (reference-id, description)
        ├── result (Passed | Failed | Partially Passed | Inconclusive | Not Applicable | Not Run)
        ├── start & end (ISO 8601 timestamps)
        ├── steps-executed (integer)
        ├── evidence: list of #Evidence (payload or source)
        └── observations (optional findings or metrics)
```

### Step 1: Setting Up Metadata and Target

Every evaluation log begins with artifact metadata and the target being assessed:

```yaml
metadata:
  id: "EVAL-2026-03-K8S-PROD-01"
  type: "EvaluationLog"
  gemara-version: "1.2.0"
  description: "Automated admission control and identity verification audit for production cluster"
  author:
    id: "sec-ops-automation"
    name: "Platform Security Engine"
    type: "Tool"
  created: "2026-03-31T09:45:00Z"
target:
  type: "Software"
  name: "ingress-gateway-service"
  version: "v2.14.0"
  uri: "git+https://github.com/example-org/ingress-gateway.git#refs/tags/v2.14.0"
```

| Field | Type | Description |
|-------|------|-------------|
| `metadata.id` | String | Unique run identifier for tracking and correlation across pipelines. |
| `metadata.type` | String | Must be `"EvaluationLog"`. |
| `metadata.gemara-version` | String | Specification version (e.g. `"1.2.0"`). |
| `target.type` | Enum | The kind of target being assessed: `Software`, `Infrastructure`, `Organization`, `Hardware`, or `Service`. |
| `target.name` | String | Unique name or identifier of the system under test. |
| `target.uri` | String (optional) | Canonical URI pointing to repository, container registry, or infrastructure endpoint. |

### Step 2: Defining Control Evaluations

An evaluation log aggregates one or more control evaluations under `evaluations`. Each control evaluation references the target control from a Layer 2 Control Catalog:

```yaml
evaluations:
  - control:
      reference-id: "CTRL-IMG-01"
      catalog-reference:
        id: "container-security-standard-v1"
        uri: "https://gemara.openssf.org/catalogs/container-security-v1.yaml"
    result: "Passed"
    summary: "Immutable container image tags and SHA256 digests strictly enforced across all deployment specs."
```

#### Supported Result States
Both `#EvaluationLog`, `#ControlEvaluation`, and `#AssessmentLog` support the standardized result enum:
* `Passed`: The target fully satisfies the requirement.
* `Failed`: The target violates the requirement.
* `Partially Passed`: Some steps or checks succeeded, but the complete requirement was not fulfilled.
* `Inconclusive`: The evaluation could not definitively determine compliance due to an environment, timeout, or sensor error.
* `Not Applicable`: The control or requirement does not apply to this target configuration.
* `Not Run`: The assessment was scheduled or defined, but skipped or not executed.

### Step 3: Granular Assessment Logs

Each control evaluation contains one or more `assessment-logs`. These capture individual requirement checks:

```yaml
    assessment-logs:
      - requirement:
          reference-id: "REQ-IMG-01.01"
          description: "Container images must be referenced using immutable SHA256 digest pinning."
        result: "Passed"
        start: "2026-03-31T09:45:02Z"
        end: "2026-03-31T09:45:06Z"
        steps-executed: 2
        confidence: "High"
        steps:
          - id: "STEP-SCAN-MANIFEST"
            name: "Parse deployment manifests"
            result: "Passed"
          - id: "STEP-VERIFY-DIGEST"
            name: "Verify sha256 prefix and hex length"
            result: "Passed"
```

> [!NOTE]
> The `control.reference-id` in the parent `#ControlEvaluation` must match the prefix or association of `requirement.reference-id` in your assessment catalogs.

### Step 4: Recording Evidence

Evidence proves the evaluation result. Gemara requires at least one of `payload` (inlined artifact/logs) or `source` (cryptographic reference or external digest):

#### Option A: Inlined Payload
Useful for structured command outputs, policy engine audit events, or short test traces:
```yaml
evidence:
  - id: "EVID-IMG-SCAN-01"
    name: "Admission Controller Policy Evaluation"
    description: "Audit evaluation trace from cluster admission controller."
    created: "2026-03-31T09:45:05Z"
    type: "Log"
    payload:
      format: "application/json"
      data: '{"enforced":true,"rejected_tags":["latest"],"validated_digest":"sha256:d891b7d59b20b22a0156d1373"}'
```

#### Option B: External Digest Reference
Useful for large artifacts, container SBOMs, or remote audit storage:
```yaml
evidence:
  - id: "EVID-MFA-API-RESPONSE"
    name: "IAM API Audit Snapshot"
    description: "Exported IAM audit log proving hardware token requirement."
    created: "2026-03-31T09:45:12Z"
    type: "Log"
    source:
      digest: "sha256:4b227777d4dd1fc61c6f884f48641d02b4d121d3fd328cb08b5531fcacdabf8a"
      uri: "https://logs.internal.net/audits/iam/2026-03-31-mfa-report.json"
```

### Step 5: Handling Unexecuted and Not Applicable Evaluations

Under Gemara's `#_AssessmentLogStrict` validation rules:
1. When an assessment runs (`Passed`, `Failed`, `Partially Passed`), the `start` timestamp is **mandatory**.
2. When an assessment is `Not Applicable` or `Not Run` and execution was bypassed before running any steps (`steps-executed: 0`), the `start` timestamp **must be omitted**:

```yaml
  - control:
      reference-id: "CTRL-DBG-03"
      catalog-reference:
        id: "container-security-standard-v1"
    result: "Not Applicable"
    summary: "Target runs in production profile where debug containers are disabled by cluster-level RBAC."
    assessment-logs:
      - requirement:
          reference-id: "REQ-DBG-03.01"
          description: "Debug ephemeral containers must require dual-custody authorization."
        result: "Not Applicable"
        steps-executed: 0
        justification: "Production ingress deployment does not enable ephemeral container capabilities (disabled at API server)."
```

---

## Full Example

Below is a complete, validated `EvaluationLog` artifact demonstrating multi-control assessments, evidence payloads, external digests, and not-applicable handling:

```yaml
metadata:
  id: "EVAL-2026-03-K8S-PROD-01"
  type: "EvaluationLog"
  gemara-version: "1.2.0"
  description: "Automated admission control and identity verification audit for production cluster"
  author:
    id: "sec-ops-automation"
    name: "Platform Security Engine"
    type: "Tool"
  created: "2026-03-31T09:45:00Z"
target:
  type: "Software"
  name: "ingress-gateway-service"
  version: "v2.14.0"
  uri: "git+https://github.com/example-org/ingress-gateway.git#refs/tags/v2.14.0"
result: "Passed"
evaluations:
  - control:
      reference-id: "CTRL-IMG-01"
      catalog-reference:
        id: "container-security-standard-v1"
        uri: "https://gemara.openssf.org/catalogs/container-security-v1.yaml"
    result: "Passed"
    summary: "Immutable container image tags and SHA256 digests strictly enforced across all deployment specs."
    assessment-logs:
      - requirement:
          reference-id: "REQ-IMG-01.01"
          description: "Container images must be referenced using immutable SHA256 digest pinning."
        result: "Passed"
        start: "2026-03-31T09:45:02Z"
        end: "2026-03-31T09:45:06Z"
        steps-executed: 2
        confidence: "High"
        steps:
          - id: "STEP-SCAN-MANIFEST"
            name: "Parse deployment manifests"
            result: "Passed"
          - id: "STEP-VERIFY-DIGEST"
            name: "Verify sha256 prefix and hex length"
            result: "Passed"
        evidence:
          - id: "EVID-IMG-SCAN-01"
            name: "Admission Controller Policy Evaluation"
            description: "Audit evaluation trace from cluster admission controller."
            created: "2026-03-31T09:45:05Z"
            type: "Log"
            payload:
              format: "application/json"
              data: '{"enforced":true,"rejected_tags":["latest"],"validated_digest":"sha256:d891b7d59b20b22a0156d1373"}'
  - control:
      reference-id: "CTRL-AUTH-02"
      catalog-reference:
        id: "identity-governance-standard-v2"
        uri: "https://gemara.openssf.org/catalogs/identity-governance-v2.yaml"
    result: "Passed"
    summary: "Multi-factor authentication enforced on administrative endpoints."
    assessment-logs:
      - requirement:
          reference-id: "REQ-AUTH-02.01"
          description: "Administrative access must require hardware-backed WebAuthn/FIDO2 MFA."
        result: "Passed"
        start: "2026-03-31T09:45:07Z"
        end: "2026-03-31T09:45:14Z"
        steps-executed: 1
        confidence: "High"
        evidence:
          - id: "EVID-MFA-API-RESPONSE"
            name: "IAM API Audit Snapshot"
            description: "Exported IAM audit log proving hardware token requirement."
            created: "2026-03-31T09:45:12Z"
            type: "Log"
            source:
              digest: "sha256:4b227777d4dd1fc61c6f884f48641d02b4d121d3fd328cb08b5531fcacdabf8a"
              uri: "https://logs.internal.net/audits/iam/2026-03-31-mfa-report.json"
  - control:
      reference-id: "CTRL-DBG-03"
      catalog-reference:
        id: "container-security-standard-v1"
        uri: "https://gemara.openssf.org/catalogs/container-security-v1.yaml"
    result: "Not Applicable"
    summary: "Target runs in production profile where debug containers are disabled by cluster-level RBAC."
    assessment-logs:
      - requirement:
          reference-id: "REQ-DBG-03.01"
          description: "Debug ephemeral containers must require dual-custody authorization."
        result: "Not Applicable"
        steps-executed: 0
        justification: "Production ingress deployment does not enable ephemeral container capabilities (disabled at API server)."
```

You can download the raw fixture file: [evaluation-log-example.yaml](evaluation-log-example.yaml).

---

## Validating Your Evaluation Log

To validate your evaluation log against the official Gemara CUE schema, use the CUE CLI:

```bash
cue vet -c -d '#EvaluationLog' github.com/gemaraproj/gemara@latest evaluation-log-example.yaml
```

If developing locally in a clone of the `gemara` repository:

```bash
cue vet -c -d '#EvaluationLog' . path/to/your-evaluation-log.yaml
```

When validation succeeds, `cue vet` exits silently with exit status `0`.
