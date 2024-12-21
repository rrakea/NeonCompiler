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
	CodeTree   *ParseTree
}

type Variable struct {
	Vartype string
	Name    string
}

type TypeCheckerInfo struct {
	Main       Function
	Functions  map[string]Function            // Name -> Function
	GlobalVars map[string]Variable            // Name -> Variable
	LocalVar   map[string]map[string]Variable // Func Name -> Var Name -> Variable
	Code       map[string]ParseTree
}

func Typecheck(tree ParseTree) (TypeCheckerInfo, bool) {
	// TODO
	// You could make this way more efficient by going through the tree once and marking down all the important stuff
	// But cant be bothered

	const locationOfNameInFuncDec = 2
	const locationOfRetTypeInFuncDec = 1
	const locationOfInputInFuncDec = 4

	info := TypeCheckerInfo{}

	// Determine Function Signatures
	// Return & Input types -> Function Map
	main := treeSearch(tree, "main")
	if len(main) > 1 {
		TypeCheckError("More than 1 main functions in sourcefile")
		return info, false
	}
	mainFunc := Function{Name: "main", ReturnType: "void", InputTypes: make(map[string]string), CodeTree: &main[0]}
	info.Main = mainFunc

	functions := make(map[string]Function)
	funcArr := treeSearch(tree, "FUNC")
	for _, f := range funcArr {
		name := f.Branches[locationOfNameInFuncDec].Leaf.Value.(string)
		fmt.Println(name)
		fmt.Println(len(f.Branches))
		returnType := f.Branches[locationOfRetTypeInFuncDec].Branches[0].Leaf.Name
		input, err := detFuncInput(f.Branches[4])
		if err != nil {
			TypeCheckError(err.Error())
			return info, false
		}
		functions[name] = Function{Name: name, ReturnType: returnType, InputTypes: input, CodeTree: &f}
	}
	info.Functions = functions

	// Determine Global scoped vars
	globals := treeSearch(tree, "GLOBALVARBLOCK")
	globalvars := make(map[string]Variable)
	for _, globalvar := range globals {
		if len(globalvar.Branches) == 1 {
			break
		}
		successful, err := determineVariables(globalvar, globalvars, 1, "", info)
		if err != nil {
			successful = false
			TypeCheckError(err.Error())
		}
		if !successful {
			return info, false
		}
	}
	info.GlobalVars = globalvars

	// Determine locally scoped vars per function
	localVars := make(map[string]map[string]Variable)
	for _, f := range info.Functions {
		localVars[f.Name] = make(map[string]Variable)
		functiontree := *f.CodeTree
		locals := treeSearch(functiontree, "VIRTUALVARBLOCK")
		for _, l := range locals {
			if len(l.Branches) == 1 {
				break
			}

			successful, err := determineVariables(l, localVars[f.Name], 1, f.Name, info)
			if err != nil {
				successful = false
				TypeCheckError(err.Error())
			}
			if !successful {
				return info, false
			}
		}
	}
	info.LocalVar = localVars

	// Determine set of a var is correct

	for _, f := range info.Functions {
		functiontree := f.CodeTree
		assigns := treeSearch(*functiontree, "VARASSIGN")
		for _, assign := range assigns {
			name := assign.Branches[0].Leaf.Value.(string)
			actualtype := ""
			localvartype, ok := info.LocalVar[f.Name][name]
			if !ok {
				globalvartype, ok := info.GlobalVars[name]
				if !ok {
					TypeCheckError("Variable has not been initialized. Var name: " + name)
				} else {
					actualtype = globalvartype.Name
				}
			} else {
				actualtype = localvartype.Vartype
			}
			expressionType, err := typeCheckExpression(assign.Branches[2], f.Name, info)
			if actualtype != expressionType || err != nil {
				TypeCheckError(err.Error() + "\nVariable assing does not type check. Variable has been declared as type " + actualtype + " while Expression has type " + expressionType)
				return info, false
			}
		}
	}

	// Determine each function call is correct
	for _, call := range treeSearch(tree, "FUNCCALL") {
		_ = call
		// TODO
	}

	// Determine if / while / return expressions are correct
	for _, f := range info.Functions {
		functree := f.CodeTree
		returnarr := treeSearch(*functree, "RETURN")
		for _, r := range returnarr {
			if len(r.Branches) == 0 {
				if f.ReturnType == "void" {
					continue
				} else {
					TypeCheckError("Void function returns a value")
					return info, false
				}
			}
			extype, err := typeCheckExpression(r.Branches[0], f.Name, info)
			if extype != f.ReturnType || err != nil {
				TypeCheckError(err.Error() + "\nReturned value of type " + extype + " does not match return type in function signature (" + f.ReturnType)
				return info, false
			}
		}

		ifarr := treeSearch(*functree, "IF")

		for _, r := range ifarr {
			extype, err := typeCheckExpression(r.Branches[2], f.Name, info)
			if extype != "bool" || err != nil {
				TypeCheckError(err.Error() + "\nExpression in if statement does not evluate to a bool. It evaluates to " + extype)
			}
		}

		whilearr := treeSearch(*functree, "WHILE")

		for _, r := range whilearr {
			extype, err := typeCheckExpression(r.Branches[2], f.Name, info)
			if extype != "bool" || err != nil {
				TypeCheckError(err.Error() + "\nExpression in while statement does not evluate to a bool. It evaluates to " + extype)
				return info, false
			}
		}
	}
	fmt.Println(info.Main)
	fmt.Println(info.Functions)
	fmt.Println(info.GlobalVars)
	fmt.Println(info.LocalVar)
	return info, true
}

func TypeCheckError(s string) {
	fmt.Print(s)
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

func determineVariables(tree ParseTree, vars map[string]Variable, static int, funcname string, info TypeCheckerInfo) (bool, error) {
	name := tree.Branches[2+static].Leaf.Value.(string)
	vartype := tree.Branches[1+static].Branches[0].Leaf.Name
	newVar := Variable{Name: name, Vartype: vartype}
	expressiontype, err := typeCheckExpression(tree.Branches[4+static], funcname, info)
	if err != nil {
		return false, err
	}
	if expressiontype != vartype {
		return false, errors.New("The variable declaration of " + name + " does not typecheck.\nThe variable is declared as type " + vartype + " while the expression evaluates to " + expressiontype)
	}
	if vars[name].Name != "" {
		return false, errors.New("The variable with name " + newVar.Name + " is declared twice")
	}
	vars[name] = newVar
	return true, nil
}

func typeCheckExpression(expression ParseTree, funcName string, info TypeCheckerInfo) (string, error) {
	switch len(expression.Branches) {
	case 0:
		return "", errors.New("Expression has 0 Children: " + expression.Leaf.Name)
	case 1:
		switch expression.Branches[0].Leaf.Name {
		case "EL1", "EL2", "EL3", "EL4", "EL5", "EL6", "EL7":
			return typeCheckExpression(expression.Branches[0], funcName, info)

		case "name":
			name := expression.Branches[0].Leaf.Value.(string)
			local, ok := info.LocalVar[funcName][name]
			if !ok {
				global, ok := info.GlobalVars[name]
				if !ok {
					return "", errors.New("Variable " + name + " was not initialized")
				}
				return global.Vartype, nil
			}
			return local.Vartype, nil

		case "LITERAL":
			switch expression.Branches[0].Branches[0].Leaf.Name {
			case "intliteral", "doubleliteral":
				return "num", nil
			case "boolliteral":
				return "bool", nil
			case "stringliteral":
				return "string", nil
			default:
				return "", errors.New("Literal error ~ Most likely error in compiler :)")
			}

		case "FUNCCALL":
			return info.Functions[expression.Branches[0].Branches[0].Leaf.Value.(string)].ReturnType, nil

		default:
			return "", errors.New("Compiler Error, Expression without covered case has only one child")
		}

	case 2:
		switch expression.Branches[0].Leaf.Name {
		case "oplv5":
			ex, err := typeCheckExpression(expression.Branches[1], funcName, info)
			if err != nil {
				return "", err
			}
			if ex == "num" {
				return "num", nil
			}
			return "", errors.New(expression.Branches[0].Leaf.Value.(string) + " did not stand before a number")
		case "oplv7":
			ex, err := typeCheckExpression(expression.Branches[1], funcName, info)
			if err != nil {
				return "", err
			}
			if ex == "bool" {
				return "bool", nil
			}
			return "", errors.New("! did not stand before a boolen expression")
		}

	case 3:
		leftside, err := typeCheckExpression(expression.Branches[0], funcName, info)
		if err != nil {
			return "", err
		}
		rightside, err := typeCheckExpression(expression.Branches[2], funcName, info)
		if err != nil {
			return "", err
		}
		switch expression.Branches[1].Leaf.Name {
		case "oplv1", "oplv2":
			if rightside == "bool" && leftside == "bool" {
				return "bool", nil
			}
		case "oplv3":
			if (rightside == "bool" && leftside == "bool") || (rightside == "num" && leftside == "num") {
				return "bool", nil
			}
		case "oplv4":
			if rightside == "num" && leftside == "num" {
				return "bool", nil
			}

		case "oplv5", "oplv6":
			if rightside == "num" && leftside == "num" {
				return "num", nil
			}
		default:
			return "", errors.New("Compiler Error: Expression with 3 children does not have the correct opperators")
		}
	}
	return "", nil
}
