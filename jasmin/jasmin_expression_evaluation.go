package jasmin

// Returns jasmincode, type, stacklimit, locals used
// Leaves the result on top of the stack!!!
func expression_evaluation(expression *tree, var_info *variable_info, build *build_info, func_sigs *function_signatures) (string, string, int, []string) {
	// Local Var Maps can also be nil!!!

	switch len(expression.Branches) {
	case 0:
		panic("Internal Error: Expression has 0 children")
	case 1:
		switch expression.Branches[0].Leaf.Name {
		case "EL1", "EL2", "EL3", "EL4", "EL5", "EL6", "EL7":
			return expression_evaluation(&expression.Branches[0], var_info, build, func_sigs)
		case "name":
			// TODO
		case "LITERAL":
			// TODO
		case "FUNCCALL":
			// TODO
		default:
			panic("Internal Error: Expression has a unrecognized child. Name: " + expression.Branches[0].Leaf.Name)
		}
	case 2: // Unary Operations
		// TODO
	case 3:
		left_side_code, left_side_type, left_side_stack_limit, left_side_locals_used := expression_evaluation(&expression.Branches[0], var_info, build, func_sigs)
		right_side_code, right_side_type, right_side_stack_limit, right_side_locals_used := expression_evaluation(&expression.Branches[2], var_info, build, func_sigs)
		
		op_jasmin, type_required, return_type := operator_to_jasmin(expression.Branches[1].Leaf.Value.(string)) 

	
	}

	return "\n", "int", 1, []string{}
}

func operator_to_jasmin (op string) (string, string, string){
	switch op {
	case "+":
		return "add", "num", "num"
	case "-":
		return "sub", "num", "num"
	case  "*":
		return "mul", "num", "num"
	case ">":
		return ""
	}
}