CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE artifacts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    digest VARCHAR(255) NOT NULL,
    registry VARCHAR(255) NOT NULL,
    repository VARCHAR(255) NOT NULL,
    first_seen TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    tenant_id UUID,
    UNIQUE (digest, registry, repository)
);

CREATE TABLE sbom_documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    artifact_id UUID NOT NULL REFERENCES artifacts(id) ON DELETE CASCADE,
    format VARCHAR(50) NOT NULL, -- e.g., 'CycloneDX', 'SPDX'
    normalized_components JSONB,
    storage_ref VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_sbom_documents_artifact_id ON sbom_documents(artifact_id);

CREATE TABLE provenance_attestations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    artifact_id UUID NOT NULL REFERENCES artifacts(id) ON DELETE CASCADE,
    builder_id VARCHAR(255) NOT NULL,
    source_repo VARCHAR(255) NOT NULL,
    slsa_level INTEGER NOT NULL,
    verified BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_provenance_attestations_artifact_id ON provenance_attestations(artifact_id);

CREATE TABLE signatures (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    artifact_id UUID NOT NULL REFERENCES artifacts(id) ON DELETE CASCADE,
    subject VARCHAR(255),
    issuer VARCHAR(255),
    verified BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_signatures_artifact_id ON signatures(artifact_id);

CREATE TABLE workload_bindings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    artifact_id UUID NOT NULL REFERENCES artifacts(id) ON DELETE CASCADE,
    cluster VARCHAR(255) NOT NULL,
    namespace VARCHAR(255) NOT NULL,
    workload_kind VARCHAR(100) NOT NULL,
    workload_name VARCHAR(255) NOT NULL,
    bound_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_workload_bindings_artifact_id ON workload_bindings(artifact_id);

CREATE TABLE drift_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workload_binding_id UUID NOT NULL REFERENCES workload_bindings(id) ON DELETE CASCADE,
    type VARCHAR(100) NOT NULL,
    severity VARCHAR(50) NOT NULL,
    detected_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_drift_events_binding_id ON drift_events(workload_binding_id);

-- Hash-chained audit log for tamper-evidence
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_type VARCHAR(100) NOT NULL,
    actor VARCHAR(255) NOT NULL,
    action VARCHAR(255) NOT NULL,
    resource_id UUID,
    details JSONB,
    previous_hash VARCHAR(64),
    hash VARCHAR(64) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
