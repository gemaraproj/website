package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRootCmdStructure(t *testing.T) {
	if rootCmd.Use != "website-docs" {
		t.Errorf("expected rootCmd.Use to be 'website-docs', got %q", rootCmd.Use)
	}

	subCommands := rootCmd.Commands()
	expectedSubs := map[string]bool{
		"openapi2md": false,
		"lexicon2md": false,
		"termlinker": false,
	}

	for _, sub := range subCommands {
		if _, ok := expectedSubs[sub.Name()]; ok {
			expectedSubs[sub.Name()] = true
		}
	}

	for name, found := range expectedSubs {
		if !found {
			t.Errorf("expected subcommand %q to be registered under rootCmd", name)
		}
	}
}

func TestRunOpenAPI2MDValidation(t *testing.T) {
	tempDir := t.TempDir()

	// 1. Missing -roots when -manifest and -nav are not set
	openAPI2MDFlags.navPath = ""
	openAPI2MDFlags.manifestPath = ""
	openAPI2MDFlags.rootsFlag = ""
	openAPI2MDFlags.inputFile = "testdata/allof.yaml"
	openAPI2MDFlags.outputDir = tempDir

	err := runOpenAPI2MD(openAPI2MDCmd, nil)
	if err == nil || !strings.Contains(err.Error(), "-roots is required") {
		t.Errorf("expected -roots required error, got %v", err)
	}

	// 2. Nav mode execution failure propagation
	openAPI2MDFlags.navPath = filepath.Join(tempDir, "nonexistent-nav.yml")
	openAPI2MDFlags.manifestPath = ""
	openAPI2MDFlags.rootsFlag = ""
	err = runOpenAPI2MD(openAPI2MDCmd, nil)
	if err == nil {
		t.Error("expected error when navPath does not exist, got nil")
	}

	// 3. Manifest mode execution failure propagation
	openAPI2MDFlags.navPath = ""
	openAPI2MDFlags.manifestPath = filepath.Join(tempDir, "nonexistent-manifest.json")
	openAPI2MDFlags.rootsFlag = ""
	err = runOpenAPI2MD(openAPI2MDCmd, nil)
	if err == nil {
		t.Error("expected error when manifestPath does not exist, got nil")
	}

	// 4. Roots mode execution success
	openAPI2MDFlags.navPath = ""
	openAPI2MDFlags.manifestPath = ""
	openAPI2MDFlags.rootsFlag = "Base,Derived"
	openAPI2MDFlags.inputFile = "testdata/allof.yaml"
	openAPI2MDFlags.outputDir = filepath.Join(tempDir, "out")
	err = runOpenAPI2MD(openAPI2MDCmd, nil)
	if err != nil {
		t.Fatalf("expected successful roots mode execution, got: %v", err)
	}
	if _, err := os.Stat(filepath.Join(openAPI2MDFlags.outputDir, "schema.md")); err != nil {
		t.Errorf("expected schema.md to be created in output directory: %v", err)
	}
}

func TestRunTermLinkerExecutionAndLoadTerms(t *testing.T) {
	tempDir := t.TempDir()

	// 1. Test loadTerms with missing file
	missingLex := filepath.Join(tempDir, "missing-lexicon.yaml")
	_, err := loadTerms(missingLex)
	if err == nil {
		t.Error("expected error for missing lexicon file, got nil")
	}

	// 2. Test loadTerms with invalid YAML
	badLex := filepath.Join(tempDir, "bad-lexicon.yaml")
	if err := os.WriteFile(badLex, []byte("terms: [invalid yaml"), 0644); err != nil {
		t.Fatal(err)
	}
	_, err = loadTerms(badLex)
	if err == nil {
		t.Error("expected error for malformed lexicon YAML, got nil")
	}

	// 3. Valid lexicon file
	validLex := filepath.Join(tempDir, "valid-lexicon.yaml")
	lexContent := `
terms:
  - id: TERM-1
    title: Risk Assessment
    definition: Process of identifying security risks.
`
	if err := os.WriteFile(validLex, []byte(lexContent), 0644); err != nil {
		t.Fatal(err)
	}
	terms, err := loadTerms(validLex)
	if err != nil {
		t.Fatalf("loadTerms failed: %v", err)
	}
	if len(terms) != 1 || terms[0].Title != "Risk Assessment" {
		t.Errorf("unexpected terms parsed: %+v", terms)
	}

	// 4. runTermLinker failure on invalid lexicon path
	termLinkerFlags.lexiconFile = missingLex
	termLinkerFlags.docsDir = tempDir
	termLinkerFlags.cleanup = false
	if err := runTermLinker(termLinkerCmd, nil); err == nil {
		t.Error("expected runTermLinker to fail on missing lexicon file, got nil")
	}

	// 5. runTermLinker failure on non-existent docs directory
	termLinkerFlags.lexiconFile = validLex
	termLinkerFlags.docsDir = filepath.Join(tempDir, "does-not-exist")
	termLinkerFlags.cleanup = false
	if err := runTermLinker(termLinkerCmd, nil); err == nil {
		t.Error("expected runTermLinker to fail on missing docs directory, got nil")
	}

	// 6. runTermLinker success in linking mode
	docsDir := filepath.Join(tempDir, "docs")
	if err := os.MkdirAll(docsDir, 0755); err != nil {
		t.Fatal(err)
	}
	docFile := filepath.Join(docsDir, "guide.md")
	docContent := "Perform a Risk Assessment before deployment."
	if err := os.WriteFile(docFile, []byte(docContent), 0644); err != nil {
		t.Fatal(err)
	}

	termLinkerFlags.lexiconFile = validLex
	termLinkerFlags.docsDir = docsDir
	termLinkerFlags.cleanup = false
	if err := runTermLinker(termLinkerCmd, nil); err != nil {
		t.Fatalf("runTermLinker failed in linking mode: %v", err)
	}

	linked, err := os.ReadFile(docFile)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(linked), "[Risk Assessment]") {
		t.Errorf("expected doc to have linked term, got:\n%s", string(linked))
	}

	// 7. runTermLinker success in cleanup mode
	termLinkerFlags.cleanup = true
	if err := runTermLinker(termLinkerCmd, nil); err != nil {
		t.Fatalf("runTermLinker failed in cleanup mode: %v", err)
	}

	cleaned, err := os.ReadFile(docFile)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(cleaned), "[Risk Assessment](") {
		t.Errorf("expected doc to have links cleaned up, got:\n%s", string(cleaned))
	}
}
