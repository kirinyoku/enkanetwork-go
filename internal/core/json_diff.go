package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
)

// JSONDiff compares two JSON documents semantically and returns the first
// mismatch. An empty string means the documents are equivalent.
func JSONDiff(expected, actual []byte) string {
	return JSONDiffWithIgnoredPaths(expected, actual, nil)
}

// JSONDiffWithIgnoredPaths compares two JSON documents semantically, skipping
// exact JSON paths listed in ignoredPaths.
func JSONDiffWithIgnoredPaths(expected, actual []byte, ignoredPaths map[string]struct{}) string {
	var expectedValue any
	if err := decodeJSON(expected, &expectedValue); err != nil {
		return fmt.Sprintf("failed to decode expected JSON: %v", err)
	}

	var actualValue any
	if err := decodeJSON(actual, &actualValue); err != nil {
		return fmt.Sprintf("failed to decode actual JSON: %v", err)
	}

	return diffJSON("$", expectedValue, actualValue, ignoredPaths)
}

func decodeJSON(data []byte, value *any) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	return decoder.Decode(value)
}

func diffJSON(path string, expected, actual any, ignoredPaths map[string]struct{}) string {
	if _, ok := ignoredPaths[path]; ok {
		return ""
	}

	switch expectedValue := expected.(type) {
	case map[string]any:
		actualValue, ok := actual.(map[string]any)
		if !ok {
			return fmt.Sprintf("%s type mismatch: expected object, got %T", path, actual)
		}
		return diffJSONObjects(path, expectedValue, actualValue, ignoredPaths)
	case []any:
		actualValue, ok := actual.([]any)
		if !ok {
			return fmt.Sprintf("%s type mismatch: expected array, got %T", path, actual)
		}
		return diffJSONArrays(path, expectedValue, actualValue, ignoredPaths)
	case json.Number:
		actualValue, ok := actual.(json.Number)
		if !ok {
			return fmt.Sprintf("%s type mismatch: expected number, got %T", path, actual)
		}
		return diffJSONNumbers(path, expectedValue, actualValue)
	case string:
		actualValue, ok := actual.(string)
		if !ok {
			return fmt.Sprintf("%s type mismatch: expected string, got %T", path, actual)
		}
		if expectedValue != actualValue {
			return fmt.Sprintf("%s value mismatch: expected %q, got %q", path, expectedValue, actualValue)
		}
	case bool:
		actualValue, ok := actual.(bool)
		if !ok {
			return fmt.Sprintf("%s type mismatch: expected bool, got %T", path, actual)
		}
		if expectedValue != actualValue {
			return fmt.Sprintf("%s value mismatch: expected %t, got %t", path, expectedValue, actualValue)
		}
	case nil:
		if actual != nil {
			return fmt.Sprintf("%s type mismatch: expected null, got %T", path, actual)
		}
	default:
		if fmt.Sprint(expected) != fmt.Sprint(actual) {
			return fmt.Sprintf("%s value mismatch: expected %v, got %v", path, expected, actual)
		}
	}

	return ""
}

func diffJSONObjects(path string, expected, actual map[string]any, ignoredPaths map[string]struct{}) string {
	keys := make(map[string]struct{}, len(expected)+len(actual))
	for key := range expected {
		keys[key] = struct{}{}
	}
	for key := range actual {
		keys[key] = struct{}{}
	}

	sortedKeys := make([]string, 0, len(keys))
	for key := range keys {
		sortedKeys = append(sortedKeys, key)
	}
	sort.Strings(sortedKeys)

	for _, key := range sortedKeys {
		childPath := path + "." + key
		if _, ok := ignoredPaths[childPath]; ok {
			continue
		}

		expectedChild, expectedOK := expected[key]
		actualChild, actualOK := actual[key]
		if !expectedOK {
			return fmt.Sprintf("%s extra field in actual", childPath)
		}
		if !actualOK {
			return fmt.Sprintf("%s missing field in actual", childPath)
		}
		if diff := diffJSON(childPath, expectedChild, actualChild, ignoredPaths); diff != "" {
			return diff
		}
	}

	return ""
}

func diffJSONArrays(path string, expected, actual []any, ignoredPaths map[string]struct{}) string {
	if len(expected) != len(actual) {
		return fmt.Sprintf("%s length mismatch: expected %d, got %d", path, len(expected), len(actual))
	}

	for i := range expected {
		childPath := fmt.Sprintf("%s[%d]", path, i)
		if diff := diffJSON(childPath, expected[i], actual[i], ignoredPaths); diff != "" {
			return diff
		}
	}

	return ""
}

func diffJSONNumbers(path string, expected, actual json.Number) string {
	expectedFloat, expectedErr := strconv.ParseFloat(expected.String(), 64)
	actualFloat, actualErr := strconv.ParseFloat(actual.String(), 64)
	if expectedErr != nil || actualErr != nil {
		if expected.String() != actual.String() {
			return fmt.Sprintf("%s number mismatch: expected %s, got %s", path, expected, actual)
		}
		return ""
	}

	if expectedFloat != actualFloat {
		return fmt.Sprintf("%s number mismatch: expected %s, got %s", path, expected, actual)
	}
	return ""
}
