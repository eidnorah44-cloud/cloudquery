package transformer

import (
	"testing"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/cloudquery/plugin-sdk/v4/schema"
)

func BenchmarkTransformSchema_Typical(b *testing.B) {
	fields := make([]arrow.Field, 50)
	for i := 0; i < 50; i++ {
		var md arrow.Metadata
		if i == 0 {
			md = arrow.MetadataFrom(map[string]string{
				schema.MetadataPrimaryKey: "true",
				schema.MetadataUnique:     "true",
			})
		}
		fields[i] = arrow.Field{
			Name:     "col_" + string(rune('a'+i%26)),
			Type:     arrow.PrimitiveTypes.Int64,
			Metadata: md,
		}
	}
	sc := arrow.NewSchema(fields, nil)
	t := NewRecordTransformer(WithRemovePKs(), WithRemoveUniqueConstraints(), WithCQIDPrimaryKey())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = t.TransformSchema(sc)
	}
}

func BenchmarkTransformSchema_AllPK(b *testing.B) {
	fields := make([]arrow.Field, 50)
	for i := 0; i < 50; i++ {
		fields[i] = arrow.Field{
			Name: "col_" + string(rune('a'+i%26)),
			Type: arrow.PrimitiveTypes.Int64,
			Metadata: arrow.MetadataFrom(map[string]string{
				schema.MetadataPrimaryKey: "true",
				schema.MetadataUnique:     "true",
			}),
		}
	}
	sc := arrow.NewSchema(fields, nil)
	t := NewRecordTransformer(WithRemovePKs(), WithRemoveUniqueConstraints(), WithCQIDPrimaryKey())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = t.TransformSchema(sc)
	}
}
