package core

import (
	"strings"
	"testing"
)

func TestJSONDiffEqualDocuments(t *testing.T) {
	if diff := JSONDiff([]byte(`{"a":1,"b":[true,"x"]}`), []byte(`{"b":[true,"x"],"a":1}`)); diff != "" {
		t.Fatalf("expected no diff, got %s", diff)
	}
}

func TestJSONDiffReportsMissingFieldPath(t *testing.T) {
	diff := JSONDiff([]byte(`{"a":{"b":1}}`), []byte(`{"a":{}}`))
	if !strings.Contains(diff, "$.a.b missing field in actual") {
		t.Fatalf("unexpected diff: %s", diff)
	}
}

func TestJSONDiffReportsExtraFieldPath(t *testing.T) {
	diff := JSONDiff([]byte(`{"a":{}}`), []byte(`{"a":{"b":1}}`))
	if !strings.Contains(diff, "$.a.b extra field in actual") {
		t.Fatalf("unexpected diff: %s", diff)
	}
}

func TestJSONDiffReportsTypeMismatch(t *testing.T) {
	diff := JSONDiff([]byte(`{"a":1}`), []byte(`{"a":"1"}`))
	if !strings.Contains(diff, "$.a type mismatch") {
		t.Fatalf("unexpected diff: %s", diff)
	}
}

func TestJSONDiffReportsValueMismatch(t *testing.T) {
	diff := JSONDiff([]byte(`{"a":[1,2]}`), []byte(`{"a":[1,3]}`))
	if !strings.Contains(diff, "$.a[1] number mismatch") {
		t.Fatalf("unexpected diff: %s", diff)
	}
}

func TestJSONDiffTreatsEquivalentNumbersAsEqual(t *testing.T) {
	if diff := JSONDiff([]byte(`{"a":1}`), []byte(`{"a":1.0}`)); diff != "" {
		t.Fatalf("expected no diff, got %s", diff)
	}
}

func TestJSONDiffDistinguishesLargeIntegersExactly(t *testing.T) {
	diff := JSONDiff([]byte(`{"a":9007199254740992}`), []byte(`{"a":9007199254740993}`))
	if !strings.Contains(diff, "$.a number mismatch") {
		t.Fatalf("unexpected diff: %s", diff)
	}
}

func TestJSONDiffWithIgnoredPaths(t *testing.T) {
	ignored := map[string]struct{}{
		"$.ttl": {},
	}

	if diff := JSONDiffWithIgnoredPaths([]byte(`{"ttl":1,"a":1}`), []byte(`{"ttl":2,"a":1}`), ignored); diff != "" {
		t.Fatalf("expected ignored path, got %s", diff)
	}
}
