package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestConvertFromNavWorkflows(t *testing.T) {
	tempDir := t.TempDir()

	validOpenAPI := `
openapi: "3.0.0"
info:
  title: Test API
  version: "1.0.0"
components:
  schemas:
    Catalog:
      type: object
      description: A control catalog
      properties:
        id:
          type: string
        title:
          type: string
      required:
        - id
    AliasCatalog:
      type: string
      description: An alias type
`
	openAPIFile := filepath.Join(tempDir, "openapi.yaml")
	if err := os.WriteFile(openAPIFile, []byte(validOpenAPI), 0644); err != nil {
		t.Fatal(err)
	}

	// 1. Missing OpenAPI file error
	err := convertFromNav(filepath.Join(tempDir, "nonexistent-openapi.yaml"), tempDir, filepath.Join(tempDir, "nav.yml"))
	if err == nil || !strings.Contains(err.Error(), "failed to read OpenAPI file") {
		t.Errorf("expected read OpenAPI file error, got: %v", err)
	}

	// 2. Malformed OpenAPI YAML error
	badOpenAPIFile := filepath.Join(tempDir, "bad-openapi.yaml")
	if err := os.WriteFile(badOpenAPIFile, []byte("components: [invalid yaml"), 0644); err != nil {
		t.Fatal(err)
	}
	err = convertFromNav(badOpenAPIFile, tempDir, filepath.Join(tempDir, "nav.yml"))
	if err == nil || !strings.Contains(err.Error(), "failed to parse OpenAPI YAML") {
		t.Errorf("expected parse OpenAPI YAML error, got: %v", err)
	}

	// 3. Missing nav file error
	err = convertFromNav(openAPIFile, tempDir, filepath.Join(tempDir, "nonexistent-nav.yml"))
	if err == nil {
		t.Error("expected error for missing nav file, got nil")
	}

	// 4. Missing schema reference in OpenAPI spec
	missingSchemaNav := filepath.Join(tempDir, "missing-schema-nav.yml")
	missingNavContent := `
pages:
  - title: Unknown Page
    filename: unknown.md
    schemas:
      - NonExistentSchema
`
	if err := os.WriteFile(missingSchemaNav, []byte(missingNavContent), 0644); err != nil {
		t.Fatal(err)
	}
	err = convertFromNav(openAPIFile, tempDir, missingSchemaNav)
	if err == nil || !strings.Contains(err.Error(), "not found in OpenAPI spec") {
		t.Errorf("expected missing schema in spec error, got: %v", err)
	}

	// 5. Successful nav conversion with custom filename, fallback slugify, and alias blocks
	validNavFile := filepath.Join(tempDir, "valid-nav.yml")
	validNavContent := `
pages:
  - title: Catalog Management
    filename: catalog-spec
    schemas:
      - Catalog
  - title: Alias Definitions Page
    schemas:
      - AliasCatalog
`
	if err := os.WriteFile(validNavFile, []byte(validNavContent), 0644); err != nil {
		t.Fatal(err)
	}

	outDir := filepath.Join(tempDir, "nav-out")
	err = convertFromNav(openAPIFile, outDir, validNavFile)
	if err != nil {
		t.Fatalf("expected successful convertFromNav execution, got: %v", err)
	}

	// Verify catalog-spec.md
	catalogPath := filepath.Join(outDir, "catalog-spec.md")
	catalogContent, err := os.ReadFile(catalogPath)
	if err != nil {
		t.Fatalf("expected %s to be created: %v", catalogPath, err)
	}
	if !strings.Contains(string(catalogContent), "## `Catalog`") || !strings.Contains(string(catalogContent), "A control catalog") {
		t.Errorf("expected Catalog schema documentation in %s, got:\n%s", catalogPath, string(catalogContent))
	}

	// Verify slugified fallback filename: alias-definitions-page.md
	aliasPath := filepath.Join(outDir, "alias-definitions-page.md")
	aliasContent, err := os.ReadFile(aliasPath)
	if err != nil {
		t.Fatalf("expected slugified %s to be created: %v", aliasPath, err)
	}
	if !strings.Contains(string(aliasContent), "AliasCatalog") {
		t.Errorf("expected AliasCatalog in %s, got:\n%s", aliasPath, string(aliasContent))
	}
}

func TestConvertPerFileWorkflows(t *testing.T) {
	tempDir := t.TempDir()

	validOpenAPI := `
openapi: "3.0.0"
info:
  title: Test API
  version: "1.0.0"
components:
  schemas:
    Control:
      type: object
      description: A control object
      properties:
        id:
          type: string
    SimpleAlias:
      type: string
      description: Simple alias string
`
	openAPIFile := filepath.Join(tempDir, "openapi.yaml")
	if err := os.WriteFile(openAPIFile, []byte(validOpenAPI), 0644); err != nil {
		t.Fatal(err)
	}

	// 1. Missing manifest file error
	err := convertPerFile(openAPIFile, tempDir, filepath.Join(tempDir, "nonexistent-manifest.json"))
	if err == nil {
		t.Error("expected error for missing manifest file, got nil")
	}

	// 2. Missing OpenAPI file error
	emptyManifestFile := filepath.Join(tempDir, "manifest.json")
	if err := os.WriteFile(emptyManifestFile, []byte(`{}`), 0644); err != nil {
		t.Fatal(err)
	}
	err = convertPerFile(filepath.Join(tempDir, "nonexistent-openapi.yaml"), tempDir, emptyManifestFile)
	if err == nil || !strings.Contains(err.Error(), "failed to read OpenAPI file") {
		t.Errorf("expected read OpenAPI file error, got: %v", err)
	}

	// 3. Malformed OpenAPI YAML error
	badOpenAPIFile := filepath.Join(tempDir, "bad-openapi.yaml")
	if err := os.WriteFile(badOpenAPIFile, []byte("components: [broken yaml"), 0644); err != nil {
		t.Fatal(err)
	}
	err = convertPerFile(badOpenAPIFile, tempDir, emptyManifestFile)
	if err == nil || !strings.Contains(err.Error(), "failed to parse OpenAPI YAML") {
		t.Errorf("expected parse OpenAPI YAML error, got: %v", err)
	}

	// 4. Successful per-file conversion with both objects, aliases, and skipping empty/unmatched schemas
	manifestFile := filepath.Join(tempDir, "valid-manifest.json")
	manifestContent := `{
  "controlcatalog.cue": ["Control", "MissingSchemaInSpec"],
  "aliases.cue": ["SimpleAlias"],
  "empty.cue": []
}`
	if err := os.WriteFile(manifestFile, []byte(manifestContent), 0644); err != nil {
		t.Fatal(err)
	}

	outDir := filepath.Join(tempDir, "per-file-out")
	err = convertPerFile(openAPIFile, outDir, manifestFile)
	if err != nil {
		t.Fatalf("expected successful convertPerFile execution, got: %v", err)
	}

	// Verify controlcatalog.md generated
	controlPath := filepath.Join(outDir, "controlcatalog.md")
	controlContent, err := os.ReadFile(controlPath)
	if err != nil {
		t.Fatalf("expected %s to be created: %v", controlPath, err)
	}
	if !strings.Contains(string(controlContent), "## `Control`") {
		t.Errorf("expected Control schema in %s, got:\n%s", controlPath, string(controlContent))
	}

	// Verify aliases.md generated
	aliasesPath := filepath.Join(outDir, "aliases.md")
	aliasesContent, err := os.ReadFile(aliasesPath)
	if err != nil {
		t.Fatalf("expected %s to be created: %v", aliasesPath, err)
	}
	if !strings.Contains(string(aliasesContent), "SimpleAlias") {
		t.Errorf("expected SimpleAlias in %s, got:\n%s", aliasesPath, string(aliasesContent))
	}

	// Verify empty.cue was skipped (no empty.md created)
	emptyPath := filepath.Join(outDir, "empty.md")
	if _, err := os.Stat(emptyPath); !os.IsNotExist(err) {
		t.Errorf("expected %s to not be created for empty schema list", emptyPath)
	}
}
