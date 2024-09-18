package rd

import (
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
	outfile, err := os.OpenFile("sql.gen.go", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	require.NoError(t, err)
	_, err = io.Copy(outfile, strings.NewReader(g.b.String()))
	require.NoError(t, err)
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
			{name: "valTuple", typ: "row_tuple"},
			{name: "valTuple", typ: "tuple_or_empty"},
			{name: "expr", typ: "tuple_expression"},
		},
		tokens: nil,
		start:  "any_command",
		defs: []*def{
			{
				name: "any_command",
				rules: []*rule{
					{
						name:  "command",
						start: true,
						body:  []string{"setParseTree(yylex, $1)"},
					},
					{
						name:  "command ';'",
						start: true,
						body:  []string{"setParseTree(yylex, $1)", "statementSeen(yylex)"},
					},
				},
			},
			{
				name: "command",
				rules: []*rule{
					{
						name:  "select_statement",
						start: true,
						body:  []string{"$$ = $1"},
					},
					{
						name:  "values_select_statement",
						start: true,
						body:  []string{"$$ = $1"},
					},
					{
						name: "stream_statement",
					},
					{
						name: "insert_statement",
					},
					{
						name:  "/*empty*/",
						start: true,
						body:  []string{"setParseTree(yylex, nil)"},
					},
				},
			},
			{
				name: "set_opt",
				rules: []*rule{
					{
						name:  "",
						start: true,
						body:  []string{"$$ = nil"},
					},
					{
						name:  "SET assignment_list",
						start: true,
						body:  []string{"$$ = $2"},
					},
				},
			},
			{
				name: "load_statement",
				rules: []*rule{
					{
						name:  "LOAD DATA local_opt infile_opt ignore_or_replace_opt load_into_table_name opt_partition_clause charset_opt fields_opt lines_opt ignore_number_opt column_list_opt set_opt",
						start: true,
						body: []string{
							"$$ = &Load{Local: $3, Infile: $4, IgnoreOrReplace: $5, Table: $6, Partition: $7, Charset: $8, Fields: $9, Lines: $10, IgnoreNum: $11, Columns: $12, SetExprs: $13}"},
					},
				},
			},
			{
				name: "select_statement",
				rules: []*rule{
					{
						name:  "with_select order_by_opt limit_opt lock_opt into_opt",
						start: true,
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
						name:  "SELECT comment_opt query_opts NEXT num_val for_from table_name",
						start: true,
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
	}

	require.Equal(t, exp, cmp)
}
