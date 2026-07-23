// SPDX-License-Identifier: Apache-2.0

package cmd

// These types mirror the OpenAPI document emitted by the spec repo's
// `gemara-docs cue2openapi` command (github.com/gemaraproj/gemara, cmd/).

type OpenAPISpec struct {
	OpenAPI    string            `yaml:"openapi" json:"openapi"`
	Info       OpenAPIInfo       `yaml:"info" json:"info"`
	Components OpenAPIComponents `yaml:"components" json:"components"`
}

type OpenAPIInfo struct {
	Title       string `yaml:"title" json:"title"`
	Version     string `yaml:"version" json:"version"`
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
}

type OpenAPIComponents struct {
	Schemas map[string]interface{} `yaml:"schemas" json:"schemas"`
}
