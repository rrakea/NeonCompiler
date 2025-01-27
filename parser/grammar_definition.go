package parser

func defGrammar(test bool) []Rule {
	_ = test
	grammar := []Rule{}
	grammar = append(grammar, testGrammar()...)
	return grammar
}

func testGrammar() []Rule {
	rules := []Rule{
		MakeRule("TYPE", []string{"double"}),
		MakeRule("TYPE", []string{"int"}),
		MakeRule("TYPE", []string{"string"}),
		MakeRule("TYPE", []string{"bool"}),
		//
		MakeRule("RETURNTYPE", []string{"void"}),
		MakeRule("RETURNTYPE", []string{"TYPE"}),
		//
		MakeRule("LITERAL", []string{"stringliteral"}),
		MakeRule("LITERAL", []string{"boolliteral"}),
		MakeRule("LITERAL", []string{"intliteral"}),
		MakeRule("LITERAL", []string{"doubleliteral"}),

		MakeRule("START", []string{"USING"}),
		MakeRule("USING", []string{"using", "name", ";", "USING"}),
		MakeRule("USING", []string{"NAMESPACE"}),
		MakeRule("NAMESPACE", []string{"namespace", "name", "{", "CLASS", "}"}),
		MakeRule("CLASS", []string{"class", "name", "{", "GLOBALVARBLOCK"}),
		MakeRule("GLOBALVARBLOCK", []string{"static", "TYPE", "name", "=", "EXPRESSION", ";", "GLOBALVARBLOCK"}),
		MakeRule("GLOBALVARBLOCK", []string{"MAIN"}),
		MakeRule("MAIN", []string{"static", "void", "main", "(", "string", "[", "]", "name", ")", "{", "VIRTUALVARBLOCK", "FUNCBLOCK"}),
		MakeRule("FUNCBLOCK", []string{"FUNC", "FUNCBLOCK"}),
		MakeRule("FUNCBLOCK", []string{"}"}),

		MakeRule("FUNC", []string{"static", "RETURNTYPE", "name", "(", "INPUTBLOCK", "{", "VIRTUALVARBLOCK"}),
		MakeRule("VIRTUALVARBLOCK", []string{"TYPE", "name", "=", "EXPRESSION", ";", "VIRTUALVARBLOCK"}),
		MakeRule("VIRTUALVARBLOCK", []string{"STATEMENTBLOCK"}),
		MakeRule("STATEMENTBLOCK", []string{"STATEMENT", "STATEMENTBLOCK"}),
		MakeRule("STATEMENTBLOCK", []string{"}"}),
		MakeRule("STATEMENT", []string{"FUNCCALL", ";"}),
		MakeRule("STATEMENT", []string{"RETURN"}),
		MakeRule("STATEMENT", []string{"VARASSIGN"}),
		MakeRule("STATEMENT", []string{"IF"}),
		MakeRule("STATEMENT", []string{"WHILE"}),

		MakeRule("INPUTBLOCK", []string{")"}),
		MakeRule("INPUTBLOCK", []string{"PARAMETER", "INPUTCONTINUED"}),
		MakeRule("INPUTCONTINUED", []string{",", "PARAMETER", "INPUTCONTINUED"}),
		MakeRule("PARAMETER", []string{"TYPE", "name"}),
		MakeRule("INPUTCONTINUED", []string{")"}),

		MakeRule("ARGBLOCK", []string{")"}),
		MakeRule("ARGBLOCK", []string{"ARG", "ARGCONTINUED"}),
		MakeRule("ARG", []string{"EXPRESSION"}),
		MakeRule("ARGCONTINUED", []string{")"}),
		MakeRule("ARGCONTINUED", []string{",", "ARG", "ARGCONTINUED"}),
		MakeRule("FUNCCALL", []string{"name", "(", "ARGBLOCK"}),
		MakeRule("RETURN", []string{"return", "EXPRESSION", ";"}),
		MakeRule("RETURN", []string{"return", ";"}),
		MakeRule("VARASSIGN", []string{"name", "=", "EXPRESSION", ";"}),

		// TODO else without {}
		MakeRule("IF", []string{"if", "(", "EXPRESSION", ")", "{", "STATEMENTBLOCK"}),
		MakeRule("IF", []string{"if", "(", "EXPRESSION", ")", "{", "STATEMENTBLOCK", "ELSE"}),
		MakeRule("IF", []string{"if", "(", "EXPRESSION", ")", "STATEMENT"}),
		//MakeRule("IF", []string{"if", "(", "EXPRESSION", ")", "STATEMENT", "ELSE"}),
		MakeRule("ELSE", []string{"else", "{", "STATEMENTBLOCK"}),
		MakeRule("ELSE", []string{"else", "STATEMENT"}),
		MakeRule("WHILE", []string{"while", "(", "EXPRESSION", ")", "{", "STATEMENTBLOCK"}),
		MakeRule("WHILE", []string{"while", "(", "EXPRESSION", ")", "STATEMENT"}),

		MakeRule("EXPRESSION", []string{"EL1"}),
		MakeRule("EL1", []string{"EL1", "oplv1", "EL2"}),
		MakeRule("EL1", []string{"EL2"}),
		MakeRule("EL2", []string{"EL2", "oplv2", "EL3"}),
		MakeRule("EL2", []string{"EL3"}),
		MakeRule("EL3", []string{"EL3", "oplv3", "EL4"}),
		MakeRule("EL3", []string{"EL4"}),
		MakeRule("EL4", []string{"EL4", "oplv4", "EL5"}),
		MakeRule("EL4", []string{"EL5"}),
		MakeRule("EL5", []string{"EL5", "oplv5", "EL6"}),
		MakeRule("EL5", []string{"EL6"}),
		MakeRule("EL6", []string{"EL6", "oplv6", "EL7"}),
		MakeRule("EL6", []string{"EL7"}),
		// +/ - before value
		MakeRule("EL7", []string{"oplv5", "EL7"}),
		MakeRule("EL7", []string{"oplv7", "EL7"}),

		MakeRule("EL7", []string{"(", "EL1", ")"}),
		MakeRule("EL7", []string{"LITERAL"}),
		MakeRule("EL7", []string{"FUNCCALL"}),
		MakeRule("EL7", []string{"name"}),
	}
	return rules
}
