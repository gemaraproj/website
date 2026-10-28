---
layout: page
title: Python SDK
---

The Python SDK (`gemara-python`) provides generated Pydantic v2 models for Gemara v1 documents. It enables loading, validating, manipulating, and inspecting Gemara JSON and YAML files with full static type hints and runtime validation. Types are generated from CUE schemas (Schema v1.5.0), ensuring consistency with the Gemara core specification.

**[Python Repository →](https://github.com/gemaraproj/gemara-python)** &nbsp;|&nbsp; **[PyPI Package →](https://pypi.org/project/gemara-python/)** &nbsp;|&nbsp; **[Issue Tracker →](https://github.com/gemaraproj/gemara-python/issues)**

## Installation

`gemara-python` requires Python 3.11+ (supports Python 3.11, 3.12, 3.13, and 3.14).

```bash
pip install gemara-python
```

Or using [uv](https://github.com/astral-sh/uv):

```bash
uv add gemara-python
```

## Loading Documents

The SDK supports both YAML and JSON document formats (YAML is handled as a superset of JSON). You can load documents either via **dynamic dispatch** or **direct concrete type loading**.

### Dynamic Dispatch Loading (`load` / `loads`)

When working with arbitrary Gemara documents where the type is determined at runtime by the document's `metadata.type` discriminator, use `load` (for files) or `loads` (for strings/bytes):

```python
from gemara.v1 import ControlCatalog, Lexicon, load

# Automatically inspects metadata.type and parses into the appropriate model
doc = load("catalog.yaml")

# Pattern match or inspect the concrete model
match doc:
    case ControlCatalog():
        print(f"Loaded Control Catalog with {len(doc.controls or [])} controls.")
    case Lexicon():
        print(f"Loaded Lexicon with {len(doc.terms)} terms.")
```

- **`load(source)`**: Accepts a file path (`str`, `Path`) or an open file object (`IO[str]`, `IO[bytes]`).
- **`loads(text)`**: Accepts raw text (`str`) or UTF-8 encoded bytes (`bytes`, `bytearray`, `memoryview`) containing JSON or YAML.

### Direct Concrete Type Loading (`from_file` / `from_text`)

When the document kind is known in advance, load directly through the model class to receive that concrete type without dynamic dispatch:

```python
from gemara.v1 import GuidanceCatalog

# Directly validates and loads guidance.yaml into a GuidanceCatalog instance
guidance = GuidanceCatalog.from_file("guidance.yaml")
print(f"Guidance ID: {guidance.metadata.id}")
```

- **`Model.from_file(source)`**: Reads and parses a file directly into the model.
- **`Model.from_text(text)`**: Parses a JSON/YAML string or bytes directly into the model.

## Core Concepts & Type Hierarchy

### Union Return Type

Functions `load()` and `loads()` return `GemaraDocument`, a `TypeAlias` union representing all 13 supported Gemara v1 document types:

- **Catalogs:** `CapabilityCatalog`, `ControlCatalog`, `GuidanceCatalog`, `PrincipleCatalog`, `RiskCatalog`, `ThreatCatalog`, `VectorCatalog`
- **Logs:** `AuditLog`, `EnforcementLog`, `EvaluationLog`
- **Other:** `Lexicon`, `MappingDocument`, `Policy`

### Category Narrowing

All catalog documents inherit from `Catalog` and all log documents inherit from `Log`:

```python
from gemara.v1 import Catalog, Log, load

doc = load("sample.yaml")

if isinstance(doc, Catalog):
    print("This is a Gemara Catalog.")
elif isinstance(doc, Log):
    print("This is a Gemara Log.")
```

## Error Handling

All custom exceptions raised by `gemara-python` inherit from `GemaraError`:

```python
from gemara.v1 import GemaraError, UnknownDocumentTypeError, load
from pydantic import ValidationError

try:
    doc = load("document.yaml")
except UnknownDocumentTypeError as e:
    # metadata.type was missing or unrecognized
    print(f"Unsupported document type: {e}")
except ValidationError as e:
    # Document schema validation failed
    print(f"Schema validation error: {e}")
except GemaraError as e:
    # Non-UTF-8 bytes, unparseable YAML/JSON, or non-mapping input
    print(f"Gemara parsing error: {e}")
except OSError as e:
    # File I/O issues
    print(f"Filesystem error: {e}")
```

## Schema Compatibility & Limitations

- **Additive Changes:** Within Gemara v1, changes are backwards compatible. The models read any v1.x document. Unrecognized properties are accepted and ignored.
- **Structural Validation vs. CUE:** `gemara-python` performs structural schema validation (required fields, enums, regex patterns, length bounds). For deep cross-field relational validation (such as referential integrity or uniqueness enforced via CUE comprehensions), validate using `cue vet` alongside the core Gemara CUE schemas.

## Useful Exports

- **`SCHEMA_VERSION`**: The Gemara schema version used to generate the models (e.g. `"1.5.0"`).
- **`__version__`**: The installed `gemara-python` library version.
- **`DOCUMENT_TYPES`**: Dictionary mapping `metadata.type` string values to model classes.

## Relationship to Other Components

### [The Model](../model)
Provides the conceptual foundation. Python SDK types correspond to elements in the model.

### [The Schemas](../schema/)
Python SDK types are generated from the CUE schemas, ensuring consistency between validation and programmatic access. The schema version supported is exposed as `gemara.v1.SCHEMA_VERSION`.
