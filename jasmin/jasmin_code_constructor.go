package jasmin

import (
	"os"
	"strconv"
)

func create_jasmin_file(name string) *os.File {
	jasminfile, err := os.Create(name + ".j")
	if err != nil {
		panic("File could not be created")
	}
	return jasminfile
}

// Returns current line count ~ Line where next instruction will be written
// 0 indexed
func write_default(jasmin_file *os.File, class_name string) int {
	boiler_plate := "" +
		".class <public> <" + class_name + ">\n" +
		".super <java/lang/object> \n"
	_, err := jasmin_file.WriteString(boiler_plate)
	if err != nil {
		panic("Could not write to file")
	}
	return 2
}

func add_global_var(jasmin_file *os.File, name string, vartype string, expression string, expression_length int) int {
	global_var_dec := "" +
		expression +
		".field <public> <" + name + "> <" + vartype + "> " // Expression value
	jasmin_file.WriteString(global_var_dec)
	return 1 + expression_length
}

func add_method(jasmin_file *os.File, method_name string, return_type string, stack_limit int, local_limit int, statements string, statement_length int) int {
	stack_limit_string := strconv.Itoa(stack_limit)
	local_limit_string := strconv.Itoa(local_limit)
	func_dec := "" +
		".method <public> <" + method_name + " " + return_type + ">\n" +
		".limit stack <" + stack_limit_string + ">\n" +
		".limit locals <" + local_limit_string + ">\n" +
		method_name + ":\n" +
		statements +
		method_name + "_end:\n" +
		".end method\n"

	_, err := jasmin_file.WriteString(func_dec)
	if err != nil {
		panic("Could not write to file")
	}
	return 6 + statement_length
}

func local_var_dec(name string, vartype string, var_count int, expression string, expression_length int, func_name string) (string, int) {
	var_count_string := strconv.Itoa(var_count)
	var_dec := "" +
		expression +
		"var <" + var_count_string + "> is <" + name + "> <" + vartype + "> from <" + func_name + "> to <" + func_name + "_end>"// Expression
	return var_dec, 1 + expression_length
}
