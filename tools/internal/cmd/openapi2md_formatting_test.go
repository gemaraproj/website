// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"strings"
	"testing"
)

func TestResolveSchemaRefErrors(t *testing.T) {
	spec := OpenAPISpec{
		Components: OpenAPIComponents{
			Schemas: map[string]interface{}{
				"ValidSchema": map[string]interface{}{
					"type": "string",
				},
			},
		},
	}

	// 1. Invalid prefix format
	invalidRefs := []string{
		"#/definitions/ValidSchema",
		"#/components/parameters/Param",
		"ValidSchema",
		"",
	}
	for _, ref := range invalidRefs {
		_, err := resolveSchemaRef(ref, spec)
		if err == nil || !strings.Contains(err.Error(), "invalid ref format") {
			t.Errorf("resolveSchemaRef(%q) expected 'invalid ref format' error, got %v", ref, err)
		}
	}

	// 2. Non-existent schema target
	_, err := resolveSchemaRef("#/components/schemas/NonExistentSchema", spec)
	if err == nil || !strings.Contains(err.Error(), "schema not found: NonExistentSchema") {
		t.Errorf("expected 'schema not found: NonExistentSchema' error, got %v", err)
	}

	// 3. Valid schema ref resolves successfully
	schema, err := resolveSchemaRef("#/components/schemas/ValidSchema", spec)
	if err != nil {
		t.Fatalf("unexpected error resolving valid ref: %v", err)
	}
	if schema.Type != "string" {
		t.Errorf("resolved schema type = %q, want string", schema.Type)
	}
}

func TestFormatFieldType(t *testing.T) {
	spec := OpenAPISpec{
		Components: OpenAPIComponents{
			Schemas: map[string]interface{}{
				"CustomActor": map[string]interface{}{
					"type": "object",
				},
			},
		},
	}

	schemaToFile := map[string]string{
		"CustomActor": "entities",
	}

	tests := []struct {
		name         string
		fieldSchema  Schema
		schemaToFile map[string]string
		want         string
	}{
		{
			name: "primitive string",
			fieldSchema: Schema{
				Type: "string",
			},
			schemaToFile: schemaToFile,
			want:         "string",
		},
		{
			name: "ref with matching schemaToFile generates markdown link without md ext",
			fieldSchema: Schema{
				Ref: "#/components/schemas/CustomActor",
			},
			schemaToFile: schemaToFile,
			want:         "[CustomActor](entities#customactor)",
		},
		{
			name: "ref without mapping in schemaToFile returns bare type name",
			fieldSchema: Schema{
				Ref: "#/components/schemas/UnmappedType",
			},
			schemaToFile: schemaToFile,
			want:         "UnmappedType",
		},
		{
			name: "array of primitive types",
			fieldSchema: Schema{
				Type: "array",
				Items: map[string]interface{}{
					"type": "integer",
				},
			},
			schemaToFile: schemaToFile,
			want:         "array[integer]",
		},
		{
			name: "array of custom ref types with linking",
			fieldSchema: Schema{
				Type: "array",
				Items: map[string]interface{}{
					"$ref": "#/components/schemas/CustomActor",
				},
			},
			schemaToFile: schemaToFile,
			want:         "array[[CustomActor](entities#customactor)]",
		},
		{
			name: "array of unmapped ref types without linking",
			fieldSchema: Schema{
				Type: "array",
				Items: map[string]interface{}{
					"$ref": "#/components/schemas/UnmappedItem",
				},
			},
			schemaToFile: schemaToFile,
			want:         "array[UnmappedItem]",
		},
		{
			name: "empty schema returns empty string",
			fieldSchema: Schema{},
			schemaToFile: schemaToFile,
			want:         "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatFieldType(tt.fieldSchema, spec, tt.schemaToFile)
			if got != tt.want {
				t.Errorf("formatFieldType() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFormatFieldInlineAndNested(t *testing.T) {
	spec := OpenAPISpec{
		Components: OpenAPIComponents{
			Schemas: map[string]interface{}{
				"TargetSchema": map[string]interface{}{
					"type":        "string",
					"description": "Fallback target description",
				},
			},
		},
	}
	schemaToFile := map[string]string{
		"TargetSchema": "targets",
	}

	// 1. Required field with prefix
	fieldLine, desc := formatFieldInline("id", Schema{
		Type:        "string",
		Description: "Primary unique identifier",
	}, spec, "metadata", true, schemaToFile)

	if fieldLine != "`metadata.id` **string** _Required_" {
		t.Errorf("unexpected fieldLine: %q", fieldLine)
	}
	if desc != "Primary unique identifier" {
		t.Errorf("unexpected description: %q", desc)
	}

	// 2. Optional field without prefix inheriting ref description
	fieldLine, desc = formatFieldInline("target", Schema{
		Ref: "#/components/schemas/TargetSchema",
	}, spec, "", false, schemaToFile)

	if fieldLine != "`target` **[TargetSchema](targets#targetschema)**" {
		t.Errorf("unexpected fieldLine: %q", fieldLine)
	}
	if desc != "Fallback target description" {
		t.Errorf("expected inherited description from ref target, got %q", desc)
	}

	// 3. formatFieldWithNested produces formatted block with trailing newlines
	nestedBlock := formatFieldWithNested("entry", Schema{
		Type:        "string",
		Description: "Entry description",
	}, spec, true, schemaToFile)

	if !strings.Contains(nestedBlock, "`entry` **string** _Required_\n\n") {
		t.Errorf("nested block missing field definition: %q", nestedBlock)
	}
	if !strings.Contains(nestedBlock, "Entry description\n") {
		t.Errorf("nested block missing description: %q", nestedBlock)
	}
}

func TestMergeSchemasEdgeCases(t *testing.T) {
	target := Schema{
		Type:        "object",
		Description: "Base description",
		Required:    []string{"id", "title"},
		Properties: map[string]interface{}{
			"id": map[string]interface{}{"type": "string"},
		},
		XStatus: "stable",
	}

	source := Schema{
		Description: "Overridden description",
		Pattern:     "^[a-z]+$",
		Format:      "uri",
		XStatus:     "experimental",
		Required:    []string{"id", "extra"}, // "id" is duplicate, "extra" is new
		Properties: map[string]interface{}{
			"extra": map[string]interface{}{"type": "integer"},
		},
	}

	mergeSchemas(&target, source)

	if target.Description != "Overridden description" {
		t.Errorf("target.Description = %q, want overridden", target.Description)
	}
	if target.Pattern != "^[a-z]+$" {
		t.Errorf("target.Pattern = %q, want ^[a-z]+$", target.Pattern)
	}
	if target.Format != "uri" {
		t.Errorf("target.Format = %q, want uri", target.Format)
	}
	if target.XStatus != "experimental" {
		t.Errorf("target.XStatus = %q, want experimental", target.XStatus)
	}

	// Check deduplicated required slice: should be ["id", "title", "extra"]
	if len(target.Required) != 3 {
		t.Fatalf("target.Required length = %d, want 3 (deduplicated)", len(target.Required))
	}
	expectedRequired := []string{"id", "title", "extra"}
	for i, req := range expectedRequired {
		if target.Required[i] != req {
			t.Errorf("target.Required[%d] = %q, want %q", i, target.Required[i], req)
		}
	}

	// Check properties merged
	if len(target.Properties) != 2 {
		t.Fatalf("target.Properties length = %d, want 2", len(target.Properties))
	}
	if _, ok := target.Properties["id"]; !ok {
		t.Errorf("missing property 'id'")
	}
	if _, ok := target.Properties["extra"]; !ok {
		t.Errorf("missing property 'extra'")
	}
}
