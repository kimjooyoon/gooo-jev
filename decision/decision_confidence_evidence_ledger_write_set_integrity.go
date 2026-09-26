package decision

// DecisionConfidenceEvidenceLedgerWriteSetIntegrity compares an observed write-set
// identity with an independently supplied expectation before any application.
type DecisionConfidenceEvidenceLedgerWriteSetIntegrity struct {
	ObservedWriteSetDigest string `json:"observed_write_set_digest"`
	ExpectedWriteSetDigest string `json:"expected_write_set_digest"`
	Status                 string `json:"status"`
	NonAuthorizing         bool   `json:"non_authorizing"`
}

func VerifyDecisionConfidenceEvidenceLedgerWriteSetIntegrity(observation DecisionConfidenceEvidenceLedgerWriteSetObservation, expectedWriteSetDigest string) DecisionConfidenceEvidenceLedgerWriteSetIntegrity {
	integrity := DecisionConfidenceEvidenceLedgerWriteSetIntegrity{
		ObservedWriteSetDigest: observation.WriteSetDigest,
		ExpectedWriteSetDigest: expectedWriteSetDigest,
		Status:                 "unknown",
		NonAuthorizing:         true,
	}
	if observation.Status != "write-set-observed" || observation.WriteSetDigest == "" || expectedWriteSetDigest == "" {
		return integrity
	}
	if observation.WriteSetDigest != expectedWriteSetDigest {
		integrity.Status = "mismatch"
		return integrity
	}
	integrity.Status = "verified"
	return integrity
}
