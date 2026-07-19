# TRUSTCHAIN Observability Runbook

## Key Service Level Objectives (SLOs)

### 1. Admission Controller Latency
**Target:** p95 < 250ms
**Description:** The Kubernetes API server waits for our webhook to approve or deny a pod. If we exceed the timeout, deployments will hang.
**Mitigation:** 
- Check `trustchain-policy-engine` CPU usage.
- Ensure Redis caching is active for policy decisions.
- If Redis is down, the Admission Controller is configured to `fail-open` (allow deployment, flag async) based on `ADMISSION_FAIL_MODE` settings to preserve cluster availability.

### 2. Runtime Agent Drift Detection SLA
**Target:** Drift Detected within 120 seconds.
**Description:** The maximum time between a malicious container starting and the Orchestrator quarantining it.
**Mitigation:**
- Check Kafka consumer lag on `trustchain-correlator-group`.
- Scale up `services/correlator` if lag exceeds 500 messages.

### 3. Ingestion Success Rate
**Target:** 99.9% success on `/api/v1/ingest/*`
**Description:** CI pipelines rely on this endpoint. If it fails, developer builds will fail.
**Mitigation:**
- Check PostgreSQL connection pool exhaustion. Max connections should be tuned according to the number of concurrent CI builds.

## Tracing
All services are instrumented with OpenTelemetry.
- Trace ID propagation: Headers `X-B3-TraceId` and `traceparent` are passed from the API Gateway down to the Policy Engine.
- View traces in Jaeger/Zipkin.
