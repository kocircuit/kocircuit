package macros

import (
	"sync"

	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
)

type memory struct {
	sync.Mutex
	seen map[string]Symbol
}

func newMemory() *memory {
	return &memory{seen: map[string]Symbol{}}
}

func (m *memory) Remember(key string, value Symbol) Symbol {
	m.Lock()
	defer m.Unlock()
	old := m.seen[key]
	m.seen[key] = value
	if old == nil {
		return EmptySymbol{}
	} else {
		return old
	}
}

func (m *memory) Recall(key string) (Symbol, bool) {
	m.Lock()
	defer m.Unlock()
	value, found := m.seen[key]
	return value, found
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
	m := newMemory()
	return MakeStructSymbol(
		FieldSymbols{
			{
				Name:  "Remember",
				Value: MakeVarietySymbol(&evalRememberMacro{m}, nil),
			},
			{
				Name:  "Recall",
				Value: MakeVarietySymbol(&evalRecallMacro{m}, nil),
			},
		},
	), nil, nil
}

type evalRememberMacro struct {
	memory *memory
}

func (m evalRememberMacro) MacroID() string { return m.Help() }

func (m evalRememberMacro) Label() string { return "remember" }

func (m evalRememberMacro) MacroSheathString() *string { return PtrString("Remember") }

func (m evalRememberMacro) Help() string { return "Remember" }

func (m evalRememberMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	a := arg.(*StructSymbol)
	return m.memory.Remember(a.Walk("key").Hash(), a.Walk("value")), nil, nil
}

type evalRecallMacro struct {
	memory *memory
}

func (m evalRecallMacro) MacroID() string { return m.Help() }

func (m evalRecallMacro) Label() string { return "recall" }

func (m evalRecallMacro) MacroSheathString() *string { return PtrString("Recall") }

func (m evalRecallMacro) Help() string { return "Recall" }

func (m evalRecallMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	a := arg.(*StructSymbol)
	keyHash := a.Walk("key").Hash()
	if value, found := m.memory.Recall(keyHash); found {
		return value, nil, nil
	} else {
		return a.Walk("otherwise"), nil, nil
	}
}
