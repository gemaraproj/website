---
layout: page
title: "Specification: OCI Bundle Distribution & Packaging"
---

## Status

* **Status:** Draft / Active
* **Type:** Tool Specification
* **Domain:** Distribution & Transport
* **Applies To:** Gemara Tooling, SDKs (`go-gemara`), Registries, CLI tools (`grcli`)

---

## 1. Abstract

Gemara artifacts (Policies, Control Catalogs, Guidance Catalogs, and Mapping Documents) must be distributed portably across local workspaces and enterprise remote registries.

This specification defines the normative standard for how Gemara artifacts and composite multi-catalog bundles are packaged into OCI (Open Container Initiative) image manifests, assigned standard media types, tagged, referenced via URI schemes, and cryptographically verified using OCI registries and ORAS (OCI Registry As Storage).

---

## 2. Terminology & Core Concepts

* **Primary Artifact (Source):** The anchor Gemara document for a bundle (e.g., a corporate security policy or a primary control catalog).
* **Imported Artifact (Dependency):** A referenced Gemara document transitively imported or cited by the primary artifact (e.g., imported baseline catalogs, guidance frameworks, or threat catalogs).
* **Gemara Bundle:** An OCI artifact manifest encapsulating a primary artifact, zero or more imported artifacts, and an OCI configuration blob detailing internal file roles and dependency relationships.
* **ORAS:** OCI Registry As Storage client library implementing OCI image and artifact packaging specifications.

---

## 3. OCI Media Types

Compliant packaging tools and SDKs **MUST** use the following registered MIME types when authoring or inspecting Gemara OCI manifests and layers:

| Media Type | Target OCI Descriptor Field | Description |
| :--- | :--- | :--- |
| `application/vnd.gemara.bundle.v1` | `artifactType` / `config.mediaType` | Identifies the OCI manifest as a conforming Gemara artifact bundle. |
| `application/vnd.gemara.manifest.v1+json` | `config.mediaType` | Carries the JSON bundle configuration and artifact dependency manifest. |
| `application/vnd.gemara.artifact.v1+yaml` | `layers[].mediaType` | Media type for individual Gemara YAML document layers. |

### Conformance Rules:
1. All document layers **MUST** carry `mediaType: application/vnd.gemara.artifact.v1+yaml`.
2. Image manifests **MUST** set `artifactType: application/vnd.gemara.bundle.v1` (for OCI Image Spec v1.1+) or carry this value in the config descriptor for backward compatibility with v1.0 registries.
3. Consumers encountering unrecognized media types in a Gemara bundle manifest **MUST** reject the bundle as incompatible.

---

## 4. Manifest Configuration Blob (`Manifest`)

The OCI config descriptor points to a JSON payload conforming to `application/vnd.gemara.manifest.v1+json`. This payload defines the bundle structure and dependency metadata:

```json
{
  "bundle-version": "1.0.0",
  "gemara-version": "1.1.0",
  "revision": "2026-11-02",
  "metadata": {
    "organization": "OpenSSF",
    "project": "osps-baseline"
  },
  "artifacts": [
    {
      "name": "policy.yaml",
      "type": "Policy",
      "id": "CORP-SEC-POLICY",
      "role": "artifact",
      "dependencies": [
        "osps-baseline.yaml"
      ]
    },
    {
      "name": "osps-baseline.yaml",
      "type": "ControlCatalog",
      "id": "osps-baseline",
      "role": "import"
    }
  ]
}
```

### Fields:
* **`bundle-version` (string, REQUIRED):** Semantic version of the bundle packaging format.
* **`gemara-version` (string, REQUIRED):** Target Gemara specification version (e.g. `1.1.0`).
* **`artifacts` (array, REQUIRED):** An ordered list of descriptors for each file stored in the bundle layers.
  * **`name` (string, REQUIRED):** Bundle-relative file path (e.g. `policy.yaml`).
  * **`type` (string, REQUIRED):** Valid Gemara artifact type (`Policy`, `ControlCatalog`, `GuidanceCatalog`, `MappingDocument`, `ThreatCatalog`, `RiskCatalog`, etc.).
  * **`id` (string, REQUIRED):** Must match the `metadata.id` declared within the artifact YAML content.
  * **`role` (string, REQUIRED):** Must be either `artifact` (for the primary source) or `import` (for dependencies).
  * **`dependencies` (array of strings, OPTIONAL):** Names of other artifacts in this bundle that this document directly cites.

---

## 5. Layer Descriptors & Annotations

Each Gemara YAML artifact layer **MUST** include OCI standard and Gemara-specific annotations on its layer descriptor in the manifest:

```json
{
  "mediaType": "application/vnd.gemara.artifact.v1+yaml",
  "digest": "sha256:7f83b1657ff1fc53b92dc18148a1d65dfc2d4b1fa3d677284addd200126d9069",
  "size": 4210,
  "annotations": {
    "org.opencontainers.image.title": "controls.yaml",
    "org.gemara.artifact.role": "artifact",
    "org.gemara.artifact.type": "ControlCatalog"
  }
}
```

### Standard Annotations:
* `org.opencontainers.image.title`: **MUST** match the relative filename (e.g., `controls.yaml`).
* `org.gemara.artifact.role`: **MUST** be either `artifact` (primary) or `import` (transitive dependency). Exactly one layer in a bundle **SHALL** have role `artifact`.
* `org.gemara.artifact.type`: **MUST** declare the Gemara artifact type (e.g., `ControlCatalog`, `Policy`).

---

## 6. Referencing & URI Syntax

Gemara tools **MUST** support referencing OCI artifacts using the `oci://` scheme.

### URI Grammar:
```text
oci://<registry>/<repository>[:<tag>][@<digest>]
```

### Examples:
* **Tagged release:** `oci://ghcr.io/gemaraproj/osps-baseline:1.0.0`
* **Digest pinned:** `oci://registry.example.org/compliance/policy@sha256:45b0e2634d243c6697d8487045760980d0ba964e1f69629f7042c37ecfe32c9f`
* **Tagged with digest pinning (RECOMMENDED for production audits):** `oci://ghcr.io/gemaraproj/osps-baseline:1.0.0@sha256:45b0e2634d...`

### Resolution Rules:
1. When both tag and digest are supplied, the digest **MUST** take precedence for integrity verification. Tools **MUST** fail if the remote manifest digest does not match the pinned digest.
2. If neither tag nor digest is specified, tools **SHOULD** default to tag `latest` or report a configuration error.

---

## 7. Size Limits & DoS Protection

To prevent memory exhaustion and zip-bomb style attacks from untrusted registries:
1. Conformant unpackers **MUST** enforce a configurable maximum bundle size limit.
2. The default maximum bundle size limit **SHOULD** be **256 MB** (`268,435,456` bytes).
3. Reads exceeding the size limit **MUST** immediately terminate with an explicit error.

---

## 8. Failure Behavior & Error Handling

Conformant tools and SDKs **MUST** handle transport and verification failures deterministically:

| Failure Condition | Expected Tool Behavior |
| :--- | :--- |
| **Missing or Mismatched Digest** | Fail execution with an integrity validation error; do not evaluate artifacts. |
| **Unsupported Layer Media Type** | Refuse unpacking and report unrecognized media type. |
| **Missing Primary Source Layer** | Reject bundle (a valid bundle must have exactly one layer with role `artifact`). |
| **Mismatched Artifact ID** | If layer annotation or manifest `id` differs from the embedded YAML `metadata.id`, report validation failure. |
| **Registry Authentication Error (401/403)** | Output actionable error message indicating credential configuration or missing registry permissions. |

---

## 9. Reference Implementation

* The canonical reference implementation for assembling, packing, and unpacking Gemara OCI bundles is available in the Go SDK: [`gemaraproj/go-gemara/bundle`](https://github.com/gemaraproj/go-gemara/tree/main/bundle).
* CLI interaction and publication pipelines are provided via `grcli` and `gemara-publish-action`.
