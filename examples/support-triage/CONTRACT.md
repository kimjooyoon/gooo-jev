# Support triage contract projection

This contract sketch is adapted into the supported declaration pipeline without
executing a route or authorizing a change.

From the repository root, run:

    go run ./cmd/gooo-contract examples/support-triage/contract.gooo

A ready result includes the original contract digest plus declaration, IR,
generation, and reverse-observation evidence. If a capability boundary is
missing, the command keeps the result UNKNOWN and exits with status 2.
