package model

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/kit/hash"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

type ReserveShaper struct {
	Shaping  `ko:"name=shaping"`
	Solution Shaper `ko:"name=solution"`
}

func (sh *ReserveShaper) String() string {
	return Sprint(sh)
}

func (sh *ReserveShaper) ShaperID() string {
	return sh.Solution.ShaperID()
}

func (sh *ReserveShaper) CircuitEffect() *GoCircuitEffect {
	return sh.Solution.CircuitEffect()
}

func (sh *ReserveShaper) RenderExprShaping(fileCtx GoFileContext, expr GoExpr) string {
	return sh.Solution.RenderExprShaping(fileCtx, expr)
}

func (sh *ReserveShaper) ShaperVerify(span *Span) error {
	return sh.Solution.ShaperVerify(span)
}

// AssignCache
type AssignCache struct {
	Parent []*AssignCache    `ko:"name=parent"`
	Seen   map[string]Shaper `ko:"name=seen"`
}

func (cache *AssignCache) Splay() Tree {
	return Quote{"AssignCache"}
}

func (cache *AssignCache) SheathID() *ID {
	return nil
}

func (cache *AssignCache) SheathLabel() *string {
	return nil
}

func (cache *AssignCache) SheathString() *string {
	return nil
}

func RefineAssignCache(span *Span, cache *AssignCache) *Span {
	return span.Refine(cache)
}

// AssignCache fulfills GoCtx interface.
func (cache *AssignCache) AssignCache() *AssignCache { return cache }

func CompressCacheUnion(parent ...*AssignCache) *AssignCache {
	return AssignCacheUnion(parent...).Compress()
}

func AssignCacheUnion(parent ...*AssignCache) *AssignCache {
	nonNilParents := []*AssignCache{}
	for _, cache := range parent {
		if cache != nil {
			nonNilParents = append(nonNilParents, cache)
		}
	}
	return &AssignCache{
		Parent: nonNilParents,
		Seen:   map[string]Shaper{},
	}
}

func (cache *AssignCache) Compress() *AssignCache {
	return &AssignCache{
		Seen: cache.compress(map[string]Shaper{}),
	}
}

func (cache *AssignCache) compress(into map[string]Shaper) map[string]Shaper {
	if cache == nil {
		return into
	}
	for _, parent := range cache.Parent {
		parent.compress(into)
	}
	for p, q := range cache.Seen {
		into[p] = q
	}
	return into
}

func assignCacheID(from, to GoType) string {
	return Mix(from.TypeID(), to.TypeID())
}

func (cache *AssignCache) Lookup(from, to GoType) Shaper {
	if cache == nil {
		return nil
	} else if result := cache.Seen[assignCacheID(from, to)]; result != nil {
		return result
	} else {
		for _, parent := range cache.Parent {
			if result := parent.Lookup(from, to); result != nil {
				return result
			}
		}
	}
	return nil
}

func (cache *AssignCache) Reserve(span *Span, from, to GoType) *ReserveShaper {
	reserve := &ReserveShaper{
		Shaping: Shaping{Origin: span, From: from, To: to},
	}
	cache.Seen[assignCacheID(from, to)] = reserve
	return reserve
}

func (cache *AssignCache) Unreserve(reserve *ReserveShaper) {
	id := assignCacheID(reserve.Shaping.From, reserve.Shaping.To)
	if r0 := cache.Seen[id]; r0 != reserve {
		panic("o")
	}
	delete(cache.Seen, id)
}
