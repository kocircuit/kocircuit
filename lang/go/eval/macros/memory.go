package macros

import (
	"fmt"
	"sync"

	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
)

type memory struct {
	origin *Span
	sync.Mutex
	seen map[string]*keyValue
}

type keyValue struct {
	key   Symbol
	value Symbol
}

func newMemory(origin *Span) *memory {
	return &memory{origin: origin, seen: map[string]*keyValue{}}
}

func (m *memory) ID() ID {
	return m.origin.SpanID()
}

func (m *memory) Remember(span *Span, key, value Symbol) Symbol {
	m.Lock()
	defer m.Unlock()
	keyHash := key.Hash(span)
	old := m.seen[keyHash]
	if IsEmptySymbol(value) {
		delete(m.seen, keyHash)
	} else {
		m.seen[keyHash] = &keyValue{key: key, value: value}
	}
	if old == nil {
		return EmptySymbol{}
	} else {
		return old.value
	}
}

func (m *memory) Recall(span *Span, key Symbol) Symbol {
	m.Lock()
	defer m.Unlock()
	if keyValue, found := m.seen[key.Hash(span)]; found {
		return keyValue.value
	} else {
		return EmptySymbol{}
	}
}

func (m *memory) Flush(span *Span) (Symbol, error) {
	m.Lock()
	defer m.Unlock()
	kvElems := make(Symbols, 0, len(m.seen))
	for _, kv := range m.seen {
		kvElems = append(kvElems,
			MakeStructSymbol(
				FieldSymbols{
					&FieldSymbol{Name: "key", Value: kv.key},
					&FieldSymbol{Name: "value", Value: kv.value},
				},
			),
		)
	}
	return MakeSeriesSymbol(span, kvElems)
}

func init() {
	RegisterEvalMacro("Memory", new(EvalMemoryMacro))
}

type EvalMemoryMacro struct{}

func (m EvalMemoryMacro) MacroID() string { return m.Help() }

func (m EvalMemoryMacro) Label() string { return "memory" }

func (m EvalMemoryMacro) MacroSheathString() *string { return PtrString("Memory") }

func (m EvalMemoryMacro) Help() string { return "Memory" }

func (EvalMemoryMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	m := newMemory(span)
	return MakeStructSymbol(
		FieldSymbols{
			{
				Name:  "name",
				Value: MakeBasicSymbol(span, m.ID().String()),
			},
			{
				Name:  "Remember",
				Value: MakeVarietySymbol(&evalRememberMacro{m}, nil),
			},
			{
				Name:  "Recall",
				Value: MakeVarietySymbol(&evalRecallMacro{m}, nil),
			},
			{
				Name:  "Flush",
				Value: MakeVarietySymbol(&evalFlushMacro{m}, nil),
			},
		},
	), nil, nil
}

// Remember
type evalRememberMacro struct {
	memory *memory
}

func (m evalRememberMacro) MacroID() string { return m.Help() }

func (m evalRememberMacro) Label() string { return "remember" }

func (m evalRememberMacro) MacroSheathString() *string { return PtrString("Remember") }

func (m evalRememberMacro) Help() string {
	return fmt.Sprintf("%v_Remember", m.memory.ID())
}

func (m evalRememberMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	a := arg.(*StructSymbol)
	return m.memory.Remember(span, a.Walk("key"), a.Walk("value")), nil, nil
}

// Recall
type evalRecallMacro struct {
	memory *memory
}

func (m evalRecallMacro) MacroID() string { return m.Help() }

func (m evalRecallMacro) Label() string { return "recall" }

func (m evalRecallMacro) MacroSheathString() *string { return PtrString("Recall") }

func (m evalRecallMacro) Help() string {
	return fmt.Sprintf("%v_Recall", m.memory.ID())
}

func (m evalRecallMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	a := arg.(*StructSymbol)
	if recalled := m.memory.Recall(span, a.Walk("key")); IsEmptySymbol(recalled) {
		return a.Walk("otherwise"), nil, nil
	} else {
		return recalled, nil, nil
	}
}

// Flush
type evalFlushMacro struct {
	memory *memory
}

func (m evalFlushMacro) MacroID() string { return m.Help() }

func (m evalFlushMacro) Label() string { return "flush" }

func (m evalFlushMacro) MacroSheathString() *string { return PtrString("Flush") }

func (m evalFlushMacro) Help() string {
	return fmt.Sprintf("%v_Flush", m.memory.ID())
}

func (m evalFlushMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	if flushed, err := m.memory.Flush(span); err != nil {
		return nil, nil, span.Errorf(err, "flushing key-value memory")
	} else {
		return flushed, nil, nil
	}
}
