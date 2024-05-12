package utils

import (
	"fmt"
	"github.com/fatih/color"
	"gorm.io/gorm/clause"
	"strings"
	"time"
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

var t = time.Now()
var greenBg = color.New(color.BgBlue).SprintFunc()

func CalculateTime(s int) {
	if s != 0 {
		fmt.Println(greenBg(fmt.Sprintln("calculate time:", s, ":", time.Now().Sub(t).Milliseconds())))
	}
	s = s + 1
	t = time.Now()
}
