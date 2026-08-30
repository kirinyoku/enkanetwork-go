package core

import "encoding/json"

// IsValidUID checks if the provided UID is a valid 9-digit number.
// Genshin and HSR UID can only be 9 digits (e.g., "618285856").
// This function is used internally to validate UIDs before making requests.
//
// Parameters:
//   - uid: The UID string to validate.
//
// Returns:
//   - true if the UID is a 9-digit number, false otherwise.
func IsValidUID(uid string) bool {
	if len(uid) != 9 {
		return false
	}
	for _, r := range uid {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// RemoveTTLField removes the TTL field from the JSON response.
// This is used for tests to ensure the response is consistent.
func RemoveTTLField(jsonBytes []byte) []byte {
	var profile map[string]any
	if err := json.Unmarshal(jsonBytes, &profile); err != nil {
		return jsonBytes
	}

	delete(profile, "ttl")

	newJSON, err := json.Marshal(profile)
	if err != nil {
		return jsonBytes
	}
	return newJSON
}
