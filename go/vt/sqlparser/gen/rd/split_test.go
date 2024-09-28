package rd

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"strings"
	"testing"
)

var testY = `

%{
package sqlparser

import "fmt"
import "strings"
%}

%union {
  empty         struct{}
  statement     Statement
  selStmt       SelectStatement
  ddl           *DDL
}

// comment

%token <bytes> NEXT VALUE SHARE MODE
%token <bytes> SQL_NO_CACHE SQL_CACHE

%type <valTuple> row_tuple tuple_or_empty
%type <expr> tuple_expression

%start any_command

%%

any_command:
  command
  {
    setParseTree(yylex, $1)
  }
| command ';'
  {
    setParseTree(yylex, $1)
    statementSeen(yylex)
  }


command:
  select_statement
  {
    $$ = $1
  }
| values_select_statement
  {
    $$ = $1
  }
| stream_statement
| insert_statement
| /*empty*/
{
  setParseTree(yylex, nil)
}

set_opt:
  {
    $$ = nil
  }
| SET assignment_list
  {
    $$ = $2
  }
load_statement:
  LOAD DATA local_opt infile_opt ignore_or_replace_opt load_into_table_name opt_partition_clause charset_opt fields_opt lines_opt ignore_number_opt column_list_opt set_opt
  {
    $$ = &Load{Local: $3, Infile: $4, IgnoreOrReplace: $5, Table: $6, Partition: $7, Charset: $8, Fields: $9, Lines: $10, IgnoreNum: $11, Columns: $12, SetExprs: $13}
  }

from_or_using:
  FROM {}
| USING {}
| OTHER

select_statement:
  with_select order_by_opt limit_opt lock_opt into_opt
  {
    $1.SetOrderBy($2)
    $1.SetLimit($3)
    $1.SetLock($4)
    if err := $1.SetInto($5); err != nil {
    	yylex.Error(err.Error())
    	return 1
    }
    $$ = $1
  }
| SELECT comment_opt query_opts NEXT num_val for_from table_name
  {
    $$ = &Select{
    	Comments: Comments($2),
    	QueryOpts: $3,
    	SelectExprs: SelectExprs{Nextval{Expr: $5}},
    	From: TableExprs{&AliasedTableExpr{Expr: $7}},
    }
  }

join_condition:
  { $$ = JoinCondition{On: $2} }
  ON expression
  { $$ = JoinCondition{On: $2} }
| USING '(' column_list ')'
  { $$ = JoinCondition{Using: $3} }

func_parens_opt:
  /*empty*/
| openb closeb
`

type emptyWriter struct{}

var _ io.Writer = (*emptyWriter)(nil)

func (e emptyWriter) Write(p []byte) (n int, err error) {
	return 0, nil
}

func TestGen(t *testing.T) {
	inp, err := os.Open("../../sql.y")
	require.NoError(t, err)
	//inp := strings.NewReader(testY)
	g := newRecursiveGen(inp, emptyWriter{})
	err = g.init()
	require.NoError(t, err)
	err = g.gen()
	require.NoError(t, err)
	outfile, err := os.OpenFile("../../sql.gen.go", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	require.NoError(t, err)
	_, err = io.Copy(outfile, strings.NewReader(g.b.String()))
	require.NoError(t, err)
}

func TestLeftRecursion(t *testing.T) {
	y := `
%union {
  expr     Expr
}

%token <bytes> null_opt OR AND NOT TOKEN_A TOKEN_B tokens row_opt ROW
%token <expr> value_expr

tokens:
  TOKEN_A
| TOKEN_B

null_opt:
  {
    $$ = nil
  }
| NULL

row_opt:
  {}
| ROW
  {}

value_expr:
|  value_expr OR value_expr 
  {
    $$ = &Or{$1, $3}
  }
|  value_expr AND value_expr 
  {
    $$ = &And{$1, $3}
  }
|  NOT value_expr
  {
    $$ = &Not{$2}
  }
|  value_expr NOT AND value_expr 
  {
    $$ = &Or{&Not{$1}, &Not{$4}}
  }
|  value_expr NOT OR value_expr 
  {
    $$ = &And{&Not{$1}, &Not{$4}}
  }
|  NOT EXISTS value_expr
  {
    $$ = &Not{&Exists{$3}}
  }
`
	inp := strings.NewReader(y)
	g := newRecursiveGen(inp, emptyWriter{})
	err := g.init()
	require.NoError(t, err)
	err = g.gen()
	fmt.Println(g.b.String())
}

func TestConflicts(t *testing.T) {
	y := `
%union {
  expr     Expr
}

%token <bytes> null_opt OR AND NOT TOKEN_A TOKEN_B tokens row_opt ROW ID
%token <expr> value_expr condition expression

expression:
  condition
  {
    $$ = $1
  }
| expression AND expression
  {
    $$ = &AndExpr{Left: $1, Right: $3}
  }
| NOT expression
  {
    $$ = &NotExpr{Expr: $2}
  }
| value_expression
  {
    $$ = $1
  }
| DEFAULT default_opt
  {
    $$ = &Default{ColName: $2}
  }

value_expression:
|  value_expression OR value_expression 
  {
    $$ = &Or{$1, $3}
  }
|  value_expression AND value_expression 
  {
    $$ = &And{$1, $3}
  }
|  NOT value_expression
  {
    $$ = &Not{$2}
  }
|  value_expression NOT AND value_expression 
  {
    $$ = &Or{&Not{$1}, &Not{$4}}
  }
|  value_expression NOT OR value_expression 
  {
    $$ = &And{&Not{$1}, &Not{$4}}
  }
|  NOT EXISTS value_expr
  {
    $$ = &Not{&Exists{$3}}
  }

condition:
  value_expression compare value_expression
  {
    $$ = &ComparisonExpr{Left: $1, Operator: $2, Right: $3}
  }
| value_expression IN col_tuple
  {
    $$ = &ComparisonExpr{Left: $1, Operator: InStr, Right: $3}
  }

default_opt:
  /*empty*/
  {
    $$ = ""
  }
| openb ID closeb
  {
    $$ = string($2)
  }
`
	inp := strings.NewReader(y)
	g := newRecursiveGen(inp, emptyWriter{})
	err := g.init()
	require.NoError(t, err)
	err = g.gen()
	fmt.Println(g.b.String())
}

func TestSplit(t *testing.T) {
	inp := strings.NewReader(testY)
	cmp, err := split(inp)
	require.NoError(t, err)
	exp := &yaccFileContents{
		prefix: []string{
			"package sqlparser",
			"import \"fmt\"",
			"import \"strings\"",
		},
		goTypes: []goType{
			{name: "empty", typ: "struct{}"},
			{name: "statement", typ: "Statement"},
			{name: "selStmt", typ: "SelectStatement"},
			{name: "ddl", typ: "*DDL"},
		},
		yaccTypes: []yaccType{
			{name: "row_tuple", typ: "valTuple"},
			{name: "tuple_or_empty", typ: "valTuple"},
			{name: "tuple_expression", typ: "expr"},
		},
		tokens: nil,
		start:  "any_command",
		defs: map[string]*def{
			"any_command": {
				name: "any_command",
				rules: &rulePrefix{
					prefix: "",
					term: []*rule{
						{
							name: "command",
							set:  true,
							body: []string{"setParseTree(yylex, $1)"},
						},
						{
							name: "command ';'",
							set:  true,
							body: []string{"setParseTree(yylex, $1)", "statementSeen(yylex)"},
						},
					},
				},
			},
			"command": {
				name: "command",
				rules: &rulePrefix{
					prefix: "",
					term: []*rule{
						{
							name: "select_statement",
							set:  true,
							body: []string{"$$ = $1"},
						},
						{
							name: "values_select_statement",
							set:  true,
							body: []string{"$$ = $1"},
						},
						{
							name: "stream_statement",
						},
						{
							name: "insert_statement",
						},
						{
							name: "/*empty*/",
							set:  true,
							body: []string{"setParseTree(yylex, nil)"},
						},
					},
				},
			},
			"set_opt": {
				name: "set_opt",
				rules: &rulePrefix{
					prefix: "",
					term: []*rule{
						{
							name: "",
							set:  true,
							body: []string{"$$ = nil"},
						},
						{
							name: "SET assignment_list",
							set:  true,
							body: []string{"$$ = $2"},
						},
					},
				},
			},
			"load_statement": {
				name: "load_statement",
				rules: &rulePrefix{
					prefix: "",
					term: []*rule{
						{
							name: "LOAD DATA local_opt infile_opt ignore_or_replace_opt load_into_table_name opt_partition_clause charset_opt fields_opt lines_opt ignore_number_opt column_list_opt set_opt",
							set:  true,
							body: []string{
								"$$ = &Load{Local: $3, Infile: $4, IgnoreOrReplace: $5, Table: $6, Partition: $7, Charset: $8, Fields: $9, Lines: $10, IgnoreNum: $11, Columns: $12, SetExprs: $13}"},
						},
					},
				},
			},
			"from_or_using": {
				name: "from_or_using",
				rules: &rulePrefix{
					prefix: "",
					term: []*rule{
						{
							name: "FROM",
							set:  false,
						},
						{
							name: "USING",
							set:  false,
						},
						{
							name: "OTHER",
							set:  false,
						},
					},
				},
			},
			"select_statement": {
				name: "select_statement",
				rules: &rulePrefix{
					prefix: "",
					term: []*rule{
						{
							name: "with_select order_by_opt limit_opt lock_opt into_opt",
							set:  true,
							body: []string{
								"$1.SetOrderBy($2)",
								"$1.SetLimit($3)",
								"$1.SetLock($4)",
								"if err := $1.SetInto($5); err != nil {",
								"yylex.Error(err.Error())",
								"return 1",
								"}",
								"$$ = $1",
							},
						},
						{
							name: "SELECT comment_opt query_opts NEXT num_val for_from table_name",
							set:  true,
							body: []string{
								"$$ = &Select{",
								"Comments: Comments($2),",
								"QueryOpts: $3,",
								"SelectExprs: SelectExprs{Nextval{Expr: $5}},",
								"From: TableExprs{&AliasedTableExpr{Expr: $7}},",
								"}",
							},
						},
					},
				},
			},
			"join_condition": {
				name: "join_condition",
				rules: &rulePrefix{
					prefix: "",
					term: []*rule{
						{
							name: "",
							set:  false,
							body: []string{"$$ = JoinCondition{On: $2}"},
						},
						{
							name: "ON expression",
							set:  false,
							body: []string{"$$ = JoinCondition{On: $2}"},
						},
						{
							name: "USING '(' column_list ')'",
							set:  false,
							body: []string{"$$ = JoinCondition{Using: $3}"},
						},
					},
				},
			},
			"func_parens_opt": {
				name: "func_parens_opt",
				rules: &rulePrefix{
					prefix: "",
					term: []*rule{
						{
							name: "/*empty*/",
							set:  false,
						},
						{
							name: "openb closeb",
							set:  false,
						},
					},
				},
			},
		},
	}

	require.Equal(t, exp, cmp)
}
