package client

import (
	"testing"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/stretchr/testify/require"
)

func TestSanitizeColumn(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "simple_col",
			expected: `"SIMPLE_COL"`,
		},
		{
			input:    `col"name`,
			expected: `"COL""NAME"`,
		},
		{
			input:    `col"; DROP TABLE users; --`,
			expected: `"COL""; DROP TABLE USERS; --"`,
		},
	}

	for _, tc := range tests {
		require.Equal(t, tc.expected, sanitizeColumn(tc.input))
	}
}

func TestSQLGenerationSanitization(t *testing.T) {
	table := &schema.Table{
		Name: `test_table"; DROP TABLE users; --`,
		Columns: []schema.Column{
			{Name: `id"; --`, Type: arrow.BinaryTypes.String, PrimaryKey: true},
			{Name: `data"col`, Type: arrow.BinaryTypes.String},
		},
	}

	colsList := createColumnsList(table)
	require.Equal(t, `$1:"ID""; --"::text as "ID""; --",$1:"DATA""COL"::text as "DATA""COL"`, colsList)

	pkList := createPrimaryKeyList(table)
	require.Equal(t, `source."ID""; --"=dest."ID""; --"`, pkList)

	updateList := updateColumnsList(table)
	require.Equal(t, ` UPDATE SET "ID""; --"=source."ID""; --","DATA""COL"=source."DATA""COL" `, updateList)

	insertList := insertColumnsList(table)
	require.Equal(t, `INSERT ("ID""; --", "DATA""COL") VALUES (source."ID""; --", source."DATA""COL")`, insertList)
}
