package lexer

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"unicode"
)

type Token struct {
	Identifier string
	Value      any
}

// Runs as go routine; called by the parser
func GetNext(tokenChannel chan Token) (*Token, error) {
	// Wait until the channel with tokens has a value inside
	select {
	case newToken, ok := <-tokenChannel:
		// If channel is closed -> File is done
		if !ok {
			return nil, errors.New("Lexer Error: File Ended")
		}
		return &newToken, nil
	}
}

func Lex(path string, tokenChannel chan Token) {
	fmt.Println("Started Lexing...")
	fmt.Println()
	// Open File
	file, err := os.Open(path)
	if err != nil {
		panic("Lexer Error: File not able to be opened. Likely to be the wrong path. Path given: " + path)
	}
	defer file.Close()

	// Scan over the file
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	lineNumber := 1

	isMultiLineComment := false

	for scanner.Scan() {
		line := scanner.Text()
		// Split String, removing whitespace etc.
		tokens := []string{"\n"}
		buffer := ""
		isString := false
		isSymbolString := false
		isSingleLineComment := false

		for _, c := range line {
			switch {
			case isSingleLineComment:
				continue

			case isMultiLineComment:
				buf := string(c)
				if buf == "*" {
					buffer = buf
				} else if buf == "/" && buffer == "*" {
					buffer = ""
					isMultiLineComment = false
				} else {
					buffer = ""
					continue
				}

			case c == '"':
				if !isString {
					isString = true
				} else {
					isString = false
					tokens = append(tokens, buffer)
					buffer = ""
				}

			case isString:
				buffer = buffer + string(c)

			case isSymbolString:
				if !isSymbol(c) {
					tokens = append(tokens, buffer)
					buffer = ""
					if (!unicode.IsSpace(c)){
						buffer = string(c)
					}
					isSymbolString = false
					continue
				}

				// Is symbol -> Can be concatonated to // /* etc.
				concSymbol := concatonateSymbols([]rune(buffer)[0], c)

				// Symbols cannot be concatonated
				if concSymbol == "" {
					tokens = append(tokens, buffer)
					buffer = string(c)
					continue
				}

				// Symbols can be concatonated -> check for if comment
				isSymbolString = false
				if concSymbol == "//" {
					isSingleLineComment = true
				} else if concSymbol == "/*" {
					isMultiLineComment = true
				} else {
					tokens = append(tokens, concSymbol)
					buffer = ""
					continue
				}

				// If this line is reached a comment has started
				buffer = ""

			case isSymbol(c):
				if buffer != ""{
					tokens = append(tokens, buffer)
				}
				isSymbolString = true
				buffer = string(c)

			case unicode.IsSpace(c):
				if buffer != "" {
					tokens = append(tokens, buffer)
				}
				buffer = ""

			default:
				buffer = buffer + string(c)
			}
		}

		if buffer != "" {
			tokens = append(tokens, buffer)
		}

		// Determine Identifier
		for _, token := range tokens {
			identifier := ""
			var tokenVal any

			tmpdigit, intConvErr := strconv.Atoi(token)

			// If could be converted to int
			if intConvErr == nil {
				sendToken("INTEGER_LITERAL", tmpdigit, tokenChannel)
				continue
			}

			tmpbool, boolConvErr := strconv.ParseBool(token)

			// If could be converted to bool
			if boolConvErr == nil {
				sendToken("BOOL_LITERAL", tmpbool, tokenChannel)
				continue
			}

			// Check for the different symbols
			switch token {
			case "\n":
				identifier = "LINE"
				tokenVal = strconv.Itoa(lineNumber)
				lineNumber++
			case "namespace":
				identifier = "NAMESPACE"
			case "using":
				identifier = "USING"
			case "class":
				identifier = "CLASS"
			case "void":
				identifier = "VOID"
			case "static":
				identifier = "STATIC"
			case "Main":
				identifier = "MAIN"
			case "int":
				identifier = "INT"
			case "bool":
				identifier = "BOOL"
			case "string":
				identifier = "STRING"
			case "double":
				identifier = "DOUBLE"
			case "if":
				identifier = "IF"
			case "else":
				identifier = "ELSE"
			case "while":
				identifier = "WHILE"
			case "return":
				identifier = "RETURN"
			case ".":
				identifier = "DOT"
			case ",":
				identifier = "COMMA"
			case "|":
				identifier = "OR"
			case "=":
				identifier = "EQUALS"
			case ">":
				identifier = "BIGGER_THAN"
			case "<":
				identifier = "SMALLER_THAN"
			case "+":
				identifier = "PLUS"
			case "*":
				identifier = "STAR"
			case ";":
				identifier = "SEMICOLON"
			case "{":
				identifier = "LEFT_BRACE"
			case "}":
				identifier = "RIGHT_BRACE"
			case "(":
				identifier = "LEFT_PARENTHESIS"
			case ")":
				identifier = "RIGHT_PARENTHESIS"
			case "[":
				identifier = "LEFT_BRACKET"
			case "]":
				identifier = "RIGHT_BRACKET"
			default:
				identifier = "NAME"
				tokenVal = token
			}
			if isMultiLineComment || isSingleLineComment {
				if isSingleLineComment {
				}
				isSingleLineComment = false
			} else {
				sendToken(identifier, tokenVal, tokenChannel)
			}
		}
	}
	close(tokenChannel)
}

func sendToken(identifier string, value any, channel chan Token){
	// Make return token and add to channel
	returnToken := new(Token)
	returnToken.Identifier = identifier
	returnToken.Value = value
	channel <- *returnToken
}

func isSymbol(r rune) bool {
	symbols := []rune{';', '.', '-', '+', '*', '>', '<', '=', '{', '}', '(', ')', '[', ']', '|', ',', '/'}
	for _, symbol := range symbols {
		if r == symbol {
			return true
		}
	}
	return false
}

func concatonateSymbols(s1 rune, s2 rune) string {
	if s2 == '=' {
		if s1 == '>' || s1 == '<' || s1 == '!' || s1 == '=' {
			return string(s1) + string(s2)
		}
	}
	if s1 == '/' {
		if s2 == '/' || s2 == '*' {
			return string(s1) + string(s2)
		}
	}
	if s1 == '*' && s2 == '/' {
		return "*/"
	}
	if s1 == '|' && s2 == '|' {
		return "||"
	}
	if s1 == '&' && s2 == '&' {
		return "&&"
	}
	return ""
}
