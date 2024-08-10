package main

import "fmt"

func RunGor(text, file string, isModuleImport, printTokens, printNodes, printVars, printVarsEachCycle bool) (ModuleImport, error) {
	lexer := NewLexer(text)
	tokens, lexerErr := lexer.Lex()
	if lexerErr != nil {
		return ModuleImport{}, lexerErr
	} else if printTokens {
		fmt.Println(tokens)
	}

	nodes, parseErr := Parse(tokens)
	if parseErr != nil {
		return ModuleImport{}, parseErr
	} else if printNodes {
		fmt.Println(nodes)
	}

	mod, interpretErr := Interpret(nodes, file, printVars, printVarsEachCycle)
	if interpretErr != nil {
		return ModuleImport{}, interpretErr
	}

	if isModuleImport {
		return mod, nil
	}
	return ModuleImport{}, nil
}
