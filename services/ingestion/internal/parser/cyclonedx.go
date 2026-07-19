package parser

import (
	"bytes"
	"fmt"
	"io"

	cdx "github.com/CycloneDX/cyclonedx-go"
	"github.com/trustchain/ingestion/internal/domain"
)

// ParseCycloneDX parses a CycloneDX JSON SBOM and extracts normalized components.
func ParseCycloneDX(r io.Reader) ([]domain.NormalizedComponent, error) {
	// We read everything into a buffer because the decoder needs it, or just use the stream
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, r); err != nil {
		return nil, fmt.Errorf("reading SBOM: %w", err)
	}

	bom := new(cdx.BOM)
	decoder := cdx.NewBOMDecoder(bytes.NewReader(buf.Bytes()), cdx.BOMFileFormatJSON)
	if err := decoder.Decode(bom); err != nil {
		// Fallback to XML if JSON fails
		decoder = cdx.NewBOMDecoder(bytes.NewReader(buf.Bytes()), cdx.BOMFileFormatXML)
		if errXML := decoder.Decode(bom); errXML != nil {
			return nil, fmt.Errorf("failed to decode CycloneDX as JSON or XML: %v", err)
		}
	}

	if bom.Components == nil {
		return []domain.NormalizedComponent{}, nil
	}

	var components []domain.NormalizedComponent
	for _, c := range *bom.Components {
		comp := domain.NormalizedComponent{
			Name:    c.Name,
			Version: c.Version,
			PURL:    c.PackageURL,
		}
		
		if c.Licenses != nil && len(*c.Licenses) > 0 {
			lic := (*c.Licenses)[0]
			if lic.License != nil {
				comp.License = lic.License.ID
				if comp.License == "" {
					comp.License = lic.License.Name
				}
			}
		}
		
		components = append(components, comp)
	}

	return components, nil
}
