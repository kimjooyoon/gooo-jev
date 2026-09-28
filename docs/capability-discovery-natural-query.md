# Natural Capability Discovery Queries

The LSP accepts a plain-language question together with the current `.gooo` source. The result is a bounded explanation of what the declaration can support now, what evidence is missing, and what can be inspected next.

## Request shape

Send one JSON object per line to `gooo-lsp`:

```json
{"id":"q1","source":"module usage_replan\ninput usage_observation\nstatus BOUND\nfield status\nconstraint no_execution","query":"What can gooo do with this declaration?"}
```

The `query` may be omitted when the default question is sufficient. A query is evidence-bound to the source and declaration; it is not an instruction to execute the declaration.

## Safe interpretation

- `AVAILABLE` means the catalog and declaration support a capability projection. Inspect the returned next operation; do not treat it as an authorization grant.
- `DEFERRED` means the capability is recognized but an external boundary or additional evidence must be supplied.
- `UNKNOWN` means the answer is not established. Preserve `first_mismatch` or `missing_stage` and ask one of the returned clarifying questions.
- Empty, malformed, or unrelated questions remain bounded by the same declaration and provenance digests.

## Consumer rule

An editor may render `capability_hover` and expose `capability_actions` as non-executing suggestions. A client must not invoke a provider, mutate a declaration, authorize a user, or infer completeness from a capability count. Later JEV observations should carry the returned evidence digest so the discovery answer can be replayed, confirmed, or refuted.

## Example questions

- `What can gooo do with this declaration?`
- `Which declared capability is available next?`
- `What evidence is missing before generation?`
- `What should I clarify before using this capability?`

The answer is useful only when its source, declaration, catalog, schema, evidence, and toolchain identities remain bound.
