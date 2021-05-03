package mysqldialect

import (
	"github.com/uptrace/bun/dialect"
	"github.com/uptrace/bun/dialect/feature"
	"github.com/uptrace/bun/dialect/sqltype"
	"github.com/uptrace/bun/schema"
	"github.com/uptrace/bun/sqlfmt"
)

type Dialect struct {
	tables   *schema.Tables
	features feature.Feature
}

func New() *Dialect {
	d := new(Dialect)
	d.tables = schema.NewTables(d)
	d.features = feature.Backticks |
		feature.AutoIncrement |
		feature.DefaultPlaceholder |
		feature.ValuesRow |
		feature.DropTableCascade
	return d
}

func (d *Dialect) Name() string {
	return dialect.MySQL
}

func (d *Dialect) Tables() *schema.Tables {
	return d.tables
}

func (d *Dialect) Append(fmter sqlfmt.QueryFormatter, b []byte, value interface{}) []byte {
	return sqlfmt.Append(fmter, b, value)
}

func (d *Dialect) OnField(field *schema.Field) {
	field.DiscoveredSQLType = sqlType(field)
}

func (d *Dialect) OnTable(table *schema.Table) {}

func (d *Dialect) Features() feature.Feature {
	return d.features
}

func sqlType(field *schema.Field) string {
	switch field.DiscoveredSQLType {
	case sqltype.VarChar:
		return field.DiscoveredSQLType + "(255)"
	case sqltype.Timestamp:
		return "DATETIME"
	}
	return field.DiscoveredSQLType
}
