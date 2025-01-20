package jasmin

import "strconv"

// Returns jasmincode, type, stacklimit, locals used
// Leaves the result on top of the stack!!!
func expression_evaluation(expression *tree, var_info *variable_info, build *build_info, func_sigs *function_signatures) (string, string, int, []string) {
	// Local Var Maps can also be nil!!!

	child := expression.Branches[0]
	switch len(expression.Branches) {
	case 0:
		panic("Internal Error: Expression has 0 children")
	case 1:
		switch expression.Branches[0].Leaf.Name {
		case "EL1", "EL2", "EL3", "EL4", "EL5", "EL6", "EL7":
			return expression_evaluation(&child, var_info, build, func_sigs)
		case "name":
			name := child.Leaf.Value.(string)
			var_type, ok := var_info.local_vars_type[name] 
			if ok {
				return var_type + "load " + name, var_type, 1, []string{name} 
			}
			var_type, ok = var_info.global_vars[name]
			if ok {
				return build.file_name + "/" + name, var_type, 1, []string{} 
			}
			panic("Internal Error: Var lookup failed. " + name + " not found in local or global var map")
		case "LITERAL":
			switch child.Leaf.Name {
			case "stringliteral":
				return "ldc " + child.Leaf.Value.(string), "Ljava/lang/String;", 1, []string{}
			case "boolliteral":
				return "ldc " + strconv.FormatBool(child.Leaf.Value.(bool)), "Z", 1, []string{}
			case "intliteral":
				return "ldc " + strconv.Itoa(child.Leaf.Value.(int)), "I", 1, []string{}
			case "doubleliteral":
				return "ldc2_w " + strconv.FormatFloat(child.Leaf.Value.(float64), 'f', -1, 64), "D", 1, []string{}  
			}
		case "FUNCCALL":
			func_name := child.Leaf.Value.(string)
			return_type, ok := func_sigs.return_type[func_name]
			if !ok {
				panic("Function name " + func_name + " not recognized")
			}

			// Evaluate args:
			args_code := ""
			arg_total_stack_limit := 0
			arg_total_locals_used := map[string]bool{}
			args := find_closest_children(&child, "arg")
			for i, arg := range args {
				arg_code, arg_type, arg_stack_limit, arg_locals_used := expression_evaluation(arg, var_info, build, func_sigs)
				_ = arg_type
				args_code += arg_code + "\n"
				if arg_stack_limit + i > arg_total_stack_limit {
					arg_total_stack_limit = arg_stack_limit + i
				}
				for _, local := range arg_locals_used {
					arg_total_locals_used[local] = true
				}
			}
			total_locals_used := []string{}
			for local := range arg_total_locals_used {
				total_locals_used = append(total_locals_used, local)
			} 
			return args_code + "invocestatic " + build.file_name + "/" + func_name + "()" + return_type, return_type, arg_total_stack_limit, total_locals_used
		default:
			panic("Internal Error: Expression has a unrecognized child. Name: " + expression.Branches[0].Leaf.Name)
		}
	case 2: // Unary Operations
		switch child.Leaf.Value {
		case "+":
			return expression_evaluation(&expression.Branches[1], var_info, build, func_sigs)  
		case "-":
			code, var_type, stack_limit, locals_used :=expression_evaluation(&expression.Branches[1], var_info, build, func_sigs)
			switch var_type {
			case "D":
				code = "dneg " + code
			case "I":
				code = "ineg " + code
			default:
				panic("Unary Operator \"-\" used before non numeric value")
			}
			return code, var_type, stack_limit, locals_used
		case "!":
			
		}
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