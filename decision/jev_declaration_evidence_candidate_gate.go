package decision

import "fmt"

const (
	jevDeclarationEvidenceCandidateGateBound   = "BOUND"
	jevDeclarationEvidenceCandidateGateUnknown = "UNKNOWN"
	jevDeclarationEvidenceCandidateGateReview  = "REVIEW"
)

// JEVDeclarationEvidenceCandidateGateInput is the declaration-facing
// boundary between a complete evidence cycle and an improvement candidate.
// The metric is flattened so this package does not execute or import the
// declaration pipeline.
type JEVDeclarationEvidenceCandidateGateInput struct {
	MetricStatus            string
	MetricStageCount        int
	MetricStageTotal        int
	MetricCoverage          float64
	MetricFirstMissingStage string
	MetricEvidenceDigest    string
	MetricDigest            string
	MetricNonExecuting      bool
	MetricNonAuthorizing    bool
	Candidate               JEVImprovementRevisionCandidate
	NonAuthorizing          bool
}

// JEVDeclarationEvidenceCandidateGateBinding binds a complete declaration
// evidence cycle to a review-only candidate. It never applies or authorizes a
// revision.
type JEVDeclarationEvidenceCandidateGateBinding struct {
	Status                  string
	Decision                string
	MetricStatus            string
	MetricStageCount        int
	MetricStageTotal        int
	MetricCoverage          float64
	MetricFirstMissingStage string
	MetricEvidenceDigest    string
	MetricDigest            string
	CandidateDigest         string
	CandidateEvidenceDigest string
	CandidateSourceDigest   string
	MissingStage            string
	BindingDigest           string
	NonExecuting            bool
	NonAuthorizing          bool
}

// BindJEVDeclarationEvidenceToCandidate admits a candidate for review only
// when declaration, IR/generation, and reverse observation are all bound.
// Deferred, unknown, invalid, or tampered evidence remains UNKNOWN.
func BindJEVDeclarationEvidenceToCandidate(input JEVDeclarationEvidenceCandidateGateInput) JEVDeclarationEvidenceCandidateGateBinding {
	output := JEVDeclarationEvidenceCandidateGateBinding{
		Status:                  jevDeclarationEvidenceCandidateGateUnknown,
		Decision:                jevDeclarationEvidenceCandidateGateReview,
		MetricStatus:            input.MetricStatus,
		MetricStageCount:        input.MetricStageCount,
		MetricStageTotal:        input.MetricStageTotal,
		MetricCoverage:          input.MetricCoverage,
		MetricFirstMissingStage: input.MetricFirstMissingStage,
		MetricEvidenceDigest:    input.MetricEvidenceDigest,
		MetricDigest:            input.MetricDigest,
		CandidateDigest:         input.Candidate.CandidateDigest,
		CandidateEvidenceDigest: input.Candidate.EvidenceDigest,
		NonExecuting:            true,
		NonAuthorizing:          true,
	}

	finalize := func() {
		output.BindingDigest = digestJEVDeclarationEvidenceCandidateGate(output)
	}
	finalize()

	if !input.NonAuthorizing || !input.MetricNonAuthorizing || !input.Candidate.NonAuthorizing {
		if !input.NonAuthorizing {
			output.NonAuthorizing = false
		}
		output.MissingStage = "authorization-boundary"
		finalize()
		return output
	}
	if !input.MetricNonExecuting || !input.Candidate.NonExecuting {
		output.MissingStage = "execution-boundary"
		finalize()
		return output
	}
	if !validDigest(input.MetricEvidenceDigest) || !validDigest(input.MetricDigest) {
		output.MissingStage = "declaration-evidence-cycle-digest"
		finalize()
		return output
	}
	if input.MetricStatus != "BOUND" {
		output.MissingStage = input.MetricFirstMissingStage
		if output.MissingStage == "" {
			output.MissingStage = "declaration-evidence-cycle"
		}
		finalize()
		return output
	}
	if input.MetricStageTotal <= 0 ||
		input.MetricStageCount != input.MetricStageTotal ||
		input.MetricCoverage != 1 ||
		input.MetricFirstMissingStage != "" {
		output.MissingStage = "declaration-evidence-cycle-completeness"
		finalize()
		return output
	}
	if err := input.Candidate.Validate(); err != nil {
		output.MissingStage = "improvement-candidate"
		finalize()
		return output
	}
	candidateSourceDigest, err := Digest(struct {
		RevisionSource string
	}{RevisionSource: input.Candidate.RevisionSource})
	if err != nil {
		output.MissingStage = "candidate-source-digest"
		finalize()
		return output
	}

	output.Status = jevDeclarationEvidenceCandidateGateBound
	output.MissingStage = ""
	output.CandidateSourceDigest = candidateSourceDigest
	finalize()
	if err := output.Validate(); err != nil {
		output.Status = jevDeclarationEvidenceCandidateGateUnknown
		output.MissingStage = "declaration-evidence-candidate-gate"
		finalize()
	}
	return output
}

func (b JEVDeclarationEvidenceCandidateGateBinding) Validate() error {
	if b.Status != jevDeclarationEvidenceCandidateGateBound && b.Status != jevDeclarationEvidenceCandidateGateUnknown {
		return fmt.Errorf("invalid declaration evidence candidate gate status %q", b.Status)
	}
	if b.Decision != jevDeclarationEvidenceCandidateGateReview {
		return fmt.Errorf("declaration evidence candidate gate must remain review-only")
	}
	if b.Status == jevDeclarationEvidenceCandidateGateBound && b.MissingStage != "" {
		return fmt.Errorf("bound declaration evidence candidate gate has a missing stage")
	}
	if b.Status == jevDeclarationEvidenceCandidateGateUnknown && b.MissingStage == "" {
		return fmt.Errorf("unknown declaration evidence candidate gate has no missing stage")
	}
	if b.MetricStageCount < 0 || b.MetricStageTotal < 0 || b.MetricStageCount > b.MetricStageTotal {
		return fmt.Errorf("declaration evidence candidate gate metric stage counts are invalid")
	}
	if b.MetricCoverage < 0 || b.MetricCoverage > 1 {
		return fmt.Errorf("declaration evidence candidate gate metric coverage is invalid")
	}
	if !validDigest(b.BindingDigest) {
		return fmt.Errorf("declaration evidence candidate gate binding digest is invalid")
	}
	if b.Status == jevDeclarationEvidenceCandidateGateBound {
		if b.MetricStatus != "BOUND" || b.MetricStageCount != b.MetricStageTotal || b.MetricCoverage != 1 || b.MetricFirstMissingStage != "" {
			return fmt.Errorf("bound declaration evidence candidate gate has incomplete metric evidence")
		}
		for name, digest := range map[string]string{
			"metric evidence":    b.MetricEvidenceDigest,
			"metric":             b.MetricDigest,
			"candidate":          b.CandidateDigest,
			"candidate evidence": b.CandidateEvidenceDigest,
			"candidate source":   b.CandidateSourceDigest,
		} {
			if !validDigest(digest) {
				return fmt.Errorf("bound declaration evidence candidate gate %s digest is invalid", name)
			}
		}
	}
	if !b.NonExecuting || !b.NonAuthorizing {
		return fmt.Errorf("declaration evidence candidate gate crossed a capability boundary")
	}
	if digestJEVDeclarationEvidenceCandidateGate(b) != b.BindingDigest {
		return fmt.Errorf("declaration evidence candidate gate binding digest mismatch")
	}
	return nil
}

func digestJEVDeclarationEvidenceCandidateGate(binding JEVDeclarationEvidenceCandidateGateBinding) string {
	digest, err := Digest(struct {
		Status                  string
		Decision                string
		MetricStatus            string
		MetricStageCount        int
		MetricStageTotal        int
		MetricCoverage          float64
		MetricFirstMissingStage string
		MetricEvidenceDigest    string
		MetricDigest            string
		CandidateDigest         string
		CandidateEvidenceDigest string
		CandidateSourceDigest   string
		MissingStage            string
		NonExecuting            bool
		NonAuthorizing          bool
	}{
		Status:                  binding.Status,
		Decision:                binding.Decision,
		MetricStatus:            binding.MetricStatus,
		MetricStageCount:        binding.MetricStageCount,
		MetricStageTotal:        binding.MetricStageTotal,
		MetricCoverage:          binding.MetricCoverage,
		MetricFirstMissingStage: binding.MetricFirstMissingStage,
		MetricEvidenceDigest:    binding.MetricEvidenceDigest,
		MetricDigest:            binding.MetricDigest,
		CandidateDigest:         binding.CandidateDigest,
		CandidateEvidenceDigest: binding.CandidateEvidenceDigest,
		CandidateSourceDigest:   binding.CandidateSourceDigest,
		MissingStage:            binding.MissingStage,
		NonExecuting:            binding.NonExecuting,
		NonAuthorizing:          binding.NonAuthorizing,
	})
	if err != nil {
		return ""
	}
	return digest
}
