package typechecker

import (
	"compiler/lexer"
	"compiler/parser"
	"errors"
	"fmt"
	"strconv"
	"sync"
)

// Type alias
type ParseTree = parser.ParseTree
type Linenumber = lexer.LineNumber

type Function struct {
	Name       string
	ReturnType string
	InputTypes map[string]InputType // Name -> Type
	CodeTree   *ParseTree
}

type Variable struct {
	Vartype    string
	Name       string
	Expression ParseTree
}

type InputType struct {
	Inputtype string
	Index     int
}

type TypeCheckerInfo struct {
	Main       Function
	Functions  map[string]Function   // Name -> Function
	GlobalVars map[string]Variable   // Name -> Variable
	LocalVar   map[string][]Variable // Func Name -> []Variables
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
	main := Parse_tree_search(tree, "MAIN")
	if len(main) != 1 {
		TypeCheckError("Wrong amount of main function in source file")
		return info, false
	}
	mainFunc := Function{Name: "main", ReturnType: "void", InputTypes: make(map[string]InputType), CodeTree: &main[0].Branches[10]}
	args_name := main[0].Branches[7].Leaf.Value.(string)
	mainFunc.InputTypes[args_name] = InputType{"string[]", 0}
	info.Main = mainFunc

	functions := make(map[string]Function)
	funcArr := Parse_tree_search(tree, "FUNC")
	for _, f := range funcArr {
		name := f.Branches[locationOfNameInFuncDec].Leaf.Value.(string)
		//fmt.Println(name)
		//fmt.Println(len(f.Branches))
		returnType := f.Branches[locationOfRetTypeInFuncDec].Branches[0].Branches[0].Leaf.Name
		//fmt.Print(returnType)
		input, err := detFuncInput(f.Branches[4])
		if err != nil {
			TypeCheckError(err.Error())
			return info, false
		}
		functions[name] = Function{Name: name, ReturnType: returnType, InputTypes: input, CodeTree: &f.Branches[6]}
	}
	functions["main"] = info.Main
	info.Functions = functions

	// Determine Global scoped vars
	globals := Parse_tree_search(tree, "GLOBALVARBLOCK")
	globalvars := make(map[string]Variable)
	for _, globalvar := range globals {
		if len(globalvar.Branches) == 1 {
			break
		}
		successful, err := determineVariables(globalvar, globalvars, 1, "", info)
		if err != nil {
			successful = false
			TypeCheckError(err.Error())
			return info, false
		}
		if !successful {
			return info, false
		}
	}
	info.GlobalVars = globalvars

	// Determine locally scoped vars per function
	localVars := map[string]map[string]Variable{}
	local_var_array := map[string][]Variable{}
	for _, f := range info.Functions {
		localVars[f.Name] = map[string]Variable{}
		functiontree := *f.CodeTree
		locals := Parse_tree_search(functiontree, "VIRTUALVARBLOCK")
		for _, l := range locals {
			if len(l.Branches) == 1 {
				break
			}
			successful, err := determineVariables(l, localVars[f.Name], 0, f.Name, info)
			if err != nil {
				successful = false
				TypeCheckError(err.Error())
				return info, false
			}
			if !successful {
				return info, false
			}
		}
		for vname, vtype := range info.Functions[f.Name].InputTypes {
			localVars[f.Name][vname] = Variable{Name: vname, Vartype: vtype.Inputtype}
		}
		for _, v := range localVars[f.Name] {
			local_var_array[f.Name] = append(local_var_array[f.Name], v)
		}
	}
	info.LocalVar = local_var_array

	// Determine set of a var is correct

	for _, f := range info.Functions {
		functiontree := f.CodeTree
		assigns := Parse_tree_search(*functiontree, "VARASSIGN")
		for _, assign := range assigns {
			name := assign.Branches[0].Leaf.Value.(string)
			actualtype := ""
			ok := false
			localvartype := Variable{}
			for _, l := range info.LocalVar[f.Name] {
				if l.Name == name {
					ok = true
					break
				}
			}
			if !ok {
				globalvartype, ok := info.GlobalVars[name]
				if !ok {
					TypeCheckError("Variable has not been initialized. Variable name: " + name)
					return info, false
				} else {
					actualtype = globalvartype.Name
				}
			} else {
				actualtype = localvartype.Vartype
			}
			expressionType, err := typeCheckExpression(assign.Branches[2], f.Name, info)
			if actualtype != expressionType || err != nil {
				TypeCheckError(err.Error() + "\nVariable assingnment did not type check.\nVariable has been declared as type " + actualtype + " while Expression has type: " + expressionType)
				return info, false
			}
		}
	}

	// Determine if / while / return expressions are correct
	for _, f := range info.Functions {

		// Determine each function call is correct
		for _, call := range Parse_tree_search(tree, "FUNCCALL") {
			if len(call.Branches) == 1 {
				// Console.Log
			}
			name := call.Branches[0].Leaf.Value.(string)
			start := call.Branches[2].Branches[0]
			if len(start.Branches) <= 1 {
				if len(info.Functions[f.Name].InputTypes) == 0 {
					continue
				} else {
					TypeCheckError("Function \"" + f.Name + "\" called with input values (needs 0)")
					return info, false
				}
			}
			starttype, err := typeCheckExpression(start.Branches[0], f.Name, info)
			if err != nil {
				TypeCheckError("In funccall expression " + name)
			}
			calltype := []string{starttype}
			calc, err := calcArgContinue(start.Branches[1], f.Name, info)
			if err != nil {
				TypeCheckError(err.Error() + "In funccall expression " + name)
				return info, false
			}
			calltype = append(calltype, calc...)

			// Check return type links up
			if len(calltype) != len(info.Functions[name].InputTypes) {
				TypeCheckError("Wrong number of Inputs for function call of function " + name)
				return info, false
			}
			/*for i, calltype := range calltype {
				for _, input := range f.InputTypes {
					if input.Index == i{
						if calltype != input.Inputtype {
							TypeCheckError("Function \"" + f.Name + "\" called with wrong parameter at index " + strconv.Itoa(i))
							return info, false
						}
					}
				}
			}*/
		}

		functree := f.CodeTree
		returnarr := Parse_tree_search(*functree, "RETURN")
		for _, r := range returnarr {
			if len(r.Branches) <= 1 {
				if f.ReturnType == "void" {
					continue
				} else {
					TypeCheckError("Void function returns a value")
					return info, false
				}
			}
			extype, err := typeCheckExpression(r.Branches[1], f.Name, info)
			if extype != f.ReturnType || err != nil {
				if err != nil {
					extype = err.Error()
				} else {
					err = errors.New("")
				}
				TypeCheckError(err.Error() + "\nReturned value of type " + extype + " does not match return type \nin function signature (" + f.ReturnType + ")")
				return info, false
			}
		}

		ifarr := Parse_tree_search(*functree, "IF")

		for _, r := range ifarr {
			extype, err := typeCheckExpression(r.Branches[2], f.Name, info)
			if extype != "bool" || err != nil {
				if err != nil {
					extype = err.Error()
				} else {
					err = errors.New("")
				}
				TypeCheckError(err.Error() + "\nIf statement header: Line: " + strconv.Itoa(r.Branches[0].Leaf.Value.(Linenumber).Line-1) + " is not bool, is " + extype)
				return info, false
			}
		}

		whilearr := Parse_tree_search(*functree, "WHILE")

		for _, r := range whilearr {
			extype, err := typeCheckExpression(r.Branches[2], f.Name, info)
			if extype != "bool" || err != nil {
				if err != nil {
					extype = err.Error()
				} else {
					err = errors.New("")
				}
				TypeCheckError(err.Error() + "\nExpression in while statement header does not evaluate to bool. It evaluates to " + extype)
				return info, false
			}
		}
	}
	return info, true
}

func TypeCheckError(s string) {
	fmt.Print("TYPE ERROR: ")
	fmt.Println(s)
}

func Parse_tree_search(tree ParseTree, name string) []ParseTree {
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
	}
	//fmt.Print(len(tree.Branches))
	for _, t := range tree.Branches {
		wg.Add(1)
		go treeSearchRoutine(t, name, channel, wg)
	}
	return
}

func detFuncInput(tree ParseTree) (map[string]InputType, error) {
	retMap := make(map[string]InputType)
	if tree.Leaf.Name != "INPUTBLOCK" {
		fmt.Print("TYPE CHECKER ERROR: Function declaration not parsed correctly")
		return retMap, nil
	}
	if tree.Branches[0].Leaf.Name == ")" {
		return retMap, nil
	}
	start := tree.Branches[0]
	retMap[start.Branches[1].Leaf.Value.(string)] = InputType{start.Branches[0].Branches[0].Leaf.Name, 0}

	err := detFuncInputRec(start.Branches[2], &retMap, 0)
	if err != nil {
		return retMap, err
	}
	return retMap, nil
}

func detFuncInputRec(tree ParseTree, inputs *map[string]InputType, index int) error {
	if len(tree.Branches) == 0 {
		fmt.Print(tree)
	}

	if tree.Branches[0].Leaf.Name != "," {
		return nil
	}
	if (*inputs)[tree.Branches[2].Leaf.Value.(string)].Inputtype != "" {
		return errors.New("Variable name declared twice in function Signature " + (*inputs)[tree.Branches[2].Leaf.Name].Inputtype + " " + tree.Branches[1].Branches[0].Leaf.Name)
	}
	(*inputs)[tree.Branches[2].Leaf.Value.(string)] = InputType{tree.Branches[1].Branches[0].Leaf.Name, index}
	return detFuncInputRec(tree.Branches[3], inputs, index+1)
}

func determineVariables(tree ParseTree, vars map[string]Variable, static int, funcname string, info TypeCheckerInfo) (bool, error) {
	name := tree.Branches[1+static].Leaf.Value.(string)
	vartype := tree.Branches[static].Branches[0].Leaf.Name
	newVar := Variable{Name: name, Vartype: vartype}
	expressiontype, err := typeCheckExpression(tree.Branches[3+static], funcname, info)
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
		return "undefined", errors.New("Expression has 0 Children: " + expression.Leaf.Name)
	case 1:
		switch expression.Branches[0].Leaf.Name {
		case "EL1", "EL2", "EL3", "EL4", "EL5", "EL6", "EL7":
			return typeCheckExpression(expression.Branches[0], funcName, info)

		case "name":
			name := expression.Branches[0].Leaf.Value.(string)
			locals := info.LocalVar[funcName]
			ok := false
			local := Variable{}
			for _, l := range locals {
				if l.Name == name {
					ok = true
					local = l
					break
				}
			}
			if !ok {
				global, ok := info.GlobalVars[name]
				if !ok {
					return "undefined", errors.New("Variable " + name + " was not initialized.")
				}
				return global.Vartype, nil
			}
			return local.Vartype, nil

		case "LITERAL":
			switch expression.Branches[0].Branches[0].Leaf.Name {
			case "intliteral":
				return "int", nil
			case "doubleliteral":
				return "double", nil
			case "boolliteral":
				return "bool", nil
			case "stringliteral":
				return "string", nil
			default:
				return "undefined", errors.New("Literal error ~ Mt likely error in compiler :).\n Calculated Type: " + expression.Branches[0].Branches[0].Leaf.Name)
			}

		case "FUNCCALL":
			return info.Functions[expression.Branches[0].Branches[0].Leaf.Value.(string)].ReturnType, nil

		default:
			return "undefined", errors.New("Compiler Error, Expression without covered case has only one child")
		}

	case 2:
		switch expression.Branches[0].Leaf.Name {
		case "oplv5":
			ex, err := typeCheckExpression(expression.Branches[1], funcName, info)
			if err != nil {
				return "undefined", err
			}
			if ex == "int" || ex == "double" {
				return ex, nil
			}
			return "undefined", errors.New(expression.Branches[0].Leaf.Value.(string) + " did not stand before a number")
		case "oplv7":
			ex, err := typeCheckExpression(expression.Branches[1], funcName, info)
			if err != nil {
				return "undefined", err
			}
			if ex == "bool" {
				return "bool", nil
			}
			return "undefined", errors.New("! did not stand before a boolen expression")
		}

	case 3:
		leftside, err := typeCheckExpression(expression.Branches[0], funcName, info)
		if err != nil {
			return "undefined", err
		}
		rightside, err := typeCheckExpression(expression.Branches[2], funcName, info)
		if err != nil {
			return "undefined", err
		}
		switch expression.Branches[1].Leaf.Name {
		case "oplv1", "oplv2":
			if rightside == "bool" && leftside == "bool" {
				return "bool", nil
			}
		case "oplv3":
			if (rightside == "bool" && leftside == "bool") || ((rightside == "int" || rightside == "double") && (leftside == "int" || leftside == "double")) {
				return "bool", nil
			}
		case "oplv4":
			if (rightside == "int" || rightside == "double") && (leftside == "int" || leftside == "double") {
				return "bool", nil
			}

		case "oplv5", "oplv6":
			if rightside == "int" && leftside == "int" {
				return "int", nil
			}

			if (rightside == "int" || rightside == "double") && (leftside == "int" || leftside == "double") {
				return "double", nil
			}
		default:
			return "undefined", errors.New("Compiler Error: Expression with 3 children does not have the correct opperators")
		}
	}
	return "undefined", nil
}

func calcArgContinue(start ParseTree, fname string, info TypeCheckerInfo) ([]string, error) {
	if len(start.Branches) <= 1 {
		return []string{}, nil
	}
	argtype, err := typeCheckExpression(start.Branches[1], fname, info)
	if err != nil {
		return []string{}, err
	}
	calc, err := calcArgContinue(start.Branches[2], fname, info)
	if err != nil {
		return []string{}, err
	}
	return append([]string{argtype}, calc...), nil
}
