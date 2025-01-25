package routines

import (
	"compiler/parser"
	"compiler/typechecker"
	"compiler/liveness"
	"fmt"
)

func Liveness(path string) {
	tree, file_name, parsingSuccesful := parser.Parse(path)
	_ = file_name

	if !parsingSuccesful {
		return
	}

	typeinfo, typechecksuccesful := typechecker.Typecheck(tree)
	_ = typeinfo
	if !typechecksuccesful {
		return
	}

	fmt.Print("Starting Liveness Analysis...")

	liveness.Liveness(&tree)
}