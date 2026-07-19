# TRUSTCHAIN Security Audit Checklist

This document is intended for internal security teams or third-party penetration testers evaluating the TRUSTCHAIN deployment.

## 1. Admission Controller Webhook Security
- [ ] **TLS Configuration**: Is the ValidatingWebhookConfiguration strictly enforcing TLS 1.3?
- [ ] **mTLS (Mutual TLS)**: Does the `services/admission-controller` require a client certificate from the Kubernetes API Server?
- [ ] **Fail-Mode Override**: If `ADMISSION_FAIL_MODE=open` is triggered, is the event correctly logged with a high severity alert?

## 2. API Gateway (BFF) Authentication
- [ ] **JWT Validation**: Ensure that the `/api/v1/auth/login` endpoint issues short-lived JWTs (e.g., 15 minutes).
- [ ] **RBAC Enforcement**: Verify that only users with the `Admin` role can invoke the `/api/v1/policies` endpoints.
- [ ] **Rate Limiting**: Attempt to flood the `/api/v1/ingest/sbom` endpoint to confirm rate limiters return HTTP 429.

## 3. Cryptographic Integrity (Cosign / Sigstore)
- [ ] **Root of Trust**: Verify that the Policy Engine is configured with the correct Sigstore Fulcio Root CA.
- [ ] **Tampered Signatures**: Inject a modified Cosign payload into the database and verify that the Verification Service correctly rejects it upon re-evaluation.

## 4. Tamper-Evident Audit Log (Enforcement Orchestrator)
- [ ] **Hash Chain Breakage**: Manually execute an `UPDATE` statement on the `audit_logs` table via direct database access. Run the verification script to confirm that the hash chain breakage is immediately detected.
- [ ] **PostgreSQL RBAC**: Ensure the microservice database user (`trustchain_user`) does NOT have `UPDATE` or `DELETE` privileges on the `audit_logs` table (Append-Only design).
