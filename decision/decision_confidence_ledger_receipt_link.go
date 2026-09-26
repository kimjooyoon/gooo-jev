package decision

// DecisionConfidenceEvidenceLedgerReceiptLink binds an observed ledger transition
// to an existing execution receipt by digest only. It never issues a receipt.
type DecisionConfidenceEvidenceLedgerReceiptLink struct {
	TransitionEntryDigest string `json:"transition_entry_digest"`
	ExecutionReceiptDigest string `json:"execution_receipt_digest"`
	Status                string `json:"status"`
	NonAuthorizing        bool   `json:"non_authorizing"`
}

func LinkDecisionConfidenceEvidenceLedgerTransitionToExecutionReceipt(transition DecisionConfidenceEvidenceLedgerTransition, receipt ExecutionReceipt) (DecisionConfidenceEvidenceLedgerReceiptLink, error) {
	receiptDigest, err := Digest(receipt)
	if err != nil {
		return DecisionConfidenceEvidenceLedgerReceiptLink{}, err
	}
	link := DecisionConfidenceEvidenceLedgerReceiptLink{
		TransitionEntryDigest:  transition.CurrentEntryDigest,
		ExecutionReceiptDigest: receiptDigest,
		Status:                "unknown",
		NonAuthorizing:        true,
	}
	if transition.Status == "extended" && transition.CurrentEntryDigest != "" {
		link.Status = "linked"
	}
	return link, nil
}
