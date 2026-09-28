# gooo-jev
Provider-neutral typed decision receipts for Go and .gooo execution plans

## Typed decision JSON boundary

Use the provider-neutral boundary when a JEV-like provider returns a bounded
decision. The input contains a declaration of the question, a state digest,
and a typed result. The result is validated before a non-authorizing receipt
is emitted.

```sh
go run ./cmd/gooo-decision-receipt observation.json
cat observation.json | go run ./cmd/gooo-decision-receipt -
```

The boundary accepts `choice`, `score`, and `noul` observations. Unknown fields,
trailing JSON, invalid choices, and incomplete provenance are rejected. A
receipt never grants execution, changes source, or replaces the later
capability envelope and reverse-observation checks.
