package decision

import (
	"strings"
	"testing"
)

func declarationEvidenceCandidateGateCandidate() JEVImprovementRevisionCandidate {
	candidate := JEVImprovementRevisionCandidate{
		Status:                jevImprovementRevisionCandidateReady,
		ParentCandidateDigest: strings.Repeat("a", 64),
		RevisionSource:        "revision-source",
		RevisionChangeDigest:  strings.Repeat("b", 64),
		DirectiveEvidenceDigest: strings.Repeat("c", 64),
		NonExecuting:          true,
		NonAuthorizing:        true,
	}
	candidate.CandidateDigest = digestJEVImprovementRevisionCandidate(
		candidate.ParentCandidateDigest,
		candidate.RevisionSource,
		candidate.RevisionChangeDigest,
		candidate.DirectiveEvidenceDigest,
	)
	candidate.EvidenceDigest = digestJEVImprovementRevisionCandidateEvidence(
		candidate.Status,
		candidate.CandidateDigest,
		candidate.DirectiveEvidenceDigest,
		candidate.RevisionChangeDigest,
	)
	return candidate
}

func declarationEvidenceCandidateGateInput() JEVDeclarationEvidenceCandidateGateInput {
	return JEVDeclarationEvidenceCandidateGateInput{
		MetricStatus:            "BOUND",
		MetricStageCount:        3,
		MetricStageTotal:        3,
		MetricCoverage:          1,
		MetricEvidenceDigest:    strings.Repeat("d", 64),
		MetricDigest:            strings.Repeat("e", 64),
		MetricNonExecuting:      true,
		MetricNonAuthorizing:    true,
		Candidate:               declarationEvidenceCandidateGateCandidate(),
		NonAuthorizing:          true,
	}
}

func TestBindJEVDeclarationEvidenceToCandidateBindsCompleteCycle(t *testing.T) {
	binding := BindJEVDeclarationEvidenceToCandidate(declarationEvidenceCandidateGateInput())
	if binding.Status != jevDeclarationEvidenceCandidateGateBound || binding.Decision != jevDeclarationEvidenceCandidateGateReview {
		t.Fatalf("expected review binding, got %#v", binding)
	}
	if binding.MissingStage != "" || binding.CandidateSourceDigest == "" {
		t.Fatalf("expected complete binding, got %#v", binding)
	}
	if err := binding.Validate(); err != nil {
		t.Fatalf("binding should validate: %v", err)
	}
}

func TestBindJEVDeclarationEvidenceToCandidatePreservesDeferredStage(t *testing.T) {
	input := declarationEvidenceCandidateGateInput()
	input.MetricStatus = "DEFERRED"
	input.MetricStageCount = 2
	input.MetricCoverage = 2.0 / 3.0
	input.MetricFirstMissingStage = "reverse_observation"

	binding := BindJEVDeclarationEvidenceToCandidate(input)
	if binding.Status != jevDeclarationEvidenceCandidateGateUnknown || binding.MissingStage != "reverse_observation" {
		t.Fatalf("expected deferred stage to remain unknown, got %#v", binding)
	}
	if binding.CandidateSourceDigest != "" {
		t.Fatalf("deferred evidence must not admit a candidate: %#v", binding)
	}
}

func TestBindJEVDeclarationEvidenceToCandidateRejectsBoundaryAndTampering(t *testing.T) {
	input := declarationEvidenceCandidateGateInput()
	input.NonAuthorizing = false
	binding := BindJEVDeclarationEvidenceToCandidate(input)
	if binding.NonAuthorizing || binding.MissingStage != "authorization-boundary" {
		t.Fatalf("expected authorization boundary rejection, got %#v", binding)
	}

	input = declarationEvidenceCandidateGateInput()
	input.MetricCoverage = 0.9
	binding = BindJEVDeclarationEvidenceToCandidate(input)
	if binding.Status != jevDeclarationEvidenceCandidateGateUnknown || binding.MissingStage != "declaration-evidence-cycle-completeness" {
		t.Fatalf("expected incomplete metric rejection, got %#v", binding)
	}

	input = declarationEvidenceCandidateGateInput()
	binding = BindJEVDeclarationEvidenceToCandidate(input)
	binding.MetricDigest = strings.Repeat("f", 64)
	if err := binding.Validate(); err == nil {
		t.Fatal("expected tampered binding digest to be rejected")
	}
}
