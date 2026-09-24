package pgarrow

import "testing"

func BenchmarkPg10ToArrow(b *testing.B) {
	types := []string{
		"boolean",
		"bigint",
		"timestamp(3) with time zone",
		"numeric(38,15)",
		"text[]",
		"varchar(50)[][]",
		"TIMESTAMPTZ USING timestamp(6)",
		"unknown_custom_type",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, t := range types {
			_ = Pg10ToArrow(t)
		}
	}
}
