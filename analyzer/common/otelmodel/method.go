package otelmodel

import (
	"cmp"
	"slices"
)

func MatchPattern(pattern *PatternTree, trace *TraceDataTree) bool {
	// 对trace中的每个子树，看看是否匹配pattern
	pRoot, tRoot := pattern.Root, trace.Root
	if pRoot == nil {
		return true
	}
	var dfs func(t *TraceDataNode) bool
	dfs = func(t *TraceDataNode) bool {
		if t == nil {
			return false
		}
		if matchPatternTreeNode(pRoot, t) {
			// 匹配成功
			return true
		}
		for _, child := range t.Children {
			if dfs(child) {
				return true
			}
		}
		return false
	}
	return dfs(tRoot)

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

func GetCybertwinInfoFromTraceData(trace *TraceData) string {
	resourceSpans := trace.Trace.ResourceSpans
	if len(resourceSpans) == 0 {
		return ""
	}
	scopeSpans := resourceSpans[0].ScopeSpans
	if len(scopeSpans) == 0 {
		return ""
	}
	spans := scopeSpans[0].Spans
	if len(spans) == 0 {
		return ""
	}
	// 这里假设cybertwin信息在第一个span中
	attributes := spans[0].Attributes
	for _, attr := range attributes {
		if attr.Key == "cybertwin_id" {
			return attr.Value.StringValue
		}
	}
	return ""
}
