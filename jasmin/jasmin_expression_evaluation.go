package jasmin

import "strconv"

// Returns jasmincode, type, stacklimit, locals used
// Leaves the result on top of the stack!!!
func expression_evaluation(expression *tree, var_info *variable_info, build *build_info, func_sigs *function_signatures, labels *label_info) (string, string, int, []string) {
	// Local Var Maps can also be nil!!!

	child := expression.Branches[0]
	switch len(expression.Branches) {
	case 0:
		panic("Internal Error: Expression has 0 children")
	case 1:
		switch expression.Branches[0].Leaf.Name {
		case "EL1", "EL2", "EL3", "EL4", "EL5", "EL6", "EL7":
			return expression_evaluation(&child, var_info, build, func_sigs, labels)
		case "name":
			name := child.Leaf.Value.(string)
			var_type, ok := var_info.local_vars_type[name]
			index := var_info.local_vars_index[name]
			if ok {
				return jasmin_type_prefix_converter(var_type) + "load " + strconv.Itoa(index) + "\n", var_type, 1, []string{name}
			}
			var_type, ok = var_info.global_vars[name]
			if ok {
				return "getstatic " + build.class + "/" + name + " " + var_type + "\n", var_type, 1, []string{}
			}
			panic("Internal Error: Var lookup failed. " + name + " not found in local or global var map")
		case "LITERAL":
			switch child.Branches[0].Leaf.Name {
			case "stringliteral":
				return "ldc " + child.Leaf.Value.(string) + "\n", "Ljava/lang/String;", 1, []string{}
			case "boolliteral":
				return "ldc " + strconv.FormatBool(child.Leaf.Value.(bool)) + "\n", "Z", 1, []string{}
			case "intliteral":
				return "ldc " + strconv.Itoa(child.Branches[0].Leaf.Value.(int)) + "\n", "I", 1, []string{}
			case "doubleliteral":
				return "ldc2_w " + strconv.FormatFloat(child.Leaf.Value.(float64), 'f', -1, 64) + "\n", "D", 1, []string{}
			default:
				panic("Invalid Literal name " + child.Leaf.Name)
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
				arg_code, arg_type, arg_stack_limit, arg_locals_used := expression_evaluation(arg, var_info, build, func_sigs, labels)
				_ = arg_type
				args_code += arg_code + "\n"
				if arg_stack_limit+i > arg_total_stack_limit {
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
			return args_code + "invocestatic " + build.class + "/" + func_name + "()" + return_type, return_type, arg_total_stack_limit, total_locals_used
		default:
			panic("Internal Error: Expression has a unrecognized child. Name: " + expression.Branches[0].Leaf.Name)
		}
	case 2: // Unary Operations
		switch child.Leaf.Value {
		case "+":
			return expression_evaluation(&expression.Branches[1], var_info, build, func_sigs, labels)
		case "-":
			code, var_type, stack_limit, locals_used := expression_evaluation(&expression.Branches[1], var_info, build, func_sigs, labels)
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
			code, var_type, stack_limit, locals_used := expression_evaluation(&expression.Branches[1], var_info, build, func_sigs, labels)
			if var_type != "B" {
				panic("Invalic Operator \"!\" before non bool value")
			}
			code = "" +
				code +
				"iconst_0 \n" +
				"ifeq bool_ex_false_" + strconv.Itoa(labels.bool_jump_count) + "\n" +
				"iconst_0\n" +
				"goto bool_ex_end_" + strconv.Itoa(labels.bool_jump_count) + ":\n" +
				"bool_ex_false_" + strconv.Itoa(labels.bool_jump_count) + "\n" +
				"iconst_1\n" +
				"bool_ex_end_" + strconv.Itoa(labels.bool_jump_count) + ":\n"
			labels.bool_jump_count += 1
			return code, "B", stack_limit + 1, locals_used
		}
	case 3:
		left_side_code, left_side_type, left_side_stack_limit, left_side_locals_used := expression_evaluation(&expression.Branches[0], var_info, build, func_sigs, labels)
		right_side_code, right_side_type, right_side_stack_limit, right_side_locals_used := expression_evaluation(&expression.Branches[2], var_info, build, func_sigs, labels)

		op_code := ""
		potential_cast_left, potential_cast_right, res_type, op_code_prefix := check_for_cast(left_side_type, right_side_type)
		op := expression.Branches[1].Leaf.Value.(string)
		switch op {
		case "+":
			op_code = op_code_prefix + "add\n"
		case "*":
			op_code = op_code_prefix + "mul\n"
		case "/":
			op_code = op_code_prefix + "div\n"
		case "%":
			op_code = op_code_prefix + "mod\n"
		case "-":
			op_code = op_code_prefix + "sub\n"
		case ">":
			res_type = "Z"
			op_code = op_to_bool("if"+op_code_prefix+"cmpgt", labels)
		case "<":
			res_type = "Z"
			op_code = op_to_bool_negated("if"+op_code_prefix+"cmpgt", labels)
		case ">=":
			switch res_type {
			case "I":
				op_code = op_to_bool("ificmpge", labels)
			case "D":
				op_code =
					"dcmpge\n"
				res_type = "Z"
			}
		case "<=":
			switch res_type {
			case "I":
				op_code = op_to_bool_negated("ificmpge", labels)
			case "D":
				res_type = "Z"
				op_code =
					"dcmpge\n" +
						"iconst_0\n" +
						op_to_bool_negated("ifeq", labels)
			}
		case "==":
			res_type = "Z"
			op_code = op_to_bool("ifeq", labels)
		case "!=":
			res_type = "Z"
			op_code = op_to_bool("ifne", labels)
		case "&&":
			if res_type != "Z" {
				panic("&& Used with 2 values that are not booleans")
			}
			op_code = "iand\n"
		case "||":
			if res_type != "Z" {
				panic("&& Used with 2 values that are not booleans")
			}
			op_code = "ior\n"

		default:
			panic("Unknown operator in expression: " + expression.Branches[0].Leaf.Value.(string))
		}
		code := potential_cast_left + left_side_code + potential_cast_right + right_side_code + op_code

		total_locals_used := deduplicate_locals_used(append(left_side_locals_used, right_side_locals_used...))

		return code, res_type, max(left_side_stack_limit, right_side_stack_limit), total_locals_used
	}
	panic("I dont even know how you got here")
}

func deduplicate_locals_used(locals []string) []string {
	tmp_map := map[string]bool{}
	for _, local := range locals {
		tmp_map[local] = true
	}
	ret_locals := []string{}
	for local := range tmp_map {
		ret_locals = append(ret_locals, local)
	}
	return ret_locals
}

func check_for_cast(left_side_type string, right_side_type string) (string, string, string, string) {
	if left_side_type == "Z" {
		if right_side_type != "Z" {
			panic("A bool and a non bool are being compared")
		}
		return "Z", "", "", "i"
	}
	res_type := ""
	potential_cast_left := ""
	potential_cast_right := ""
	op_code_type := ""

	if left_side_type == "I" && right_side_type == "I" {
		res_type = "I"
		op_code_type = "i"
		return potential_cast_left, potential_cast_right, res_type, op_code_type
	}
	if left_side_type == "D" {
		res_type = "D"
		potential_cast_right = "i2d"
	}
	if right_side_type == "D" {
		res_type = "D"
		if potential_cast_right == "i2d" {
			potential_cast_right = ""
		} else {
			potential_cast_left = "i2d"
		}
	}
	op_code_type = "d"
	return potential_cast_left, potential_cast_right, res_type, op_code_type
}

func op_to_bool(op_code string, labels *label_info) string {
	code := op_code +
		"bool_ex_false_" + strconv.Itoa(labels.bool_jump_count) + "\n" +
		"iconst_1\n" +
		"goto bool_ex_end_" + strconv.Itoa(labels.bool_jump_count) + ":\n" +
		"bool_ex_false_" + strconv.Itoa(labels.bool_jump_count) + "\n" +
		"iconst_0\n" +
		"bool_ex_end_" + strconv.Itoa(labels.bool_jump_count) + ":\n"
	labels.bool_jump_count += 1
	return code
}

func op_to_bool_negated(op_code string, labels *label_info) string {
	code := op_code +
		"bool_ex_false_" + strconv.Itoa(labels.bool_jump_count) + "\n" +
		"iconst_1\n" +
		"goto bool_ex_end_" + strconv.Itoa(labels.bool_jump_count) + ":\n" +
		"bool_ex_false_" + strconv.Itoa(labels.bool_jump_count) + "\n" +
		"iconst_0\n" +
		"bool_ex_end_" + strconv.Itoa(labels.bool_jump_count) + ":\n"
	labels.bool_jump_count += 1
	return code
}
