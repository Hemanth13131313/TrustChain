package verifier

import (
	"context"
	"fmt"

	"github.com/sigstore/cosign/v2/pkg/cosign"
	"github.com/trustchain/verification/internal/domain"
)

type CosignVerifier struct {
	// For MVP, we will use the default Sigstore public good instance
}

func NewCosignVerifier() *CosignVerifier {
	return &CosignVerifier{}
}

// VerifySignature attempts to verify keyless signatures on the given OCI image reference
func (v *CosignVerifier) VerifySignature(ctx context.Context, imageRef string, artifactID string) ([]domain.SignatureInfo, error) {
	// 1. Resolve the image reference
	// For testing without live registry, we stub this out if imageRef == "mock-image"
	if imageRef == "mock-image" {
		return []domain.SignatureInfo{{
			ArtifactID: artifactID,
			Subject:    "test@trustchain.dev",
			Issuer:     "https://github.com/login/oauth",
			Verified:   true,
		}}, nil
	}

	// In a real implementation, we would call cosign.VerifyImageSignatures
	// options := &cosign.CheckOpts{
	// 	ClaimVerifier: cosign.SimpleClaimVerifier,
	// 	// OIDC issuer constraints would be passed here based on policy
	// }
	// checkedSignatures, bundleVerified, err := cosign.VerifyImageSignatures(ctx, ref, options)

	return nil, fmt.Errorf("live registry verification not implemented in MVP scaffold, use mock-image")
}
