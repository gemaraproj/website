---
layout: page
title: React SDK
---

The React SDK (`@gemara/react`) provides headless React renderers and TypeScript components for Gemara GRC artifacts. Built to be server-component-clean by default, it renders Gemara documents in modern React frameworks (such as Next.js App Router, Remix, or Astro) with full static typing and zero unnecessary client JavaScript.

**[React Repository →](https://github.com/gemaraproj/gemara-react)** &nbsp;|&nbsp; **[npm Package →](https://www.npmjs.com/package/@gemara/react)** &nbsp;|&nbsp; **[Issue Tracker →](https://github.com/gemaraproj/gemara-react/issues)**

## Installation

`@gemara/react` requires React 18 or newer (`react >= 18`).

Install via `npm`:

```bash
npm install @gemara/react
```

Or using `pnpm` or `yarn`:

```bash
pnpm add @gemara/react
# or
yarn add @gemara/react
```

## Core Design Principles

- **Headless & Unstyled:** The library renders clean, accessible, semantic HTML without injecting opinionated CSS, giving you complete styling freedom.
- **Server Component Clean by Default:** Catalog and document renderers contain no `useState`, no `useEffect`, and no `"use client"` directives. In React Server Component (RSC) environments, documents render entirely on the server with zero client bundle weight.
- **Tree-Shakeable Subpath Imports:** Each artifact renderer can be imported granularly via subpaths (e.g. `@gemara/react/control-catalog`) to keep bundle sizes minimal.
- **Isolated Interactivity:** Interactive components (such as collapsible control sections) are segregated into dedicated client subpaths so server-rendered pages remain lean.

## Rendering Documents

Pass your parsed Gemara JSON/YAML document object directly to the corresponding renderer component.

### Server Component Example (Next.js App Router)

```tsx
import { ControlCatalog } from '@gemara/react/control-catalog';
import fs from 'node:fs/promises';
import yaml from 'js-yaml';

export default async function CatalogPage() {
  // Read and parse your Gemara YAML/JSON artifact on the server
  const fileContents = await fs.readFile('./controls.yaml', 'utf8');
  const catalogData = yaml.load(fileContents);

  return (
    <main className="container mx-auto p-6">
      <h1 className="text-2xl font-bold mb-4">Security Controls</h1>
      {/* Server-component-clean renderer */}
      <ControlCatalog data={catalogData} />
    </main>
  );
}
```

## Styling with `data-gemara-part`

Because `@gemara/react` is headless, internal sub-components and elements emit semantic HTML data attributes (`data-gemara-part`). This allows you to style any part of the rendered output using standard CSS attribute selectors or Tailwind CSS variant plugins.

### Common `data-gemara-part` Identifiers

- `data-gemara-part="references"` — Wrapping cross-document references and mappings.
- `data-gemara-part="contact-affiliation"` / `"contact-email"` — Author and contact metadata.
- `data-gemara-part="values"` / `"non-compliance"` — Metric thresholds, evidence values, and compliance logs.
- `data-gemara-part="document-references"` — External links and framework cross-references.

### CSS Example

```css
/* Style references block */
[data-gemara-part="references"] {
  border-left: 3px solid #3b82f6;
  padding-left: 1rem;
  margin-top: 0.5rem;
}

/* Style compliance status badges */
[data-gemara-part="non-compliance"] {
  background-color: #fee2e2;
  color: #991b1b;
  border-radius: 4px;
  padding: 0.25rem 0.5rem;
}
```

## Interactive Components (`CollapsibleGroup`)

When interactive disclosure widgets (accordion/toggle behavior) are desired, use interactive components such as `CollapsibleGroup`. These components can be used in Client Component environments:

```tsx
'use client';

import { CollapsibleGroup } from '@gemara/react/interactive';
import { ControlCatalog } from '@gemara/react/control-catalog';

export default function InteractiveCatalog({ catalogData }) {
  return (
    <CollapsibleGroup defaultExpanded={false}>
      <ControlCatalog data={catalogData} />
    </CollapsibleGroup>
  );
}
```

## Supported Renderers

`@gemara/react` provides renderers for Gemara document types, including:

- **Control Catalogs (`ControlCatalog`):** Controls, parameters, objectives, and implementation guidance.
- **Guidance Catalogs (`GuidanceCatalog`):** Best practices, procedures, and implementation steps.
- **Risk Catalogs (`RiskCatalog`):** Risk definitions, severity vectors, and threat associations.
- **Threat Catalogs (`ThreatCatalog`):** Threat agents, attack vectors, and exploit scenarios.
- **Mapping Documents (`MappingDocument`):** Cross-framework mappings between controls and requirements.
- **Evaluation Logs (`EvaluationLog`):** Automated assessment results, test outcomes, and audit trails.

## Relationship to Other Components

### [The Model](../model)
Provides the conceptual foundation for document types and UI hierarchies.

### [The Schemas](../schema/)
React renderers format documents adhering to the core Gemara CUE schemas.
