package jasmin

import (
	"strconv"
)

// Returns: The code as string, stack limit
func Statement_block_evaluate(function_body *tree, var_info *variable_info, func_sigs *function_signatures, build *build_info, labels *label_info) (string, int) {
	block_stack_limit := 0
	code := ""

	statements := function_body.Search_top_occurences("STATEMENT")

	for _, statement := range statements {
		statement_code, statement_stack_limit := Statement_evaluate(statement, func_sigs, var_info, build, labels)
		if statement_stack_limit > block_stack_limit {
			block_stack_limit = statement_stack_limit
		}
		code += statement_code
	}
	return code, block_stack_limit
}

// Returns code, stack limit
func Statement_evaluate(statement_tree *tree, func_sigs *function_signatures, var_info *variable_info, build *build_info, labels *label_info) (string, int) {
	statement := statement_tree.Branches[0]
	switch statement.Leaf.Name {
	case "FUNCCALL":
		name_tree := statement.Search_first_child("name")
		name := name_tree.Leaf.Value.(string)
		argblock := statement.Search_first_child("ARGBLOCK")
		return func_call_evaluate(name, argblock, func_sigs, var_info, build, labels)
	case "VARASSIGN":
		name_tree := statement.Search_first_child("name")
		name := name_tree.Leaf.Value.(string)
		expression := statement.Search_first_child("EXPRESSION")
		return var_assign_evaluate(name, expression, var_info, build, func_sigs, labels)

	case "RETURN":
		expression := statement.Search_first_child("EXPRESSION")
		return return_evaluate(expression, var_info, build, func_sigs, labels)
	case "IF":
		condition := statement.Search_first_child("EXPRESSION")
		if_block := statement.Search_first_child("STATEMENTBLOCK")
		else_block := statement.Search_first_child("ELSE")
		return if_evaluate(condition, if_block, else_block, var_info, build, labels, func_sigs)
	case "WHILE":
		condition := statement.Search_first_child("EXPRESSION")
		statement_block := statement.Search_first_child("STATEMENTBLOCK")
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
	args := arg_block.Search_top_occurences("ARG")
	arg_stack_limit := 0
	arg_code := ""
	first_arg_type := ""

	for i, arg := range args {
		ex_code, ex_type, ex_stack_limit, _ := expression_evaluation(&arg.Branches[0], var_info, build, func_sigs, labels)
		if i == 0 {
			first_arg_type = ex_type
		}
		_ = ex_type
		arg_code += ex_code
		arg_stack_limit = max(arg_stack_limit, ex_stack_limit+len(args))
	}

	if func_name == "Console.WriteLine" {
		call_type := jasmin_type_converter(first_arg_type)
		call := "" +
			"getstatic java/lang/System/out Ljava/io/PrintStream;\n" +
			arg_code +
			"invokevirtual java/io/PrintStream/println(" + call_type + ")V\n"
		return call, arg_stack_limit
	}

	call := "" +
		arg_code +
		"invokestatic " + build.class + "/" + func_name + "(" + func_sigs.parameter_type[func_name] + ")" + func_sigs.return_type[func_name] + "\n"
	return call, arg_stack_limit
}

func if_evaluate(condition *tree, if_block *tree, else_block *tree, var_info *variable_info, build *build_info, labels *label_info, func_sigs *function_signatures) (string, int) {
	// TODO else
	cond_code, cond_type, cond_stack_limit, _ := expression_evaluation(condition, var_info, build, func_sigs, labels)
	if cond_type != "z" { // "Z" is bool in jasmin for some reason
		panic("Internal error: Typecheck passed, but conditional expression does not evaluate to bool")
	}
	if_label := strconv.Itoa(labels.if_count)
	labels.if_count += 1

	if_statement_block, if_statement_block_stack_limit := Statement_block_evaluate(if_block, var_info, func_sigs, build, labels)
	else_statement_block := ""
	else_statement_block_stack_limit := 0
	if else_block != nil {
		else_statement_block, else_statement_block_stack_limit = Statement_block_evaluate(else_block, var_info, func_sigs, build, labels)
	}

	if_code := "" +
		cond_code + "\n" +
		"ifeq ELSE_LABEL_" +  if_label + "\n" +
		if_statement_block +
		"goto END_IF_ELSE_" + if_label + "\n" +
		"ELSE_LABEL_" + if_label + ":\n" +
		else_statement_block +
		"END_IF_ELSE_" + if_label + ":\n"

	return if_code, max(if_statement_block_stack_limit, else_statement_block_stack_limit) + cond_stack_limit + 1
}

func while_evaluate(condition *tree, statement_block *tree, var_info *variable_info, build *build_info, labels *label_info, func_sigs *function_signatures) (string, int) {
	while_label := strconv.Itoa(labels.while_count)
	cond_code, cond_type, cond_stack_limit, _ := expression_evaluation(condition, var_info, build, func_sigs, labels)
	if cond_type != "z" { // "Z" is bool in jasmin for some reason
		panic("Internal error: Typecheck passed, but conditional expression does not evaluate to bool")
	}

	while_statement_block, while_statement_stack_limit := Statement_block_evaluate(statement_block, var_info, func_sigs, build, labels)
	labels.while_count += 1
	
	while_code := "" +
		"WHILE_BEGIN" +  while_label + ":\n" +
		cond_code + "\n" +
		"ifeq WHILE_END_" + while_label + "\n" + 
		while_statement_block +
		"goto WHILE_BEGIN" +
		"WHILE_END_" + while_label

	while_statement_block_stack_limit := cond_stack_limit + 1 + while_statement_stack_limit
	return while_code, while_statement_block_stack_limit
}
