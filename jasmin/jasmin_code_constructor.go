package jasmin

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func create_jasmin_file(name string) *os.File {
	file_name := strings.TrimSuffix(name, filepath.Ext(name))
	jasminfile, err := os.Create(file_name + ".j")
	if err != nil {
		panic("File could not be created")
	}
	return jasminfile
}

// Returns current line count ~ Line where next instruction will be written
// 0 indexed
func (build *build_info) add_header() {
	boiler_plate := "" +
		".class public " + build.file_name + "\n" +
		".super java/lang/Object \n\n"
	_, err := build.jasmin_file.WriteString(boiler_plate)
	if err != nil {
		panic("Could not write to file")
	}
}

// Add global vars. Values get set in clinit
func (build *build_info) add_global_var(var_name string, var_type string) {
	global_var_dec := "" +
		".field public static " + var_name + " " + var_type + "\n"
	build.jasmin_file.WriteString(global_var_dec)
}

func local_var_dec(name string, var_type string, var_count int, expression string) string {
	_ = name
	var_count_string := strconv.Itoa(var_count)
	var_dec := "" +
		expression +
		jasmin_type_prefix_converter(var_type) + "store " + var_count_string + "\n"
	return var_dec
}

// To initialize fields on the class initialisation
func (build *build_info) add_clinit(global_var_code map[string]string, global_var_type map[string]string, stack_limit int, local_limit int) {
	file_name := build.jasmin_file.Name()
	local_limit_string := strconv.Itoa(local_limit)
	stack_limit_string := strconv.Itoa(stack_limit)
	clinit := "" +
		".method static <clinit>()V \n" +
		".limit locals " + local_limit_string + "\n" +
		".limit stack " + stack_limit_string + "\n"

	for name, statement_block := range global_var_code {
		clinit += statement_block
		clinit += "putstatic " + file_name + "/" + name + " " + global_var_type[name] + "\n"
	}

	clinit += "return \n"
	clinit += ".end method\n\n"
	build.jasmin_file.WriteString(clinit)
}

func (build *build_info) add_function(method_name string, return_type string, arg_types []string, stack_limit int, local_limit int, statements string) {
	stack_limit_string := strconv.Itoa(stack_limit)
	local_limit_string := strconv.Itoa(local_limit)
	void_return := ""
	if return_type == "void" {
		void_return = "return \n"
	}
	arg_type_string := ""
	for _, arg := range arg_types {
		arg_type_string += jasmin_type_converter(arg)
	}

	func_dec := "" +
		".method public static " + method_name + "(" + arg_type_string + ") " + return_type + "\n" +
		".limit stack " + stack_limit_string + "\n" +
		".limit locals " + local_limit_string + "\n" +
		statements +
		void_return +
		".end method\n\n"

	_, err := build.jasmin_file.WriteString(func_dec)
	if err != nil {
		panic("Could not write to file")
	}
}
