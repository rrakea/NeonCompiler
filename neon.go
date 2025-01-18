package main

import (
	"compiler/parser"
	"compiler/typechecker"
	"compiler/jasmin"
	"flag"
	"fmt"
	"os"
)

func main() {
	help := flag.Bool("help", false, "Display this Info")
	compile := flag.Bool("compile", false, "Compile the code")
	liveness := flag.Bool("liveness", false, "Start liveness analysis")
	constants := flag.Bool("constants", false, "Start constant propagation analysis")
	
	flag.Parse()

	fmt.Println()

	if *help {
		flag.PrintDefaults()
		fmt.Println()
		return
	}

	if len(os.Args) < 2 {
		fmt.Println("Please provide a code path and a flag")
		fmt.Println()
		flag.PrintDefaults()
		fmt.Println()
		return
	}

	if len(os.Args) > 3 {
		fmt.Println("To many arguments")
		fmt.Println()
		flag.PrintDefaults()
		fmt.Println()
		return
	}

	if !*compile && !*liveness && !*constants {
		fmt.Println("Please specify what the program should do.")
		fmt.Println()
		flag.PrintDefaults()
		fmt.Println()
		return
	}

	var path string

	if *compile {
		if len(os.Args) != 3 {
			fmt.Println("No path provided")
			fmt.Println()
			return
		}
		path = os.Args[2]
		// Send code to tokenizer
		tree, file_name, parsingSuccesful := parser.Parse(path, true)
		fmt.Println()

		if !parsingSuccesful{
			return
		}

		typeinfo, typechecksuccesful := typechecker.Typecheck(tree)
		if !typechecksuccesful {
			return
		}
		fmt.Println("Type Check Succesful!")
		fmt.Println()
		jasmin.Compile(&tree, &typeinfo, file_name)
		
		fmt.Println("Compilation completed!")
		fmt.Println("Jasmin file created!")
	}

	if *liveness {
		// Only main function!!
		fmt.Println("Not implemented")
		fmt.Println()
		return
	}

	if *constants {
		fmt.Println("Not implemented")
		fmt.Println()
		return
	}
}
