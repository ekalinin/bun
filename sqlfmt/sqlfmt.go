package sqlfmt

import (
	"github.com/uptrace/bun/internal"
)

type QueryAppender interface {
	AppendQuery(fmter QueryFormatter, b []byte) ([]byte, error)
}

type ColumnsAppender interface {
	AppendColumns(fmter QueryFormatter, b []byte) ([]byte, error)
}

//------------------------------------------------------------------------------

// Safe represents a safe SQL query.
type Safe string

var _ QueryAppender = (*Safe)(nil)

func (s Safe) AppendQuery(fmter QueryFormatter, b []byte) ([]byte, error) {
	return append(b, s...), nil
}

//------------------------------------------------------------------------------

// Ident represents a SQL identifier, for example, table or column name.
type Ident string

var _ QueryAppender = (*Ident)(nil)

func (s Ident) AppendQuery(fmter QueryFormatter, b []byte) ([]byte, error) {
	return AppendIdent(fmter, b, string(s)), nil
}

//------------------------------------------------------------------------------

type QueryWithArgs struct {
	Query string
	Args  []interface{}
}

var _ QueryAppender = QueryWithArgs{}

func SafeQuery(query string, args []interface{}) QueryWithArgs {
	if query != "" && args == nil {
		args = make([]interface{}, 0)
	}
	return QueryWithArgs{Query: query, Args: args}
}

func UnsafeIdent(ident string) QueryWithArgs {
	return QueryWithArgs{Query: ident}
}

func (q QueryWithArgs) IsZero() bool {
	return q.Query == "" && q.Args == nil
}

func (q QueryWithArgs) AppendQuery(fmter QueryFormatter, b []byte) ([]byte, error) {
	if q.Args == nil {
		return AppendIdent(fmter, b, q.Query), nil
	}
	return fmter.FormatQuery(b, q.Query, q.Args...), nil
}

func (q QueryWithArgs) Value() Safe {
	b, err := q.AppendQuery(defaultFmter, nil)
	if err != nil {
		return Safe(err.Error())
	}
	return Safe(internal.String(b))
}

//------------------------------------------------------------------------------

type QueryWithSep struct {
	QueryWithArgs
	Sep string
}

func SafeQueryWithSep(query string, args []interface{}, sep string) QueryWithSep {
	return QueryWithSep{
		QueryWithArgs: SafeQuery(query, args),
		Sep:           sep,
	}
}
