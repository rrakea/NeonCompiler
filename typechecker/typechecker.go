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
	ParameterOrder []string
	CodeTree       *ParseTree
	LocalTree      *ParseTree
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
	info.Functions = functions

	funcArr := tree.Search_tree("FUNC")
	for _, f := range funcArr {
		name := f.Search_first_child("name").Leaf.Value.(string)
		rettype := f.Search_first_child("RETURNTYPE")
		returnType := rettype.Branches[0].Branches[0].Leaf.Name
		input, parameter_order, err := det_func_parameters(*f.Search_first_child("INPUTBLOCK"))
		code := f.Search_first_occurence_depth("STATEMENTBLOCK")
		localvars := f.Search_first_occurence_depth("VIRTUALVARBLOCK")
		// Find the actual code after the local vars
		if err != nil {
			TypeCheckError(err.Error())
			return info, false
		}
		if name == "main" || name == "Main" {
			TypeCheckError("Main function declared twice")
			return info, false
		}
		// Check against double entries
		_, exist := functions[name]
		if exist {
			TypeCheckError("Function " + name + " declared twice")
			return info, false
		}
		functions[name] = Function{Name: name, ReturnType: returnType, ParameterTypes: input, ParameterOrder: parameter_order, CodeTree: code, LocalTree: localvars}
	}
	// Add the main Function
	maintree := tree.Search_tree("MAIN")[0]
	argname := maintree.Search_direct_children("name")[0].Leaf.Value.(string)
	code := maintree.Search_first_occurence_depth("STATEMENTBLOCK")
	mainlocals := maintree.Search_first_occurence_depth("VIRTUALVARBLOCK")
	mainfunc := Function{Name: "main", ReturnType: "void", ParameterTypes: map[string]string{argname: "string[]"}, ParameterOrder: []string{"string[]"}, CodeTree: code, LocalTree: mainlocals}
	functions["main"] = mainfunc

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
		for _, l := range info.Functions[f.Name].LocalTree.Search_tree("VIRTUALVARBLOCK") {
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
		func_calls := functiontree.Search_tree("FUNCCALL")
		for _, call := range func_calls {
			args := call.Search_top_occurences("ARG")
			called_function := call.Search_first_child("name").Leaf.Value.(string)
			input_amount_wanted := len(info.Functions[called_function].ParameterOrder)
			if called_function == "Console.WriteLine" {
				input_amount_wanted = 1
			}
			if len(args) != input_amount_wanted {
				TypeCheckError("Function \"" + called_function + "\" called with the incorect amount of inputs (needs " + strconv.Itoa(input_amount_wanted) + " got " + strconv.Itoa(len(args)))
				return info, false
			}
			for i, arg := range args {
				ex := arg.Branches[0]
				ex_type, err := typeCheckExpression(ex, f.Name, info)
				if err != nil {
					TypeCheckError(err.Error() + "\nIn funccall expression " + called_function)
					return info, false
				}
				if called_function == "Console.WriteLine" {
					break
				}
				if ex_type != info.Functions[called_function].ParameterOrder[i] {
					TypeCheckError("Parameter " + strconv.Itoa(i) + " for function call \"" + called_function + "\" is misplaced. Function does not take a " + ex_type + " at that index")
					return info, false
				}
			}
		}

		for _, r := range functiontree.Search_tree("RETURN") {
			if len(r.Branches) == 2 {
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

func det_func_parameters(tree ParseTree) (map[string]string, []string, error) {
	retMap := map[string]string{}
	// Has no inputs
	if tree.Branches[0].Leaf.Name == ")" {
		return retMap, []string{}, nil
	}
	parameters_type_arr := []string{}
	for _, parameter := range tree.Search_tree("PARAMETER") {
		name := parameter.Branches[1].Leaf.Value.(string)
		paratype := det_func_input_type(parameter)
		parameters_type_arr = append(parameters_type_arr, paratype)
		_, exists := retMap[name]
		if exists {
			return retMap, parameters_type_arr, errors.New("Variable name declared twice in function signature " + name)
		}
		retMap[name] = paratype
	}
	return retMap, parameters_type_arr, nil
}

func det_func_input_type(tree *ParseTree) string {
	if len(tree.Branches) == 0 {
		return tree.Leaf.Name
	}
	return det_func_input_type(&tree.Branches[0])
}

func type_check_declaration(tree *ParseTree, vars map[string]Variable, funcname string, info *TypeCheckerInfo) error {
	name := tree.Search_first_child("name").Leaf.Value.(string)
	vartype := tree.Search_first_child("TYPE").Branches[0].Leaf.Name
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
