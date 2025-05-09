package otelmodel

import (
	"encoding/json"
	"errors"
	"strings"
)

type PatternGraph struct {
	Edges []Edge `json:"edges"`
	Nodes []Node `json:"nodes"`
}

type Edge struct {
	Label  string `json:"label"`
	Source string `json:"source"`
	Target string `json:"target"`
}

type Node struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

type PatternTree struct {
	Root *PatternTreeNode
	Name string
	ID   int
}

type PatternTreeNode struct {
	Value    string
	Service  string
	Children []*PatternTreeNode
}

// Pattern2Tree converts a pattern string to a PatternTree.
// pattern is a json string
func Pattern2Tree(id int, name string, pattern []byte) (*PatternTree, error) {
	var patternGraph PatternGraph
	if err := json.Unmarshal(pattern, &patternGraph); err != nil {
		return nil, err
	}
	nodeMap := make(map[string]*PatternTreeNode)
	for _, node := range patternGraph.Nodes {
		splits := strings.Split(node.Label, "~") // e.g: svc3~/tosvc4
		if len(splits) != 2 {
			return nil, errors.New("invalid pattern format")
		}

		nodeMap[node.ID] = &PatternTreeNode{
			Value:    splits[1],
			Service:  splits[0],
			Children: make([]*PatternTreeNode, 0),
		}
	}
	for _, edge := range patternGraph.Edges {
		// 构建节点中的children，并且找到root(没有作为dst的节点一定是root)
		src, dst := edge.Source, edge.Target
		if _, ok := nodeMap[src]; !ok {
			return nil, errors.New("source node not found")
		}
		if _, ok := nodeMap[dst]; !ok {
			return nil, errors.New("target node not found")
		}
		nodeMap[src].Children = append(nodeMap[src].Children, nodeMap[dst])
		nodeMap[dst] = nil
	}
	for _, node := range nodeMap {
		if node != nil {
			return &PatternTree{
				Root: node,
				Name: name,
				ID:   id,
			}, nil
		}
	}
	return nil, errors.New("construct error, no root found")
}
