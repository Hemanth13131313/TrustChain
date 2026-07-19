package parser

import (
	"strings"
	"testing"
)

func TestParseCycloneDX(t *testing.T) {
	mockSBOM := `{
		"bomFormat": "CycloneDX",
		"specVersion": "1.4",
		"version": 1,
		"components": [
			{
				"type": "library",
				"name": "lodash",
				"version": "4.17.21",
				"purl": "pkg:npm/lodash@4.17.21",
				"licenses": [
					{
						"license": {
							"id": "MIT"
						}
					}
				]
			}
		]
	}`

	components, err := ParseCycloneDX(strings.NewReader(mockSBOM))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(components) != 1 {
		t.Fatalf("expected 1 component, got %d", len(components))
	}

	c := components[0]
	if c.Name != "lodash" {
		t.Errorf("expected name 'lodash', got '%s'", c.Name)
	}
	if c.Version != "4.17.21" {
		t.Errorf("expected version '4.17.21', got '%s'", c.Version)
	}
	if c.PURL != "pkg:npm/lodash@4.17.21" {
		t.Errorf("expected PURL 'pkg:npm/lodash@4.17.21', got '%s'", c.PURL)
	}
	if c.License != "MIT" {
		t.Errorf("expected license 'MIT', got '%s'", c.License)
	}
}
