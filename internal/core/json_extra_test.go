package core

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestPreserveUnknownJSON(t *testing.T) {
	data := []byte(`{"known":"value","future":{"enabled":true},"count":3}`)

	raw, extra, err := PreserveUnknownJSON(data, "known")
	if err != nil {
		t.Fatalf("PreserveUnknownJSON() error = %v", err)
	}
	if string(raw) != string(data) {
		t.Fatalf("raw = %s, want %s", raw, data)
	}
	if _, ok := extra["known"]; ok {
		t.Fatal("known field should not be preserved in Extra")
	}
	assertJSONEqual(t, []byte(`{"enabled":true}`), extra["future"])
	assertJSONEqual(t, []byte(`3`), extra["count"])
}

func TestPreserveUnknownJSONWithoutExtra(t *testing.T) {
	raw, extra, err := PreserveUnknownJSON([]byte(`{"known":"value"}`), "known")
	if err != nil {
		t.Fatalf("PreserveUnknownJSON() error = %v", err)
	}
	if string(raw) != `{"known":"value"}` {
		t.Fatalf("raw = %s", raw)
	}
	if extra != nil {
		t.Fatalf("extra = %v, want nil", extra)
	}
}

func TestPreserveUnknownJSONForStruct(t *testing.T) {
	type model struct {
		Known    string `json:"known,omitempty"`
		Untagged string
		Ignored  string `json:"-"`
	}

	raw, extra, err := PreserveUnknownJSONForStruct([]byte(`{"known":"value","Untagged":"tagless","Ignored":"skip","unexported":"skip","future":true}`), model{})
	if err != nil {
		t.Fatalf("PreserveUnknownJSONForStruct() error = %v", err)
	}
	if string(raw) == "" {
		t.Fatal("expected raw JSON to be preserved")
	}
	if _, ok := extra["future"]; !ok {
		t.Fatal("expected future in Extra")
	}
	if _, ok := extra["Ignored"]; !ok {
		t.Fatal("expected ignored JSON field to remain extra")
	}
	if _, ok := extra["unexported"]; !ok {
		t.Fatal("expected unexported JSON field to remain extra")
	}
	if _, ok := extra["known"]; ok {
		t.Fatal("known field should not be preserved in Extra")
	}
	if _, ok := extra["Untagged"]; ok {
		t.Fatal("untagged exported field should not be preserved in Extra")
	}
}

func TestJSONFieldNamesReturnsIndependentSlice(t *testing.T) {
	type model struct {
		Known string `json:"known,omitempty"`
	}

	first := JSONFieldNames(model{})
	first[0] = "changed"

	second := JSONFieldNames(model{})
	if second[0] != "known" {
		t.Fatalf("JSONFieldNames() returned shared mutable data: %q", second[0])
	}
}

func TestMergeKnownAndExtraJSON(t *testing.T) {
	known := struct {
		Known string `json:"known"`
	}{
		Known: "value",
	}
	extra := map[string]json.RawMessage{
		"future": json.RawMessage(`{"enabled":true}`),
	}

	got, err := MergeKnownAndExtraJSON(known, extra)
	if err != nil {
		t.Fatalf("MergeKnownAndExtraJSON() error = %v", err)
	}

	assertJSONEqual(t, []byte(`{"known":"value","future":{"enabled":true}}`), got)
}

func TestMergeKnownAndExtraJSONKnownFieldsWin(t *testing.T) {
	known := struct {
		Known string `json:"known"`
	}{
		Known: "typed",
	}
	extra := map[string]json.RawMessage{
		"known":  json.RawMessage(`"extra"`),
		"future": json.RawMessage(`true`),
	}

	got, err := MergeKnownAndExtraJSON(known, extra)
	if err != nil {
		t.Fatalf("MergeKnownAndExtraJSON() error = %v", err)
	}

	assertJSONEqual(t, []byte(`{"known":"typed","future":true}`), got)
}

func TestMergeKnownAndExtraJSONWithoutExtra(t *testing.T) {
	known := struct {
		Known string `json:"known"`
	}{
		Known: "value",
	}

	got, err := MergeKnownAndExtraJSON(known, nil)
	if err != nil {
		t.Fatalf("MergeKnownAndExtraJSON() error = %v", err)
	}

	assertJSONEqual(t, []byte(`{"known":"value"}`), got)
}

func TestMergeKnownExtraAndRawJSONRestoresExplicitZeroLikeValues(t *testing.T) {
	known := struct {
		Count int      `json:"count,omitempty"`
		Name  string   `json:"name,omitempty"`
		Items []string `json:"items,omitempty"`
		Flag  bool     `json:"flag,omitempty"`
	}{
		Count: 0,
		Name:  "",
		Items: nil,
		Flag:  false,
	}
	raw := json.RawMessage(`{"count":0.00,"name":"","items":[],"flag":false,"future":true}`)
	extra := map[string]json.RawMessage{
		"future": json.RawMessage(`true`),
	}

	got, err := MergeKnownExtraAndRawJSON(known, extra, raw)
	if err != nil {
		t.Fatalf("MergeKnownExtraAndRawJSON() error = %v", err)
	}

	assertJSONEqual(t, []byte(`{"count":0.00,"name":"","items":[],"flag":false,"future":true}`), got)
}

func TestMergeKnownExtraAndRawJSONDoesNotRestoreNonZeroRawValues(t *testing.T) {
	known := struct {
		Count int `json:"count,omitempty"`
	}{
		Count: 0,
	}
	raw := json.RawMessage(`{"count":5}`)

	got, err := MergeKnownExtraAndRawJSON(known, nil, raw)
	if err != nil {
		t.Fatalf("MergeKnownExtraAndRawJSON() error = %v", err)
	}

	assertJSONEqual(t, []byte(`{}`), got)
}

func TestMergeKnownExtraAndRawJSONKnownFieldsWin(t *testing.T) {
	known := struct {
		Count int `json:"count,omitempty"`
	}{
		Count: 7,
	}
	raw := json.RawMessage(`{"count":0}`)

	got, err := MergeKnownExtraAndRawJSON(known, nil, raw)
	if err != nil {
		t.Fatalf("MergeKnownExtraAndRawJSON() error = %v", err)
	}

	assertJSONEqual(t, []byte(`{"count":7}`), got)
}

func TestIsZeroLikeJSON(t *testing.T) {
	for _, value := range []string{"null", "false", `""`, "[]", "{}", "0", "-0", "0.00", "0e10", "-0E-2", "0e999999"} {
		if !isZeroLikeJSON(json.RawMessage(value)) {
			t.Errorf("isZeroLikeJSON(%s) = false, want true", value)
		}
	}

	for _, value := range []string{"true", `"0"`, "1", "-1", "0.1", "1e-1000", "00", "0.", "0e", "-", `{"value":0}`, "[0]"} {
		if isZeroLikeJSON(json.RawMessage(value)) {
			t.Errorf("isZeroLikeJSON(%s) = true, want false", value)
		}
	}
}

func assertJSONEqual(t *testing.T, wantJSON, gotJSON []byte) {
	t.Helper()

	var want any
	if err := json.Unmarshal(wantJSON, &want); err != nil {
		t.Fatalf("failed to unmarshal expected JSON: %v", err)
	}

	var got any
	if err := json.Unmarshal(gotJSON, &got); err != nil {
		t.Fatalf("failed to unmarshal actual JSON: %v", err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("JSON mismatch\nwant: %s\ngot:  %s", wantJSON, gotJSON)
	}
}
