package decision

import (
	"strings"
	"time"
)

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionGuardInput
// carries a revision candidate into the existing typed action guard.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionGuardInput struct {
	RevisionCandidate   ExecutionEnvelopeGoooEvidenceFullProvenanceDirectionRevisionCandidateBinding
	CandidateID         string
	ObservationDigest   string
	ObservationAt       time.Time
	Now                 time.Time
	MaxAge              time.Duration
	RequestedPath       string
	AllowedPathPrefix   string
	WriteRequested      bool
	WriteAllowed        bool
	ConfirmationRequired bool
	ConfirmationPresent  bool
	NonAuthorizing      bool
}

// ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionGuardBinding
// keeps revision provenance attached to the guard disposition.
type ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionGuardBinding struct {
	Status                    string
	MissingStage              string
	CandidateID               string
	RevisionCandidateDigest   string
	RevisionCandidateEvidenceDigest string
	GuardStatus               string
	GuardEvidenceDigest       string
	EvidenceDigest            string
	NonExecuting              bool
	NonAuthorizing            bool
}

func (b ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionGuardBinding) Validate() error {
	if b.Status != "admitted" && b.Status != "review" && b.Status != "rejected" {
		return errIncompleteGoooRevisionActionGuardBinding()
	}
	if b.MissingStage == "" && b.GuardStatus != b.Status {
		return errIncompleteGoooRevisionActionGuardBinding()
	}
	if b.CandidateID == "" ||
		b.RevisionCandidateDigest == "" ||
		b.RevisionCandidateEvidenceDigest == "" ||
		b.GuardStatus == "" ||
		b.GuardEvidenceDigest == "" ||
		b.EvidenceDigest == "" {
		return errIncompleteGoooRevisionActionGuardBinding()
	}
	if !b.NonExecuting || !b.NonAuthorizing {
		return errIncompleteGoooRevisionActionGuardBinding()
	}
	expected := digestJEVGoooEvidenceFullProvenanceRevisionActionGuard(
		b.Status,
		b.CandidateID,
		b.RevisionCandidateDigest,
		b.RevisionCandidateEvidenceDigest,
		b.GuardStatus,
		b.GuardEvidenceDigest,
	)
	if b.EvidenceDigest != expected {
		return errIncompleteGoooRevisionActionGuardBinding()
	}
	return nil
}

func errIncompleteGoooRevisionActionGuardBinding() error {
	return fmt.Errorf("incomplete Gooo revision action guard binding")
}

// GuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionAction
// validates candidate provenance before exposing the typed guard result.
func GuardExecutionEnvelopeGoooEvidenceFullProvenanceRevisionAction(input ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionGuardInput) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionGuardBinding {
	unknown := func(stage string) ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionGuardBinding {
		if strings.TrimSpace(stage) == "" {
			stage = "revision-action-guard"
		}
		return ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionGuardBinding{
			Status: "UNKNOWN", MissingStage: stage,
			NonExecuting: true, NonAuthorizing: true,
		}
	}
	if !input.NonAuthorizing || !input.RevisionCandidate.NonAuthorizing {
		return unknown("authorization-boundary")
	}
	if !input.RevisionCandidate.NonExecuting {
		return unknown("execution-boundary")
	}
	if err := input.RevisionCandidate.Validate(); err != nil {
		return unknown("revision-candidate-validation")
	}
	guard := GuardJEVActionCandidate(JEVActionCandidateGuardInput{
		CandidateID:          input.CandidateID,
		ObservationDigest:    input.ObservationDigest,
		ObservationAt:        input.ObservationAt,
		Now:                  input.Now,
		MaxAge:               input.MaxAge,
		RequestedPath:        input.RequestedPath,
		AllowedPathPrefix:    input.AllowedPathPrefix,
		WriteRequested:       input.WriteRequested,
		WriteAllowed:         input.WriteAllowed,
		ConfirmationRequired: input.ConfirmationRequired,
		ConfirmationPresent:  input.ConfirmationPresent,
		NonAuthorizing:       true,
	})
	if guard.EvidenceDigest == "" {
		return unknown(guard.MissingStage)
	}
	output := ExecutionEnvelopeGoooEvidenceFullProvenanceRevisionActionGuardBinding{
		Status:                          guard.Status,
		MissingStage:                    guard.MissingStage,
		CandidateID:                     guard.CandidateID,
		RevisionCandidateDigest:         input.RevisionCandidate.RevisionCandidateDigest,
		RevisionCandidateEvidenceDigest: input.RevisionCandidate.RevisionCandidateEvidenceDigest,
		GuardStatus:                     guard.Status,
		GuardEvidenceDigest:             guard.EvidenceDigest,
		NonExecuting:                    true,
		NonAuthorizing:                  true,
	}
	output.EvidenceDigest = digestJEVGoooEvidenceFullProvenanceRevisionActionGuard(
		output.Status,
		output.CandidateID,
		output.RevisionCandidateDigest,
		output.RevisionCandidateEvidenceDigest,
		output.GuardStatus,
		output.GuardEvidenceDigest,
	)
	if err := output.Validate(); err != nil {
		return unknown("revision-action-guard-evidence")
	}
	return output
}

func digestJEVGoooEvidenceFullProvenanceRevisionActionGuard(
	status,
	candidateID,
	revisionCandidateDigest,
	revisionCandidateEvidenceDigest,
	guardStatus,
	guardEvidenceDigest string,
) string {
	digest, err := Digest(struct {
		Status                       string
		CandidateID                  string
		RevisionCandidateDigest      string
		RevisionCandidateEvidenceDigest string
		GuardStatus                  string
		GuardEvidenceDigest          string
	}{
		Status:                          status,
		CandidateID:                     candidateID,
		RevisionCandidateDigest:         revisionCandidateDigest,
		RevisionCandidateEvidenceDigest: revisionCandidateEvidenceDigest,
		GuardStatus:                     guardStatus,
		GuardEvidenceDigest:             guardEvidenceDigest,
	})
	if err != nil {
		return ""
	}
	return digest
}