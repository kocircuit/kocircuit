package model

import (
	"testing"

	. "github.com/kocircuit/kocircuit/lang/go/kit/subset"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

func TestAssign(t *testing.T) {
	for i, test := range testAssign {
		if !test.Enable {
			continue
		}
		bridge, _, err := Assign(NewCacheSpan(nil), test.From, test.To)
		if test.Error {
			if err == nil {
				t.Errorf("test %d: expecting error, but got none", i)
				continue
			}
		} else {
			if err != nil {
				t.Errorf("test %d: assignment error (%v)", i, err)
				continue
			}
			if bridge.Shape(test.From).TypeID() != test.To.TypeID() {
				t.Errorf("test %d: bridge %v is broken", i, bridge)
			}
			if !IsSubset(test.Bridge, bridge) {
				t.Errorf("test %d: expecting bridge %v, got %v", i, Sprint(test.Bridge), Sprint(bridge))
				continue
			}
		}
	}
}

var testAssign = []struct {
	Enable bool
	From   GoType
	To     GoType
	Bridge Shaper
	Error  bool
}{
	{Enable: true,
		From: GoUint32,
		To:   GoUint64,
		Bridge: &ReserveShaper{
			Shaping: Shaping{From: GoUint32, To: GoUint64},
			Solution: &ConvertTypeShaper{
				Shaping: Shaping{From: GoUint32, To: GoUint64},
			},
		},
	},
	{Enable: true,
		From: NewGoIntegerNumber(33),
		To:   GoInt64,
		Bridge: &ReserveShaper{
			Shaping: Shaping{
				Origin: nil,
				From:   NewGoIntegerNumber(33),
				To:     GoInt64,
			},
			Solution: &NumberShaper{
				Shaping: Shaping{
					Origin: nil,
					From:   NewGoIntegerNumber(33),
					To:     GoInt64,
				},
			},
		},
	},
	{Enable: true,
		From: NewGoNeverNilPtr(NewGoStruct()),
		To:   NewGoNeverNilPtr(NewGoStruct()),
		Bridge: &ReserveShaper{
			Shaping: Shaping{
				Origin: nil,
				From:   NewGoNeverNilPtr(NewGoStruct()),
				To:     NewGoNeverNilPtr(NewGoStruct()),
			},
			Solution: &ReShaper{
				Shaping: Shaping{
					Origin: nil,
					From:   NewGoNeverNilPtr(NewGoStruct()),
					To:     NewGoNeverNilPtr(NewGoStruct()),
				},
			},
		},
	},
}
