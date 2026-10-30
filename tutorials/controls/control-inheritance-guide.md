---
layout: page
title: Control Inheritance and Shared Responsibility Guide
description: Practical guide to expressing multi-party control responsibility and decoupled validation using Gemara conventions.
---

## What This Is

Modern software systems depend on shared infrastructure, cloud services, and third-party vendors. In cloud and enterprise security frameworks (such as FedRAMP, NIST SP 800-53, and FINOS Common Cloud Controls), **controls are rarely implemented by a single team or system alone**.

This guide demonstrates how to model **Control Inheritance** and **Shared Responsibility** in Gemara using existing schema primitives and standardized conventions.

By following this guide, you will learn how to:
1. Express multi-party responsibility across **5 standardized tiers**.
2. Scope assessment requirements using the **`resp:*` applicability-groups convention** in Layer 2 Control Catalogs.
3. Decouple validation methods in Layer 3 Policies so automated checks, vendor attestations, and risk acceptance each apply to their appropriate tier.

---

## The 5 Responsibility Tiers

When authoring controls for multi-tenant or cloud-hosted systems, requirements fall into five distinct categories:

| Tier | Convention ID | Definition | Real-World Example |
| :--- | :--- | :--- | :--- |
| **Native** | `resp:native` | The component meets the requirement by default without consumer configuration. | Managed cloud endpoints enforcing TLS 1.2+ out of the box. |
| **Configurable** | `resp:configurable` | The component provides the capability, but requires explicit consumer configuration to satisfy the control. | Enabling customer-managed keys (CMEK) and setting minimum key length to AES-256. |
| **Inherited** | `resp:inherited` | The component consumes the capability entirely from an upstream source and depends on that implementation. | Physical data center security or hardware-level storage volume encryption provided by a cloud provider. |
| **Shared (Inherited Configurable)** | `resp:shared` | Upstream platform provides the base infrastructure, but the consumer must configure and operate it correctly. | Cloud Key Management Service (KMS) where the provider hosts the HSM, but the tenant configures key rotation and access policies. |
| **Custom** | `resp:custom` | The platform does not provide the capability. The implementer must satisfy the requirement externally or accept the risk. | Custom business transaction logging implemented via a sidecar proxy because the native storage API does not log application payloads. |

---

## Walkthrough: Cloud Storage Security

We demonstrate these conventions through a real-world scenario: securing **Cloud Object Storage**.

### Step 1: Upstream Provider Capabilities (Layer 2)

Upstream platform providers publish a **`CapabilityCatalog`** documenting the security features their platform makes available:

```yaml
title: Cloud Platform Storage Capabilities
metadata:
  id: cloud-platform-storage-caps
  type: CapabilityCatalog
  gemara-version: "1.2.0"
  description: Capabilities provided natively by the cloud storage platform.
  author:
    id: cloud-provider
    name: Cloud Service Provider
    type: Software

groups:
  - id: encryption-services
    title: Encryption Services
    description: Cryptographic capabilities provided by the platform.

capabilities:
  - id: CAP-DEFAULT-ENCRYPTION
    title: Default Server-Side Encryption
    description: Platform transparently encrypts all stored objects at rest.
    group: encryption-services
  - id: CAP-CMEK
    title: Customer-Managed Encryption Keys
    description: Platform supports integration with Key Management Service for customer keys.
    group: encryption-services
```

---

### Step 2: Consumer Control Catalog with `resp:*` Conventions (Layer 2)

As a consumer, you define your security controls in a **`ControlCatalog`**. You import the upstream capabilities and define the standardized `resp:*` applicability groups in `metadata`:

```yaml
title: Cloud Storage Security Control Catalog
metadata:
  id: org-cloud-storage-controls
  type: ControlCatalog
  gemara-version: "1.2.0"
  author:
    id: infosec-team
    name: Information Security Team
    type: Human
  mapping-references:
    - id: cloud-platform-storage-caps
      title: Cloud Platform Storage Capabilities
      version: "1.0.0"
      url: https://example.com/catalogs/cloud-platform-storage-caps.yaml
  applicability-groups:
    - id: resp:native
      title: Native Responsibility
      description: Satisfied by component defaults without configuration.
    - id: resp:configurable
      title: Configurable Responsibility
      description: Provided by component; requires consumer configuration.
    - id: resp:inherited
      title: Inherited Responsibility
      description: Satisfied entirely by upstream platform or provider.
    - id: resp:shared
      title: Shared Responsibility
      description: Provided upstream, but requires consumer configuration.
    - id: resp:custom
      title: Custom Responsibility
      description: Not provided by platform; requires external implementation or risk acceptance.
    - id: production
      title: Production Environments
      description: Storage buckets holding sensitive production data.
```

Next, each control tags its `assessment-requirements` with the appropriate `resp:*` group:

```yaml
controls:
  # 1. INHERITED: Consumed from cloud provider
  - id: CTRL-STORAGE-ENCRYPT-DEFAULT
    title: Default Server-Side Storage Encryption
    objective: Ensure data at rest is encrypted via inherited provider capability.
    group: storage-data-protection
    assessment-requirements:
      - id: CTRL-STORAGE-ENCRYPT-DEFAULT-AR1
        text: The cloud provider MUST provide default server-side encryption for all objects.
        applicability: ["resp:inherited", "production"]

  # 2. SHARED / CONFIGURABLE: Requires consumer configuration
  - id: CTRL-STORAGE-CMEK-CONFIG
    title: Customer-Managed Encryption Key Configuration
    objective: Ensure customer-managed keys meet organizational strength requirements.
    group: storage-data-protection
    assessment-requirements:
      - id: CTRL-STORAGE-CMEK-CONFIG-AR1
        text: Storage buckets holding sensitive data MUST configure CMEK with AES-256.
        applicability: ["resp:shared", "production"]

  # 3. NATIVE: Out-of-the-box baseline
  - id: CTRL-STORAGE-BLOCK-PUBLIC
    title: Restrict Public Storage Exposure
    objective: Ensure storage resources disallow anonymous public read access by default.
    group: storage-data-protection
    assessment-requirements:
      - id: CTRL-STORAGE-BLOCK-PUBLIC-AR1
        text: Storage buckets MUST block public access configurations by default.
        applicability: ["resp:native", "production"]

  # 4. CUSTOM: Unmet by cloud platform
  - id: CTRL-STORAGE-APP-AUDIT
    title: Application-Layer Data Access Auditing
    objective: Maintain application-specific transaction logs for object operations.
    group: storage-auditability
    assessment-requirements:
      - id: CTRL-STORAGE-APP-AUDIT-AR1
        text: Application data access events MUST be captured in an external log repository.
        applicability: ["resp:custom", "production"]
```

---

### Step 3: Decoupled Validation in Policy (Layer 3)

In Layer 3, an organizational **`Policy`** imports the control catalog and defines how compliance is evaluated. 

Control inheritance decouples validation into the right inspection methods and executors:
* **`resp:configurable` and `resp:shared`**: Evaluated continuously by **Software** (CI/CD linters, Terraform policy checks, CSPM scanners).
* **`resp:inherited`**: Evaluated periodically by a **Human** auditor reviewing vendor attestations (e.g., SOC 2 Type II reports, ISO certifications).
* **`resp:custom`**: Handled via documented **risk acceptance** or custom compensating tool inspection.

```yaml
title: Cloud Storage Security Policy
metadata:
  id: org-cloud-storage-policy
  type: Policy
  gemara-version: "1.2.0"
  author:
    id: security-compliance
    name: Security & Compliance Office
    type: Human
  mapping-references:
    - id: org-cloud-storage-controls
      title: Cloud Storage Security Control Catalog
      version: "1.0.0"

imports:
  catalogs:
    - reference-id: org-cloud-storage-controls

adherence:
  assessment-plans:
    # 1. Inherited Control -> Verified via Vendor Attestation (Annual Human Review)
    - id: plan-inherited-default-encryption
      requirement-id: CTRL-STORAGE-ENCRYPT-DEFAULT-AR1
      frequency: annual
      evidence-requirements: "Cloud Provider SOC 2 Type II report confirming encryption"
      evaluation-methods:
        - id: em-vendor-attestation
          type: Intent
          mode: Manual
          required: true
          description: Verify cloud provider maintains third-party SOC 2 attestation.
          executor:
            id: compliance-analyst
            name: Senior Compliance Analyst
            type: Human

    # 2. Shared/Configurable Control -> Verified via Automated CI/CD Scanner
    - id: plan-shared-cmek-configuration
      requirement-id: CTRL-STORAGE-CMEK-CONFIG-AR1
      frequency: continuous
      evidence-requirements: "Automated scan output verifying Terraform bucket encryption"
      evaluation-methods:
        - id: em-automated-iac-check
          type: Behavioral
          mode: Automated
          required: true
          description: Automated CI pipeline check validates minimum encryption strength in IaC.
          executor:
            id: terraform-policy-scanner
            name: Terraform Policy Scanner
            type: Software
          parameters:
            - id: min-key-length
              label: Minimum Key Length
              accepted-values: ["AES-256"]

# Custom / Unmet Control -> Formally Documented Risk Acceptance
risks:
  accepted:
    - id: risk-accept-app-audit-logging
      risk:
        entry-id: RISK-DATA-TRACEABILITY-1
      justification: >
        Native cloud storage does not support tenant-level custom business context logging.
        Compensating logging implemented via external application middleware sidecar.
        Residual risk formally accepted per annual review.
      scope:
        in:
          technologies: ["Cloud Storage"]
```

---

## Why Use Conventions Instead of New Schema Fields?

1. **Zero Breaking Changes**: The `resp:*` naming convention works immediately with published versions of the Gemara CUE schema without waiting for schema releases.
2. **Flexible Composability**: Requirements can combine responsibility tags with environmental tags (e.g. `applicability: ["resp:shared", "production", "pci-scope"]`).
3. **Automated Tooling Compatibility**: Security scanners can filter evaluation targets by `resp:configurable` to run automated tests, while compliance dashboards can group `resp:inherited` controls into vendor attestation matrices.

For full architectural discussion and alternatives considered, see [ADR 0023: Define Control Inheritance Conventions for Multi-Party Governance](../../adrs/0023-control-inheritance-conventions).

---

## Ready-to-Use Artifact

A complete, multi-document YAML artifact combining the upstream capability catalog, consumer control catalog, and decoupled policy is available at [control-inheritance-example.yaml](control-inheritance-example.yaml).
