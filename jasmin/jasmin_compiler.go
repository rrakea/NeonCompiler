package jasmin

import (
	"compiler/parser"
	"compiler/typechecker"
)

type tree = parser.ParseTree

func Compile(parsetree *tree, info *typechecker.TypeCheckerInfo, file_name string) {
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
		ex_code, ex_type, ex_stack_limit, ex_locals_used := expression_evaluation(&global_var.Expression)
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

		// Which locals are used by name, so that we dont set a local limit that is too high
		locals_used_map := make(map[string]bool)

		// Maps the var name to its local var number
		var_map := make(map[string]int)

		// Local Variables
		local_var_code := ""
		for var_index, local_var := range info.LocalVar[function.Name] {
			ex_code, ex_type, ex_stack_limit, ex_locals_used := expression_evaluation(&local_var.Expression)
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
			var_map[local_var.Name] = var_index
		}
		func_local_limit := len(locals_used_map)

		statements, statement_local_limit, statement_stack_limit := Statement_block_evaluate(function.CodeTree, file_name, var_map)
		func_code := local_var_code + statements
		func_stack_limit += statement_stack_limit
		func_local_limit += statement_local_limit

		add_function(jasmin_file, function.Name, function.ReturnType, func_stack_limit, func_local_limit, func_code)
	}
}
