package parser

func defGrammar(test bool) []Rule {
	grammar := []Rule{}
	if test {
		grammar = append(grammar, testGrammar()...)
	} else {
		grammar = append(grammar, compilerGrammar()...)
	}
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
		MakeRule("LITERAL", []string{"NUMLITERAL"}),
		MakeRule("NUMLITERAL", []string{"intliteral"}),
		MakeRule("NUMLITERAL", []string{"doubleliteral"}),

		MakeRule("START", []string{"USINGBLOCK"}),
		MakeRule("USINGBLOCK", []string{"using", "name", ";", "USINGBLOCK"}),
		MakeRule("USINGBLOCK", []string{"NAMESPACE"}),
		MakeRule("NAMESPACE", []string{"namespace", "name", "{", "CLASS"}),
		MakeRule("CLASS", []string{"class", "name", "{", "GLOBALVARBLOCK", "}"}),
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
		MakeRule("STATEMENT", []string{"ARRAYASSIGN"}),
		MakeRule("STATEMENT", []string{"IF"}),
		MakeRule("STATEMENT", []string{"WHILE"}),

		MakeRule("INPUTBLOCK", []string{")"}),
		MakeRule("INPUTBLOCK", []string{"string", "[", "]", "name", ")"}),
		MakeRule("INPUTBLOCK", []string{"INPUTSTART"}),
		MakeRule("INPUTSTART", []string{"TYPE", "name", "INPUTCONTINUED"}),
		MakeRule("INPUTCONTINUED", []string{",", "TYPE", "name"}),
		MakeRule("INPUTCONTINUED", []string{")"}),
		MakeRule("ARGBLOCK", []string{")"}),
		MakeRule("ARGBLOCK", []string{"ARGSSTART"}),
		MakeRule("ARGSSTART", []string{"EXPRESSION", "ARGCONTINUED"}),
		MakeRule("ARGCONTINUED", []string{")"}),
		MakeRule("ARGCONTINUED", []string{",", "EXPRESSION", "ARGCONTINUED"}),
		MakeRule("FUNCCALL", []string{"name", "(", "ARGBLOCK"}),
		MakeRule("FUNCCALL", []string{"name", ".", "name", "(", "ARGBLOCK"}),
		MakeRule("ARRAYACCESS", []string{"name", "[", "EXPRESSION", "]"}),
		MakeRule("RETURN", []string{"return", "EXPRESSION", ";"}),
		MakeRule("RETURN", []string{"return", ";"}),
		MakeRule("VARASSIGN", []string{"name", "=", "EXPRESSION", ";"}),
		MakeRule("ARRAYASSIGN", []string{"name", "[", "EXPRESSION", "]", "=", "EXPRESSION", ";"}),
		MakeRule("IF", []string{"if", "(", "EXPRESSION", ")", "{", "STATEMENTBLOCK"}),
		MakeRule("IF", []string{"if", "(", "EXPRESSION", ")", "{", "STATEMENTBLOCK", "ELSE"}),
		MakeRule("IF", []string{"if", "(", "EXPRESSION", ")", "STATEMENT"}),
		MakeRule("ELSE", []string{"else", "{", "STATEMENTBLOCK"}),
		MakeRule("WHILE", []string{"while", "(", "EXPRESSION", ")", "{", "STATEMENTBLOCK"}),
		MakeRule("WHILE", []string{"while", "(", "EXPRESSION", ")", "STATEMENT"}),

		MakeRule("EXPRESSION", []string{"EXPRESSION", "logicaloperator", "TERM"}),
		//MakeRule("EXPRESSION", []string{"!", "EXPRESSION"}),
		MakeRule("EXPRESSION", []string{"TERM"}),
		MakeRule("TERM", []string{"TERM", "unaryoperator", "FACTOR"}),
		MakeRule("TERM", []string{"FACTOR"}),
		MakeRule("FACTOR", []string{"FACTOR", "multoperator", "PRIMARY"}),
		MakeRule("FACTOR", []string{"PRIMARY"}),
		MakeRule("PRIMARY", []string{"FUNCCALL"}),
		MakeRule("PRIMARY", []string{"ARRAYACCESS"}),
		MakeRule("PRIMARY", []string{"LITERAL"}),
		MakeRule("PRIMARY", []string{"name"}),
		MakeRule("PRIMARY", []string{"!", "PRIMARY"}),
		MakeRule("PRIMARY", []string{"unaryoperator", "PRIMARY"}),
		MakeRule("PRIMARY", []string{"(", "EXPRESSION", ")"}),
	}
	return rules
}

func compilerGrammar() []Rule {
	return []Rule{}
}
