package gooo

import "fmt"

type SemanticToken struct {
	Name     string
	Kind     SymbolKind
	Position Position
	Length   int
	Digest   string
}

type SemanticTokenResult struct {
	Status         string
	MissingStage   string
	SourceDigest   string
	IRDigest       string
	Tokens         []SemanticToken
	TokensDigest    string
	NonExecuting   bool
	NonAuthorizing bool
}

func SemanticTokens(source string) SemanticTokenResult {
	snapshot := Analyze(source)
	result := SemanticTokenResult{
		Status:         "UNKNOWN",
		MissingStage:   "lsp-semantic-tokens",
		SourceDigest:   snapshot.SourceDigest,
		IRDigest:       snapshot.IRDigest,
		NonExecuting:   true,
		NonAuthorizing: true,
	}
	if snapshot.Status != "BOUND" {
		result.MissingStage = snapshot.MissingStage
		if result.MissingStage == "" {
			result.MissingStage = "lsp-semantic-analysis"
		}
		return result
	}
	for _, symbol := range snapshot.Symbols {
		token := SemanticToken{
			Name: symbol.Name, Kind: symbol.Kind, Position: symbol.Position,
			Length: len(symbol.Name),
		}
		token.Digest = digestSemanticToken(token)
		result.Tokens = append(result.Tokens, token)
	}
	result.Status = "BOUND"
	result.MissingStage = ""
	result.TokensDigest = digestSemanticTokens(result)
	return result
}

func (s SemanticTokenResult) Validate() error {
	if s.Status != "BOUND" && s.Status != "UNKNOWN" {
		return fmt.Errorf("semantic token status %q is invalid", s.Status)
	}
	if !validDigest(s.SourceDigest) {
		return fmt.Errorf("semantic token source digest is invalid")
	}
	if !s.NonExecuting || !s.NonAuthorizing {
		return fmt.Errorf("semantic tokens must remain non-executing and non-authorizing")
	}
	if s.Status == "UNKNOWN" {
		if s.MissingStage == "" {
			return fmt.Errorf("unknown semantic tokens must retain a missing stage")
		}
		return nil
	}
	if s.MissingStage != "" || len(s.Tokens) == 0 || !validDigest(s.TokensDigest) {
		return fmt.Errorf("bound semantic tokens are incomplete")
	}
	for _, token := range s.Tokens {
		if token.Name == "" || token.Kind == "" || token.Position.Line < 1 || token.Position.Column < 1 || token.Length < 1 {
			return fmt.Errorf("semantic token is incomplete")
		}
		if !validDigest(token.Digest) || digestSemanticToken(token) != token.Digest {
			return fmt.Errorf("semantic token digest does not match its fields")
		}
	}
	if digestSemanticTokens(s) != s.TokensDigest {
		return fmt.Errorf("semantic tokens digest does not match its fields")
	}
	return nil
}

func digestSemanticToken(token SemanticToken) string {
	return digestString(fmt.Sprintf("%s|%s|%d|%d|%d", token.Name, token.Kind, token.Position.Line, token.Position.Column, token.Length))
}

func digestSemanticTokens(result SemanticTokenResult) string {
	value := fmt.Sprintf("%s|%s", result.SourceDigest, result.IRDigest)
	for _, token := range result.Tokens {
		value += "|" + token.Digest
	}
	return digestString(value)
}
