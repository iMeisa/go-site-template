package funcs

import (
	"github.com/iMeisa/go-site-template/internal/tools"
	"html/template"
)

var Functions = template.FuncMap{
	"add":       add,
	"contains":  tools.Contains,
	"makeRange": makeRange,
	"subtract":  subtract,
}

func add(num1, num2 int) int {
	return num1 + num2
}

func makeRange(count int) []int {
	var seq []int

	for i := 0; i < count; i++ {
		seq = append(seq, i)
	}

	return seq
}

func subtract(num1, num2 int) int {
	return num1 - num2
}
