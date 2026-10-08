---
layout: page
title: Contributing
nav-title: Contributing
---

Thanks for your interest in contributing to Gemara! Contributions of all
kinds are welcome including docs, schemas, SDKs, examples, and ideas. This page
covers a few things worth knowing before you open a pull request, including
how we like AI-assisted work to be attributed.

## Getting started

- Sign off your commits with `git commit -s` (see [DCO](#dco) below).
- Use [Conventional Commits](https://www.conventionalcommits.org/) for PR
  titles, for example `docs: fix typo on the model page`.
- Run `make cleanup` before committing to clear out generated files.

Planning something larger? Come say hello first. We're happy to help shape
the idea:

- **Slack:** [#gemara](https://openssf.slack.com/archives/C09A9PP765Q) on OpenSSF Slack
- **Meetings:** every other Thursday — see the [OpenSSF calendar](https://calendar.google.com/calendar/u/0?cid=czYzdm9lZmhwNWk5cGZsdGI1cTY3bmdwZXNAZ3JvdXAuY2FsZW5kYXIuZ29vZ2xlLmNvbQ)

## DCO

Gemara uses the [Developer Certificate of Origin](https://developercertificate.org/),
so please sign off your commits with `git commit -s`. The sign-off says that
you wrote the change, or otherwise have the right to submit it. So it should
be yours, even when an AI tool helped you write it.

## Working with AI

AI-assisted contributions are welcome. We just ask that you be open about
where AI was involved, so reviewers and future maintainers understand what
they're looking at. Please refer to the Linux Foundation's [Generative AI Policy](https://www.linuxfoundation.org/legal/generative-ai)
for additional guidance.

### Commits

When AI helped author a commit, add an `Assisted-by:` trailer naming the tool:

```
Assisted-by: Claude Code <noreply@anthropic.com>
```

Please don't use `Co-authored-by:` for AI tools (this is blocked by our branch protections).
The DCO sign-off is a statement made by a person, so co-authorship should reflect the humans
involved.

### PRs, issues, and comments

If an AI generated or largely drafted a PR description, issue, or comment,
add a 🤖 so it's clear which text is AI-generated and which is yours.

### Understanding what you contribute

Whether a change is hand-written or AI-assisted, the goal is the same: you
should understand it well enough to walk someone through it. A good rule of
thumb is to imagine explaining the change months from now. You might not
remember every function name, but the reasoning, the design, and the
packages it relies on should still feel like yours.

If a change gets that kind of understanding, it's in good shape. If it
doesn't yet, it's worth spending a bit more time with it before asking for a
review.
