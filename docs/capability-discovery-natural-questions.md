# Natural capability questions

The capability discovery interface starts with a broad question and narrows it using signals found in the declaration.

| Declaration signal | Useful next question |
| --- | --- |
| module | What module boundary can gooo inspect here? |
| input | What inputs are visible and what can be checked next? |
| output | What output contract is described? |
| status | Which status transitions are observable? |
| field | Which fields can be inspected or transformed? |
| constraint | Which constraints can be checked next? |
| observe | What observation step is available? |
| transform | What transformation can be inspected? |
| contract | Which contract evidence is still missing? |
| policy | Which policy boundary can be inspected? |
| workflow | What workflow stage is available next? |

## Recommended interaction

Start with:

    What can gooo do with this declaration?

Then select a focused question from the suggested list. The original question remains preserved in the capability trail, so a later replay can distinguish the user's intent from the declaration-driven refinement.

## What the interface does not imply

A suggested question is not a promise that the feature is implemented. The response must retain AVAILABLE, DEFERRED, or UNKNOWN, the first missing stage, and its evidence digest. Discovery remains non-executing and non-authorizing.
