# TRUSTCHAIN Disaster Recovery Runbook

## Overview
This runbook defines the procedures for recovering the TRUSTCHAIN control plane from a catastrophic failure.

## 1. PostgreSQL Database Recovery
The PostgreSQL database is the most critical stateful component, holding the `artifacts`, `signatures`, and the tamper-evident `audit_logs`.

### 1.1 Backup Procedure
Backups should be executed daily via a Kubernetes CronJob.
```bash
pg_dump -h $DB_HOST -U trustchain -Fc trustchain > /backup/trustchain_$(date +%Y%m%d).dump
```

### 1.2 Restore Procedure
> [!CAUTION]
> Restoring the database will overwrite current state. Ensure all microservices (Ingestion, Verification) are scaled to zero before proceeding.

```bash
# 1. Scale down writers
kubectl scale deployment trustchain-ingestion --replicas=0
kubectl scale deployment trustchain-orchestrator --replicas=0

# 2. Restore database
pg_restore -h $DB_HOST -U trustchain -d trustchain -1 /backup/trustchain_latest.dump

# 3. Verify Audit Log Integrity
# You must manually verify the hash chain of the latest record to ensure it matches the final state of the backup.

# 4. Scale up writers
kubectl scale deployment trustchain-ingestion --replicas=3
kubectl scale deployment trustchain-orchestrator --replicas=3
```

## 2. Kafka Recovery
Kafka holds ephemeral observation and incident streams.
- If a Kafka broker goes down, rely on the replication factor (RF=3, MinISR=2).
- If the entire Kafka cluster is lost, no historical restore is required. Recreate the topics (`trustchain.runtime.observations`, `trustchain.incidents`). The Runtime Agent will automatically publish fresh observations on its next polling cycle.

## 3. High Availability (HA) Considerations
To minimize the need for DR, ensure the following deployments are running in HA:
- `trustchain-admission-controller`: Min Replicas 3 (Critical for cluster operations).
- `trustchain-policy-engine`: Min Replicas 3 (Cache enabled via Redis).
