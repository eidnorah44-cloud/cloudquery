package client

import (
	"testing"

	"github.com/apache/arrow-go/v18/arrow"
	cqtypes "github.com/cloudquery/plugin-sdk/v4/types"
)

func TestSnowflakeToSchemaType(t *testing.T) {
	cases := []struct {
		snowflakeType string
		want          arrow.DataType
	}{
		{"boolean", arrow.FixedWidthTypes.Boolean},
		{"tinyint", arrow.PrimitiveTypes.Int64},
		{"smallint", arrow.PrimitiveTypes.Int64},
		{"integer", arrow.PrimitiveTypes.Int64},
		{"int", arrow.PrimitiveTypes.Int64},
		{"bigint", arrow.PrimitiveTypes.Int64},
		{"float", arrow.PrimitiveTypes.Float64},
		{"double", arrow.PrimitiveTypes.Float64},
		{"date", arrow.FixedWidthTypes.Date32},
		{"binary", arrow.BinaryTypes.Binary},
		{"variant", cqtypes.ExtensionTypes.JSON},
		{"timestamp_ltz", arrow.FixedWidthTypes.Timestamp_ns},
		{"timestamp_tz", arrow.FixedWidthTypes.Timestamp_ns},
		{"timestamp_ntz", arrow.FixedWidthTypes.Timestamp_ns},
		{"timestamp(0)", arrow.FixedWidthTypes.Timestamp_s},
		{"timestamp(3)", arrow.FixedWidthTypes.Timestamp_ms},
		{"timestamp(6)", arrow.FixedWidthTypes.Timestamp_us},
		{"timestamp(9)", arrow.FixedWidthTypes.Timestamp_ns},
		{"time(0)", arrow.FixedWidthTypes.Time32s},
		{"time(3)", arrow.FixedWidthTypes.Time32ms},
		{"time(6)", arrow.FixedWidthTypes.Time64us},
		{"time(9)", arrow.FixedWidthTypes.Time64ns},
		{"number(18,0)", arrow.PrimitiveTypes.Int64},
		{"number(20,0)", arrow.PrimitiveTypes.Uint64},
		{"decimal(38,15)", &arrow.Decimal128Type{Precision: 38, Scale: 15}},
		{"numeric(50,25)", &arrow.Decimal256Type{Precision: 50, Scale: 25}},
		{"varchar", arrow.BinaryTypes.String},
		{"text[]", arrow.ListOf(arrow.BinaryTypes.String)},
	}

	for _, c := range cases {
		t.Run(c.snowflakeType, func(t *testing.T) {
			got := SnowflakeToSchemaType(c.snowflakeType)
			if !arrow.TypeEqual(got, c.want) {
				t.Errorf("SnowflakeToSchemaType(%q) = %v, want %v", c.snowflakeType, got, c.want)
			}
		})
	}
}

func BenchmarkSnowflakeToSchemaType(b *testing.B) {
	typesToTest := []string{
		"boolean", "bigint", "variant", "text[]",
		"timestamp(3)", "number(38,15)", "varchar",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, t := range typesToTest {
			_ = SnowflakeToSchemaType(t)
		}
	}
}
