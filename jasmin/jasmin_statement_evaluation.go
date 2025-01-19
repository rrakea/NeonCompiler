package jasmin

import (
	"compiler/typechecker"
	"strconv"
)

type Function = typechecker.Function

// Returns: The code as string, stack limit
func Statement_block_evaluate(function_body *tree, file_name string, var_info *variable_info, functions map[string]Function, if_count int, while_count int) (string, int) {
	block_stack_limit := 0
	code := ""

	statements := find_top_level_statements(function_body)

	
	for _, statement := range statements {
		if len(statement.Branches) == 1 {
			continue
		}
		statement_code, statement_stack_limit := Statement_evaluate(statement, file_name, var_info.local_vars_index, var_info.local_vars_type, var_info.global_vars, functions, if_count, while_count)
		if statement_stack_limit > block_stack_limit {
			block_stack_limit = statement_stack_limit
		}
		code += statement_code
	}
	return code, block_stack_limit
}

func find_top_level_statements(block *tree) []*tree {
	statement_chan := make(chan *tree)
	go find_routine(block, statement_chan)
	statements := []*tree{}
	select {
	case statement, ok := <- statement_chan:
		if !ok {
			break
		}
		statements = append(statements, statement)
	}
	return statements
}

func find_routine(block *tree,stat_chan chan *tree) {
	for _, branch := range block.Branches {
		if branch.Leaf.Name == "STATEMENT" {
			stat_chan <- &branch
		} else {
			go find_routine(&branch, stat_chan)
		}
	} 
}

// Returns code, stack limit
func Statement_evaluate(statement_tree *tree, file_name string, var_map_index map[string]int, var_map_type map[string]string, global_var_type map[string]string, functions map[string]Function) (string, int) {
	statement := statement_tree.Branches[0]
	switch statement.Leaf.Name {
	case "FUNCCALL":
		name_tree, err := statement.Find_child("name")
		name := name_tree.Leaf.Value.(string)
		argblock, err := statement.Find_child("ARGBLOCK")
		if err != nil {
			panic("Invalid funccall found in parse tree")
		}
		return func_call_evaluate(name, argblock, var_map_index, var_map_type, global_var_type, file_name, functions)
	case "VARASSIGN":
		name_tree, err := statement.Find_child("name")
		name := name_tree.Leaf.Value.(string)
		expression, err := statement.Find_child("EXPRESSION")
		if err != nil {
			panic("Invalid var assign in parse tree")
		}
		return var_assign_evaluate(name, expression, var_map_index, var_map_type, global_var_type, file_name, functions)

	case "RETURN":
		expression, err := statement.Find_child("expression")
		if err != nil {
			expression = nil
		}
		return return_evaluate(expression, var_map_index, var_map_type, global_var_type, file_name, functions)
	case "IF":
		// TODO If/ Else

	case "WHILE":
		expression, err := statement.Find_child("expression")
		if err != nil {
			panic("Invalid while found in parse tree")
		}
		return while_evaluate(expression, var_map_index, var_map_type, global_var_type, file_name, functions)
	default:
		panic("Unrecognized statement in parse tree")
	}
	return "\n", 1
}

func var_assign_evaluate(var_name string, expression *tree, var_map_index map[string]int, var_map_type map[string]string, global_var_type map[string]string, file_name string, functions map[string]Function) (string, int) {
	ex_code, ex_type, ex_stack_limit, _ := expression_evaluation(expression, var_map_index, var_map_type, global_var_type, file_name)
	location, ok := var_map_index[var_name]
	if !ok {
		panic("Unitialized Variable " + var_name)
	}
	location_string := strconv.Itoa(location)
	
	retstring := "" +
		ex_code +
		ex_type + "load " + location_string + "\n"
	return retstring, ex_stack_limit
}

func return_evaluate(expression *tree, var_map_index map[string]int, var_map_type map[string]string, global_var_type map[string]string, file_name string, functions map[string]Function) (string, int) {
	if expression == nil {
		// Void Return
		return "return\n", 0
	}

	ex_code, ex_type, ex_stack_limit, _ := expression_evaluation(expression, var_map_index, var_map_type, global_var_type, file_name)

	retstring := "" +
		ex_code +
		ex_type + "return"
	return retstring, ex_stack_limit
}

func func_call_evaluate(func_name string, arg_block *tree, var_map_index map[string]int, var_map_type map[string]string, global_var_type map[string]string, file_name string, functions map[string]Function) (string, int) {
	args := typechecker.Parse_tree_search(*arg_block, "ARG")
	arg_stack_limit := 0
	arg_code := ""
	func_input_type := ""
	for _, arg_type:= range functions[func_name].InputTypes {
		func_input_type += jasmin_type_converter(arg_type.Inputtype) 
	}

	for _, arg := range args {
		ex_code, ex_type, ex_stack_limit, _ := expression_evaluation(&arg.Branches[0], var_map_index, var_map_type, global_var_type, file_name)
		_ = ex_type
		arg_code += ex_code + "\n"
		arg_stack_limit = max(arg_stack_limit, ex_stack_limit + len(args))
	}
	call := "" +
		arg_code +
		"invokestatic " + file_name + "/" + func_name + "(" + func_input_type + ")" + jasmin_type_converter(functions[func_name].ReturnType)
	return call, arg_stack_limit
}

func if_evaluate(if_statement *tree, var_map_index map[string]int, var_map_type map[string]string, global_var_type map[string]string, file_name string, functions map[string]Function) (string int) {
	condition, err := if_statement.Find_child("EXPRESSION")
	if err != nil {
		panic("Internal Error: If statement has no expressions as children")
	}

	cond_code, cond_type, cond_stack_limit, _ := expression_evaluation(condition, var_map_index, var_map_type, global_var_type, file_name)
	if cond_type != "Z" { // "Z" is bool in jasmin for some reason
		panic("Internal error: Typecheck passed, but conditional expression does not evaluate to bool")
	}
	
	if_code := "" +
	cond_code + "\n" +
	"iconst_0\n" +
	"if_icmpeq else_label" + if_count + "\n"

	return if_code
}

func while_evaluate(expression *tree, var_map_index map[string]int, var_map_type map[string]string, global_var_type map[string]string, file_name string, functions map[string]Function) (string, int) {
	ex_code, ex_type, ex_stack_limit, _ := expression_evaluation(expression, var_map_index, var_map_type, global_var_type, file_name)
}
