# capability discovery

The discovery command explains what a .gooo declaration can do now and
which capabilities remain UNKNOWN. It does not use a completion score as a
semantic answer and never executes or authorizes a declaration.

~~~text
go run ./cmd/gooo-discover examples/support-triage/partial.gooo pro
go run ./cmd/gooo-discover examples/support-triage/usecase.gooo pro
go run ./cmd/gooo-discover examples/support-triage/partial.gooo pro > /tmp/gooo-discovery.json
go run ./cmd/gooo-plan /tmp/gooo-discovery.json
~~~

A partial declaration can still offer syntax completion and diagnostics.
Generation, round-trip observation, reverse observation, and symbol
completion remain UNKNOWN until the declaration reaches the required stage.
The discovery JSON can be passed directly to the planner; completion
digests and diagnostics survive the JSON round trip.

The complete example uses this repository's decision declaration dialect:
`property name type` belongs to the preceding entity, and an activity can
have ordered, named inputs such as `ticket Ticket, confidence number`.
Discovery preserves those parts in its source-bound IR. The `Route` return
type is declared as an entity, so discovery can report BOUND. An unresolved
input or return type remains UNKNOWN at `ir-validation`; it is not treated as
an available generation capability. This subset is separate from the full
syntax implemented by the `meta-ontology-go` compiler.
