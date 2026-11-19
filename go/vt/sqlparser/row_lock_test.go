/*
Copyright 2026 Dolthub, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package sqlparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRowLockItems(t *testing.T) {
	tests := []struct {
		query    string
		expected []RowLockItem
	}{
		{query: "select * from t"},
		{
			query:    "select * from t for update",
			expected: []RowLockItem{{Strength: RowLockForUpdate}},
		},
		{
			query:    "select * from t lock in share mode",
			expected: []RowLockItem{{Strength: RowLockForShare}},
		},
		{
			query:    "select * from t for update skip locked",
			expected: []RowLockItem{{Strength: RowLockForUpdate, WaitPolicy: RowLockSkipLocked}},
		},
		{
			query:    "select * from t for update nowait",
			expected: []RowLockItem{{Strength: RowLockForUpdate, WaitPolicy: RowLockNoWait}},
		},
		{
			query: "select * from t1 join t2 on t1.a = t2.a for update of t1, t2",
			expected: []RowLockItem{{
				Strength: RowLockForUpdate,
				Tables:   TableNames{{Name: NewTableIdent("t1")}, {Name: NewTableIdent("t2")}},
			}},
		},
		{
			query: "select * from t1 join t2 on t1.a = t2.a for update of t2 skip locked",
			expected: []RowLockItem{{
				Strength:   RowLockForUpdate,
				WaitPolicy: RowLockSkipLocked,
				Tables:     TableNames{{Name: NewTableIdent("t2")}},
			}},
		},
		{
			query: "select * from t1 join t2 on t1.a = t2.a for update of t1 nowait",
			expected: []RowLockItem{{
				Strength:   RowLockForUpdate,
				WaitPolicy: RowLockNoWait,
				Tables:     TableNames{{Name: NewTableIdent("t1")}},
			}},
		},
	}

	for _, test := range tests {
		t.Run(test.query, func(t *testing.T) {
			stmt, err := Parse(test.query)
			require.NoError(t, err)
			sel, ok := stmt.(*Select)
			require.True(t, ok)
			assert.Equal(t, test.expected, sel.Lock.RowLockItems())
		})
	}

	t.Run("items set by a caller take precedence", func(t *testing.T) {
		items := []RowLockItem{
			{Strength: RowLockForKeyShare, Tables: TableNames{{Name: NewTableIdent("t1")}}},
			{Strength: RowLockForNoKeyUpdate, WaitPolicy: RowLockNoWait, Tables: TableNames{{Name: NewTableIdent("t2")}}},
		}
		lock := &Lock{Type: ForUpdateStr, Items: items}
		assert.Equal(t, items, lock.RowLockItems())
	})
}
