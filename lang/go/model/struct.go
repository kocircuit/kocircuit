package model

import (
	"bytes"
	"fmt"

	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/gate"
	. "github.com/kocircuit/kocircuit/lang/go/kit/hash"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

type Structure interface {
	StructureField() []*GoField
}

func StructureMonadicField(s Structure) *GoField {
	for _, f := range s.StructureField() {
		if f.IsMonadic() {
			return f
		}
	}
	return nil
}

type GoStructure interface {
	GoType
	Structure
}

// GoStruct describes the defition of a Go composite type.
type GoStruct struct {
	ID      string     `ko:"name=id"`
	Comment string     `ko:"name=comment"`
	Field   []*GoField `ko:"name=field"`
}

func NewGoStruct(ff ...*GoField) *GoStruct {
	mix := []string{"struct"}
	for _, f := range ff {
		mix = append(mix, f.Name, f.Type.TypeID())
	}
	s := &GoStruct{
		ID: Mix(mix...),
		// Field: ff, // go compiler gets confused with gate.Field
	}
	s.Field = ff
	return s
}

func (strct *GoStruct) TypeID() string { return strct.ID }

func (strct *GoStruct) Len() int { return len(strct.Field) }

func (strct *GoStruct) StructureField() []*GoField {
	return strct.Field
}

func (strct *GoStruct) Doc() string { return strct.Comment }

func (strct *GoStruct) String() string { return Sprint(strct) }

func (strct *GoStruct) Sketch(ctx *GoSketchCtx) interface{} {
	k := K{}
	for _, f := range strct.Field {
		k[f.KoName()] = f.Type.Sketch(ctx)
	}
	return k
}

// GoField describes a go struct field.
type GoField struct {
	Comment string   `ko:"name=comment" ctx:"expand"`
	Name    string   `ko:"name=name"` // go name
	Type    GoType   `ko:"name=type"`
	Tag     []*GoTag `ko:"name=tag"` // ko name and monadicity
}

func BuildGoField(span *Span, name string, typ GoType, monadic bool) *GoField {
	comment := ""
	if span != nil {
		comment = span.SourceLine()
	}
	return &GoField{
		Comment: comment,
		Name:    GoNameFor(name), // ko name -> go name
		Type:    typ,
		Tag:     KoTags(name, monadic),
	}
}

func SelectFieldByGoName(name string, field []*GoField) *GoField {
	return FilterFieldByGoName(name, field)[0]
}

func FilterFieldByGoName(name string, field []*GoField) (filtered []*GoField) {
	for _, f := range field {
		if f.Name == name {
			filtered = append(filtered, f)
		}
	}
	return
}

func FilterFieldByKoName(name string, field []*GoField) (filtered []*GoField) {
	for _, f := range field {
		if f.KoName() == name {
			filtered = append(filtered, f)
		}
	}
	return
}

// GoTag models a key value in a `ko:"key=value"` field tag
type GoTag struct {
	Key   string `ko:"name=key"`
	Value string `ko:"name=value"`
}

func KoTags(name string, monadic bool) []*GoTag {
	if monadic {
		return []*GoTag{{Key: "name", Value: name}, {Key: Monadic}}
	} else {
		return []*GoTag{{Key: "name", Value: name}}
	}
}

func (strct *GoStruct) FieldNames() []string {
	n := make([]string, len(strct.Field))
	for i, f := range strct.Field {
		n[i] = f.Name
	}
	return n
}

func (strct *GoStruct) Tag() []*GoTag { return nil }

func (strct *GoStruct) SelectKoField(field string) (goField string, goType GoType) {
	for _, f := range strct.Field {
		if f.KoName() == field {
			return f.Name, f.Type
		}
	}
	return "", nil
}

// RenderDef returns a type definition of the form: struct{...}
func (strct *GoStruct) RenderDef(fileCtx GoFileContext) string {
	if len(strct.Field) == 0 {
		return "struct{}"
	}
	var w bytes.Buffer
	fmt.Fprintf(&w, "struct{\n")
	for _, f := range strct.Field {
		fmt.Fprintf(&w, "\t%s %s %s\n", f.Name, f.Type.RenderRef(fileCtx), f.RenderTag(f.Type.Tag()))
	}
	fmt.Fprintf(&w, "}")
	return w.String()
}

// RenderRef returns a type reference of the form: struct{...}
func (strct *GoStruct) RenderRef(fileCtx GoFileContext) string {
	var w bytes.Buffer
	fmt.Fprintf(&w, "struct{")
	for i, f := range strct.Field {
		if i > 0 {
			fmt.Fprint(&w, "; ")
		}
		fmt.Fprintf(&w, "%s %s %s", f.Name, f.Type.RenderRef(fileCtx), f.RenderTag(f.Type.Tag()))
	}
	fmt.Fprintf(&w, "}")
	return w.String()
}

// RenderZero returns a zero value of the form: struct{...}{}
func (strct *GoStruct) RenderZero(fileCtx GoFileContext) string {
	return fmt.Sprintf("%s{}", strct.RenderRef(fileCtx))
}

func (goField *GoField) RenderTag(typeTag []*GoTag) string {
	name := goField.KoName()
	koTags := append(append([]*GoTag{}, goField.Tag...), typeTag...)
	return RenderTagLabel(
		[]*GoTagLabel{
			{Label: "ko", Tag: koTags},
			{Label: "json", Tag: []*GoTag{{Value: name}}},
		},
	)
}

func RenderTagLabel(tagLabel []*GoTagLabel) string {
	var w bytes.Buffer
	w.WriteString("`")
	for i, l := range tagLabel {
		if i > 0 {
			w.WriteString(" ")
		}
		w.WriteString(l.Render())
	}
	w.WriteString("`")
	return w.String()
}

type GoTagLabel struct {
	Label string   `ko:"name=label"`
	Tag   []*GoTag `ko:"name=tag"`
}

func (gtl *GoTagLabel) Render() string {
	var w bytes.Buffer
	for i, t := range gtl.Tag {
		if i > 0 {
			fmt.Fprint(&w, ",")
		}
		if t.Key != "" {
			fmt.Fprintf(&w, "%s=%s", t.Key, t.Value)
		} else {
			fmt.Fprintf(&w, "%s", t.Value)
		}
	}
	return fmt.Sprintf("%s:%q", gtl.Label, w.String())
}

func (goField *GoField) KoName() string {
	for _, tag := range goField.Tag {
		if tag.Key == "name" {
			return tag.Value
		}
	}
	return ""
}

func (goField *GoField) IsMonadic() bool {
	for _, tag := range goField.Tag {
		if tag.Key == Monadic {
			return true
		}
	}
	return false
}
