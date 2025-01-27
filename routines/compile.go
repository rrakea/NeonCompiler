package routines

import (
	"compiler/jasmin"
	"compiler/parser"
	"compiler/typechecker"
	"fmt"
)

func Compile(path string) {
	tree, file_name, parsingSuccesful := parser.Parse(path)

	if !parsingSuccesful {
		return
	}

	typeinfo, typechecksuccesful := typechecker.Typecheck(tree)
	if !typechecksuccesful {
		return
	}

	fmt.Println("Type Check Succesful!")

	jasmin.Build_jasmin(&tree, &typeinfo, file_name)

	fmt.Println("Compilation completed!")
	fmt.Println("Jasmin file created!")
}