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
func add_header(jasmin_file *os.File, file_name string) {
	boiler_plate := "" +
		".class <public> <" + file_name + ">\n" +
		".super <java/lang/object> \n"
	_, err := jasmin_file.WriteString(boiler_plate)
	if err != nil {
		panic("Could not write to file")
	}
}

// Add global vars. Values get set in clinit
func add_global_var(jasmin_file *os.File, name string, vartype string) {
	global_var_dec := "" +
		".field public static " + name + " " + vartype
	jasmin_file.WriteString(global_var_dec)
}

func local_var_dec(name string, var_type string, var_count int, expression string) string {
	var_count_string := strconv.Itoa(var_count)
	var_dec := "" +
		expression +
		var_type + "store " + var_count_string
	return var_dec
}

// To initialize fields on the class initialisation
func add_clinit(jasmin_file *os.File, static_vars map[string]string, var_type map[string]string, stack_limit int, local_limit int) {
	file_name := jasmin_file.Name()
	local_limit_string := strconv.Itoa(local_limit)
	stack_limit_string := strconv.Itoa(stack_limit)
	clinit := "" +
		".method static <clinit>()V \n" +
		".limit locals " + local_limit_string +
		".limit stack " + stack_limit_string

	for name, statement_block := range static_vars {
		clinit += statement_block
		clinit += "putstatic " + file_name + "/" + name + var_type[name] + "\n"
	}

	clinit += "return \n"

	jasmin_file.WriteString(clinit)
}

func add_function(jasmin_file *os.File, method_name string, return_type string, arg_types string, stack_limit int, local_limit int, statements string) {
	stack_limit_string := strconv.Itoa(stack_limit)
	local_limit_string := strconv.Itoa(local_limit)
	func_dec := "" +
		".method public static " + method_name + "(" + arg_types + ")" + return_type + "\n" +
		".limit stack " + stack_limit_string + "\n" +
		".limit locals " + local_limit_string + "\n" +
		method_name + ":\n" +
		statements +
		method_name + "_end:\n" +
		".end method\n"

	_, err := jasmin_file.WriteString(func_dec)
	if err != nil {
		panic("Could not write to file")
	}
}
