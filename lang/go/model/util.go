package model

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

type M map[string]interface{}

func ParseTmpl(tmpl string) *template.Template {
	return template.Must(template.New("").Funcs(goTemplateFunc).Parse(tmpl))
}

func ApplyTmpl(tmpl *template.Template, data interface{}) string {
	var w bytes.Buffer
	if err := tmpl.Execute(&w, data); err != nil {
		panic(err)
	}
	return w.String()
}

var goTemplateFunc = template.FuncMap{
	"join":    strings.Join,
	"lines":   Lines,
	"comment": StringComment,
	"indent":  Indent,
}

func Lines(text string) []string {
	var r []string
	scanner := bufio.NewScanner(bytes.NewBufferString(text))
	for scanner.Scan() {
		r = append(r, scanner.Text())
	}
	return r
}

func StringComment(c string) string {
	if c == "" {
		return ""
	}
	var w bytes.Buffer
	for _, l := range Lines(c) {
		fmt.Fprintln(&w, "//", l)
	}
	return w.String()
}

// Indent indents the text following a new line by a tab character.
func Indent(s interface{}) string {
	switch t := s.(type) {
	case string:
		return strings.Replace(t, "\n", "\n\t", -1)
	case stringer:
		return strings.Replace(t.String(), "\n", "\n\t", -1)
	}
	panic("o")
}

type stringer interface {
	String() string
}
