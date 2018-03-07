package model

func ExtractGroupPaths(types []*GoAlias) (groups []GoGroupPath) {
	seen := map[GoGroupPath]bool{}
	for _, t := range types {
		if !seen[t.Address.GroupPath] {
			groups = append(groups, t.Address.GroupPath)
			seen[t.Address.GroupPath] = true
		}
	}
	return
}

// ExtractDuctTypes does not return duplicates.
func ExtractDuctTypes(circuit []*GoCircuit) (duct []*GoAlias) {
	ctx := &extractAliasCtx{
		Seen: map[string]bool{},
	}
	for _, c := range circuit {
		for _, t := range c.DuctType() {
			duct = append(duct, ctx.Extract(t)...)
		}
	}
	return
}

func ExtractAliases(top []GoType) (all []*GoAlias) {
	ctx := &extractAliasCtx{
		Seen: map[string]bool{},
	}
	for _, start := range top {
		all = append(all, ctx.Extract(start)...)
	}
	return
}

type extractAliasCtx struct {
	Seen map[string]bool `ko:"name=seen"`
}

func (ctx *extractAliasCtx) Extract(start GoType) (all []*GoAlias) {
	for {
		switch u := start.(type) {
		case *GoPtr:
			start = u.Elem
		case *GoNeverNilPtr:
			start = u.Elem
		case *GoAlias:
			if u.Address.IsHereditary() {
				return nil
			}
			if ctx.Seen[u.TypeID()] {
				return nil //already visiting
			} else {
				ctx.Seen[u.TypeID()] = true
			}
			return append([]*GoAlias{u}, ctx.Extract(u.For)...) // deep extraction
		case *GoSlice:
			return ctx.Extract(u.Elem)
		case *GoArray:
			return ctx.Extract(u.Elem)
		case *GoStruct:
			for _, f := range u.Field {
				all = append(all, ctx.Extract(f.Type)...)
			}
			return all
		case GoReal:
			start = u.Real()
		default:
			return nil
		}
	}
}
