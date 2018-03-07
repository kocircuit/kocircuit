package model

import (
	"sort"
)

type GoDirectivePkg struct {
	GroupPath GoGroupPath    `ko:"name=groupPath"`
	Directive []*GoDirective `ko:"name=directive"`
}

type GoDirectiveIndex []*GoDirective

// PkgIndex assumes index is sorted.
func (index GoDirectiveIndex) PkgIndex() (pkgIndex []*GoDirectivePkg) {
	directivePkg := &GoDirectivePkg{}
	for _, directive := range index {
		if directive.GroupPath != directivePkg.GroupPath {
			directivePkg = &GoDirectivePkg{GroupPath: directive.GroupPath}
			pkgIndex = append(pkgIndex, directivePkg)
		}
		directivePkg.Directive = append(directivePkg.Directive, directive)
	}
	return
}

func (index GoDirectiveIndex) Len() int {
	return len(index)
}

func (index GoDirectiveIndex) Less(i, j int) bool {
	return index[i].GroupPath.String() < index[j].GroupPath.String()
}

func (index GoDirectiveIndex) Swap(i, j int) {
	index[i], index[j] = index[j], index[i]
}

func IndexDirectives(directive []*GoDirective) []*GoDirectivePkg {
	index := GoDirectiveIndex(DedupDirectives(directive))
	sort.Sort(index)
	return index.PkgIndex()
}

func DedupDirectives(directive []*GoDirective) (dedup []*GoDirective) {
	seen := map[string]bool{}
	for _, d := range directive {
		if id := d.DirectiveID(); !seen[id] {
			dedup = append(dedup, d)
			seen[id] = true
		}
	}
	return
}

type GoCircuitPkg struct {
	GroupPath GoGroupPath  `ko:"name=groupPath"`
	Circuit   []*GoCircuit `ko:"name=circuit"`
}

type GoCircuitIndex []*GoCircuit

// PkgIndex assumes index is sorted.
func (index GoCircuitIndex) PkgIndex() (pkgIndex []*GoCircuitPkg) {
	circuitPkg := &GoCircuitPkg{}
	for _, circuit := range index {
		if circuit.Valve.Address.GroupPath != circuitPkg.GroupPath {
			circuitPkg = &GoCircuitPkg{GroupPath: circuit.Valve.Address.GroupPath}
			pkgIndex = append(pkgIndex, circuitPkg)
		}
		circuitPkg.Circuit = append(circuitPkg.Circuit, circuit)
	}
	return
}

func (index GoCircuitIndex) Len() int {
	return len(index)
}

func (index GoCircuitIndex) Less(i, j int) bool {
	return index[i].Valve.Address.GroupPath.String() < index[j].Valve.Address.GroupPath.String()
}

func (index GoCircuitIndex) Swap(i, j int) {
	index[i], index[j] = index[j], index[i]
}

func IndexCircuits(circuit []*GoCircuit) []*GoCircuitPkg {
	index := GoCircuitIndex(DedupCircuits(circuit))
	sort.Sort(index)
	return index.PkgIndex()
}

func DedupCircuits(circuit []*GoCircuit) (dedup []*GoCircuit) {
	seen := map[string]bool{}
	for _, c := range circuit {
		if id := c.Valve.TypeID(); !seen[id] {
			dedup = append(dedup, c)
			seen[id] = true
		}
	}
	return
}

type GoTypePkg struct {
	GroupPath GoGroupPath   `ko:"name=groupPath"`
	File      []*GoTypeFile `ko:"name=file"`
}

type GoTypeFile struct {
	Filebase string     `ko:"name=filebase"`
	Type     []*GoAlias `ko:"name=type"`
}

type GoTypeIndex []*GoAlias

// PkgIndex assumes index is sorted.
func (index GoTypeIndex) PkgIndex() (pkgIndex []*GoTypePkg) {
	typePkg, typeFile := &GoTypePkg{}, &GoTypeFile{}
	for _, alias := range index {
		if alias.Address.GroupPath != typePkg.GroupPath {
			typePkg = &GoTypePkg{GroupPath: alias.Address.GroupPath}
			pkgIndex = append(pkgIndex, typePkg)
			typeFile = &GoTypeFile{}
		}
		if alias.Address.Filebase() != typeFile.Filebase {
			typeFile = &GoTypeFile{Filebase: alias.Address.Filebase()}
			typePkg.File = append(typePkg.File, typeFile)
		}
		typeFile.Type = append(typeFile.Type, alias)
	}
	return
}

func (index GoTypeIndex) Len() int {
	return len(index)
}

func (index GoTypeIndex) Less(i, j int) bool {
	igroup, jgroup := index[i].Address.GroupPath.String(), index[j].Address.GroupPath.String()
	switch {
	case igroup < jgroup:
		return true
	case igroup == jgroup:
		return index[i].Address.Filebase() < index[j].Address.Filebase()
	default: // igroup > jgroup:
		return false
	}
}

func (index GoTypeIndex) Swap(i, j int) {
	index[i], index[j] = index[j], index[i]
}

func IndexTypes(types []*GoAlias) []*GoTypePkg {
	index := GoTypeIndex(DedupTypes(types))
	sort.Sort(index)
	return index.PkgIndex()
}

func DedupTypes(types []*GoAlias) (dedup []*GoAlias) {
	seen := map[string]bool{}
	for _, t := range types {
		if id := t.TypeID(); !seen[id] {
			dedup = append(dedup, t)
			seen[id] = true
		}
	}
	return
}
