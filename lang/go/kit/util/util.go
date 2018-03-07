package util

import (
	"sort"
)

func OptString(opt *string, otherwise string) string {
	if opt != nil {
		return *opt
	} else {
		return otherwise
	}
}

func OptInt64(opt *int64, otherwise int64) int64 {
	if opt != nil {
		return *opt
	} else {
		return otherwise
	}
}

func OptInt(opt *int, otherwise int) int {
	if opt != nil {
		return *opt
	} else {
		return otherwise
	}
}

func PtrInterface(i interface{}) *interface{} { return &i }

func PtrString(s string) *string { return &s }

func PtrBool(b bool) *bool { return &b }

func SortStringTable(ss [][]string) [][]string {
	sort.Sort(sortStringTable(ss))
	return ss
}

type sortStringTable [][]string

func (ss sortStringTable) Len() int { return len(ss) }

func (ss sortStringTable) Less(i, j int) bool { return lessStrings(ss[i], ss[j]) }

func lessStrings(u, v []string) bool {
	if len(u) != len(v) {
		panic("o")
	}
	for i := range u {
		switch {
		case u[i] < v[i]:
			return true
		case u[i] > v[i]:
			return false
		}
	}
	return false
}

func (ss sortStringTable) Swap(i, j int) { ss[i], ss[j] = ss[j], ss[i] }
