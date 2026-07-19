# TRUSTCHAIN

Continuous software supply chain verification and enforcement platform.

## Local Development

### Prerequisites
- Docker & Docker Compose
- Go 1.21+
- Node.js 18+ (for frontend)
- `golang-migrate` (for DB migrations)
- OPA (for policy testing)

### Starting Local Infrastructure

Spin up the local dependencies (PostgreSQL, Redis, Kafka, ClickHouse):
```bash
cd trustchain
docker-compose up -d
```

### Running Migrations

Apply the initial database schema:
```bash
migrate -path migrations -database "postgres://trustchain:trustchain_password@localhost:5432/trustchain?sslmode=disable" up
```

### Services

- `ingestion`: SBOM/Provenance/VEX ingestion
- `verification`: Signature & SLSA provenance verification
- `policy-engine`: OPA-based Rego policy evaluation
- `admission-controller`: K8s Webhook
- `runtime-agent`: K8s DaemonSet for re-attestation
- `correlator`: Maps running containers to artifacts, detects drift
- `enforcement-orchestrator`: Applies quarantine/labeling
- `api-gateway`: REST API entrypoint
- `web/dashboard`: React/TypeScript frontend

### Policies

Rego policies are stored in the `policies/` directory. Run tests with:
```bash
opa test policies/
```
