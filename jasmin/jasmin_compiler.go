package jasmin

import (
	"compiler/parser"
	"compiler/typechecker"
	"os"
)

type tree = parser.ParseTree

type variable_info struct {
	global_vars      map[string]string // Name -> Type
	local_vars_index map[string]int    // Name -> Index
	local_vars_type  map[string]string // Name -> Type
}

type build_info struct {
	class       string
	jasmin_file *os.File
	parse_info  *typechecker.TypeCheckerInfo
}

type function_signatures struct {
	return_type    map[string]string
	parameter_type map[string]string
}

type label_info struct {
	if_count        int
	while_count     int
	bool_jump_count int
}

func Build_jasmin(parsetree *tree, info *typechecker.TypeCheckerInfo, file_name string) {
	build := new(build_info)
	jasmin_file := create_jasmin_file(file_name, build)
	defer jasmin_file.Close()

	build.jasmin_file = jasmin_file
	build.parse_info = info

	build.add_header()

	labels := label_info{0, 0, 0}

	func_sigs := evaluate_func_signatures(info)

	// Name -> Code
	global_var_code := make(map[string]string)
	global_var_type := make(map[string]string)
	global_var_stack_limit := 0
	global_var_locals_used := make(map[string]bool)
	var_info_only_for_globals := variable_info{nil, nil, global_var_type}

	// Global Variable Definition
	for _, global_var := range info.GlobalVars {
		ex_code, ex_type, ex_stack_limit, ex_locals_used := expression_evaluation(&global_var.Expression, &var_info_only_for_globals, build, &func_sigs, &labels)
		if ex_type != jasmin_type_converter(global_var.Vartype) {
			panic("Internal Error: Type Checked Expression does not equal actual type of expression")
		}
		for _, locals_used := range ex_locals_used {
			ok := global_var_locals_used[locals_used]
			if !ok {
				global_var_locals_used[locals_used] = true
			}
		}

		global_var_stack_limit += ex_stack_limit
		build.add_global_var(global_var.Name, global_var.Vartype)
		global_var_code[global_var.Name] = ex_code
		global_var_type[global_var.Name] = ex_type
	}
	build.jasmin_file.WriteString("\n")

	// The global var initialisation is in <clinit>
	global_var_local_limit := len(global_var_locals_used)
	build.add_clinit(global_var_code, global_var_type, global_var_stack_limit, global_var_local_limit)

	// Functions
	for _, function := range info.Functions {
		func_stack_limit := 0
		func_arg_type := []string{}

		// Which locals are used by name, so that we dont set a local limit that is too high
		locals_used_map := make(map[string]bool)

		// Maps the var name to its local var number
		var_map_count := make(map[string]int)
		var_map_type := make(map[string]string)
		var_info := variable_info{local_vars_index: var_map_count, local_vars_type: var_map_type, global_vars: global_var_type}

		// Function Arguments
		arg_count := 0
		for arg_name, arg_type := range info.Functions[function.Name].InputTypes {
			func_arg_type = append(func_arg_type, arg_type.Inputtype)
			var_map_count[arg_name] = arg_count
			var_map_type[arg_name] = arg_type.Inputtype
			arg_count++
		}

		// Local Variables
		local_var_code := ""
		for var_index, local_var := range info.LocalVar[function.Name] {
			// Check if the var is a parameter
			if len(local_var.Expression.Branches) != 0 {
				ex_code, ex_type, ex_stack_limit, ex_locals_used := expression_evaluation(&local_var.Expression, &var_info, build, &func_sigs, &labels)
				if ex_type != jasmin_type_converter(local_var.Vartype) {
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
			}
			var_map_count[local_var.Name] = var_index + arg_count
			var_map_type[local_var.Name] = local_var.Vartype
		}

		statements, statement_stack_limit := Statement_block_evaluate(function.CodeTree, &var_info, &func_sigs, build, &labels)
		func_code := local_var_code + statements
		func_stack_limit += statement_stack_limit
		func_local_limit := len(locals_used_map)

		build.add_function(function.Name, jasmin_type_converter(function.ReturnType), func_arg_type, func_stack_limit, func_local_limit, func_code)
	}
}

func evaluate_func_signatures(info *typechecker.TypeCheckerInfo) function_signatures {
	func_sig := function_signatures{map[string]string{}, map[string]string{}}
	for func_name, func_struct := range info.Functions {
		func_sig.return_type[func_name] = jasmin_type_converter(func_struct.ReturnType)
		parameters := ""
		for _, parameter := range func_struct.InputTypes {
			parameters += jasmin_type_converter(parameter.Inputtype)
		}
		func_sig.parameter_type[func_name] = parameters
	}
	return func_sig
}

func jasmin_type_converter(var_type string) string {
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

func jasmin_type_prefix_converter(var_type string) string {
	switch var_type {
	case "int":
		return "i"
	case "double":
		return "d"
	case "bool":
		return "z"
	case "string":
		return "a"
	case "string[]":
		return "[a"
	default:
		panic("Internal Error: Invalid Type used " + var_type)
	}
}
