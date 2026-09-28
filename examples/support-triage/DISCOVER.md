# capability discovery

The discovery command explains what a .gooo declaration can do now and
which capabilities remain UNKNOWN. It does not use a completion score as a
semantic answer and never executes or authorizes a declaration.

~~~text
go run ./cmd/gooo-discover examples/support-triage/partial.gooo pro
~~~

A partial declaration can still offer syntax completion and diagnostics.
Generation, round-trip observation, reverse observation, and symbol
completion remain UNKNOWN until the declaration reaches the required stage.