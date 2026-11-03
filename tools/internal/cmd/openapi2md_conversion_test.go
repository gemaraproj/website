package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateAliasBlock(t *testing.T) {
	aliasSchema := Schema{
		Type:        "string",
		Description: "A unique identifier string.\n\n(Enforced by the CUE schema; not by this OpenAPI projection.)",
		Format:      "uri",
		Pattern:     "^[a-z0-9-]+$",
	}

	// Test top-level heading (##)
	block := generateAliasBlock("ArtifactID", aliasSchema, false)
	if !strings.HasPrefix(block, "## `ArtifactID`\n\n") {
		t.Errorf("expected top-level heading '## `ArtifactID`', got:\n%s", block)
	}
	if !strings.Contains(block, "A unique identifier string.") {
		t.Errorf("expected description to be present in alias block:\n%s", block)
	}
	if strings.Contains(block, "(Enforced by the CUE schema; not by this OpenAPI projection.)") {
		t.Errorf("expected CUE projection note to be stripped from alias block:\n%s", block)
	}
	if !strings.Contains(block, "- **Type**: `string`") {
		t.Errorf("missing type in alias block:\n%s", block)
	}
	if !strings.Contains(block, "- **Format**: `uri`") {
		t.Errorf("missing format in alias block:\n%s", block)
	}
	if !strings.Contains(block, "- **Value**: `^[a-z0-9-]+$`") {
		t.Errorf("missing pattern in alias block:\n%s", block)
	}
	if !strings.HasSuffix(block, "\n---\n\n") {
		t.Errorf("expected divider suffix, got:\n%s", block)
	}

	// Test subheading (###)
	subBlock := generateAliasBlock("SubID", aliasSchema, true)
	if !strings.HasPrefix(subBlock, "### `SubID`\n\n") {
		t.Errorf("expected subheading '### `SubID`', got:\n%s", subBlock)
	}
}

func TestConvertFromNav(t *testing.T) {
	tempDir := t.TempDir()
	outDir := filepath.Join(tempDir, "docs")
	validNav := filepath.Join(tempDir, "valid_nav.yml")
	if err := os.WriteFile(validNav, []byte("pages: []\n"), 0644); err != nil {
		t.Fatal(err)
	}

	// Missing OpenAPI spec file
	if err := convertFromNav("nonexistent.yaml", outDir, validNav); err == nil {
		t.Error("expected error for missing OpenAPI input file, got nil")
	}

	// Invalid OpenAPI YAML
	badYAML := filepath.Join(tempDir, "bad_openapi.yaml")
	if err := os.WriteFile(badYAML, []byte("invalid: yaml: ["), 0644); err != nil {
		t.Fatal(err)
	}
	if err := convertFromNav(badYAML, outDir, validNav); err == nil {
		t.Error("expected error for invalid OpenAPI YAML, got nil")
	}

	// Missing nav file
	if err := convertFromNav("testdata/allof.yaml", outDir, filepath.Join(tempDir, "missing_nav.yml")); err == nil {
		t.Error("expected error for missing nav file, got nil")
	}

	// Nav references non-existent schema
	brokenNav := filepath.Join(tempDir, "broken_nav.yml")
	brokenNavContent := `
pages:
  - title: Bad Page
    filename: bad-page.md
    schemas:
      - NoSuchSchema
`
	if err := os.WriteFile(brokenNav, []byte(brokenNavContent), 0644); err != nil {
		t.Fatal(err)
	}
	if err := convertFromNav("testdata/allof.yaml", outDir, brokenNav); err == nil {
		t.Error("expected error when nav references nonexistent schema, got nil")
	}

	// Valid nav execution with schemas from allof.yaml (Derived and Metadata alias)
	goodNav := filepath.Join(tempDir, "good_nav.yml")
	goodNavContent := `
pages:
  - title: Model Specs
    filename: model-specs
    schemas:
      - Derived
      - Metadata
  - title: Auto Slug Page
    schemas:
      - Base
`
	if err := os.WriteFile(goodNav, []byte(goodNavContent), 0644); err != nil {
		t.Fatal(err)
	}

	if err := convertFromNav("testdata/allof.yaml", outDir, goodNav); err != nil {
		t.Fatalf("convertFromNav failed: %v", err)
	}

	// Verify generated file 1
	p1, err := os.ReadFile(filepath.Join(outDir, "model-specs.md"))
	if err != nil {
		t.Fatalf("failed to read generated page 1: %v", err)
	}
	p1Content := string(p1)
	if !strings.Contains(p1Content, "`Derived`") {
		t.Errorf("model-specs.md missing Derived schema section:\n%s", p1Content)
	}
	if !strings.Contains(p1Content, "`Metadata`") {
		t.Errorf("model-specs.md missing Metadata alias block:\n%s", p1Content)
	}

	// Verify auto-slug generated file 2 (auto_slug_page -> auto-slug-page.md)
	p2, err := os.ReadFile(filepath.Join(outDir, "auto-slug-page.md"))
	if err != nil {
		t.Fatalf("failed to read auto-slug generated page 2: %v", err)
	}
	if !strings.Contains(string(p2), "`Base`") {
		t.Errorf("auto-slug-page.md missing Base schema section:\n%s", string(p2))
	}
}

func TestConvertPerFile(t *testing.T) {
	tempDir := t.TempDir()
	outDir := filepath.Join(tempDir, "out")
	manifestPath := filepath.Join(tempDir, "manifest.json")

	// Missing manifest
	if err := convertPerFile("testdata/allof.yaml", outDir, filepath.Join(tempDir, "nonexistent.json")); err == nil {
		t.Error("expected error for missing manifest file, got nil")
	}

	// Missing OpenAPI input file
	goodManifestContent := `{"base.cue": ["Base"], "empty.cue": []}`
	if err := os.WriteFile(manifestPath, []byte(goodManifestContent), 0644); err != nil {
		t.Fatal(err)
	}
	if err := convertPerFile("nonexistent.yaml", outDir, manifestPath); err == nil {
		t.Error("expected error for missing OpenAPI input file, got nil")
	}

	// Valid execution
	if err := convertPerFile("testdata/allof.yaml", outDir, manifestPath); err != nil {
		t.Fatalf("convertPerFile failed: %v", err)
	}

	// Verify generated file base.md
	baseFile, err := os.ReadFile(filepath.Join(outDir, "base.md"))
	if err != nil {
		t.Fatalf("failed to read generated base.md: %v", err)
	}
	if !strings.Contains(string(baseFile), "`Base`") {
		t.Errorf("base.md missing Base schema section:\n%s", string(baseFile))
	}

	// Empty cue file in manifest should not create a markdown file
	if _, err := os.Stat(filepath.Join(outDir, "empty.md")); !os.IsNotExist(err) {
		t.Errorf("expected empty.md not to be created, but it exists")
	}
}

func TestConvertOpenAPIToMarkdownRoots(t *testing.T) {
	tempDir := t.TempDir()
	outDir := filepath.Join(tempDir, "out")

	// Missing OpenAPI input file
	if err := convertOpenAPIToMarkdown("nonexistent.yaml", outDir, []string{"Base"}); err == nil {
		t.Error("expected error for missing input file, got nil")
	}

	// Non-existent root schema
	if err := convertOpenAPIToMarkdown("testdata/allof.yaml", outDir, []string{"NoSuchRoot"}); err == nil {
		t.Error("expected error for nonexistent root schema, got nil")
	}

	// Valid conversion with Base and Derived
	if err := convertOpenAPIToMarkdown("testdata/allof.yaml", outDir, []string{"Base", "Derived"}); err != nil {
		t.Fatalf("convertOpenAPIToMarkdown failed: %v", err)
	}

	// Verify generated schema.md containing roots and aliases
	schemaFile, err := os.ReadFile(filepath.Join(outDir, "schema.md"))
	if err != nil {
		t.Fatalf("failed to read schema.md: %v", err)
	}
	content := string(schemaFile)
	if !strings.Contains(content, "`Base`") {
		t.Errorf("schema.md missing Base section:\n%s", content)
	}
	if !strings.Contains(content, "`Derived`") {
		t.Errorf("schema.md missing Derived section:\n%s", content)
	}
	if !strings.Contains(content, "## Aliases") {
		t.Errorf("schema.md missing Aliases section:\n%s", content)
	}
}
