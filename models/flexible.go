package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

// StringNumber is a string-facing scalar that accepts either a JSON string or
// number. It is useful for API hash fields that sometimes change representation.
type StringNumber string

// UnmarshalJSON accepts JSON strings, numbers, and null.
func (s *StringNumber) UnmarshalJSON(data []byte) error {
	trimmed := bytes.TrimSpace(data)
	if bytes.Equal(trimmed, []byte("null")) {
		*s = ""
		return nil
	}

	if len(trimmed) > 0 && trimmed[0] == '"' {
		var asString string
		if err := json.Unmarshal(trimmed, &asString); err == nil {
			*s = StringNumber(asString)
			return nil
		}
	}

	if isJSONNumber(trimmed) {
		*s = StringNumber(string(trimmed))
		return nil
	}

	return fmt.Errorf("StringNumber: expected string or number, got %s", data)
}

// MarshalJSON serializes the value as a JSON string.
func (s StringNumber) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(s))
}

// String returns the normalized string value.
func (s StringNumber) String() string {
	return string(s)
}

// Int64 parses the value as an int64.
func (s StringNumber) Int64() (int64, error) {
	return strconv.ParseInt(string(s), 10, 64)
}

// IsZero reports whether the value is empty.
func (s StringNumber) IsZero() bool {
	return s == ""
}

// IntString is an integer-facing scalar that accepts either a JSON number or
// string. It is kept for API fields that are numeric in Go but unstable on the wire.
type IntString int64

// UnmarshalJSON accepts JSON strings, numbers, and null.
func (i *IntString) UnmarshalJSON(data []byte) error {
	trimmed := bytes.TrimSpace(data)
	if bytes.Equal(trimmed, []byte("null")) {
		*i = 0
		return nil
	}

	if isJSONNumber(trimmed) {
		parsed, parseErr := strconv.ParseInt(string(trimmed), 10, 64)
		if parseErr != nil {
			return parseErr
		}
		*i = IntString(parsed)
		return nil
	}

	if len(trimmed) > 0 && trimmed[0] == '"' {
		var asString string
		if err := json.Unmarshal(trimmed, &asString); err == nil {
			parsed, parseErr := strconv.ParseInt(asString, 10, 64)
			if parseErr != nil {
				return parseErr
			}
			*i = IntString(parsed)
			return nil
		}
	}

	return fmt.Errorf("IntString: expected string or number, got %s", data)
}

// MarshalJSON serializes the value as a JSON number.
func (i IntString) MarshalJSON() ([]byte, error) {
	return strconv.AppendInt(nil, int64(i), 10), nil
}

// String returns the normalized decimal string value.
func (i IntString) String() string {
	return strconv.FormatInt(int64(i), 10)
}

// Int64 returns the normalized int64 value.
func (i IntString) Int64() (int64, error) {
	return int64(i), nil
}

// IsZero reports whether the value is zero.
func (i IntString) IsZero() bool {
	return i == 0
}

func isJSONNumber(data []byte) bool {
	if len(data) == 0 || (data[0] != '-' && (data[0] < '0' || data[0] > '9')) {
		return false
	}
	return json.Valid(data)
}
