package console

import (
	"fmt"

	"github.com/fatih/color"
)

type Console interface {
	Frame(...interface{}) string
	KoFile(...interface{}) string
	GoFile(...interface{}) string
	Ref(...interface{}) string
	Func(...interface{}) string
	Step(...interface{}) string
	String(...interface{}) string
	Field(...interface{}) string
	TaskStderr(...interface{}) string
}

var (
	Mono = &ansiPrinter{
		frame:      chromeless,
		koFile:     chromeless,
		goFile:     chromeless,
		ref:        chromeless,
		func_:      chromeless,
		step:       chromeless,
		string_:    chromeless,
		field:      chromeless,
		taskStderr: chromeless,
	}
	ANSI = &ansiPrinter{
		frame:      color.New(color.FgCyan).SprintFunc(),
		koFile:     color.New(color.FgGreen).SprintFunc(),
		goFile:     color.New(color.FgRed).SprintFunc(),
		ref:        color.New(color.Bold).SprintFunc(),
		func_:      color.New(color.Underline).SprintFunc(),
		step:       color.New(color.Bold).SprintFunc(),
		string_:    color.New(color.FgBlue).SprintFunc(),
		field:      color.New(color.Bold).SprintFunc(),
		taskStderr: color.New(color.FgGreen).SprintFunc(),
	}
)

type ansiPrinter struct {
	frame      func(...interface{}) string
	koFile     func(...interface{}) string
	goFile     func(...interface{}) string
	ref        func(...interface{}) string
	func_      func(...interface{}) string
	step       func(...interface{}) string
	string_    func(...interface{}) string
	field      func(...interface{}) string
	taskStderr func(...interface{}) string
}

func (ansi *ansiPrinter) Frame(a ...interface{}) string      { return ansi.frame(a...) }
func (ansi *ansiPrinter) KoFile(a ...interface{}) string     { return ansi.koFile(a...) }
func (ansi *ansiPrinter) GoFile(a ...interface{}) string     { return ansi.goFile(a...) }
func (ansi *ansiPrinter) Ref(a ...interface{}) string        { return ansi.ref(a...) }
func (ansi *ansiPrinter) Func(a ...interface{}) string       { return ansi.func_(a...) }
func (ansi *ansiPrinter) Step(a ...interface{}) string       { return ansi.step(a...) }
func (ansi *ansiPrinter) String(a ...interface{}) string     { return ansi.string_(a...) }
func (ansi *ansiPrinter) Field(a ...interface{}) string      { return ansi.field(a...) }
func (ansi *ansiPrinter) TaskStderr(a ...interface{}) string { return ansi.taskStderr(a...) }

func chromeless(a ...interface{}) string { return fmt.Sprint(a...) }
