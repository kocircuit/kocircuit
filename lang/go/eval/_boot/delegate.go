package boot

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
)

// delegate forwards invocation panics to its caller.
// errors in types are forwarded in the return value.
func (b *Booter) delegate(delegated *Span, vty *VarietySymbol, fields Fields) (result *BootStepResult, err error) {
	result, _, err := vty.Evoke(delegated, fields)
	if err != nil {
		return nil, err
	} else {
		returned, _, err := result.Select(delegated, Path{"returns"})
		if err != nil {
			return nil, delegated.Errorf(err, "expecting a structure with a returns field")
		}
		effect, _, err := result.Select(delegated, Path{"effect"})
		if err != nil {
			return nil, delegated.Errorf(err, "expecting a structure with an effect field")
		}
		return &BootStepResult{Returned: returned, Effect: effect}, nil
	}
}
