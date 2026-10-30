---
layout: page
title: Define Control Inheritance Conventions for Multi-Party Governance
---

- **ADR:** 0023
- **Proposal Author(s):** @harshsomankar123-tech, @jpower432
- **Status:** Proposed

## Context

In modern cloud and multi-tier enterprise systems, security controls are rarely satisfied by a single component or organization in isolation. Organizations routinely operate within shared responsibility models (such as those defined by cloud providers, FedRAMP, NIST SP 800-53, and the FINOS Common Cloud Controls). Responsibilities are partitioned among cloud service providers, platform teams, infrastructure administrators, application developers, and external vendors.

Gemara's CUE schema defines primitives capable of expressing multi-party responsibility, including:
- `metadata.applicability-groups` and requirement `applicability` tags (Layer 2 Control Catalogs)
- `imports` and `mapping-references` for referencing upstream catalogs and capabilities
- `CapabilityCatalog` for detailing platform-provided functionality
- `executor` and `evaluation-methods` within `adherence.assessment-plans` (Layer 3 Policies)

However, Gemara has lacked standardized conventions for using these primitives together to express **control inheritance**. Without established conventions:
- Authors describe responsibility boundaries using disparate, ad-hoc terms or free-text prose.
- Downstream evaluators and automated tooling cannot distinguish between controls satisfied out of the box by an upstream platform versus those requiring tenant-side configuration.
- Policy adherence definitions struggle to decouple validation methods: testing a tenant's configuration requires automated behavioral inspection, whereas verifying an upstream provider's platform controls requires reviewing third-party attestations (e.g., SOC 2 reports).

To address this, we must define conventions for five canonical tiers of control responsibility:
1. **Native**: The component satisfies the control objective by default without user configuration.
2. **Configurable**: The component provides the capability to satisfy the control, but requires explicit consumer configuration.
3. **Inherited**: The component consumes the capability from an upstream source and depends on that implementation.
4. **Inherited Configurable (Shared)**: Functionality is inherited from an upstream source, but requires consumer configuration.
5. **Custom**: The component does not provide the capability. The implementer must satisfy the control outside the component scope (e.g., via compensating controls) or document risk acceptance.

## Decision

Standardize conventions for control inheritance by leveraging existing schema primitives—specifically **namespaced applicability groups** in Layer 2 Control Catalogs combined with **decoupled assessment plans** in Layer 3 Policies.

### 1. Applicability Group Namespace (`resp:*`)

When defining applicability groups in catalog metadata, authors SHOULD use the `resp:` namespace prefix to represent responsibility and inheritance tiers:

- `resp:native` — Satisfied by the component/platform by default.
- `resp:configurable` — Provided by the component, but requires consumer configuration.
- `resp:inherited` — Fully provided by an upstream provider or platform.
- `resp:shared` — Provided by an upstream provider, but requires consumer configuration (inherited configurable).
- `resp:custom` — Not provided by the platform; requires compensating safeguards or risk acceptance.

Example in a `ControlCatalog`:

```yaml
metadata:
  applicability-groups:
    - id: resp:native
      title: Native Responsibility
      description: Component satisfies the control by default without configuration.
    - id: resp:configurable
      title: Configurable Responsibility
      description: Component provides capability, but requires explicit configuration.
    - id: resp:inherited
      title: Inherited Responsibility
      description: Functionality is inherited from an upstream provider or platform.
    - id: resp:shared
      title: Shared Responsibility
      description: Functionality is inherited upstream, but requires consumer configuration.
    - id: resp:custom
      title: Custom Responsibility
      description: Not provided by the component; requires external implementation or risk acceptance.
```

### 2. Tagging Assessment Requirements

Within controls, `assessment-requirements` declare applicability using the appropriate `resp:*` identifier, optionally alongside environment or workload applicability tags:

```yaml
controls:
  - id: CTRL-ENCRYPTION-AT-REST
    title: Encryption at Rest
    objective: Ensure stored data is encrypted using approved algorithms.
    group: data-protection
    assessment-requirements:
      - id: CTRL-ENCRYPTION-AR1
        text: Default server-side encryption MUST be active on all storage buckets.
        applicability: ["resp:inherited", "production"]
      - id: CTRL-ENCRYPTION-AR2
        text: Customer-managed keys MUST enforce a minimum key length of AES-256.
        applicability: ["resp:configurable", "production"]
```

### 3. Decoupled Policy Validation Patterns

In Layer 3 Policies, assessment plans decouple validation based on the control responsibility:
- **`resp:configurable` requirements** map to automated behavioral evaluation methods executed by software (e.g., CI/CD linters, Terraform policy checks, runtime CSPM scanners).
- **`resp:inherited` requirements** map to manual intent evaluation methods executed by humans (e.g., verifying vendor SOC 2 Type II reports or ISO certifications on an annual frequency).
- **`resp:custom` requirements** that cannot be fulfilled within scope are documented under `risks.accepted` with an explicit risk reference and technical justification.

## Consequences

### Positive

- **Zero Breaking Schema Changes**: Works immediately within existing Gemara schemas (`metadata.applicability-groups`, `applicability`, `imports`, and `adherence`). No modification to CUE type definitions is required.
- **Interoperability**: Downstream evaluators (e.g., automated scanning engines) can filter assessment requirements by `resp:configurable` to execute automated code scans, while ignoring `resp:inherited` checks that belong in audit pipelines.
- **Clarity across Tiers**: Clear alignment with industry standards (FedRAMP CRM, NIST SP 800-53 Rev. 5 control allocation, and FINOS Common Cloud Controls).

### Negative

- Because `applicability-groups` is an open list of `#Group` entries, schema validation alone cannot strictly enforce the spelling of `resp:*` tags unless tooling or linting rules specifically validate this convention.

## Alternatives Considered

### Dedicated `responsibility` enum in `#AssessmentRequirement`

Adding a mandatory or optional field `responsibility: #ResponsibilityTier` directly to `#AssessmentRequirement` in `controlcatalog.cue`.

**Rejected because:** 
1. It introduces schema changes and breaking differences across Gemara versions.
2. Applicability is already the canonical Gemara mechanism for scoping when and where an assessment requirement applies. Adding a redundant responsibility field fragments scoping metadata.

### Adding a `recommended-executor` field at the ControlCatalog layer

Adding a field to suggest whether an evaluation should be performed by `Human` or `Software` at Layer 2.

**Rejected because:** 
Layer 2 defines technology controls independently of how a particular organization chooses to evaluate them. Executor assignment belongs to Layer 3 Policy (`adherence.assessment-plans.evaluation-methods.executor`), maintaining separation of concerns between definition and measurement.
