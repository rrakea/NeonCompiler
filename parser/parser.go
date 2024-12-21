package parser

import (
	"compiler/lexer"
	"fmt"
	"strconv"
)

func createParser(test bool) (*SLR_parsing_Table, *Grammar) {
	rules := defGrammar(test)
	grammar := MakeGrammar(rules, "START")
	grammar.Augment()
	first := grammar.FIRST()
	first["S"] = append(first["S"], "namespace") // TODO
	first["START"] = append(first["START"], "namespace")
	//PrintFirst(first)
	follow := grammar.FOLLOW(first)
	grammar.follow = follow
	grammar.CalcClosure()

	automata := grammar.CreateSLRAutomata()
	table := automata.CreateSLRTable(grammar)

	return table, grammar
}

func Parse(path string, test bool) (ParseTree, bool) {
	parseTreeChannel := make(chan any)
	go createParseTree(parseTreeChannel)

	tokenChannel := make(chan lexer.Token)
	slrTable, grammar := createParser(test)

	go lexer.Lex(path, tokenChannel)

	linecount := 0
	stack := makeStack(0)
	accepts := false

	for true {
		token := lexer.GetNext(tokenChannel)

		if token.Identifier == "" {
			if accepts {
				break
			} else {
				parseError(*token, linecount, *stack, slrTable, parseTreeChannel)
			}
		}
		if token.Identifier == "LINE" {
			linecount, _ = strconv.Atoi(token.Value.(string))
			continue
		}

		// Do only once, unless reduce is found
		for i := 0; i < 1; i++ {
			stackVal := stack.peek().(*any)
			res, err := slrTable.GetAction((*stackVal).(int), token.Identifier)
			if err != nil {
				parseError(*token, linecount, *stack, slrTable, parseTreeChannel)
				return ParseTree{}, false
			}
			switch res.actionType {
			case "Shift":
				parseTreeChannel <- *token
				stack.add(token)
				stack.add(res.value)
			case "Reduce":
				// Redo the loop ~ Dont get another input symbol
				i--
				// Get the Rule that we reduce by
				reductionRule := grammar.rules[res.value]
				parseTreeChannel <- reductionRule
				for range reductionRule.production {
					stack.pop()
					stack.pop()
				}
				stateBefore := stack.peek().(*any)
				gotoVal, err := slrTable.GetGoto((*stateBefore).(int), reductionRule.nonTerminal)
				if err != nil {
					parseError(*token, linecount, *stack, slrTable, parseTreeChannel)
					return ParseTree{}, false
				}
				stack.add(reductionRule.nonTerminal)
				stack.add(gotoVal.val)
			case "Accept":
				accepts = true
			}
		}
	}
	if accepts {
		parseTreeChannel <- true
		fmt.Println("Code passed parser")
		tree := <- parseTreeChannel
		return tree.(ParseTree), true
	}
	parseError(lexer.Token{}, linecount, *stack, slrTable, parseTreeChannel)
	return ParseTree{}, false
}

func parseError(token lexer.Token, linecount int, stack Stack, table *SLR_parsing_Table, parseTreeChannel chan any) {
	parseTreeChannel <- false
	select {
	case tree := <-parseTreeChannel:
		PrintTree(tree.(ParseTree))
	}
	fmt.Println("DID NOT PASS")

	lineString := strconv.Itoa(linecount)
	state := stack.peek().(*any)
	next := table.getNextExpectedTokens((*state).(int))

	conv, err := token.Value.(int)
	if err && conv != 0 {
		next = append(next, fmt.Sprintf("%v", token.Value))
	}

	nextString := formatNext(next)

	if token.Identifier == "$" {
		fmt.Println("Unexpected end of file reached. At line: " + lineString + ".\nExpecting: " + nextString)
		return
	}
	unexpected := formatToken(token)

	fmt.Println("Syntax Error: Unexpected: \"" + unexpected + "\" at line " + lineString + ".\nExpecting: " + nextString)
}

func formatNext(next []string) string {
	returnstring := ""
	for _, n := range next {
		nextString := ""
		switch n {
		case "logicaloperator":
			nextString = "==, ||, &&, >=, <=, =="
		case "multoperator":
			nextString = "*, /"
		case "unaryoperator":
			nextString = "+, -"
		case "name":
			nextString = "descriptor"
		case "intliteral", "doubleliteral":
			nextString = "num"
		case "boolliteral":
			nextString = "bool"
		case "$":
			nextString = "End of file"
		default:
			nextString = n
		}
		returnstring += nextString + "," + " "
		nextString = ""
	}
	if len(returnstring) >= 2{
		return returnstring[:len(returnstring)-2]
	}
	return ""
}

func formatToken(token lexer.Token) string {
	switch token.Identifier {
	case "name", "logicaloperator", "multoperator", "unaryoperator", "boolliteral":
		return token.Value.(string)
	case "intliteral":
		return strconv.Itoa(token.Value.(int))
	case "doubleliteral":
		return strconv.FormatFloat(token.Value.(float64), 'f', -1, 64)
	case "stringliteral":
		return "string"
	default:
		return token.Identifier
	}
}
