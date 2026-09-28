# .gooo generation workflow

The generation command writes canonical .gooo only after the declaration
has valid source, IR, generation, round-trip, and reverse-observation evidence.
It never executes or authorizes the declaration.

~~~text
go run ./cmd/gooo-generate examples/support-triage/usecase.gooo /tmp/support-triage.generated.gooo
~~~

An incomplete declaration remains UNKNOWN and does not produce an output
file. The command prints the evidence digests on stderr so a build or editor
adapter can bind the generated file to its exact provenance.