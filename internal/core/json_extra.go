package core

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"
)

// PreserveUnknownJSON returns a copy of the original JSON document and a map of
// fields not listed in knownFields. It is intended for drift-tolerant API models.
func PreserveUnknownJSON(data []byte, knownFields ...string) (json.RawMessage, map[string]json.RawMessage, error) {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return nil, nil, err
	}

	for _, field := range knownFields {
		delete(fields, field)
	}

	raw := append(json.RawMessage(nil), data...)
	if len(fields) == 0 {
		return raw, nil, nil
	}
	return raw, fields, nil
}

// PreserveUnknownJSONForStruct returns a copy of the original JSON document and
// fields not listed in the JSON tags of model. It is useful for drift-tolerant
// structs with many fields because the known field list follows the struct tags.
func PreserveUnknownJSONForStruct(data []byte, model any) (json.RawMessage, map[string]json.RawMessage, error) {
	return PreserveUnknownJSON(data, JSONFieldNames(model)...)
}

// JSONFieldNames returns JSON object field names declared on a struct type.
func JSONFieldNames(model any) []string {
	typ := reflect.TypeOf(model)
	for typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return nil
	}

	fields := make([]string, 0, typ.NumField())
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.PkgPath != "" {
			continue
		}

		tag := field.Tag.Get("json")
		name, _, _ := strings.Cut(tag, ",")
		switch {
		case name == "-":
			continue
		case name != "":
			fields = append(fields, name)
		default:
			fields = append(fields, field.Name)
		}
	}
	return fields
}

// MergeKnownAndExtraJSON marshals known fields and then adds preserved unknown
// fields. Known fields win if Extra contains a duplicate key.
func MergeKnownAndExtraJSON(known any, extra map[string]json.RawMessage) ([]byte, error) {
	knownJSON, err := json.Marshal(known)
	if err != nil {
		return nil, err
	}
	if len(extra) == 0 {
		return knownJSON, nil
	}

	var fields map[string]json.RawMessage
	if err := json.Unmarshal(knownJSON, &fields); err != nil {
		return nil, err
	}

	for key, value := range extra {
		if _, exists := fields[key]; !exists {
			fields[key] = value
		}
	}

	return json.Marshal(fields)
}

// MergeKnownExtraAndRawJSON marshals known fields, restores explicitly present
// zero-like known values from raw JSON when omitempty dropped them, and then adds
// preserved unknown fields. It is intended for API response models that keep Raw
// and Extra for drift-tolerant decode/encode round trips.
func MergeKnownExtraAndRawJSON(known any, extra map[string]json.RawMessage, raw json.RawMessage) ([]byte, error) {
	knownJSON, err := json.Marshal(known)
	if err != nil {
		return nil, err
	}

	var fields map[string]json.RawMessage
	if err := json.Unmarshal(knownJSON, &fields); err != nil {
		return nil, err
	}

	if len(raw) > 0 {
		var rawFields map[string]json.RawMessage
		if err := json.Unmarshal(raw, &rawFields); err != nil {
			return nil, err
		}
		for key, value := range rawFields {
			if _, exists := fields[key]; !exists && isZeroLikeJSON(value) {
				fields[key] = value
			}
		}
	}

	for key, value := range extra {
		if _, exists := fields[key]; !exists {
			fields[key] = value
		}
	}

	return json.Marshal(fields)
}

func isZeroLikeJSON(value json.RawMessage) bool {
	trimmed := bytes.TrimSpace(value)
	if bytes.Equal(trimmed, []byte("null")) ||
		bytes.Equal(trimmed, []byte("false")) ||
		bytes.Equal(trimmed, []byte(`""`)) ||
		bytes.Equal(trimmed, []byte("[]")) ||
		bytes.Equal(trimmed, []byte("{}")) {
		return true
	}

	var number json.Number
	decoder := json.NewDecoder(bytes.NewReader(trimmed))
	decoder.UseNumber()
	if err := decoder.Decode(&number); err != nil {
		return false
	}
	asFloat, err := number.Float64()
	return err == nil && asFloat == 0
}
