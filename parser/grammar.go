package parser

type Grammar struct{
	start string
	nonTerminals []string
	terminals []string
	rules []Rule
	follow map[string][]string
	closure map[string][]Rule
}

/*type GrammarFollow struct{
	nullable map[string]bool
	first map[string][]string
	follow map[string][]string
}

*/

// Epsilon -> " "
type Rule struct{
	nonTerminal string
	production []string
}

func (grammar *Grammar) addSymbol(s string){
	grammar.nonTerminals = append(grammar.nonTerminals, s)
}

func MakeRule (nonTerminal string, production []string) Rule{
	newRule := new(Rule)
	newRule.nonTerminal = nonTerminal
	newRule.production = append(newRule.production, production...)
	return *newRule
}

func (grammar *Grammar) AddRule (nonTerminal string, production []string) Rule{
	newRule := MakeRule(nonTerminal, production)
	grammar.rules = append(grammar.rules, newRule)
	return newRule
}

func (grammar *Grammar) NULLABLE() *map[string]bool{
	// TODO: CALC NULLABLE
}

func (grammar *Grammar) FIRST(nonTerminal string, nullable map[string]bool) []string{
	// TODO: CALC FIRST
}

func (grammar *Grammar) FOLLOW(nonTerminal string, nullable map[string]bool, first map[string][]string) []string{
	// TODO: CALC FOLLOW
}


func makeStandardGrammar() *Grammar{
	newGrammar := new(Grammar)
	newGrammar.start = "S"
	newGrammar.nonTerminals = append(newGrammar.nonTerminals, "S")
	return newGrammar
}


