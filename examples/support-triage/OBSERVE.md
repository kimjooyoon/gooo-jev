# Usage observation

The discovery slice first reports what a partial declaration can do, then plans the next non-executing action.
A READY action may be confirmed or refuted. A DEFERRED action remains UNKNOWN until its declaration boundary is complete.

Flow:
1. gooo-discover produces capability evidence.
2. gooo-plan produces provenance-bound next actions.
3. gooo-observe records a bounded observation without executing code.
