package parser

import (
	"compiler/lexer"
	"slices"

	"github.com/pterm/pterm"
)

type ParseTree struct {
	Leaf     ParseLeaf
	Branches []ParseTree
}

type ParseLeaf struct {
	Name  string
	Value any
}

func createParseTree(parseChan chan any) {
	//typeCheckerChan := make(chan any)
	//go typechecker.typecheck(typeCheckerChan)

	Trees := []ParseTree{}
	for true {
		newItem := <-parseChan
		switch newItem.(type) {
		case lexer.Token:
			token := newItem.(lexer.Token)
			newLeaf := ParseLeaf{Name: token.Identifier, Value: token.Value}
			newTree := ParseTree{Leaf: newLeaf, Branches: []ParseTree{}}
			Trees = append(Trees, newTree)
		case Rule:
			rule := newItem.(Rule)
			newTree := ParseTree{}
			newBranches := []ParseTree{}
			newTree.Leaf = ParseLeaf{Name: rule.nonTerminal, Value: 0}
			for i := range rule.production {
				len := len(Trees)
				if len == 0 {
					panic("Parse Tree Error, no new node possible")
				}
				newBranches = append(newBranches, Trees[len-i-1])
			}
			Trees = Trees[:len(Trees)-len(rule.production)]
			slices.Reverse(newBranches)
			newTree.Branches = newBranches
			Trees = append(Trees, newTree)
		case bool:
			// Change for debug
			if true {
				for _, t := range Trees {
					PrintTree(t)
				}
			}
			parseChan <- Trees[0]
			close(parseChan)
		}
	}
	for _, t := range Trees {
		parseChan <- t
	}
	close(parseChan)
}

func PrintTree(tree ParseTree) {
	ptree := makePTree(tree)
	renderTree := pterm.DefaultTree.WithRoot(ptree)
	renderTree.Render()
}

func makePTree(tree ParseTree) pterm.TreeNode {
	root := pterm.TreeNode{Text: tree.Leaf.Name, Children: []pterm.TreeNode{}}
	for _, t := range tree.Branches {
		root.Children = append(root.Children, makePTree(t))
	}
	return root
}
