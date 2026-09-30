package otel

import (
	"time"

	"github.com/cloudquery/cloudquery/cli/v6/internal/utils"
)

type TableDurations struct {
	data *utils.ConcurrentMap[string, time.Duration]
}

func NewTableDurations() *TableDurations {
	return &TableDurations{
		data: utils.NewConcurrentMap[string, time.Duration](),
	}
}

func (td *TableDurations) Set(table string, duration time.Duration) {
	// Use SetIfGreater to perform the lookup, comparison, and update atomically under a single write lock.
	// This avoids acquiring an RLock followed by a write Lock and redundant map hashing.
	td.data.SetIfGreater(table, duration, func(newVal, existingVal time.Duration) bool {
		return newVal > existingVal
	})
}

func (td *TableDurations) GetAll() map[string]time.Duration {
	return td.data.GetAll()
}
