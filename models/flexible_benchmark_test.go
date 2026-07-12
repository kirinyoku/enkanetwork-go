package models

import (
	"encoding/json"
	"testing"
)

var (
	benchmarkStringNumber StringNumber
	benchmarkIntString    IntString
	benchmarkScalarJSON   []byte
)

func BenchmarkStringNumberUnmarshalNumber(b *testing.B) {
	data := []byte(`1311116816978843608`)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if err := json.Unmarshal(data, &benchmarkStringNumber); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkStringNumberUnmarshalString(b *testing.B) {
	data := []byte(`"1311116816978843608"`)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if err := json.Unmarshal(data, &benchmarkStringNumber); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkIntStringUnmarshalNumber(b *testing.B) {
	data := []byte(`1311116816978843608`)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if err := json.Unmarshal(data, &benchmarkIntString); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkIntStringUnmarshalString(b *testing.B) {
	data := []byte(`"1311116816978843608"`)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if err := json.Unmarshal(data, &benchmarkIntString); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkIntStringMarshalJSON(b *testing.B) {
	value := IntString(1311116816978843608)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		encoded, err := value.MarshalJSON()
		if err != nil {
			b.Fatal(err)
		}
		benchmarkScalarJSON = encoded
	}
}
