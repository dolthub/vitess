package sqlparser

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParser(t *testing.T) {
	tests := []struct {
		q   string
		exp Statement
		ok  bool
	}{
		{
			q:  "insert into xy values (0,'0', .0), (1,'1', 1.0)",
			ok: true,
		},
		{
			q:  "insert into xy (x,y,z) values (0,'0', .0), (1,'1', 1.0)",
			ok: true,
		},
		{
			q:  "insert into db.xy values (0,'0', .0), (1,'1', 1.0)",
			ok: true,
		},
		{
			q:  "select * from xy where x = 1",
			ok: true,
		},
		{
			q:  "select id from sbtest1 where id = 1000",
			ok: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.q, func(t *testing.T) {
			p := new(parser)
			p.tok = NewStringTokenizer(tt.q)
			res, ok := p.any_command(p.tok)
			require.Equal(t, tt.ok, ok)
			require.Equal(t, tt.exp, res)
		})
	}
}
