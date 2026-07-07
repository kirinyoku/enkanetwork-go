package models

import (
	"encoding/json"
	"testing"
)

func TestStringNumberAcceptsStringAndNumber(t *testing.T) {
	var fromString StringNumber
	if err := json.Unmarshal([]byte(`"1311116816978843608"`), &fromString); err != nil {
		t.Fatalf("failed to unmarshal string: %v", err)
	}
	if fromString.String() != "1311116816978843608" {
		t.Fatalf("unexpected string value: %q", fromString.String())
	}
	if got, err := fromString.Int64(); err != nil || got != 1311116816978843608 {
		t.Fatalf("unexpected int64 value: got %d err %v", got, err)
	}

	var fromNumber StringNumber
	if err := json.Unmarshal([]byte(`12345`), &fromNumber); err != nil {
		t.Fatalf("failed to unmarshal number: %v", err)
	}
	if fromNumber.String() != "12345" {
		t.Fatalf("unexpected number value: %q", fromNumber.String())
	}

	gotJSON, err := json.Marshal(fromNumber)
	if err != nil {
		t.Fatalf("failed to marshal StringNumber: %v", err)
	}
	if string(gotJSON) != `"12345"` {
		t.Fatalf("expected JSON string, got %s", gotJSON)
	}
}

func TestStringNumberNullIsZero(t *testing.T) {
	value := StringNumber("123")
	if err := json.Unmarshal([]byte(`null`), &value); err != nil {
		t.Fatalf("failed to unmarshal null: %v", err)
	}
	if !value.IsZero() {
		t.Fatalf("expected zero value, got %q", value.String())
	}
}

func TestIntStringAcceptsStringAndNumber(t *testing.T) {
	var fromString IntString
	if err := json.Unmarshal([]byte(`"12345"`), &fromString); err != nil {
		t.Fatalf("failed to unmarshal string: %v", err)
	}
	if got, err := fromString.Int64(); err != nil || got != 12345 {
		t.Fatalf("unexpected int64 value: got %d err %v", got, err)
	}
	if fromString.String() != "12345" {
		t.Fatalf("unexpected string value: %q", fromString.String())
	}

	var fromNumber IntString
	if err := json.Unmarshal([]byte(`12345`), &fromNumber); err != nil {
		t.Fatalf("failed to unmarshal number: %v", err)
	}

	gotJSON, err := json.Marshal(fromNumber)
	if err != nil {
		t.Fatalf("failed to marshal IntString: %v", err)
	}
	if string(gotJSON) != `12345` {
		t.Fatalf("expected JSON number, got %s", gotJSON)
	}
}
