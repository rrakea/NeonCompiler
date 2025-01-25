package liveness

import (
	"compiler/parser"
	"fmt"
)

/*
Endresult: map[int][]string that use linenumber -> var that has to have a register here
You go through the var decs and set the beginning number of each var, as well as the last number to -1
You go thorugh all the expressions and increase the last used number if the var is used
Treat while loops as one line -> the variable needs to be scoped the whole loop
In the end you check the max length of the map -> max amount of registers that have to be used
*/

type liveness_var struct {
	name string
	firstuse int
	lastuse int
}

func Liveness (tree *parser.ParseTree) {
	main := tree.Search_tree("MAIN")[0]

	vars :=  map[string]*liveness_var{}
	usesInDec := 0 // Lines that use variables in the local var declaration

	for i, dec := range main.Search_tree("VIRTUALVARBLOCK") {
		name := dec.Search_first_child("name").Leaf.Value.(string)
		vars[name] = &liveness_var{name: name, firstuse: i, lastuse: -1}
		ex := dec.Search_first_child("EXPRESSION")
		usedInDecEx := LivenessCheckExpression(ex)
		for _, use := range usedInDecEx {
			vars[use].lastuse = i
		}
		usesInDec = i
	}

	// We only search the top level statements, since we dont want to search inside of while loops
	statements := tree.Search_top_occurences("STATEMENT")
	for i := 0; i < len(statements); i++  {
		stat := statements[i].Branches[0]
		switch  stat.Leaf.Name{
		
		case "VARSSIGN":
			vars[stat.Search_first_child("name").Leaf.Value.(string)].lastuse = i
			for _, use := range LivenessCheckExpression(stat.Search_first_child("EXPRESSION")) {
				vars[use].lastuse = i + usesInDec
			}

		case "WHILE":
			for _, ex := range stat.Branches[0].Search_tree("EXPRESSIOn") {
				for _, use := range LivenessCheckExpression(ex) {
					vars[use].lastuse = i + usesInDec
				}
			}

		case "IF":
			// Go Through Expression
			for _, use := range LivenessCheckExpression(stat.Search_first_child("EXPRESSION")) {
				vars[use].lastuse = i + usesInDec
			}
			// Add the values inside the if statement to the statement slice 
			ifstat := stat.Search_top_occurences("STATEMENT")
			statements = append(statements[:i], append(ifstat, statements[i:]...)...)
		
		case "FUNCCALL":
			for _, ex := range stat.Search_tree("EXPRESSION") {
				for _, use := range LivenessCheckExpression(ex) {
					vars[use].lastuse = i + usesInDec
				}
			}

		case "RETURN":
			for _, use := range LivenessCheckExpression(stat.Search_first_child("EXPRESSION")) {
				vars[use].lastuse = i + usesInDec
			}
		
		}
	}

	// Create the map from linenumber -> vars that need to be live on that line
	live_vars := map[int]int{}
	line_number := 0
	for _, localvar := range vars {
		// Also discards any var with lastused == -1
		if line_number >= localvar.firstuse && line_number < localvar.lastuse {
			live_vars[line_number] += 1
		} 
		line_number++
	}

	// Find the max size
	min_registers := 0 
	for _, size := range live_vars {
		if size > min_registers {
			min_registers = size
		}
	}
	fmt.Println("Minimum Registers: ", min_registers)
}

func LivenessCheckExpression (ex *parser.ParseTree) []string {
	vars_used := []string{}
	literals := ex.Search_tree("EL7")
	for _, literal := range literals {
		if literal.Branches[0].Leaf.Name == "name" {
			vars_used = append(vars_used, literal.Branches[0].Leaf.Value.(string))
		}
	}
	return vars_used
}