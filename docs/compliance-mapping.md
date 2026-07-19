# TRUSTCHAIN Compliance Mapping

This document maps the technical capabilities of the TRUSTCHAIN platform to major international cybersecurity regulations and frameworks.

## 1. Executive Order (EO) 14028 (Improving the Nation's Cybersecurity)
**Focus:** Securing the Software Supply Chain

| EO 14028 Requirement | TRUSTCHAIN Implementation |
|----------------------|---------------------------|
| **Sec 4(e)(i)**: Employ automated tools to maintain trusted source code supply chains. | **Phase 9 (CI/CD)**: GitHub/GitLab templates automate signing (Cosign) and SBOM generation directly at the build step. |
| **Sec 4(e)(vi)**: Provide a purchaser a Software Bill of Materials (SBOM) for each product directly or by publishing it on a public website. | **Phase 1 (Ingestion)**: Ingests, normalizes, and stores CycloneDX and SPDX SBOMs in a centralized database. |
| **Sec 4(e)(x)**: Ensuring and attesting to the integrity and provenance of open source software used within any portion of a product. | **Phase 2 (Verification)**: Cryptographically verifies SLSA Provenance (in-toto attestations) and OCI image signatures using Sigstore Fulcio/Rekor. |

## 2. EU NIS2 Directive
**Focus:** Network and Information Systems Security

| NIS2 Article | TRUSTCHAIN Implementation |
|--------------|---------------------------|
| **Article 21(2)(d)**: Supply chain security, including security-related aspects concerning the relationships between each entity and its direct suppliers. | **Phase 4 (Admission Controller)**: Enforces zero-trust at the Kubernetes boundary. Prevents the deployment of any third-party artifact that lacks a verified signature and SBOM. |
| **Article 21(2)(e)**: Security in network and information systems acquisition, development and maintenance, including vulnerability handling and disclosure. | **Phase 3 (Policy Engine)**: Evaluates Vulnerability Exploitability eXchange (VEX) documents to automatically filter out false-positive CVEs, ensuring teams focus only on exploitable vulnerabilities. |

## 3. DORA (Digital Operational Resilience Act)
**Focus:** Financial Sector IT Resilience

| DORA Article | TRUSTCHAIN Implementation |
|--------------|---------------------------|
| **Article 9**: Protection and prevention (Monitoring anomalous activity). | **Phase 5 (Runtime Agent)**: Continuously polls the container runtime interface (CRI) to detect "drift" (unauthorized changes to running workloads). |
| **Article 11**: Response and recovery (Incident management). | **Phase 6 (Enforcement Orchestrator)**: Automatically quarantines drifted pods via NetworkPolicies and generates a tamper-evident, hash-chained audit log for post-mortem forensics. |
