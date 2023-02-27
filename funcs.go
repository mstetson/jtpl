package main

import (
	"fmt"
	"html/template"
	"strings"
)

var tplFuncs = template.FuncMap{
	"csv": csvEscaper,
}

func csvEscaper(args ...interface{}) string {
	s := fmt.Sprint(args...)
	if !strings.ContainsAny(s, ",\"\r\n") {
		return s
	}
	s = strings.Replace(s, "\"", "\"\"", -1)
	return "\"" + s + "\""
}
