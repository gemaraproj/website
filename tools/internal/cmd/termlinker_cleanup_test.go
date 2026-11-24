package cmd

import (
	"regexp"
	"testing"
)

func TestCleanupContentMarkdownCleaning(t *testing.T) {
	controlRegex, _ := regexp.Compile(`(?i)\bControl\b`)
	assessmentRegex, _ := regexp.Compile(`(?i)\bAssessment\b`)

	terms := []TermInfo{
		{
			OriginalTerm: "Control",
			LowerTerm:    "control",
			Slug:         "control",
			Regex:        controlRegex,
		},
		{
			OriginalTerm: "Assessment",
			LowerTerm:    "assessment",
			Slug:         "assessment",
			Regex:        assessmentRegex,
		},
	}

	defPath := "02-definitions.md"

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "preserves front matter without changes",
			input: `---
title: Sample [Control](/specifications/02-definitions#control)
layout: default
---
This is a [Control](/specifications/02-definitions#control) in body.`,
			expected: `---
title: Sample [Control](/specifications/02-definitions#control)
layout: default
---
This is a Control in body.`,
		},
		{
			name:     "preserves fenced code blocks",
			input:    "Here is code:\n```go\nvar c = [Control](/specifications/02-definitions#control)\n```\nAnd text [Control](/specifications/02-definitions#control).",
			expected: "Here is code:\n```go\nvar c = [Control](/specifications/02-definitions#control)\n```\nAnd text Control.",
		},
		{
			name:     "strips termlinker link but leaves non-termlinker links intact",
			input:    "See [Control](/specifications/02-definitions#control) and external [Google](https://google.com).",
			expected: "See Control and external [Google](https://google.com).",
		},
		{
			name:     "handles nested parentheses in link URL gracefully",
			input:    "Review [Assessment](/docs(v2)/02-definitions#assessment) procedure.",
			expected: "Review Assessment procedure.",
		},
		{
			name:     "handles multiple term links on the same line",
			input:    "Each [Control](/specifications/02-definitions#control) requires an [Assessment](/specifications/02-definitions#assessment).",
			expected: "Each Control requires an Assessment.",
		},
		{
			name:     "leaves plain text unchanged",
			input:    "No links here at all.",
			expected: "No links here at all.",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := cleanupContent(tc.input, terms, defPath)
			if got != tc.expected {
				t.Errorf("cleanupContent() = %q, want %q", got, tc.expected)
			}
		})
	}
}

func TestProcessLineSkipsExistingLinks(t *testing.T) {
	regexControl, _ := regexp.Compile(`(?i)\bControl\b`)

	termInfos := []TermInfo{
		{
			OriginalTerm: "Control",
			LowerTerm:    "control",
			Slug:         "control",
			Regex:        regexControl,
		},
	}

	defPath := "02-definitions.md"

	linkLine := "Check [existing Control](https://example.com/control)."
	resLink := processLine(linkLine, termInfos, defPath)
	if resLink != linkLine {
		t.Errorf("processLine modified existing link: got %q, want %q", resLink, linkLine)
	}

	htmlLine := `<span class="Control">text</span>`
	resHTML := processLine(htmlLine, termInfos, defPath)
	if resHTML != htmlLine {
		t.Errorf("processLine modified HTML tag: got %q, want %q", resHTML, htmlLine)
	}
}

func TestContentStateTransitions(t *testing.T) {
	var state contentState

	// Frontmatter transition
	state.update("---")
	if !state.inFrontMatter {
		t.Errorf("expected inFrontMatter to be true")
	}
	state.update("title: Test")
	if !state.inFrontMatter {
		t.Errorf("expected inFrontMatter to remain true")
	}
	state.update("---")
	if state.inFrontMatter {
		t.Errorf("expected inFrontMatter to be false after closing delimiter")
	}

	// Code block transition
	state.update("```go")
	if !state.inCodeBlock {
		t.Errorf("expected inCodeBlock to be true")
	}
	state.update("fmt.Println()")
	if !state.inCodeBlock {
		t.Errorf("expected inCodeBlock to remain true")
	}
	state.update("```")
	if state.inCodeBlock {
		t.Errorf("expected inCodeBlock to be false after closing fence")
	}

	// Multi-line Jekyll include transition
	state.update("{% include header.html")
	if !state.inJekyllInclude {
		t.Errorf("expected inJekyllInclude to be true after opening multi-line include")
	}
	state.update("   param=\"value\"")
	if !state.inJekyllInclude {
		t.Errorf("expected inJekyllInclude to remain true inside include body")
	}
	state.update("%}")
	if state.inJekyllInclude {
		t.Errorf("expected inJekyllInclude to be false after closing multi-line include")
	}

	// Multi-line HTML block transition
	state.update(`<div class="wrapper">`)
	if !state.inHTMLBlock || state.htmlTagStack != 1 {
		t.Errorf("expected inHTMLBlock to be true with stack 1, got inHTMLBlock=%v stack=%d", state.inHTMLBlock, state.htmlTagStack)
	}
	state.update("</div>")
	if state.inHTMLBlock || state.htmlTagStack != 0 {
		t.Errorf("expected inHTMLBlock to be false with stack 0 after closing div")
	}
}
