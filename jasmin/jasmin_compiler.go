package jasmin

import (
	"compiler/parser"
	"compiler/typechecker"
)

type tree = parser.ParseTree

func Compile(parsetree *tree, info *typechecker.TypeCheckerInfo, file_name string) {
	jasmin_file := create_jasmin_file(file_name)
	write_default(jasmin_file, file_name)

	// Global Variable Definition
	for _, global_var := range info.GlobalVars {
		ex_code, ex_type, ex_length := expression_evaluation(&global_var.Expression)
		if ex_type != global_var.Vartype {
			panic("Internal Error: Type Checked Expression does not equal actual type of expression")
		}
		add_global_var(jasmin_file, global_var.Name, global_var.Vartype, ex_code, ex_length)
	}

	// Add main function
	add_main_func(info.Main.CodeTree)


	// Functions
	for _, function := range info.Functions {

		// Maps the var name to its local var number used in jasmin
		var_map := make(map[string]int)

		// Local Variables
		// You go over a map so no i :(
		var_count := 0
		local_var_code := ""
		for _, local_var := range info.LocalVar[function.Name] {
			ex_code, ex_type, ex_length := expression_evaluation(&local_var.Expression)
			if ex_type != local_var.Vartype {
				panic("Internal Error: Type Checked Expression does not equal actual type of expression")
			}
			var_code, _ := local_var_jasmin_code(local_var.Name, local_var.Vartype, var_count, ex_code, ex_length, function.Name)
			local_var_code += var_code + "\n"
			var_map[local_var.Name] = var_count
			var_count++
		}


		local_limit := var_count + 1 // 0 Indexed
		statements, statement_length := Statement_block_evaluate(function.CodeTree, file_name, var_map)
		// TODO: Calc Stack Limit
		func_code := local_var_code + "\n"+ statements
		
		add_function(jasmin_file, function.Name, function.ReturnType, stack_limit, local_limit, func_code, statement_length)
	}
}
