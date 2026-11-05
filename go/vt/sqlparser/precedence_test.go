/*
Copyright 2019 The Vitess Authors.

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
	"fmt"
	"testing"
)

func readable(node Expr) string {
	switch node := node.(type) {
	case *OrExpr:
		return fmt.Sprintf("(%s or %s)", readable(node.Left), readable(node.Right))
	case *XorExpr:
		return fmt.Sprintf("(%s xor %s)", readable(node.Left), readable(node.Right))
	case *AndExpr:
		return fmt.Sprintf("(%s and %s)", readable(node.Left), readable(node.Right))
	case *BinaryExpr:
		return fmt.Sprintf("(%s %s %s)", readable(node.Left), node.Operator, readable(node.Right))
	case *IsExpr:
		return fmt.Sprintf("(%s %s)", readable(node.Expr), node.Operator)
	default:
		return String(node)
	}
}

func TestLogicOperatorPrecedence(t *testing.T) {
	validSQL := []struct {
		input  string
		output string
	}{{
		input:  "select * from a where a=b and c=d or e=f",
		output: "((a = b and c = d) or e = f)",
	}, {
		input:  "select * from a where a=b or c=d and e=f",
		output: "(a = b or (c = d and e = f))",
	}, {
		input:  "select * from a where f=g xor a=b and e=f",
		output: "(f = g xor (a = b and e = f))",
	}, {
		input:  "select * from a where f=g xor a=b or c=d and e=f",
		output: "((f = g xor a = b) or (c = d and e = f))",
	}, {
		input:  "select * from t where a = 5 or b xor d = 3 and c is not null",
		output: "(a = 5 or (b xor (d = 3 and (c is not null))))",
	}}
	for _, tcase := range validSQL {
		tree, err := Parse(tcase.input)
		if err != nil {
			t.Error(err)
			continue
		}
		expr := readable(tree.(*Select).Where.Expr)
		if expr != tcase.output {
			t.Errorf("Parse: \n%s, want: \n%s", expr, tcase.output)
		}
	}
}

func TestPlusStarPrecedence(t *testing.T) {
	validSQL := []struct {
		input  string
		output string
	}{{
		input:  "select 1+2*3 from a",
		output: "(1 + (2 * 3))",
	}, {
		input:  "select 1*2+3 from a",
		output: "((1 * 2) + 3)",
	}}
	for _, tcase := range validSQL {
		tree, err := Parse(tcase.input)
		if err != nil {
			t.Error(err)
			continue
		}
		expr := readable(tree.(*Select).SelectExprs[0].(*AliasedExpr).Expr)
		if expr != tcase.output {
			t.Errorf("Parse: \n%s, want: \n%s", expr, tcase.output)
		}
	}
}

func TestIsPrecedence(t *testing.T) {
	validSQL := []struct {
		input  string
		output string
	}{{
		input:  "select * from a where a+b is true",
		output: "((a + b) is true)",
	}, {
		input:  "select * from a where a=1 and b=2 is true",
		output: "(a = 1 and (b = 2 is true))",
	}, {
		input:  "select * from a where (a=1 and b=2) is true",
		output: "((a = 1 and b = 2) is true)",
	}}
	for _, tcase := range validSQL {
		tree, err := Parse(tcase.input)
		if err != nil {
			t.Error(err)
			continue
		}
		expr := readable(tree.(*Select).Where.Expr)
		if expr != tcase.output {
			t.Errorf("Parse: \n%s, want: \n%s", expr, tcase.output)
		}
	}
}

func TestComparisonPrecedence(t *testing.T) {
	// https://github.com/dolthub/dolt/issues/12345
	validSQL := []struct {
		name   string
		input  string
		output string
	}{{
		name:   "chained scalar comparison",
		input:  "select 1 < 2 = 1",
		output: "((1 < 2) = 1)",
	}, {
		name:   "chained row comparison",
		input:  "select (1,2) < (1,3) = 1",
		output: "(((1, 2) < (1, 3)) = 1)",
	}, {
		name:   "grouped scalar comparison",
		input:  "select (1 < 2) = 1",
		output: "((1 < 2) = 1)",
	}, {
		name:   "grouped row comparison",
		input:  "select ((1,2) < (1,3)) = 1",
		output: "(((1, 2) < (1, 3)) = 1)",
	}, {
		name:   "chained comparisons group left",
		input:  "select 2 = 1 = 0",
		output: "((2 = 1) = 0)",
	}, {
		name:   "chained comparisons of one operator",
		input:  "select 1 < 2 < 3 < 4",
		output: "(((1 < 2) < 3) < 4)",
	}, {
		name:   "null-safe equal in chain",
		input:  "select 1 <=> null = 0",
		output: "((1 <=> null) = 0)",
	}, {
		name:   "chain under not",
		input:  "select not 1 < 2 = 1",
		output: "(not ((1 < 2) = 1))",
	}, {
		name:   "chain under and",
		input:  "select 1 < 2 = 1 and 0",
		output: "(((1 < 2) = 1) and 0)",
	}, {
		name:   "chain under or",
		input:  "select 0 = 1 or 1 < 2 = 1",
		output: "((0 = 1) or ((1 < 2) = 1))",
	}, {
		name:   "chain with arithmetic operand",
		input:  "select 1 + 1 < 3 = 1",
		output: "(((1 + 1) < 3) = 1)",
	}, {
		name:   "chain with subquery operand",
		input:  "select (select 1) < 2 = 1",
		output: "(((select 1) < 2) = 1)",
	}, {
		name:   "row chain with mismatched arity parses",
		input:  "select (1,2) < (1,3) = (1,2)",
		output: "(((1, 2) < (1, 3)) = (1, 2))",
	}}
	for _, tcase := range validSQL {
		t.Run(tcase.name, func(t *testing.T) {
			tree, err := Parse(tcase.input)
			if err != nil {
				t.Fatal(err)
			}
			expr := fmtComparison(tree.(*Select).SelectExprs[0].(*AliasedExpr).Expr)
			if expr != tcase.output {
				t.Errorf("Parse: \n%s, want: \n%s", expr, tcase.output)
			}
		})
	}
}

// fmtComparison formats |node| like readable, parenthesizing comparisons,
// AND, OR, and NOT so grouping is visible.
func fmtComparison(node Expr) string {
	switch node := node.(type) {
	case *ComparisonExpr:
		return fmt.Sprintf("(%s %s %s)", fmtComparison(node.Left), node.Operator, fmtComparison(node.Right))
	case *AndExpr:
		return fmt.Sprintf("(%s and %s)", fmtComparison(node.Left), fmtComparison(node.Right))
	case *OrExpr:
		return fmt.Sprintf("(%s or %s)", fmtComparison(node.Left), fmtComparison(node.Right))
	case *NotExpr:
		return fmt.Sprintf("(not %s)", fmtComparison(node.Expr))
	}
	return readable(node)
}

func fmtSetOp(s SelectStatement) string {
	switch s := s.(type) {
	case *SetOp:
		return fmt.Sprintf("(%s %s %s)", fmtSetOp(s.Left), s.Type, fmtSetOp(s.Right))
	case *Select:
		return String(s)
	case *ParenSelect:
		return String(s)
	}
	return ""
}

func TestSetOperatorPrecedence(t *testing.T) {
	validSQL := []struct {
		input  string
		output string
	}{
		{
			input:  "select 1 union select 2 union select 3 union select 4",
			output: "(((select 1 union select 2) union select 3) union select 4)",
		},
		{
			input:  "select 1 intersect select 2 intersect select 3 intersect select 4",
			output: "(((select 1 intersect select 2) intersect select 3) intersect select 4)",
		},
		{
			input:  "select 1 except select 2 except select 3 except select 4",
			output: "(((select 1 except select 2) except select 3) except select 4)",
		},

		{
			input:  "select 1 union select 2 intersect select 3 except select 4",
			output: "((select 1 union (select 2 intersect select 3)) except select 4)",
		},
		{
			input:  "select 1 union select 2 except select 3 intersect select 4",
			output: "((select 1 union select 2) except (select 3 intersect select 4))",
		},

		{
			input:  "select 1 intersect select 2 union select 3 except select 4",
			output: "(((select 1 intersect select 2) union select 3) except select 4)",
		},
		{
			input:  "select 1 intersect select 2 except select 3 union select 4",
			output: "(((select 1 intersect select 2) except select 3) union select 4)",
		},

		{
			input:  "select 1 except select 2 intersect select 3 union select 4",
			output: "((select 1 except (select 2 intersect select 3)) union select 4)",
		},
		{
			input:  "select 1 except select 2 union select 3 intersect select 4",
			output: "((select 1 except select 2) union (select 3 intersect select 4))",
		},

		{
			input:  "(table a) union (table b)",
			output: "((select * from a) union (select * from b))",
		},
		{
			input:  "(table a) intersect (table b)",
			output: "((select * from a) intersect (select * from b))",
		},
		{
			input:  "(table a) except (table b)",
			output: "((select * from a) except (select * from b))",
		},
		{
			input:  "select 1 intersect (select 2 union select 3)",
			output: "(select 1 intersect (select 2 union select 3))",
		},

	}
	for _, tcase := range validSQL {
		tree, err := Parse(tcase.input)
		if err != nil {
			t.Error(err)
			continue
		}
		expr := fmtSetOp(tree.(SelectStatement))
		if expr != tcase.output {
			t.Errorf("Parse: \n%s, want: \n%s", expr, tcase.output)
		}
	}
}