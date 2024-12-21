package typechecker

import (
	"compiler/parser"
	"errors"
	"fmt"
	"sync"
)

// Type alias
type ParseTree = parser.ParseTree

type Function struct {
	Name       string
	ReturnType string
	InputTypes map[string]string // Name -> Type
	Code       *ParseTree
}

type Variable struct {
	Vartype   string
	initvalue any
}

type TypeCheckerInfo struct {
	Main       ParseTree
	functions  map[string]Function            // Name -> Function
	globalVars map[string]Variable            // Name -> Variable
	localVar   map[string]map[string]Variable // Func Name -> Var Name -> Variable
}

func Typecheck(tree ParseTree) bool {
	const locationOfNameInFuncDec = 2
	const locationOfRetTypeInFuncDec = 1
	const locationOfInputInFuncDec = 4

	// Determine Function Signatures
	// Return & Input types -> Function Map
	main := treeSearch(tree, "main")
	_ = main
	functions := make(map[string]Function)
	funcArr := treeSearch(tree, "FUNC")
	for _, f := range funcArr {
		name := f.Branches[locationOfNameInFuncDec].Leaf.Value.(string)
		fmt.Println(name)
		fmt.Println(len(f.Branches))
		returnType := f.Branches[locationOfRetTypeInFuncDec].Branches[0].Leaf.Name
		input, err := detFuncInput(f.Branches[4])
		if err != nil {
			TypeCheckError(err)
			return false
		}
		functions[name] = Function{Name: name, ReturnType: returnType, InputTypes: input, Code: &f}
	}

	// Determine Global scoped vars

	// Determine locally scoped vars per function

	// Determine set of a var is correct

	// Determine each function call is correct

	// Determine every called function is existant
	return true
}

func TypeCheckError(err error) {
	fmt.Print(err)
}

func treeSearch(tree ParseTree, name string) []ParseTree {
	tokenchannel := make(chan ParseTree)
	wg := new(sync.WaitGroup)
	res := []ParseTree{}
	go func() {
		for true {
			token, ok := <-tokenchannel
			if !ok {
				return
			}
			res = append(res, token)
		}
	}()
	wg.Add(1)
	go treeSearchRoutine(tree, name, tokenchannel, wg)
	wg.Wait()
	close(tokenchannel)

	return res
}

func treeSearchRoutine(tree ParseTree, name string, channel chan ParseTree, wg *sync.WaitGroup) {
	defer wg.Done()
	if tree.Leaf.Name == name {
		channel <- tree
		return
	}
	//fmt.Print(len(tree.Branches))
	for _, t := range tree.Branches {
		wg.Add(1)
		go treeSearchRoutine(t, name, channel, wg)
	}
	return
}

func detFuncInput(tree ParseTree) (map[string]string, error) {
	retMap := make(map[string]string)
	if tree.Leaf.Name != "INPUTBLOCK" {
		fmt.Print("TYPE CHECKER ERROR: Function declaration not parsed correctly")
		return retMap, nil
	}
	if tree.Branches[0].Leaf.Name == ")" {
		return retMap, nil
	}
	start := tree.Branches[0]
	retMap[start.Branches[1].Leaf.Value.(string)] = start.Branches[0].Branches[0].Leaf.Name

	err := detFuncInputRec(start.Branches[2], &retMap)
	if err != nil {
		return retMap, err
	}
	return retMap, nil
}

func detFuncInputRec(tree ParseTree, inputs *map[string]string) error {
	if len(tree.Branches) == 0 {
		fmt.Print(tree)
	}

	if tree.Branches[0].Leaf.Name != "," {
		return nil
	}
	if (*inputs)[tree.Branches[2].Leaf.Value.(string)] != "" {
		return errors.New("Variable name declared twice in function Signature " + (*inputs)[tree.Branches[2].Leaf.Name] + " " + tree.Branches[1].Branches[0].Leaf.Name)
	}
	(*inputs)[tree.Branches[2].Leaf.Value.(string)] = tree.Branches[1].Branches[0].Leaf.Name
	return detFuncInputRec(tree.Branches[3], inputs)
}
