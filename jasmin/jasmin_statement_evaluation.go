package jasmin

import "compiler/parser"

type tree = parser.ParseTree

func Statement_evaluate(statement_tree tree) (string, int) {
	statement := statement_tree.Branches[0]
	switch statement.Leaf.Name {
	case "FUNCCALL":
		name_tree, err := statement.Find_child("name")
		name := name_tree.Leaf.Value.(string)
		argblock, err := statement.Find_child("ARGBLOCK")
		if err != nil {
			panic("Invalid funccall found in parse tree")
		}
		func_call_evaluate(name, argblock)

	case "VARASSIGN":
		name_tree, err := statement.Find_child("name")
		name := name_tree.Leaf.Value.(string)
		expression, err := statement.Find_child("EXPRESSION")
		if err != nil {
			panic("Invalid var assign in parse tree")
		}
		return var_assign_evaluate(name, expression)

	case "RETURN":
		expression, err := statement.Find_child("expression")
		if err != nil {
			expression = nil
		}
		return return_evaluate(expression)
	case "IF":
		// TODO If/ Else

	case "WHILE":
		expression, err := statement.Find_child("expression")
		if err != nil {
			panic("Invalid while found in parse tree")
		}
		return while_evaluate(expression)
	default:
		panic("Unrecognized statement in parse tree")
	}
	return "\n", 1
}

func var_assign_evaluate(var_name string, expression *tree) (string, int) {
	ex_string, ex_length := expression_evaluation(expression)
}

func return_evaluate(expression *tree) (string, int) {
	ex_string, ex_length := expression_evaluation(expression)
}

func func_call_evaluate(name string, arg_block *tree) (string, int) {

}

func if_evaluate() (string int) {
	ex_string, ex_length := expression_evaluation(expression)
}

func while_evaluate(expression *tree) (string, int) {
	ex_string, ex_length := expression_evaluation(expression)
}
