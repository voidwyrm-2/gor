package main

import "fmt"

func RunGor(text string, printTokens, printNodes, printVars, printVarsEachCycle bool) {
	lexer := NewLexer(text)
	tokens, lexerErr := lexer.Lex()
	if lexerErr != nil {
		fmt.Println(lexerErr.Error())
		return
	} else if printTokens {
		fmt.Println(tokens)
	}

	nodes, parseErr := Parse(tokens)
	if parseErr != nil {
		fmt.Println(parseErr.Error())
		return
	} else if printNodes {
		fmt.Println(nodes)
	}

	interpretErr := Interpret(nodes, printVars, printVarsEachCycle)
	if interpretErr != nil {
		fmt.Println(interpretErr.Error())
		return
	}
}
