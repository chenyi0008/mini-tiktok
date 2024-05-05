package utils

import (
	"gorm.io/gorm/clause"
	"strings"
)

func BuildSQLStringForIDs(ids []int) clause.OrderBy {
	var placeholders []string
	var vars []interface{}
	for _, id := range ids {
		placeholders = append(placeholders, "?")
		vars = append(vars, id)
	}
	idsStr := strings.Join(placeholders, ", ")

	// 执行查询并排序
	return clause.OrderBy{
		Expression: clause.Expr{SQL: "FIELD(id, " + idsStr + ")", Vars: vars},
	}
}
