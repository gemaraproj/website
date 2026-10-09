package cmd

import (
	"regexp"
	"testing"
)

func TestProcessContentPreservesLiquidOutput(t *testing.T) {
	termInfos := []TermInfo{{
		OriginalTerm: "organization",
		Slug:         "organization",
		Regex:        regexp.MustCompile(`(?i)\borganization\b`),
	}}
	for _, content := range []string{
		"{{ maintainer.organization }}",
		"{{ organization | append: \"}\" }}",
		"{% assign organization = \"100%\" %}",
	} {
		if got := processContent(content, termInfos, "../model/02-definitions.html"); got != content {
			t.Fatalf("processContent() = %q, want %q", got, content)
		}
	}
}
