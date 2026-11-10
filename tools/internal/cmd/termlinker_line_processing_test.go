package cmd

import (
	"strings"
	"testing"
)

func TestProcessLineComplexEdgeCases(t *testing.T) {
	terms := []Term{
		{ID: "T1", Title: "Control Catalog", Definition: "A catalog containing controls"},
		{ID: "T2", Title: "Control", Definition: "A single control specification"},
		{ID: "T3", Title: "Preventive Enforcement", Definition: "Enforcement preventing violation"},
		{ID: "T4", Title: "Enforcement", Definition: "General enforcement mechanism"},
		{ID: "T5", Title: "Risk", Definition: "Potential harm or vulnerability"},
	}
	termInfos := buildTermInfos(terms)
	defPath := "../model/02-definitions.md"

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "longer compound term matches before substring term",
			input:    "This involves Preventive Enforcement across the boundary.",
			expected: "This involves [Preventive Enforcement](../model/02-definitions.md#preventive-enforcement) across the boundary.",
		},
		{
			name:     "cleans up malformed split term nested link into clean compound slug",
			input:    "See [[Control](url#control) Catalog](url#[control](url#control)-catalog) for guidance.",
			expected: "See [Control Catalog](../model/02-definitions.md#control-catalog) for guidance.",
		},
		{
			name:     "cleans up general nested link into outer compound term slug",
			input:    "Check [[Control](some/path) Catalog](outer/path) now.",
			expected: "Check [Control Catalog](../model/02-definitions.md#control-catalog) now.",
		},
		{
			name:     "cleans up long hash fragment matching compound term suffix",
			input:    "Refer to [Control Catalog](url#unnecessary-prefix-and-extra-long-hash-control-catalog).",
			expected: "Refer to [Control Catalog](url#control-catalog).",
		},
		{
			name:     "cleans up trailing slug pattern repetition",
			input:    "See [Control Catalog](url#control-catalog)-catalog).",
			expected: "See [Control Catalog](url#control-catalog).",
		},
		{
			name:     "removes unintended extra closing parenthesis after link",
			input:    "Read [Control Catalog](url#control-catalog)) before proceeding.",
			expected: "Read [Control Catalog](url#control-catalog) before proceeding.",
		},
		{
			name:     "preserves intentional acronym parenthesized link",
			input:    "Review the ([Control](url#control)) specification.",
			expected: "Review the ([Control](url#control)) specification.",
		},
		{
			name:     "preserves balanced prose parenthesis containing link",
			input:    "(such as an access [Control](url#control)) in place.",
			expected: "(such as an access [Control](url#control)) in place.",
		},
		{
			name:     "cleans up double bracket formatting",
			input:    "Important **[[Control](url#control)]** here.",
			expected: "Important **[Control](url#control)** here.",
		},
		{
			name:     "skips term replacement inside inline code",
			input:    "Do not link inside `Control Catalog` code spans.",
			expected: "Do not link inside `Control Catalog` code spans.",
		},
		{
			name:     "skips term replacement inside liquid template syntax",
			input:    "Template output: {{ page.Control }} or {% include Control %}",
			expected: "Template output: {{ page.Control }} or {% include Control %}",
		},
		{
			name:     "skips term replacement inside html tags and content",
			input:    "<span class=\"Control\">Do not link Control inside HTML</span>",
			expected: "<span class=\"Control\">Do not link Control inside HTML</span>",
		},
		{
			name:     "skips term replacement when already a markdown link",
			input:    "Already linked [Control](https://example.com/custom) should stay unchanged.",
			expected: "Already linked [Control](https://example.com/custom) should stay unchanged.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := processLine(tt.input, termInfos, defPath)
			if got != tt.expected {
				t.Errorf("\nInput:    %s\nGot:      %s\nExpected: %s", tt.input, got, tt.expected)
			}
		})
	}
}

func TestProcessLinePluralAndBoundaryMatching(t *testing.T) {
	terms := []Term{
		{ID: "T1", Title: "Control", Definition: "A control specification"},
		{ID: "T2", Title: "Policy", Definition: "A policy document"},
	}
	termInfos := buildTermInfos(terms)
	defPath := "../model/02-definitions.md"

	// Exact whole-word match
	lineWithExactWord := "The framework has a Control defined."
	processed := processLine(lineWithExactWord, termInfos, defPath)
	if !strings.Contains(processed, "[Control](../model/02-definitions.md#control)") {
		t.Errorf("expected 'Control' to match slug 'control', got:\n%s", processed)
	}

	// Word boundary: partial word should not match (e.g. "Controller" or "Subcontrol")
	partialWord := "The Controller coordinates subcontrols."
	processedPartial := processLine(partialWord, termInfos, defPath)
	if strings.Contains(processedPartial, "[") {
		t.Errorf("partial word should not match term boundary, got:\n%s", processedPartial)
	}
}
