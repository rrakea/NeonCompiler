package main

import (
	"compiler/routines"
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

	switch {
	case *help:
		flag.PrintDefaults()
	
	case len(os.Args) < 2:
		fmt.Println("Please provide a code path and a flag")
		fmt.Println()
		flag.PrintDefaults()
	
	case len(os.Args) > 3:
		fmt.Println("To many arguments")
		fmt.Println()
		flag.PrintDefaults()
	
	case !*compile && !*liveness && !*constants:
		fmt.Println("Please specify what the program should do.")
		fmt.Println()
		flag.PrintDefaults()
	
	default:
		path := os.Args[2]
		switch {
			case *compile:
				routines.Compile(path)
			case *liveness:
				routines.Liveness(path)
			case *constants:
				routines.ConstantPropogation(path)
			default:
				fmt.Print("This should not happen. Hmmm")
		}
	}
	fmt.Println()
}