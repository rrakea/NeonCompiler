package typechecker

import (
	"compiler/lexer"
	"compiler/parser"
	"errors"
	"fmt"
	"strconv"
)

// Type alias
type ParseTree = parser.ParseTree
type Linenumber = lexer.LineNumber

type Function struct {
	Name           string
	ReturnType     string
	ParameterTypes map[string]string // Name -> Type
	CodeTree       *ParseTree
}

type Variable struct {
	Vartype    string
	Name       string
	Expression ParseTree
}

type TypeCheckerInfo struct {
	Functions  map[string]Function   // Name -> Function
	GlobalVars map[string]Variable   // Name -> Variable
	LocalVar   map[string][]Variable // Func Name -> []Variables
}

func Typecheck(tree ParseTree) (TypeCheckerInfo, bool) {
	info := TypeCheckerInfo{}

	// Determine Functions, Return Types, Paramter Types etc.
	functions := make(map[string]Function)
	funcArr := tree.Search_tree("FUNC")
	for _, f := range funcArr {
		name := f.Search_first_child("name").Leaf.Value.(string)
		returnType := f.Search_first_child("RETURNTYPE").Branches[0].Branches[0].Leaf.Name
		input, err := det_func_parameters(*f.Search_first_child("INPUTBLOCK"))
		pVarBlock := *f.Search_first_child("VIRTUALVARBLOCK")
		// Find the actual code after the local vars
		var code *ParseTree
		for true {
			if len(pVarBlock.Branches) == 1 {
				code = &pVarBlock.Branches[0]
				break
			}
			pVarBlock = *pVarBlock.Search_first_child("VIRTUALVARBLOCK")
		}
		if err != nil {
			TypeCheckError(err.Error())
			return info, false
		}
		// Check against double entries
		_, exist := functions[name]
		if exist {
			TypeCheckError("Function " + name + " declared twice")
			return info, false
		}
		functions[name] = Function{Name: name, ReturnType: returnType, ParameterTypes: input, CodeTree: code}
	}
	_, exist := functions["Main"]
	if !exist {
		TypeCheckError("No main function in code")
		return info, false
	}
	info.Functions = functions

	// Determine Global scoped vars
	globals := tree.Search_tree("GLOBALVARBLOCK")
	globalvars := map[string]Variable{}
	for _, globalvar := range globals {
		if len(globalvar.Branches) == 1 {
			continue
		}
		err := type_check_declaration(globalvar, globalvars, "", &info)
		if err != nil {
			TypeCheckError(err.Error())
			return info, false
		}
	}
	info.GlobalVars = globalvars

	// Determine locally scoped vars per function
	localVars := map[string]map[string]Variable{}
	local_var_array := map[string][]Variable{}
	info.LocalVar = local_var_array

	// Big Loop over all the functions to type check all statements
	for _, f := range info.Functions {
		functiontree := f.CodeTree

		// Local Variables per Function
		localVars[f.Name] = map[string]Variable{}
		for _, l := range functiontree.Search_tree("VIRTUALVARBLOCK") {
			if len(l.Branches) == 1 {
				break
			}
			err := type_check_declaration(l, localVars[f.Name], f.Name, &info)
			if err != nil {
				TypeCheckError(err.Error())
				return info, false
			}
		}
		for vname, vtype := range info.Functions[f.Name].ParameterTypes {
			localVars[f.Name][vname] = Variable{Name: vname, Vartype: vtype}
		}
		for _, v := range localVars[f.Name] {
			local_var_array[f.Name] = append(local_var_array[f.Name], v)
		}

		// Determine assign of a var is correct
		for _, assign := range functiontree.Search_tree("VARASSIGN") {
			name := assign.Search_first_child("name").Leaf.Value.(string)
			actualtype := ""
			ok := false
			localvartype := Variable{}
			// local var exists
			for _, l := range info.LocalVar[f.Name] {
				if l.Name == name {
					ok = true
					break
				}
			}
			if !ok {
				// Check if global var
				globalvartype, ok := info.GlobalVars[name]
				if !ok {
					TypeCheckError("Variable has not been initialized. Variable name: " + name)
					return info, false
				} else {
					actualtype = globalvartype.Vartype
				}
			} else {
				actualtype = localvartype.Vartype
			}
			expressionType, err := typeCheckExpression(*assign.Search_first_child("EXPRESSION"), f.Name, info)
			if actualtype != expressionType || err != nil {
				TypeCheckError(err.Error() + "\nVariable assingnment did not type check.\nVariable has been declared as type " + actualtype + " while Expression has type: " + expressionType)
				return info, false
			}
		}

		// Determine each function call is correct
		// TODO
		for _, call := range tree.Search_tree("FUNCCALL") {
			if len(call.Branches) == 1 {
				// Console.Log
			}
			name := call.Branches[0].Leaf.Value.(string)
			start := call.Branches[2].Branches[0]
			if len(start.Branches) <= 1 {
				if len(info.Functions[f.Name].ParameterTypes) == 0 {
					continue
				} else {
					TypeCheckError("Function \"" + f.Name + "\" called with input values (needs 0)")
					return info, false
				}
			}
			_ = name
			/*starttype, err := typeCheckExpression(start.Branches[0], f.Name, info)
			if err != nil {
				TypeCheckError("In funccall expression " + name)
			}
			_ = starttype
			*/
			/*
				//TODO Arg Type Checking
				calltype := []string{starttype}
				calc, err := calcArgContinue(start.Branches[1], f.Name, info)
				if err != nil {
					TypeCheckError(err.Error() + " In funccall expression \"" + name + "\"")
					return info, false
				}
				calltype = append(calltype, calc...)

				// Check return type links up
				if len(calltype) != len(info.Functions[name].InputTypes) {
					TypeCheckError("Wrong number of Inputs for function call of function " + name)
					return info, false
				}
			*/
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

		for _, r := range functiontree.Search_tree("RETURN") {
			if len(r.Branches) == 1 {
				if f.ReturnType == "void" {
					continue
				} else {
					TypeCheckError("Void function returns a value")
					return info, false
				}
			}
			extype, err := typeCheckExpression(*r.Search_first_child("EXPRESSION"), f.Name, info)
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

		for _, r := range functiontree.Search_tree("IF") {
			extype, err := typeCheckExpression(*r.Search_first_child("EXPRESSION"), f.Name, info)
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

		for _, r := range functiontree.Search_tree("WHILE") {
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

func det_func_parameters(tree ParseTree) (map[string]string, error) {
	retMap := map[string]string{}
	// Has no inputs
	if tree.Branches[0].Leaf.Name == ")" {
		return retMap, nil
	}
	for _, parameter := range tree.Search_tree("PARAMETER") {
		name := parameter.Branches[1].Leaf.Value.(string)
		paratype := det_func_input_type(parameter)
		_, exists := retMap[name]
		if exists {
			return retMap, errors.New("Variable name declared twice in function signature " + name)
		}
		retMap[name] = paratype
	}
	return retMap, nil
}

func det_func_input_type(tree *ParseTree) string {
	if tree.Branches == nil {
		return tree.Leaf.Name
	}
	if len(tree.Branches) > 1 {
		return "[]string"
	}
	return det_func_input_type(&tree.Branches[0])
}

func type_check_declaration(tree *ParseTree, vars map[string]Variable, funcname string, info *TypeCheckerInfo) error {
	name := tree.Search_first_child("name").Leaf.Value.(string)
	vartype := tree.Search_first_child("TYPE").Leaf.Name
	expression := tree.Search_first_child("EXPRESSION")
	newVar := Variable{Name: name, Vartype: vartype, Expression: *expression}

	expressiontype, err := typeCheckExpression(*expression, funcname, *info)
	if err != nil {
		return err
	}
	if expressiontype != vartype {
		return errors.New("The variable declaration of " + name + " does not typecheck.\nThe variable is declared as type " + vartype + " while the expression evaluates to " + expressiontype)
	}
	if vars[name].Name != "" {
		return errors.New("The variable with name " + newVar.Name + " is declared twice")
	}
	vars[name] = newVar
	return nil
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
