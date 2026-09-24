package spec

import (
	"testing"
	"time"

	"github.com/cloudquery/filetypes/v4"
)

func BenchmarkReplacePathVariables(b *testing.B) {
	s := &Spec{
		Path:     "path/to/files/{{TABLE}}/{{UUID}}.parquet",
		FileSpec: filetypes.FileSpec{Format: filetypes.FormatTypeParquet},
	}
	tm := time.Now().UTC()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.ReplacePathVariables("aws_ec2_instances", "12345678-1234-1234-1234-1234567890ab", tm, "sync-id-123")
	}
}
