# Support triage declaration

This is a small, executable use-case for the supported .gooo declaration
subset. It describes a support ticket and a routing activity without executing
the activity or authorizing a route.

From the repository root, run:

    go run ./cmd/gooo-usecase examples/support-triage/usecase.gooo

The command emits a JSON evidence report containing the source, IR, generated
declaration, and reverse-observation digests. The status is ready only when
the generated declaration reparses to the same IR. An unresolved declaration
keeps the status as UNKNOWN and exits with status 2.
