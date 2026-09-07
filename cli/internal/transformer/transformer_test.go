package transformer

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	"github.com/apache/arrow-go/v18/arrow/memory"
	"github.com/cloudquery/plugin-sdk/v4/schema"
)

var transformTestCases = []struct {
	name               string
	transformer        func() *RecordTransformer
	originalSchema     *arrow.Schema
	originalJSONRecord []byte
	expectedSchema     *arrow.Schema
	expectedJSONRecord []byte
}{
	{
		name: "no_transformation",
		transformer: func() *RecordTransformer {
			return NewRecordTransformer()
		},
		originalSchema: arrow.NewSchema([]arrow.Field{
			{Name: "id", Type: arrow.PrimitiveTypes.Int64},
		}, nil),
		originalJSONRecord: []byte(`{"id": 1}`),
		expectedSchema: arrow.NewSchema([]arrow.Field{
			{Name: "id", Type: arrow.PrimitiveTypes.Int64},
		}, nil),
		expectedJSONRecord: []byte(`{"id": 1}`),
	},
	{
		name: "add_source",
		transformer: func() *RecordTransformer {
			return NewRecordTransformer(WithSourceNameColumn("test"))
		},
		originalSchema: arrow.NewSchema([]arrow.Field{
			{Name: "id", Type: arrow.PrimitiveTypes.Int64},
		}, nil),
		originalJSONRecord: []byte(`{"id": 1}`),
		expectedSchema: arrow.NewSchema([]arrow.Field{
			{Name: "_cq_source_name", Type: arrow.BinaryTypes.String, Nullable: true},
			{Name: "id", Type: arrow.PrimitiveTypes.Int64},
		}, nil),
		expectedJSONRecord: []byte(`{"_cq_source_name": "test","id": 1}`),
	},
	{
		name: "add_sync_time",
		transformer: func() *RecordTransformer {
			t, err := time.Parse(time.RFC3339, "2023-06-21T17:54:44.488177Z")
			if err != nil {
				panic(err)
			}
			return NewRecordTransformer(WithSyncTimeColumn(t))
		},
		originalSchema: arrow.NewSchema([]arrow.Field{
			{Name: "id", Type: arrow.PrimitiveTypes.Int64},
		}, nil),
		originalJSONRecord: []byte(`{"id": 1}`),
		expectedSchema: arrow.NewSchema([]arrow.Field{
			{Name: "_cq_sync_time", Type: arrow.FixedWidthTypes.Timestamp_us, Nullable: true},
			{Name: "id", Type: arrow.PrimitiveTypes.Int64},
		}, nil),
		expectedJSONRecord: []byte(`{"_cq_sync_time": "2023-06-21 17:54:44.488177","id": 1}`),
	},
	{
		name: "add_source_and_sync_time",
		transformer: func() *RecordTransformer {
			t, err := time.Parse(time.RFC3339, "2023-06-21T17:54:44.488177Z")
			if err != nil {
				panic(err)
			}
			return NewRecordTransformer(WithSyncTimeColumn(t), WithSourceNameColumn("test"))
		},
		originalSchema: arrow.NewSchema([]arrow.Field{
			{Name: "id", Type: arrow.PrimitiveTypes.Int64},
		}, nil),
		originalJSONRecord: []byte(`{"id": 1}`),
		expectedSchema: arrow.NewSchema([]arrow.Field{
			{Name: "_cq_sync_time", Type: arrow.FixedWidthTypes.Timestamp_us, Nullable: true},
			{Name: "_cq_source_name", Type: arrow.BinaryTypes.String, Nullable: true},
			{Name: "id", Type: arrow.PrimitiveTypes.Int64},
		}, nil),
		expectedJSONRecord: []byte(`{"_cq_sync_time": "2023-06-21 17:54:44.488177","_cq_source_name": "test","id": 1}`),
	},
	{
		name: "use_cq_id_primary_key_with_remove_pks",
		transformer: func() *RecordTransformer {
			return NewRecordTransformer(WithRemovePKs(), WithCQIDPrimaryKey())
		},
		originalSchema: arrow.NewSchema([]arrow.Field{
			{Name: "id", Type: arrow.PrimitiveTypes.Int64, Metadata: arrow.MetadataFrom(map[string]string{schema.MetadataPrimaryKey: "true"})},
			{Name: "_cq_id", Type: arrow.PrimitiveTypes.Int64},
		}, nil),
		originalJSONRecord: []byte(`{"id": 1, "_cq_id": 2}`),
		expectedSchema: arrow.NewSchema([]arrow.Field{
			{Name: "id", Type: arrow.PrimitiveTypes.Int64, Metadata: arrow.MetadataFrom(map[string]string{})},
			{Name: "_cq_id", Type: arrow.PrimitiveTypes.Int64, Metadata: arrow.MetadataFrom(map[string]string{schema.MetadataPrimaryKey: "true"})},
		}, nil),
		expectedJSONRecord: []byte(`{"id": 1, "_cq_id": 2}`),
	},
	{
		name: "use_cq_id_primary_key_with_remove_unique",
		transformer: func() *RecordTransformer {
			return NewRecordTransformer(WithRemovePKs(), WithCQIDPrimaryKey(), WithRemoveUniqueConstraints())
		},
		originalSchema: arrow.NewSchema([]arrow.Field{
			{Name: "id", Type: arrow.PrimitiveTypes.Int64, Metadata: arrow.MetadataFrom(map[string]string{schema.MetadataPrimaryKey: "true", schema.MetadataUnique: "true"})},
			{Name: "_cq_id", Type: arrow.PrimitiveTypes.Int64},
		}, nil),
		originalJSONRecord: []byte(`{"id": 1, "_cq_id": 2}`),
		expectedSchema: arrow.NewSchema([]arrow.Field{
			{Name: "id", Type: arrow.PrimitiveTypes.Int64, Metadata: arrow.MetadataFrom(map[string]string{})},
			{Name: "_cq_id", Type: arrow.PrimitiveTypes.Int64, Metadata: arrow.MetadataFrom(map[string]string{schema.MetadataPrimaryKey: "true"})},
		}, nil),
		expectedJSONRecord: []byte(`{"id": 1, "_cq_id": 2}`),
	},
	{
		name: "use_with_remove_unique",
		transformer: func() *RecordTransformer {
			return NewRecordTransformer(WithRemovePKs(), WithRemoveUniqueConstraints())
		},
		originalSchema: arrow.NewSchema([]arrow.Field{
			{Name: "id", Type: arrow.PrimitiveTypes.Int64, Metadata: arrow.MetadataFrom(map[string]string{schema.MetadataPrimaryKey: "true", schema.MetadataUnique: "true"})},
			{Name: "_cq_id", Type: arrow.PrimitiveTypes.Int64},
		}, nil),
		originalJSONRecord: []byte(`{"id": 1, "_cq_id": 2}`),
		expectedSchema: arrow.NewSchema([]arrow.Field{
			{Name: "id", Type: arrow.PrimitiveTypes.Int64, Metadata: arrow.MetadataFrom(map[string]string{})},
			{Name: "_cq_id", Type: arrow.PrimitiveTypes.Int64},
		}, nil),
		expectedJSONRecord: []byte(`{"id": 1, "_cq_id": 2}`),
	},
}

func TestRecord(t *testing.T) {
	for _, tc := range transformTestCases {
		t.Run(tc.name, func(t *testing.T) {
			bldr := array.NewRecordBuilder(memory.DefaultAllocator, tc.originalSchema)
			if err := bldr.UnmarshalJSON(tc.originalJSONRecord); err != nil {
				t.Fatal(err)
			}
			record := bldr.NewRecordBatch()
			transformedRecord := tc.transformer().Transform(record)
			if transformedRecord.Schema().String() != tc.expectedSchema.String() {
				t.Fatalf("expected schema\n%v, got\n%v", tc.expectedSchema, transformedRecord.Schema())
			}
			bldr = array.NewRecordBuilder(memory.DefaultAllocator, transformedRecord.Schema())
			if err := bldr.UnmarshalJSON(tc.expectedJSONRecord); err != nil {
				t.Fatal(err)
			}
			expectedRecord := bldr.NewRecordBatch()
			if !array.RecordEqual(expectedRecord, transformedRecord) {
				b, err := json.Marshal(transformedRecord)
				if err != nil {
					t.Fatal(err)
				}
				t.Logf("expected record %v, got %v", string(tc.expectedJSONRecord), string(b))
				t.Fatalf("expected record %v, got %v", expectedRecord, transformedRecord)
			}
		})
	}
}

func TestTransformSchemaCaching(t *testing.T) {
	tr := NewRecordTransformer(
		WithSourceNameColumn("test"),
		WithRemovePKs(),
	)
	sc := arrow.NewSchema([]arrow.Field{
		{Name: "id", Type: arrow.PrimitiveTypes.Int64, Metadata: arrow.MetadataFrom(map[string]string{schema.MetadataPrimaryKey: "true"})},
	}, nil)

	s1 := tr.TransformSchema(sc)
	s2 := tr.TransformSchema(sc)

	if s1 != s2 {
		t.Fatalf("expected identical pointer for cached schema transformation, got %p vs %p", s1, s2)
	}
}

func BenchmarkTransform(b *testing.B) {
	tTime, _ := time.Parse(time.RFC3339, "2023-06-21T17:54:44.488177Z")
	tr := NewRecordTransformer(
		WithSyncTimeColumn(tTime),
		WithSourceNameColumn("test_source"),
		WithRemovePKs(),
		WithRemoveUniqueConstraints(),
		WithCQIDPrimaryKey(),
	)

	fields := make([]arrow.Field, 20)
	for i := 0; i < 19; i++ {
		fields[i] = arrow.Field{
			Name:     "col_" + string(rune('a'+i)),
			Type:     arrow.PrimitiveTypes.Int64,
			Metadata: arrow.MetadataFrom(map[string]string{"comment": "test"}),
		}
	}
	fields[0].Metadata = arrow.MetadataFrom(map[string]string{schema.MetadataPrimaryKey: "true"})
	fields[19] = arrow.Field{Name: "_cq_id", Type: arrow.BinaryTypes.String}

	sc := arrow.NewSchema(fields, nil)
	bldr := array.NewRecordBuilder(memory.DefaultAllocator, sc)
	for i := 0; i < 20; i++ {
		if i == 19 {
			bldr.Field(i).(*array.StringBuilder).Append("cq_id_val")
		} else {
			bldr.Field(i).(*array.Int64Builder).Append(int64(i))
		}
	}
	rec := bldr.NewRecordBatch()
	defer rec.Release()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = tr.Transform(rec)
	}
}
