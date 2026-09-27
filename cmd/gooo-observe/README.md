# gooo-observe

Record a provenance-bound usage observation without executing or authorizing a planned action.

Usage:
go run ./cmd/gooo-observe PLAN ACTION_ID KIND METRIC_NAME METRIC_VALUE EVIDENCE_DIGEST RECORDED_AT_RFC3339

KIND is one of confirmed, refuted, or unknown. A deferred action accepts only unknown.
The output keeps source, discovery, plan, action, and evidence digests and marks the observation as non-executing and non-authorizing.
