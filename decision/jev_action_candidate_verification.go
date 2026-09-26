package decision

import "strings"

// JEVActionCandidateVerificationInput combines a guarded candidate with the
// independent reverse observation of its expected provenance boundary.
type JEVActionCandidateVerificationInput struct {
	Guard              JEVActionCandidateGuard
	ReverseObservation ExecutionEnvelopeReverseObservationOutput
	NonAuthorizing     bool
}

// JEVActionCandidateVerification records evidence alignment without claiming
// that the candidate was executed successfully.
type JEVActionCandidateVerification struct {
	Status             string
	CandidateID        string
	GuardEvidenceDigest string
	ReverseStatus      string
	ReverseEvidenceDigest string
	FirstMismatch      string
	EvidenceDigest     string
	MissingStage       string
	NonExecuting       bool
	NonAuthorizing     bool
}

// VerifyJEVActionCandidate requires an admitted guard and a reproduced reverse
// observation before classifying the candidate boundary as verified.
func VerifyJEVActionCandidate(input JEVActionCandidateVerificationInput) JEVActionCandidateVerification {
	output := JEVActionCandidateVerification{
		Status:         "UNKNOWN",
		NonExecuting:   true,
		NonAuthorizing: true,
	}
	if !input.NonAuthorizing || !input.Guard.NonAuthorizing || !input.ReverseObservation.NonAuthorizing {
		output.NonAuthorizing = false
		output.MissingStage = "authorization-boundary"
		return output
	}
	if !input.Guard.NonExecuting {
		output.NonExecuting = false
		output.MissingStage = "execution-boundary"
		return output
	}
	if input.Guard.Status != "admitted" {
		output.Status = "review"
		output.MissingStage = "candidate-guard"
		return output
	}
	if strings.TrimSpace(input.Guard.CandidateID) == "" ||
		strings.TrimSpace(input.Guard.EvidenceDigest) == "" {
		output.MissingStage = "candidate-guard-evidence"
		return output
	}
	if strings.TrimSpace(input.ReverseObservation.EvidenceDigest) == "" {
		output.MissingStage = "reverse-evidence"
		return output
	}
	evidenceDigest, err := Digest(struct {
		CandidateID           string
		GuardEvidenceDigest   string
		ReverseStatus         string
		ReverseEvidenceDigest string
		FirstMismatch         string
	}{
		CandidateID:           input.Guard.CandidateID,
		GuardEvidenceDigest:   input.Guard.EvidenceDigest,
		ReverseStatus:         input.ReverseObservation.Status,
		ReverseEvidenceDigest: input.ReverseObservation.EvidenceDigest,
		FirstMismatch:         input.ReverseObservation.FirstMismatch,
	})
	if err != nil {
		output.MissingStage = "verification-evidence"
		return output
	}
	output.CandidateID = input.Guard.CandidateID
	output.GuardEvidenceDigest = input.Guard.EvidenceDigest
	output.ReverseStatus = input.ReverseObservation.Status
	output.ReverseEvidenceDigest = input.ReverseObservation.EvidenceDigest
	output.FirstMismatch = input.ReverseObservation.FirstMismatch
	output.EvidenceDigest = evidenceDigest
	switch input.ReverseObservation.Status {
	case "reproduced":
		if input.ReverseObservation.FirstMismatch != "" {
			output.Status = "review"
			output.MissingStage = "reverse-observation"
			return output
		}
		output.Status = "verified"
	case "counterexample":
		output.Status = "review"
		output.MissingStage = "reverse-observation"
	default:
		output.Status = "review"
		output.MissingStage = "reverse-observation"
	}
	return output
}
