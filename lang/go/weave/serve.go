package weave

import (
	"fmt"
	"strings"

	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/lex"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/circuit/syntax"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
	. "github.com/kocircuit/kocircuit/lang/go/model"
)

func init() {
	RegisterGoMacro("Serve", new(GoServeMacro))
}

type GoServeMacro struct{}

func (m GoServeMacro) MacroID() string { return m.Help() }

func (m GoServeMacro) Label() string { return "serve" }

func (m GoServeMacro) MacroSheathString() *string { return PtrString("Serve") }

func (m GoServeMacro) Help() string {
	return GoInterfaceTypeAddress(m).String()
}

func (GoServeMacro) Invoke(span *Span, arg Arg) (Return, Effect, error) {
	_, logic, err := GoSelectSimplify(span, Path{"func"}, arg.(GoStructure))
	if err != nil {
		return nil, nil, span.Errorf(err, "serve expects a func argument")
	}
	varietal, ok := logic.(GoVarietal)
	if !ok {
		return nil, nil, span.Errorf(nil, "serve func must be a variety")
	}
	_, p, err := GoSelectSimplify(span, Path{"path"}, arg.(GoStructure))
	if err != nil {
		return nil, nil, span.Errorf(err, "serve expects a path argument")
	}
	endPath, ok := p.(*GoStringNumber)
	if !ok {
		return nil, nil, span.Errorf(nil, "serve path argument must be a string")
	}
	fullName := endPath.Value_
	if err = sanitizeServiceName(fullName); err != nil {
		return nil, nil, span.Errorf(err, "service name %q not valid", fullName)
	}
	chamber := fmt.Sprintf("serve_%s", SanitizeForPkgName(fullName))
	if valve, returns, effect, err := GoFix(span, chamber, varietal); err != nil {
		return nil, nil, err
	} else {
		return returns,
			effect.(*GoMacroEffect).AggregateDirective(ServeDirective(valve, fullName)),
			nil
	}
}

func sanitizeServiceName(n string) error {
	if l, err := LexifyString("sanitize_service_name.ko", n); err != nil {
		return err
	} else if _, remain, err := ParseRef(l); err != nil {
		return err
	} else if len(remain) > 0 {
		return err
	} else {
		return nil
	}
}

func ServeDirective(valve *GoValve, fullName string) *GoDirective {
	registerEvalGateAt := &GoAddress{
		GroupPath: GoGroupPath{
			Group: GoHereditaryPkgGroup,
			Path:  "github.com/kocircuit/kocircuit/lang/go/eval",
		},
		Name: "RegisterEvalGateAt",
	}
	return &GoDirective{
		Label:     fmt.Sprintf("serve_%s", fullName),
		GroupPath: GoGroupPath{Group: KoPkgGroup, Path: ""}, // main pkg
		Inject: &GoInitExpr{
			Line: []GoExpr{
				&GoCallExpr{
					Func: registerEvalGateAt,
					Arg: []GoExpr{
						&GoQuoteExpr{String: ""},
						&GoQuoteExpr{String: strings.Join([]string{"service", fullName}, ".")},
						&GoZeroExpr{GoType: valve.Real()},
					},
				},
			},
		},
	}
}
