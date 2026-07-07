package models

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestOwnerPreservesUnknownFields(t *testing.T) {
	data := []byte(`{"id":749,"username":"Algoinde","futureField":{"enabled":true}}`)

	var owner Owner
	if err := json.Unmarshal(data, &owner); err != nil {
		t.Fatalf("failed to unmarshal owner: %v", err)
	}
	if owner.Raw == nil {
		t.Fatal("expected Raw to be set")
	}
	if _, ok := owner.Extra["futureField"]; !ok {
		t.Fatal("expected futureField in Extra")
	}

	got, err := json.Marshal(owner)
	if err != nil {
		t.Fatalf("failed to marshal owner: %v", err)
	}

	assertJSONEqual(t, data, got)
}

func assertJSONEqual(t *testing.T, wantJSON, gotJSON []byte) {
	t.Helper()

	var want, got any
	if err := json.Unmarshal(wantJSON, &want); err != nil {
		t.Fatalf("failed to unmarshal expected JSON: %v", err)
	}
	if err := json.Unmarshal(gotJSON, &got); err != nil {
		t.Fatalf("failed to unmarshal actual JSON: %v", err)
	}
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("JSON mismatch\nwant: %s\ngot:  %s", wantJSON, gotJSON)
	}
}
