# TRUSTCHAIN Developer Onboarding

Welcome to TRUSTCHAIN! To ensure your application can be successfully deployed to our Kubernetes clusters, your container images must meet our zero-trust requirements:
1. They must have an associated **SBOM** (Software Bill of Materials).
2. They must be **Cryptographically Signed**.
3. They must have **SLSA Provenance**.

We have created reusable CI templates to completely automate this process for you.

## GitHub Actions

To integrate TRUSTCHAIN into your GitHub repository, create a workflow file (e.g., `.github/workflows/build.yml`) that calls our reusable workflow:

```yaml
name: Production Build

on:
  push:
    branches: [ "main" ]

jobs:
  secure-build:
    uses: Hemanth13131313/TrustChain/.github/workflows/trustchain-secure-build.yml@main
    with:
      image_name: ghcr.io/${{ github.repository }}
      image_tag: ${{ github.sha }}
      trustchain_url: https://trustchain.internal.company.com
    secrets:
      registry_username: ${{ github.actor }}
      registry_password: ${{ secrets.GITHUB_TOKEN }}
      trustchain_token: ${{ secrets.TRUSTCHAIN_API_TOKEN }}
```

### Important Prerequisites
Your repository must grant `id-token: write` permissions so that `cosign` can securely request short-lived signing certificates from Sigstore's Fulcio identity provider. You do not need to manage long-lived GPG or RSA keys!

## GitLab CI

To integrate TRUSTCHAIN into your GitLab repository, include our template in your `.gitlab-ci.yml`:

```yaml
include:
  - project: 'trustchain-org/trustchain'
    file: '/ci-templates/gitlab/trustchain-secure-build.gitlab-ci.yml'

variables:
  TRUSTCHAIN_URL: "https://trustchain.internal.company.com"
  # Set TRUSTCHAIN_TOKEN in your GitLab CI/CD Variables
```

## How It Works
When you push code:
1. Your image is built.
2. `syft` scans the image and generates a CycloneDX SBOM.
3. `cosign` signs the image digest using an OIDC token.
4. SLSA provenance is generated.
5. All of this metadata is automatically POSTed to the TRUSTCHAIN Ingestion Service.
6. When Kubernetes attempts to pull and run your image, the TRUSTCHAIN Admission Controller will query the Policy Engine to verify these cryptographic signatures and attestations before allowing the Pod to start.
