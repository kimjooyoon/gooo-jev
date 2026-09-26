package decision

func (ledger Ledger) AppendVerified(spec Spec, state State, result Result) (Ledger, error) {
	receipt, err := Observe(spec, state, result)
	if err != nil {
		return Ledger{}, err
	}
	if err := receipt.Verify(spec, state, result); err != nil {
		return Ledger{}, err
	}
	return ledger.Append(receipt)
}
