package jasmin



import (
	"strconv"
)

func Statement_block_evaluate(statement_block *tree, class_name string, var_map map[string]int) (string, int){
	// TODO
}

func Statement_evaluate(statement_tree tree, class_name string, varmap map[string]int) (string, int) {
	statement := statement_tree.Branches[0]
	switch statement.Leaf.Name {
	case "FUNCCALL":
		name_tree, err := statement.Find_child("name")
		name := name_tree.Leaf.Value.(string)
		argblock, err := statement.Find_child("ARGBLOCK")
		if err != nil {
			panic("Invalid funccall found in parse tree")
		}
		func_call_evaluate(name, argblock, class_name)

	case "VARASSIGN":
		name_tree, err := statement.Find_child("name")
		name := name_tree.Leaf.Value.(string)
		expression, err := statement.Find_child("EXPRESSION")
		if err != nil {
			panic("Invalid var assign in parse tree")
		}
		return var_assign_evaluate(name, expression, varmap)

	case "RETURN":
		expression, err := statement.Find_child("expression")
		if err != nil {
			expression = nil
		}
		return return_evaluate(expression, varmap)
	case "IF":
		// TODO If/ Else

	case "WHILE":
		expression, err := statement.Find_child("expression")
		if err != nil {
			panic("Invalid while found in parse tree")
		}
		return while_evaluate(expression, varmap)
	default:
		panic("Unrecognized statement in parse tree")
	}
	return "\n", 1
}

func var_assign_evaluate(var_name string, expression *tree, varmap map[string]int) (string, int) {
	ex_string, extype, ex_length := expression_evaluation(expression)
	location, ok := varmap[var_name]
	if !ok {
		panic("Unitialized Variable " + var_name)
	}
	location_string := strconv.Itoa(location)

	retstring := "" + 
	ex_string +
	extype + "load " + location_string + "\n"
	return retstring, ex_length + 1
}

func return_evaluate(expression *tree, varmap map[string]int) (string, int) {
	if expression == nil {
		// Void Return
		return "return V", 1
	}

	ex_string, ex_type, ex_length := expression_evaluation(expression)

	retstring := "" + 
	ex_string + 
	ex_type + "return" 
	return retstring, ex_length + 1
}

func func_call_evaluate(func_name string, arg_block *tree, class_name string) (string, int) {

	call := "" +
		"invokestatic " + class_name + "/" + func_name + "()"
	return call, 1
}

func if_evaluate(varmap map[string]int) (string int) {
	ex_string, ex_type, ex_length := expression_evaluation(expression)
}

func while_evaluate(expression *tree, varmap map[string]int) (string, int) {
	ex_string, ex_type, ex_length := expression_evaluation(expression)
}
