# .gooo completion workflow

The completion command is intentionally useful before a declaration is
complete. It returns keyword candidates together with source and IR digests,
diagnostics, and explicit non-executing/non-authorizing flags.

From the repository root:

~~~text
go run ./cmd/gooo-complete examples/support-triage/partial.gooo pro
~~~

The example asks for the pro prefix after an unfinished property line.
The JSON response contains the property candidate and preserves the
provenance needed by an editor adapter. A partial document may remain
UNKNOWN; that is an honest completion state, not a failed completion request.
