package typechecker

import "errors"

func typeCheckExpression(expression ParseTree, funcName string, info TypeCheckerInfo) (string, error) {
	switch len(expression.Branches) {
	case 0:
		return "undefined", errors.New("Expression has 0 Children: " + expression.Leaf.Name)
	case 1:
		switch expression.Branches[0].Leaf.Name {
		case "EL1", "EL2", "EL3", "EL4", "EL5", "EL6", "EL7":
			return typeCheckExpression(expression.Branches[0], funcName, info)

		case "name":
			name := expression.Branches[0].Leaf.Value.(string)
			locals := info.LocalVar[funcName]
			ok := false
			local := Variable{}
			for _, l := range locals {
				if l.Name == name {
					ok = true
					local = l
					break
				}
			}
			if !ok {
				global, ok := info.GlobalVars[name]
				if !ok {
					return "undefined", errors.New("Variable " + name + " was not initialized.")
				}
				return global.Vartype, nil
			}
			return local.Vartype, nil

		case "LITERAL":
			switch expression.Branches[0].Branches[0].Leaf.Name {
			case "intliteral":
				return "int", nil
			case "doubleliteral":
				return "double", nil
			case "boolliteral":
				return "bool", nil
			case "stringliteral":
				return "string", nil
			default:
				return "undefined", errors.New("Literal error ~ Mt likely error in compiler :).\n Calculated Type: " + expression.Branches[0].Branches[0].Leaf.Name)
			}

		case "FUNCCALL":
			return info.Functions[expression.Branches[0].Branches[0].Leaf.Value.(string)].ReturnType, nil

		default:
			return "undefined", errors.New("Compiler Error, Expression without covered case has only one child")
		}

	case 2:
		switch expression.Branches[0].Leaf.Name {
		case "oplv5":
			ex, err := typeCheckExpression(expression.Branches[1], funcName, info)
			if err != nil {
				return "undefined", err
			}
			if ex == "int" || ex == "double" {
				return ex, nil
			}
			return "undefined", errors.New(expression.Branches[0].Leaf.Value.(string) + " did not stand before a number")
		case "oplv7":
			ex, err := typeCheckExpression(expression.Branches[1], funcName, info)
			if err != nil {
				return "undefined", err
			}
			if ex == "bool" {
				return "bool", nil
			}
			return "undefined", errors.New("! did not stand before a boolen expression")
		}

	case 3:
		leftside, err := typeCheckExpression(expression.Branches[0], funcName, info)
		if err != nil {
			return "undefined", err
		}
		rightside, err := typeCheckExpression(expression.Branches[2], funcName, info)
		if err != nil {
			return "undefined", err
		}
		switch expression.Branches[1].Leaf.Name {
		case "oplv1", "oplv2":
			if rightside == "bool" && leftside == "bool" {
				return "bool", nil
			}
		case "oplv3":
			if (rightside == "bool" && leftside == "bool") || ((rightside == "int" || rightside == "double") && (leftside == "int" || leftside == "double")) || (rightside == "string" && leftside == "string"){
				return "bool", nil
			}
		case "oplv4":
			if (rightside == "int" || rightside == "double") && (leftside == "int" || leftside == "double") {
				return "bool", nil
			}

		case "oplv5", "oplv6":
			if rightside == "int" && leftside == "int" {
				return "int", nil
			}

			if (rightside == "int" || rightside == "double") && (leftside == "int" || leftside == "double") {
				return "double", nil
			}
		default:
			return "undefined", errors.New("Compiler Error: Expression with 3 children does not have the correct opperators")
		}
	}
	return "undefined", nil
}
