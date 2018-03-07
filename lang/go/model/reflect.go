package model

import (
	"reflect"
	"strings"
)

type goTypeCtx struct {
	address map[GoAddress]*GoAlias
}

func newGoTypeCtx() *goTypeCtx {
	return &goTypeCtx{address: make(map[GoAddress]*GoAlias)}
}

func GoInterfaceTypeAddress(v interface{}) *GoAddress {
	return GoTypeAddress(reflect.TypeOf(v))
}

func GoTypeAddress(t reflect.Type) *GoAddress {
	return &GoAddress{
		Name: t.Name(),
		GroupPath: GoGroupPath{
			Group: GoHereditaryPkgGroup,
			Path:  t.PkgPath(),
		},
	}
}

func GraftGoType(t reflect.Type) GoType {
	return newGoTypeCtx().GraftGoType(t)
}

func (ctx *goTypeCtx) GraftGoType(t reflect.Type) GoType {
	if goIsPublic(t) {
		a := GoTypeAddress(t)
		if alias, ok := ctx.address[*a]; ok {
			return alias
		}
		alias := NewAliasDeferType(a)
		ctx.address[*a] = alias
		alias.For = ctx.graftGoKind(t)
		return alias
	}
	return ctx.graftGoKind(t)
}

func goIsPublic(t reflect.Type) bool { return t.PkgPath() != "" && t.Name() != "" }

func (ctx *goTypeCtx) graftGoKind(t reflect.Type) GoType {
	switch t.Kind() {
	case reflect.Invalid:
		panic("invalid")
	case reflect.Bool:
		return NewGoBuiltin(t.Kind())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return NewGoBuiltin(t.Kind())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return NewGoBuiltin(t.Kind())
	case reflect.Float32, reflect.Float64:
		return NewGoBuiltin(t.Kind())
	case reflect.Complex64, reflect.Complex128:
		return NewGoBuiltin(t.Kind())
	case reflect.String:
		return NewGoBuiltin(t.Kind())
	case reflect.UnsafePointer:
		return NewGoBuiltin(t.Kind())
	case reflect.Chan:
		return NewGoChan(ctx.GraftGoType(t.Elem()))
	case reflect.Func:
		arg := make([]GoType, t.NumIn())
		returns := make([]GoType, t.NumOut())
		for i := 0; i < t.NumIn(); i++ {
			arg[i] = ctx.GraftGoType(t.In(i))
		}
		for i := 0; i < t.NumOut(); i++ {
			returns[i] = ctx.GraftGoType(t.Out(i))
		}
		return NewGoFunc(arg, returns, t.IsVariadic())
	case reflect.Interface:
		return NewGoInterface(t)
	case reflect.Map:
		return NewGoMap(ctx.GraftGoType(t.Key()), ctx.GraftGoType(t.Elem()))
	case reflect.Ptr:
		return NewGoPtr(ctx.GraftGoType(t.Elem()))
	case reflect.Array:
		return NewGoArray(t.Len(), ctx.GraftGoType(t.Elem()))
	case reflect.Slice:
		return NewGoSlice(ctx.GraftGoType(t.Elem()))
	case reflect.Struct:
		field := make([]*GoField, t.NumField())
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			field[i] = &GoField{Name: f.Name, Type: ctx.GraftGoType(f.Type), Tag: graftGoTag(f.Tag)}
		}
		return NewGoStruct(field...)
	}
	panic("unknown go kind")
}

func graftGoTag(tag reflect.StructTag) []*GoTag {
	goTag := []*GoTag{}
	for _, p := range strings.Split(tag.Get("ko"), ",") {
		kv := strings.SplitN(p, "=", 2)
		switch len(kv) {
		case 1:
			goTag = append(goTag, &GoTag{Key: kv[0], Value: ""})
		case 2:
			goTag = append(goTag, &GoTag{Key: kv[0], Value: kv[1]})
		default:
			panic("confusing tag")
		}
	}
	return goTag
}
