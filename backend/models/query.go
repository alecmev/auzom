package models

import (
	"strings"

	sq "github.com/Masterminds/squirrel"
)

var psqlbuilder = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type QueryModifier struct {
	Offset uint64
	Limit  uint64
	Filter map[string]interface{}
	Sort   []string
}

func (modifier *QueryModifier) SetColumnFilter(
	column string, value interface{},
) {
	if modifier.Filter == nil {
		modifier.Filter = make(map[string]interface{})
	}

	modifier.Filter[column] = value
}

// TODO: We should be able to embed QueryBase when decoding url data. Sadly this
//   is currently Not supported by goji/param. Resorted to copying for now.
type QueryBase struct {
	Offset uint64            `param:"offset"`
	Limit  uint64            `param:"count"`
	Filter map[string]string `param:"filter"`
	Sort   string            `param:"sort"` // prepend "-" for descending
}

// TODO: check if sortAllowed and filterAllowed can be moved to struct tags.
func NewQueryModifier(
	queryBase QueryBase, filterAllowed []string, sortAllowed []string,
) *QueryModifier {
	var result *QueryModifier
	if queryBase.Offset > 0 {
		result = &QueryModifier{Offset: queryBase.Offset}
	}

	if queryBase.Limit > 0 {
		if result == nil {
			result = &QueryModifier{}
		}

		result.Limit = queryBase.Limit
	}

	for col, val := range queryBase.Filter {
		for _, allowedcol := range filterAllowed {
			if allowedcol == col {
				if result == nil {
					result = &QueryModifier{}
				}

				result.SetColumnFilter(allowedcol, val)
				break
			}
		}
	}

	cols := strings.Split(queryBase.Sort, ",")
	for _, col := range cols {
		if len(col) > 0 {
			order := " ASC"
			if col[:1] == "-" {
				order = " DESC"
				col = col[1:]
			}

			for _, allowedcol := range sortAllowed {
				if allowedcol == col {
					if result == nil {
						result = &QueryModifier{}
					}

					result.Sort = append(result.Sort, allowedcol+order)
					break
				}
			}
		}
	}

	return result
}

func (modifier *QueryModifier) ToSql(
	table string, columns ...string,
) (string, []interface{}, error) {
	q := psqlbuilder.Select(columns...).From(table)
	if modifier == nil {
		return q.ToSql()
	}

	if len(modifier.Filter) > 0 {
		eq := sq.Eq{}
		for column, value := range modifier.Filter {
			if value == "\x00" {
				eq[column] = nil
			} else {
				eq[column] = value
			}
		}

		q = q.Where(eq)
	}

	if len(modifier.Sort) > 0 {
		q = q.OrderBy(modifier.Sort...)
	}

	if modifier.Offset > 0 {
		q = q.Offset(modifier.Offset)
	}

	if modifier.Limit > 0 {
		q = q.Limit(modifier.Limit)
	}

	return q.ToSql()
}
