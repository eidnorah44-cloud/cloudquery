package client

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/cloudquery/plugin-sdk/v4/types"
)

var (
	reTimestamp = regexp.MustCompile(`timestamp(?:_(?:ltz|tz|ntz))?\s*(?:\(([0-9])\))?`)
	reTime      = regexp.MustCompile(`time\s*(?:\(([0-9])\))?`)
	reNumeric   = regexp.MustCompile(`(?:numeric|number|decimal)\s*(?:\(([0-9]+)\s*(?:,\s*([0-9]+))?\))?`)
)

func SchemaTypeToSnowflake(t arrow.DataType) string {
	switch t.(type) {
	case *arrow.ListType, *arrow.FixedSizeListType:
		return "array"
	case *arrow.BooleanType:
		return "boolean"
	case *arrow.Int8Type, *arrow.Uint8Type, *arrow.Int16Type, *arrow.Uint16Type,
		*arrow.Int32Type, *arrow.Uint32Type, *arrow.Int64Type, *arrow.Uint64Type:
		return "number"
	case *arrow.Float32Type, *arrow.Float64Type:
		return "float"
	case *arrow.StringType, *arrow.LargeStringType:
		return "text"
	case *arrow.BinaryType, *arrow.LargeBinaryType:
		return "binary"
	case *arrow.TimestampType:
		return "timestamp_tz"
	case *types.JSONType, *arrow.StructType:
		return "variant"
	default:
		return "text"
	}
}

func SnowflakeToSchemaType(t string) arrow.DataType {
	t = strings.ToLower(strings.TrimSpace(t))
	if strings.HasSuffix(t, "[]") {
		return arrow.ListOf(SnowflakeToSchemaType(t[:len(t)-2]))
	}

	switch t {
	case "boolean":
		return arrow.FixedWidthTypes.Boolean

	case "tinyint", "smallint", "integer", "int", "bigint":
		return arrow.PrimitiveTypes.Int64

	case "float", "float4", "float8", "double", "double precision", "real":
		return arrow.PrimitiveTypes.Float64

	case "date":
		return arrow.FixedWidthTypes.Date32

	case "binary", "varbinary":
		return arrow.BinaryTypes.Binary

	case "variant", "object", "array":
		return types.ExtensionTypes.JSON
	}

	if strings.HasPrefix(t, "timestamp") {
		if got, matched := parseTimestamp(t); matched {
			return got
		}
	}
	if strings.HasPrefix(t, "time") {
		if got, matched := parseTime(t); matched {
			return got
		}
	}
	if strings.HasPrefix(t, "numeric") || strings.HasPrefix(t, "number") || strings.HasPrefix(t, "decimal") {
		if got, matched := parseNumeric(t); matched {
			return got
		}
	}

	return arrow.BinaryTypes.String
}

func parseTimestamp(t string) (arrow.DataType, bool) {
	// Handle Snowflake TIMESTAMP_* types with optional precision
	if t == "timestamp_ltz" || t == "timestamp_tz" || t == "timestamp_ntz" || t == "timestamp" {
		return arrow.FixedWidthTypes.Timestamp_ns, true
	}

	matches := reTimestamp.FindStringSubmatch(t)
	if len(matches) == 0 {
		return nil, false
	}

	precisionStr := matches[1]
	precision := 9 // default
	if precisionStr != "" {
		p, err := strconv.Atoi(precisionStr)
		if err == nil {
			precision = p
		}
	}

	switch {
	case precision == 0:
		return arrow.FixedWidthTypes.Timestamp_s, true
	case precision >= 1 && precision <= 3:
		return arrow.FixedWidthTypes.Timestamp_ms, true
	case precision >= 4 && precision <= 6:
		return arrow.FixedWidthTypes.Timestamp_us, true
	case precision >= 7 && precision <= 9:
		return arrow.FixedWidthTypes.Timestamp_ns, true
	default:
		return arrow.FixedWidthTypes.Timestamp_ns, true
	}
}

func parseTime(t string) (arrow.DataType, bool) {
	if t == "time" {
		return arrow.FixedWidthTypes.Time64ns, true
	}

	matches := reTime.FindStringSubmatch(t)
	if len(matches) == 0 {
		return nil, false
	}

	precisionStr := matches[1]
	precision := 9 // default
	if precisionStr != "" {
		p, err := strconv.Atoi(precisionStr)
		if err == nil {
			precision = p
		}
	}

	switch {
	case precision == 0:
		return arrow.FixedWidthTypes.Time32s, true
	case precision >= 1 && precision <= 3:
		return arrow.FixedWidthTypes.Time32ms, true
	case precision >= 4 && precision <= 6:
		return arrow.FixedWidthTypes.Time64us, true
	case precision >= 7 && precision <= 9:
		return arrow.FixedWidthTypes.Time64ns, true
	default:
		return arrow.FixedWidthTypes.Time64ns, true
	}
}

func parseNumeric(t string) (arrow.DataType, bool) {
	if t == "numeric" || t == "number" || t == "decimal" {
		return arrow.PrimitiveTypes.Int64, true
	}

	matches := reNumeric.FindStringSubmatch(t)
	if len(matches) == 0 {
		return nil, false
	}

	// No precision/scale specified - default to Int64
	if len(matches) < 3 || matches[1] == "" {
		return arrow.PrimitiveTypes.Int64, true
	}

	precision, err := strconv.ParseInt(matches[1], 10, 32)
	if precision == 0 || err != nil {
		panic("precision cannot be 0")
	}
	scale, err := strconv.ParseInt(matches[2], 10, 32)
	if err != nil {
		panic("error parsing scale " + err.Error())
	}

	switch {
	case precision <= 18 && scale == 0:
		return arrow.PrimitiveTypes.Int64, true
	case precision == 20 && scale == 0:
		return arrow.PrimitiveTypes.Uint64, true
	case precision <= 38:
		return &arrow.Decimal128Type{Precision: int32(precision), Scale: int32(scale)}, true
	case precision <= 76:
		return &arrow.Decimal256Type{Precision: int32(precision), Scale: int32(scale)}, true
	default:
		return arrow.BinaryTypes.String, true
	}
}
