package client

import (
	"testing"
)

func TestSnowflakeToSchemaType(t *testing.T) {
	cases := []string{
		"timestamp_tz(3)",
		"time(0)",
		"number(38,15)",
		"boolean",
		"varchar[]",
	}
	for _, c := range cases {
		got := SnowflakeToSchemaType(c)
		if got == nil {
			t.Errorf("SnowflakeToSchemaType(%q) returned nil", c)
		}
	}
}

func BenchmarkSnowflakeToSchemaType(b *testing.B) {
	typesToTest := []string{
		"timestamp_tz(3)",
		"time(0)",
		"number(38,15)",
		"boolean",
		"varchar[]",
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for _, t := range typesToTest {
			_ = SnowflakeToSchemaType(t)
		}
	}
}
