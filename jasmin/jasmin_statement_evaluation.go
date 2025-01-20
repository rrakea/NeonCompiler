package jasmin

import (
	"compiler/typechecker"
	"strconv"
	"sync"
)

type Function = typechecker.Function

// Returns: The code as string, stack limit
func Statement_block_evaluate(function_body *tree, var_info *variable_info, func_sigs *function_signatures, build *build_info, labels *label_info) (string, int) {
	block_stack_limit := 0
	code := ""

	statements := find_closest_children(function_body, "STATEMENT")

	for _, statement := range statements {
		statement_code, statement_stack_limit := Statement_evaluate(statement, func_sigs, var_info, build, labels)
		if statement_stack_limit > block_stack_limit {
			block_stack_limit = statement_stack_limit
		}
		code += statement_code
	}
	return code, block_stack_limit
}

// Finds the children of the name, without searching the children of these nodes
func find_closest_children(block *tree, name string) []*tree {
	statement_chan := make(chan *tree)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go find_routine(block, statement_chan, name, wg)
	statements := []*tree{}
	go func() {
		wg.Wait()
		close(statement_chan)
	}()
	for stat := range statement_chan {
		statements = append(statements, stat)
	}
	return statements
}

func find_routine(block *tree, stat_chan chan *tree, name string, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, branch := range block.Branches {
		if branch.Leaf.Name == name {
			stat_chan <- &branch
		} else {
			wg.Add(1)
			go find_routine(&branch, stat_chan, name, wg)
		}
	}
}

// Returns code, stack limit
func Statement_evaluate(statement_tree *tree, func_sigs *function_signatures, var_info *variable_info, build *build_info, labels *label_info) (string, int) {
	statement := statement_tree.Branches[0]
	switch statement.Leaf.Name {
	case "FUNCCALL":
		name_tree, err := statement.Find_child("name")
		name := name_tree.Leaf.Value.(string)
		argblock, err := statement.Find_child("ARGBLOCK")
		if err != nil {
			panic("Invalid funccall found in parse tree")
		}
		return func_call_evaluate(name, argblock, func_sigs, var_info, build, labels)
	case "VARASSIGN":
		name_tree, err := statement.Find_child("name")
		name := name_tree.Leaf.Value.(string)
		expression, err := statement.Find_child("EXPRESSION")
		if err != nil {
			panic("Invalid var assign in parse tree")
		}
		return var_assign_evaluate(name, expression, var_info, build, func_sigs, labels)

	case "RETURN":
		expression, err := statement.Find_child("EXPRESSION")
		if err != nil {
			expression = nil
		}
		return return_evaluate(expression, var_info, build, func_sigs, labels)
	case "IF":
		condition, err := statement.Find_child("EXPRESSION")
		if err != nil {
			panic("Internal Error: Invalid if statements block, expression missing")
		}
		statement_block, err := statement.Find_child("STATEMENTBLOCK")
		if err != nil {
			panic("Internal Error: Invalid if statement block, no statement block")
		}
		return if_evaluate(condition, statement_block, var_info, build, labels, func_sigs)
	case "WHILE":
		condition, err := statement.Find_child("EXPRESSION")
		if err != nil {
			panic("Invalid while found in parse tree")
		}
		statement_block, err := statement.Find_child("STATEMENTBLOCK")
		if err != nil {
			panic("Internal Error: Invalid while statement block, no statement block")
		}
		return while_evaluate(condition, statement_block, var_info, build, labels, func_sigs)
	default:
		panic("Unrecognized statement in parse tree")
	}
}

func var_assign_evaluate(var_name string, expression *tree, var_info *variable_info, build *build_info, func_sigs *function_signatures, labels *label_info) (string, int) {
	ex_code, ex_type, ex_stack_limit, _ := expression_evaluation(expression, var_info, build, func_sigs, labels)
	location, ok := var_info.local_vars_index[var_name]
	if !ok {
		panic("Unitialized Variable " + var_name)
	}
	location_string := strconv.Itoa(location)

	retstring := "" +
		ex_code +
		ex_type + "load " + location_string + "\n"
	return retstring, ex_stack_limit
}

func return_evaluate(expression *tree, var_info *variable_info, build *build_info, func_sigs *function_signatures, labels *label_info) (string, int) {
	if expression == nil {
		// Void Return
		return "return\n", 0
	}

	ex_code, ex_type, ex_stack_limit, _ := expression_evaluation(expression, var_info, build, func_sigs, labels)

	retstring := "" +
		ex_code +
		ex_type + "return\n"
	return retstring, ex_stack_limit
}

func func_call_evaluate(func_name string, arg_block *tree, func_sigs *function_signatures, var_info *variable_info, build *build_info, labels *label_info) (string, int) {
	args := typechecker.Parse_tree_search(*arg_block, "ARG")
	arg_stack_limit := 0
	arg_code := ""

	for _, arg := range args {
		ex_code, ex_type, ex_stack_limit, _ := expression_evaluation(&arg.Branches[0], var_info, build, func_sigs, labels)
		_ = ex_type
		arg_code += ex_code + "\n"
		arg_stack_limit = max(arg_stack_limit, ex_stack_limit+len(args))
	}
	call := "" +
		arg_code +
		"invokestatic " + build.class + "/" + func_name + "(" + func_sigs.parameter_type[func_name] + ")" + func_sigs.return_type[func_name] + "\n"
	return call, arg_stack_limit
}

func if_evaluate(condition *tree, statement_block *tree, var_info *variable_info, build *build_info, labels *label_info, func_sigs *function_signatures) (string, int) {
	// TODO else
	cond_code, cond_type, cond_stack_limit, _ := expression_evaluation(condition, var_info, build, func_sigs, labels)
	if cond_type != "Z" { // "Z" is bool in jasmin for some reason
		panic("Internal error: Typecheck passed, but conditional expression does not evaluate to bool")
	}

	if_code := "" +
		cond_code + "\n" +
		"iconst_0\n" +
		"if_icmpeq else_label" + strconv.Itoa(labels.if_count) + "\n"
	labels.if_count += 1

	if_statement_block, statement_block_stack_limit := Statement_block_evaluate(statement_block, var_info, func_sigs, build, labels)
	if_code += if_statement_block

	if_statement_stack_limit := cond_stack_limit + 1 + statement_block_stack_limit
	return if_code, if_statement_stack_limit
}

func while_evaluate(condition *tree, statement_block *tree, var_info *variable_info, build *build_info, labels *label_info, func_sigs *function_signatures) (string, int) {
	cond_code, cond_type, cond_stack_limit, _ := expression_evaluation(condition, var_info, build, func_sigs, labels)
	if cond_type != "Z" { // "Z" is bool in jasmin for some reason
		panic("Internal error: Typecheck passed, but conditional expression does not evaluate to bool")
	}

	while_code := "" +
		"while_begin" + strconv.Itoa(labels.while_count) + ":\n" +
		cond_code + "\n" +
		"iconst_0\n" +
		"if_icmpeq else_label" + strconv.Itoa(labels.if_count) + "\n"
	labels.while_count += 1

	while_statement_block, statement_block_stack_limit := Statement_block_evaluate(statement_block, var_info, func_sigs, build, labels)
	while_code += while_statement_block
	while_code += "GOTO while_begin" + strconv.Itoa(labels.while_count)

	while_statement_stack_limit := cond_stack_limit + 1 + statement_block_stack_limit
	return while_code, while_statement_stack_limit
}
