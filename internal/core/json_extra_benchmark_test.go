package core

import (
	"encoding/json"
	"testing"
)

type benchmarkJSONModel struct {
	Count int      `json:"count,omitempty"`
	Name  string   `json:"name,omitempty"`
	Items []string `json:"items,omitempty"`
}

var benchmarkMergedJSON []byte

func BenchmarkMergeKnownExtraAndRawJSON(b *testing.B) {
	model := benchmarkJSONModel{}
	raw := json.RawMessage(`{"count":0,"name":"","items":[],"future":{"enabled":true}}`)
	extra := map[string]json.RawMessage{
		"future": json.RawMessage(`{"enabled":true}`),
	}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		encoded, err := MergeKnownExtraAndRawJSON(model, extra, raw)
		if err != nil {
			b.Fatal(err)
		}
		benchmarkMergedJSON = encoded
	}
}

func BenchmarkMergeKnownExtraAndRawJSONWithoutPreservedData(b *testing.B) {
	model := benchmarkJSONModel{
		Count: 42,
		Name:  "benchmark",
		Items: []string{"one", "two"},
	}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		encoded, err := MergeKnownExtraAndRawJSON(model, nil, nil)
		if err != nil {
			b.Fatal(err)
		}
		benchmarkMergedJSON = encoded
	}
}

func BenchmarkIsZeroLikeJSONNonNumber(b *testing.B) {
	value := json.RawMessage(`"not-zero"`)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if isZeroLikeJSON(value) {
			b.Fatal("unexpected zero-like value")
		}
	}
}

func BenchmarkIsZeroLikeJSONNumber(b *testing.B) {
	value := json.RawMessage(`0.00e10`)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if !isZeroLikeJSON(value) {
			b.Fatal("expected zero-like value")
		}
	}
}
