package otelmodel

import (
	"cmp"
	"slices"
)

func MatchPattern(pattern *PatternTree, trace *TraceDataTree) bool {
	// 对trace中的每个子树，看看是否匹配pattern
	pRoot, tRoot := pattern.Root, trace.Root
	return matchPatternTreeNode(pRoot, tRoot)

}

func matchPatternTreeNode(patternRoot *PatternTreeNode, traceRoot *TraceDataNode) bool {
	if patternRoot == nil || traceRoot == nil {
		return true
	} else if len(patternRoot.Children) == 0 {
		return patternRoot.Value == traceRoot.Value && patternRoot.Service == traceRoot.Service
	} else if len(traceRoot.Children) == 0 {
		return false
	}
	if patternRoot.Value != traceRoot.Value || patternRoot.Service != traceRoot.Service {
		return false
	}
	// children排序
	slices.SortFunc(patternRoot.Children, func(a, b *PatternTreeNode) int {
		// 先value后service排序
		//cmp.Compare(a.Value, b.Value)
		if a.Value != b.Value {
			return cmp.Compare(a.Value, b.Value)
		} else {
			return cmp.Compare(a.Service, b.Service)
		}
	})
	slices.SortFunc(traceRoot.Children, func(a, b *TraceDataNode) int {
		// 先value后service排序
		if a.Value != b.Value {
			return cmp.Compare(a.Value, b.Value)
		} else {
			return cmp.Compare(a.Service, b.Service)
		}
	})
	childMatchFlag := true
	for _, pChild := range patternRoot.Children {
		match := false
		for _, tChild := range traceRoot.Children {
			if matchPatternTreeNode(pChild, tChild) {
				// 找到一个匹配的就返回true
				match = true
				break
			}
		}
		childMatchFlag = match && childMatchFlag
	}
	return childMatchFlag

}
