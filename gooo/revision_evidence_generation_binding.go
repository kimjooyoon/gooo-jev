package gooo

import "fmt"

// RevisionEvidenceGenerationBinding links the complete evidence chain to a
// canonical GenerationReceipt without authorizing execution or persistence.
type RevisionEvidenceGenerationBinding struct {
	Status                 string
	MissingStage           string
	ChainDigest            string
	SourceDigest           string
	ProposedSourceDigest   string
	GenerationSourceDigest string
	GeneratedSourceDigest  string
	GeneratedIRDigest      string
	StructureDigest        string
	StructureMatch         bool
	ExactSourceMatch       bool
	ReverseObserved        bool
	BindingDigest          string
	NonExecuting           bool
	NonAuthorizing         bool
}

// BindRevisionEvidenceGeneration validates the existing evidence chain and
// binds its source and reverse-observation facts to canonical generation.
func BindRevisionEvidenceGeneration(chain RevisionEvidenceChain, generation GenerationReceipt) (RevisionEvidenceGenerationBinding, error) {
	binding := RevisionEvidenceGenerationBinding{
		Status:                 "UNKNOWN",
		MissingStage:           "revision-evidence-generation-binding",
		ChainDigest:            chain.ChainDigest,
		SourceDigest:           chain.SourceDigest,
		ProposedSourceDigest:   chain.ProposedSourceDigest,
		GenerationSourceDigest: generation.SourceDigest,
		GeneratedSourceDigest:  generation.GeneratedSourceDigest,
		GeneratedIRDigest:      generation.GeneratedIRDigest,
		StructureDigest:        generation.StructureDigest,
		StructureMatch:         generation.StructureMatch,
		ExactSourceMatch:       generation.ExactSourceMatch,
		ReverseObserved:        chain.ReverseObserved,
		NonExecuting:           true,
		NonAuthorizing:         true,
	}
	setBindingDigest := func() {
		binding.BindingDigest = digestRevisionEvidenceGenerationBinding(binding)
	}
	setBindingDigest()

	if err := chain.Validate(); err != nil {
		binding.MissingStage = "revision-evidence-generation-binding-chain"
		setBindingDigest()
		return binding, fmt.Errorf("revision evidence chain is not valid: %w", err)
	}
	if err := validateGenerationReceipt(generation); err != nil {
		binding.MissingStage = "revision-evidence-generation-binding-generation"
		setBindingDigest()
		return binding, fmt.Errorf("generation receipt is not valid: %w", err)
	}
	if chain.SourceDigest != generation.SourceDigest {
		binding.MissingStage = "revision-evidence-generation-binding-link"
		setBindingDigest()
		return binding, fmt.Errorf("evidence chain source digest does not match generation source digest")
	}

	binding.Status = "BOUND"
	binding.MissingStage = ""
	setBindingDigest()
	if err := binding.Validate(); err != nil {
		binding.Status = "UNKNOWN"
		binding.MissingStage = "revision-evidence-generation-binding"
		setBindingDigest()
		return binding, fmt.Errorf("revision evidence generation binding is not valid: %w", err)
	}
	return binding, nil
}

func (b RevisionEvidenceGenerationBinding) Validate() error {
	if b.Status == "" {
		return fmt.Errorf("revision evidence generation binding status is empty")
	}
	if b.Status == "BOUND" && b.MissingStage != "" {
		return fmt.Errorf("bound revision evidence generation binding has a missing stage")
	}
	if b.Status == "UNKNOWN" && b.MissingStage == "" {
		return fmt.Errorf("unknown revision evidence generation binding has no missing stage")
	}
	if !b.NonExecuting || !b.NonAuthorizing {
		return fmt.Errorf("revision evidence generation binding must remain non-executing and non-authorizing")
	}
	for name, digest := range map[string]string{
		"chain": b.ChainDigest,
		"source": b.SourceDigest,
		"proposed source": b.ProposedSourceDigest,
		"generation source": b.GenerationSourceDigest,
		"generated source": b.GeneratedSourceDigest,
		"generated IR": b.GeneratedIRDigest,
		"structure": b.StructureDigest,
		"binding": b.BindingDigest,
	} {
		if digest != "" && !validDigest(digest) {
			return fmt.Errorf("revision evidence generation binding %s digest is invalid", name)
		}
	}
	if b.Status == "BOUND" && (!b.StructureMatch || !b.ReverseObserved) {
		return fmt.Errorf("bound revision evidence generation binding lacks reverse observation")
	}
	if digestRevisionEvidenceGenerationBinding(b) != b.BindingDigest {
		return fmt.Errorf("revision evidence generation binding digest does not match its fields")
	}
	return nil
}

func digestRevisionEvidenceGenerationBinding(binding RevisionEvidenceGenerationBinding) string {
	return digestString(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%s|%s|%s|%t|%t|%t|%t|%t",
		binding.Status,
		binding.MissingStage,
		binding.ChainDigest,
		binding.SourceDigest,
		binding.ProposedSourceDigest,
		binding.GenerationSourceDigest,
		binding.GeneratedSourceDigest,
		binding.GeneratedIRDigest,
		binding.StructureDigest,
		binding.StructureMatch,
		binding.ExactSourceMatch,
		binding.ReverseObserved,
		binding.NonExecuting,
		binding.NonAuthorizing,
	))
}