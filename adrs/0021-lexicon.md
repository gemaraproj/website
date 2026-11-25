---
layout: page
title: New schema for Lexicon artifacts
description: "ADR-0021: add a Lexicon schema for documents that define the terms other Gemara artifacts use."
---

- **ADR:** 0021
- **Proposal Author(s):** @eddie-knight
- **Status:** Accepted

## Context

We currently have a field in all documents for `Metadata.Lexicon`, which is an `ArtifactMapping` type. This allows authors to link out to another document for their definitions. We do not, however, have any opinion on what that document looks like.

Gemara user [OSPS Baseline](https://baseline.openssf.org) uses a custom YAML schema to hold their definitions, which is then used when generating the project's Markdown web page. The generator inspects the text fields of every control and requirement, then wraps any lexicon terms in a link which enables users to jump directly to the definition.  

## Action

Create a new lexicon schema inspired by the OSPS Basline Lexicon. It should behave like all other artifact schemas, such as the presence of `Metadata`. 

Similar to mapping documents, this should not be affiliated with a Gemara Model Layer.

The body of the document should be approximately:

- title
- definition
- synonyms (optional, used for linking logic)
- references (optional, functions like citations for the definition)

## Consequences

Positive: Fills the gap created by `Metadata.Lexicon`. Allows official Gemara tooling to behave similar to how users are currently operating.

Negative: Yet another schema.

## Alternatives Considered

None
