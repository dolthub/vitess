// Copyright 2026 Dolthub, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sqlparser

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

// Include every comparison node in the rendering: formatting SQL alone cannot
// distinguish a left-associated chain from a right-associated one.
func comparisonTree(expr Expr) string {
	switch e := expr.(type) {
	case *ComparisonExpr:
		return fmt.Sprintf("(%s %s %s)", comparisonTree(e.Left), e.Operator, comparisonTree(e.Right))
	case *ParenExpr:
		return comparisonTree(e.Expr)
	case *BinaryExpr:
		return fmt.Sprintf("(%s %s %s)", comparisonTree(e.Left), e.Operator, comparisonTree(e.Right))
	case *AndExpr:
		return fmt.Sprintf("(%s and %s)", comparisonTree(e.Left), comparisonTree(e.Right))
	case *OrExpr:
		return fmt.Sprintf("(%s or %s)", comparisonTree(e.Left), comparisonTree(e.Right))
	case *NotExpr:
		return fmt.Sprintf("(not %s)", comparisonTree(e.Expr))
	case *IsExpr:
		return fmt.Sprintf("(%s %s)", comparisonTree(e.Expr), e.Operator)
	default:
		return String(expr)
	}
}

func TestChainedComparisonPrecedence(t *testing.T) {
	tests := []struct{ input, want string }{
		{"1 < 2 = 1", "((1 < 2) = 1)"},
		{"(1, 2) < (1, 3) = 1", "(((1, 2) < (1, 3)) = 1)"},
		{"2 < 3 = 1", "((2 < 3) = 1)"},
		{"2 < (3 = 1)", "(2 < (3 = 1))"},
		{"(2 < 3) = 1", "((2 < 3) = 1)"},
		{"1 < 2 <= 3 <> 0 = 1", "((((1 < 2) <= 3) != 0) = 1)"},
		{"3 > 2 >= 1 <=> 1", "(((3 > 2) >= 1) <=> 1)"},
		{"null < 2 = 1", "((null < 2) = 1)"},
		{"null <=> null = 1", "((null <=> null) = 1)"},
		{"1 + 1 < 3 * 2 = 1", "(((1 + 1) < (3 * 2)) = 1)"},
		{"not 2 < 3 = 1", "(not ((2 < 3) = 1))"},
		{"2 < 3 = 1 and 4 > 3 = 1", "(((2 < 3) = 1) and ((4 > 3) = 1))"},
		{"2 < 3 = 1 or 0 and 1", "(((2 < 3) = 1) or (0 and 1))"},
		{"2 < 3 = 1 is true", "(((2 < 3) = 1) is true)"},
		{"null < 3 = 1 is null", "(((null < 3) = 1) is null)"},
		{"2 < 3 = (1 in (1, 2))", "((2 < 3) = (1 in (1, 2)))"},
		{"2 < 3 = ('a' like 'a')", "((2 < 3) = ('a' like 'a'))"},
		{"2 < 3 = ('a' regexp 'a')", "((2 < 3) = ('a' regexp 'a'))"},
		{"2 < 3 = (1 between 0 and 2)", "((2 < 3) = 1 between 0 and 2)"},
		{"2 < 3 = 1 and 1 between 0 and 2", "(((2 < 3) = 1) and 1 between 0 and 2)"},
		{"2 < 3 = 1 and 1 in (1, 2)", "(((2 < 3) = 1) and (1 in (1, 2)))"},
		{"(select 2) < 3 = 1", "(((select 2) < 3) = 1)"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			stmt, err := Parse("select " + tt.input)
			require.NoError(t, err)
			expr := stmt.(*Select).SelectExprs[0].(*AliasedExpr).Expr
			require.Equal(t, tt.want, comparisonTree(expr))
			// Formatting must retain enough grouping to reconstruct the same tree.
			again, err := Parse(String(stmt))
			require.NoError(t, err)
			require.Equal(t, tt.want, comparisonTree(again.(*Select).SelectExprs[0].(*AliasedExpr).Expr))
		})
	}
}
