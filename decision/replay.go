package decision

import "fmt"

type ReplayInput struct {
	Spec   Spec
	State  State
	Result Result
}

func (ledger Ledger) Replay(inputs []ReplayInput) error {
	if err := ledger.Validate(); err != nil {
		return err
	}
	if len(inputs) != len(ledger.Entries) {
		return fmt.Errorf("replay input count %d does not match ledger entries %d", len(inputs), len(ledger.Entries))
	}
	for index, input := range inputs {
		entry := ledger.Entries[index]
		if err := entry.Receipt.Verify(input.Spec, input.State, input.Result); err != nil {
			return fmt.Errorf("replay entry %d: %w", index, err)
		}
	}
	return nil
}
