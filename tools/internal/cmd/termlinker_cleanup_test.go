package cmd

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
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
			name: "preserves fenced code blocks",
			input: "Here is code:\n```go\nvar c = [Control](/specifications/02-definitions#control)\n```\nAnd text [Control](/specifications/02-definitions#control).",
			expected: "Here is code:\n```go\nvar c = [Control](/specifications/02-definitions#control)\n```\nAnd text Control.",
		},
		{
			name:     "strips termlinker link but leaves non-termlinker links intact",
			input:    "See [Control](/specifications/02-definitions#control) and external [Google](https://google.com).",
			expected: "See Control and external [Google](https://google.com).",
		},
		{
			name:     "handles parentheses in link text gracefully",
			input:    "Review [Assessment](/specifications/02-definitions#assessment) procedure.",
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

func TestLoadTermsValidation(t *testing.T) {
	tmpDir := t.TempDir()

	// 1. Valid lexicon file
	validFile := filepath.Join(tmpDir, "lexicon.yaml")
	validContent := `terms:
  - title: Control
    definition: A security control.
  - title: Assessment
    definition: An assessment process.
`
	if err := os.WriteFile(validFile, []byte(validContent), 0o644); err != nil {
		t.Fatal(err)
	}

	terms, err := loadTerms(validFile)
	if err != nil {
		t.Fatalf("expected loadTerms to succeed, got %v", err)
	}
	if len(terms) != 2 {
		t.Fatalf("expected 2 terms, got %d", len(terms))
	}
	if terms[0].Title != "Control" || terms[1].Title != "Assessment" {
		t.Errorf("unexpected loaded term titles: %v, %v", terms[0].Title, terms[1].Title)
	}

	// 2. Non-existent file error
	_, err = loadTerms(filepath.Join(tmpDir, "does-not-exist.yaml"))
	if err == nil {
		t.Fatal("expected error for non-existent lexicon file, got nil")
	}

	// 3. Malformed YAML file
	malformedFile := filepath.Join(tmpDir, "bad.yaml")
	if err := os.WriteFile(malformedFile, []byte("terms: [unclosed list"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err = loadTerms(malformedFile)
	if err == nil {
		t.Fatal("expected parse error for malformed YAML, got nil")
	}
}

func TestProcessLineTermReplacementsAndEscaping(t *testing.T) {
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

	// 1. Replaces term in plain text
	line := "This is a Control in Gemara."
	res := processLine(line, termInfos, defPath)
	if !strings.Contains(res, "[Control](02-definitions.md#control)") {
		t.Errorf("expected termlinker markdown link, got: %s", res)
	}

	// 2. Skips inline code blocks (`code`)
	codeLine := "Execute `Control` inside code block."
	resCode := processLine(codeLine, termInfos, defPath)
	if strings.Contains(resCode, "[Control]") {
		t.Errorf("expected terms in inline code to not be converted to links, got: %s", resCode)
	}

	// 3. Skips already formed markdown links
	linkLine := "Check [existing Control](https://example.com/control)."
	resLink := processLine(linkLine, termInfos, defPath)
	if strings.Contains(resLink, "[[") {
		t.Errorf("expected no nested links created, got: %s", resLink)
	}

	// 4. Skips HTML tags (<tag>)
	htmlLine := "<span class=\"Control\">text</span>"
	resHTML := processLine(htmlLine, termInfos, defPath)
	if strings.Contains(resHTML, "<span class=\"[Control]") {
		t.Errorf("expected HTML attributes to remain unlinked, got: %s", resHTML)
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

	// Jekyll include transition
	state.update("{% include header.html %}")
	if state.inJekyllInclude {
		t.Errorf("single line include should start and end on same line")
	}

	// Multi-line HTML block transition
	state.update("<div class=\"wrapper\">")
	if !state.inHTMLBlock || state.htmlTagStack != 1 {
		t.Errorf("expected inHTMLBlock to be true with stack 1, got inHTMLBlock=%v stack=%d", state.inHTMLBlock, state.htmlTagStack)
	}
	state.update("</div>")
	if state.inHTMLBlock || state.htmlTagStack != 0 {
		t.Errorf("expected inHTMLBlock to be false with stack 0 after closing div")
	}
}
