package jasmin

import (
	"compiler/parser"
	"compiler/typechecker"
)

type tree = parser.ParseTree

func Build_jasmin(parsetree *tree, info *typechecker.TypeCheckerInfo, file_name string) {
	jasmin_file := create_jasmin_file(file_name)
	defer jasmin_file.Close()

	add_header(jasmin_file, file_name)

	// Name -> Code
	global_var_code := make(map[string]string)
	global_var_type := make(map[string]string)
	global_var_stack_limit := 0
	global_var_locals_used := make(map[string]bool)
	// Global Variable Definition
	for _, global_var := range info.GlobalVars {
		ex_code, ex_type, ex_stack_limit, ex_locals_used := expression_evaluation(&global_var.Expression, nil, nil, global_var_type, file_name)
		if ex_type != global_var.Vartype {
			panic("Internal Error: Type Checked Expression does not equal actual type of expression")
		}
		for _, locals_used := range ex_locals_used {
			ok := global_var_locals_used[locals_used]
			if !ok {
				global_var_locals_used[locals_used] = true
			}
		}

		global_var_stack_limit += ex_stack_limit
		add_global_var(jasmin_file, global_var.Name, global_var.Vartype)
		global_var_code[global_var.Name] = ex_code
		global_var_type[global_var.Name] = ex_type
	}

	// The global var initialisation is in <clinit>
	global_var_local_limit := len(global_var_locals_used)
	add_clinit(jasmin_file, global_var_code, global_var_type, global_var_stack_limit, global_var_local_limit)

	// Functions
	for _, function := range info.Functions {
		func_stack_limit := 0
		func_arg_type := ""

		// Which locals are used by name, so that we dont set a local limit that is too high
		locals_used_map := make(map[string]bool)

		// Maps the var name to its local var number
		var_map_count := make(map[string]int)
		var_map_type := make(map[string]string)

		// Function Arguments
		arg_count := 0
		for arg_name, arg_type := range info.Functions[function.Name].InputTypes {
			func_arg_type += arg_type.Inputtype
			var_map_count[arg_name] = arg_count
			var_map_type[arg_name] = arg_type.Inputtype
			arg_count++
		}

		// Local Variables
		local_var_code := ""
		for var_index, local_var := range info.LocalVar[function.Name] {
			ex_code, ex_type, ex_stack_limit, ex_locals_used := expression_evaluation(&local_var.Expression, var_map_count, var_map_type, global_var_type, file_name)
			if ex_type != local_var.Vartype {
				panic("Internal Error: Type Checked Expression does not equal actual type of expression")
			}
			// Set which local vars were used in the expression
			for _, locals_used := range ex_locals_used {
				ok := global_var_locals_used[locals_used]
				if !ok {
					global_var_locals_used[locals_used] = true
				}
			}
			func_stack_limit += ex_stack_limit

			var_code := local_var_dec(local_var.Name, local_var.Vartype, var_index, ex_code)
			local_var_code += var_code
			var_map_count[local_var.Name] = var_index + arg_count
			var_map_type[local_var.Name] = local_var.Vartype
		}

		statements, statement_stack_limit := Statement_block_evaluate(function.CodeTree, file_name, var_map_count, var_map_type, global_var_type, info.Functions)
		func_code := local_var_code + statements
		func_stack_limit += statement_stack_limit
		func_local_limit := len(locals_used_map)


		add_function(jasmin_file, function.Name, function.ReturnType, func_arg_type, func_stack_limit, func_local_limit, func_code)
	}
}

func jasmin_type_converter(var_type string) string{
	switch var_type {
	case "int":
		return "I"
	case "double":
		return "D"
	case "bool":
		return "Z"
	case "string":
		return "Ljava/lang/String;"
	case "string[]":
		return "[Ljava/lang/String;"
	case "void":
		return "V"
	default:
		panic("Internal Error: Invalid Type used")
	}
}