package genshin

import (
	"encoding/json"
	"os"
	"testing"
)

var benchmarkProfileJSON []byte

func BenchmarkProfileUnmarshalJSON(b *testing.B) {
	data := benchmarkFixture(b, "testdata/profile.json")
	b.ReportAllocs()
	b.SetBytes(int64(len(data)))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var profile Profile
		if err := json.Unmarshal(data, &profile); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkProfileMarshalJSON(b *testing.B) {
	data := benchmarkFixture(b, "testdata/profile.json")
	var profile Profile
	if err := json.Unmarshal(data, &profile); err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.SetBytes(int64(len(data)))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		encoded, err := json.Marshal(&profile)
		if err != nil {
			b.Fatal(err)
		}
		benchmarkProfileJSON = encoded
	}
}

func benchmarkFixture(b *testing.B, path string) []byte {
	b.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		b.Fatal(err)
	}
	return data
}
