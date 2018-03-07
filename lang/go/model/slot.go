package model

import (
	"fmt"

	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
)

type Slot interface {
	Label() string
}

type RootSlot struct{}

func (RootSlot) Label() string { return "" }

type KnotFieldSlot struct {
	Field Field `ko:"name=field"`
	Index int   `ko:"name=index"`
}

func (slot KnotFieldSlot) Label() string {
	return fmt.Sprintf("%s_%d", slot.Field.Name, slot.Index)
}

type NameSlot struct {
	Name string `ko:"name=name"`
}

func (slot NameSlot) Label() string { return slot.Name }
