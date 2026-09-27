# Capability assistance

The usage layer can explain what the current declaration can support without executing it or authorizing anything.

```sh
go run ./cmd/gooo-assist examples/support-triage/partial.gooo pro
go run ./cmd/gooo-assist --json examples/support-triage/partial.gooo pro
```

The human-readable form lists READY and DEFERRED capabilities, their next operation, and the reason for each state. The JSON form retains `source_digest`, `discovery_digest`, `plan_digest`, and `assistance_digest` so an editor or tool can verify that the explanation belongs to the declaration it inspected.

A DEFERRED capability is not treated as a failure or a success. It remains a bounded suggestion, normally pointing to `complete_declaration`, until the declaration reaches the evidence boundary required for generation or reverse observation. The assistant is non-executing and non-authorizing.
