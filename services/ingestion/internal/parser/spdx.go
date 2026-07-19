package parser

import (
	"fmt"
	"io"

	"github.com/spdx/tools-golang/json"
	"github.com/spdx/tools-golang/tvsaver"
	"github.com/trustchain/ingestion/internal/domain"
)

// ParseSPDX parses an SPDX JSON (or TagValue) SBOM and extracts normalized components.
func ParseSPDX(r io.Reader, format string) ([]domain.NormalizedComponent, error) {
	var components []domain.NormalizedComponent

	// SPDX tools-golang loading is a bit varied by format
	if format == "json" {
		doc, err := json.Read(r)
		if err != nil {
			return nil, fmt.Errorf("parsing SPDX JSON: %w", err)
		}
		
		for _, pkg := range doc.Packages {
			comp := domain.NormalizedComponent{
				Name:    pkg.PackageName,
				Version: pkg.PackageVersion,
			}
			// Find purl in ExternalRefs
			for _, ref := range pkg.PackageExternalReferences {
				if ref.RefType == "purl" {
					comp.PURL = ref.Locator
					break
				}
			}
			components = append(components, comp)
		}
	} else if format == "tag-value" {
		doc, err := tvsaver.Load2_2(r)
		if err != nil {
			return nil, fmt.Errorf("parsing SPDX Tag-Value: %w", err)
		}
		for _, pkg := range doc.Packages {
			comp := domain.NormalizedComponent{
				Name:    pkg.PackageName,
				Version: pkg.PackageVersion,
			}
			for _, ref := range pkg.PackageExternalReferences {
				if ref.RefType == "purl" {
					comp.PURL = ref.Locator
					break
				}
			}
			components = append(components, comp)
		}
	} else {
		return nil, fmt.Errorf("unsupported SPDX format: %s", format)
	}

	return components, nil
}
