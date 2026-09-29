package gooo

import "encoding/json"

func (item SyntaxCompletionItem) MarshalJSON() ([]byte, error) {
	type itemJSON struct {
		Label      string `json:"label"`
		Kind       string `json:"kind"`
		Detail     string `json:"detail"`
		ItemDigest string `json:"item_digest"`
	}
	return json.Marshal(itemJSON{
		Label: item.Label, Kind: item.Kind, Detail: item.Detail, ItemDigest: item.ItemDigest,
	})
}

func (item *SyntaxCompletionItem) UnmarshalJSON(data []byte) error {
	type itemJSON struct {
		Label      string `json:"label"`
		Kind       string `json:"kind"`
		Detail     string `json:"detail"`
		ItemDigest string `json:"item_digest"`
	}
	var decoded itemJSON
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}
	item.Label = decoded.Label
	item.Kind = decoded.Kind
	item.Detail = decoded.Detail
	item.ItemDigest = decoded.ItemDigest
	return nil
}

func (response SyntaxCompletionResponse) MarshalJSON() ([]byte, error) {
	type responseJSON struct {
		Status         string                 `json:"status"`
		MissingStage   string                 `json:"missing_stage"`
		SourceDigest   string                 `json:"source_digest"`
		IRDigest       string                 `json:"ir_digest"`
		Prefix         string                 `json:"prefix"`
		Items          []SyntaxCompletionItem `json:"items"`
		ItemsDigest    string                 `json:"items_digest"`
		Diagnostics    []Diagnostic           `json:"diagnostics"`
		NonExecuting   bool                   `json:"non_executing"`
		NonAuthorizing bool                   `json:"non_authorizing"`
	}
	return json.Marshal(responseJSON{
		Status: response.Status, MissingStage: response.MissingStage,
		SourceDigest: response.SourceDigest, IRDigest: response.IRDigest,
		Prefix: response.Prefix, Items: response.Items, ItemsDigest: response.ItemsDigest,
		Diagnostics: response.Diagnostics, NonExecuting: response.NonExecuting,
		NonAuthorizing: response.NonAuthorizing,
	})
}

func (response *SyntaxCompletionResponse) UnmarshalJSON(data []byte) error {
	type responseJSON struct {
		Status         string                 `json:"status"`
		MissingStage   string                 `json:"missing_stage"`
		SourceDigest   string                 `json:"source_digest"`
		IRDigest       string                 `json:"ir_digest"`
		Prefix         string                 `json:"prefix"`
		Items          []SyntaxCompletionItem `json:"items"`
		ItemsDigest    string                 `json:"items_digest"`
		Diagnostics    []Diagnostic           `json:"diagnostics"`
		NonExecuting   bool                   `json:"non_executing"`
		NonAuthorizing bool                   `json:"non_authorizing"`
	}
	var decoded responseJSON
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}
	response.Status = decoded.Status
	response.MissingStage = decoded.MissingStage
	response.SourceDigest = decoded.SourceDigest
	response.IRDigest = decoded.IRDigest
	response.Prefix = decoded.Prefix
	response.Items = decoded.Items
	response.ItemsDigest = decoded.ItemsDigest
	response.Diagnostics = decoded.Diagnostics
	response.NonExecuting = decoded.NonExecuting
	response.NonAuthorizing = decoded.NonAuthorizing
	return nil
}