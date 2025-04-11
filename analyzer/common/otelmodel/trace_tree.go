package otelmodel

import (
	"errors"
	"github.com/sirupsen/logrus"
	"otelDemo/utils/otelutils"
)

type TraceDataTree struct {
	Root    *TraceDataNode
	TraceId string
}

type TraceDataNode struct {
	SpanId       string
	ParentSpanId string
	Value        string // 方法名
	Service      string // 服务名
	Children     []*TraceDataNode
}

func NewTraceDataNodeFromSpan(span *Span, scopeName string) *TraceDataNode {
	if span == nil {
		return nil
	}
	spanId, err := otelutils.DecodeOtelID(span.SpanID)
	if err != nil {
		logrus.Errorf("[NewTraceDataNodeFromSpan] - spanId = %s, err = %v", span.SpanID, err)
		return nil
	}
	parentSpanId := ""
	if span.ParentSpanID != "" {
		parentSpanId, err = otelutils.DecodeOtelID(span.ParentSpanID)
		if err != nil {
			logrus.Errorf("[NewTraceDataNodeFromSpan] - parentSpanId = %s, err = %v", span.ParentSpanID, err)
			return nil
		}
	}

	return &TraceDataNode{
		SpanId:       spanId,
		ParentSpanId: parentSpanId,
		Value:        span.Name,
		Service:      scopeName,
		Children:     make([]*TraceDataNode, 0),
	}
}

func TransferTraceData2Tree(traceData *TraceData) (*TraceDataTree, error) {
	if traceData == nil {
		return nil, errors.New("traceData is nil")
	}
	root := &TraceDataNode{}
	elseNode := make([]*TraceDataNode, 0)
	edges := make([]*TraceDataNode, 0)         // kind类型为"SPAN_KIND_CLIENT"为边
	nodeMap := make(map[string]*TraceDataNode) // key 为spanid，这个m可以方便我们构造tree
	traceId := ""
	var err error
	for _, resourceSpan := range traceData.Trace.ResourceSpans {
		for _, scopeSpan := range resourceSpan.ScopeSpans {
			for _, span := range scopeSpan.Spans {
				node := NewTraceDataNodeFromSpan(&span, scopeSpan.Scope.Name)
				if span.ParentSpanID == "" {
					root = node
					traceId, err = otelutils.DecodeOtelID(span.TraceID)
					if err != nil {
						return nil, err
					}
				} else if span.Kind == "SPAN_KIND_SERVER" {
					elseNode = append(elseNode, node)
				} else if span.Kind == "SPAN_KIND_CLIENT" {
					edges = append(edges, node)
				} else {
					return nil, errors.New("unknow span kind")
				}
				spanID, err := otelutils.DecodeOtelID(span.SpanID)
				if err != nil {
					return nil, err
				}
				nodeMap[spanID] = node
			}
		}
	}

	// 遍历edge然后构建Tree
	for _, edge := range edges {

		parentID := edge.ParentSpanId
		parentNode, ok := nodeMap[parentID]
		if !ok {
			// 理论上这种情况不应该出现
			logrus.Errorf("[TransferTraceData2Tree] - parentID = %s, edge = %v未查询到父span，请check原因", parentID, edge)
			return nil, errors.New("parent span not found")
		}
		// 在elseNode里面找，谁的parentid是这个spanid
		childNode := mapFind(nodeMap, func(node *TraceDataNode) bool {
			return node.ParentSpanId == edge.SpanId
		})
		if childNode == nil {
			logrus.Errorf("[TransferTraceData2Tree] - childNode = nil, parentID = %s, edge = %v错误", parentID, edge)
		}
		parentNode.Children = append(parentNode.Children, childNode)
	}

	return &TraceDataTree{Root: root, TraceId: traceId}, nil
}

func slicesFind(a []*TraceDataNode, f func(node *TraceDataNode) bool) *TraceDataNode {
	for _, each := range a {
		if f(each) {
			return each
		}
	}
	return nil
}

func mapFind(m map[string]*TraceDataNode, f func(node *TraceDataNode) bool) *TraceDataNode {
	for _, each := range m {
		if f(each) {
			return each
		}
	}
	return nil
}
