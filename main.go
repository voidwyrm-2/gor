package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"slices"
	"strings"
)

func ChangeNestedMapItem(m *map[any]any, dictPath []any, value any, unwrap bool) {
	var a map[any]any
	if unwrap {
		ab := *m
		a = ab["0"].(map[any]any)
	} else {
		a = *m
	}

	if len(dictPath) > 0 {
		if len(dictPath) == 1 {
			a[dictPath[0]] = value
		} else {
			inner := a[dictPath[0]]
			ChangeNestedMapItem(&map[any]any{"0": inner}, dictPath[1:], value, true)
			a[dictPath[0]] = inner
		}
	}

	*m = a
}

func readFile(fileName string) (string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	content := ""
	for scanner.Scan() {
		content += scanner.Text() + "\n"
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return content, nil
}

func GorREPL(printTokens, printNodes, printVars, printVarsEachCycle bool) {
	var codeBuffer []string
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(">>> ")
		scanner.Scan()
		input := strings.TrimSpace(scanner.Text())

		if input == "--exit" || input == "--quit" {
			return
		} else {
			if input != "" {
				codeBuffer = append(codeBuffer, input)
			}
			RunGor(strings.Join(codeBuffer, "\n"), printTokens, printNodes, printVars, printVarsEachCycle)
		}
	}
}

func main() {
	printTokens := slices.Contains(os.Args, "-t")
	printNodes := slices.Contains(os.Args, "-n")
	printVars := slices.Contains(os.Args, "-v")
	printVarsEachCycle := slices.Contains(os.Args, "-cv")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		scanner.Scan()
		command := strings.TrimSpace(scanner.Text())

		switch command {
		case "exit", "quit":
			return
		case "repl":
			fmt.Println("Gor REPL(type '--exit' or '--quit' to end the repl)")
			GorREPL(printTokens, printNodes, printVars, printVarsEachCycle)
		case "help":
			fmt.Println(strings.Join([]string{
				"'help': shows this text",
				"'exit/quit': exits the interpreter",
				"'repl': starts the Gor repl",
				"'run [path]': runs a .gor file",
				"cli flags:",
				" -t: print lexer tokens",
				" -n: print AST nodes",
				" -v: print variables after execution of all code",
				" -cv: print variables after execution of each node in the AST(this overrides -v)",
			}, "\n"))
		default:
			if len(command) > 4 && command[:4] == "run " {
				fileName := command[4:]

				if len(fileName) < 4 {
					if path.Ext(fileName) != "" {
						fmt.Printf("file '%s' is not a .gor file\n", fileName)
						continue
					}
					fileName += ".gor"
				} else if fileName[:len(fileName)-4] != ".gor" {
					if path.Ext(fileName) != "" {
						fmt.Printf("file '%s' is not a .gor file\n", fileName)
						continue
					}
					fileName += ".gor"
				}

				content, err := readFile(fileName)
				if err != nil {
					fmt.Println(err.Error())
					continue
				}

				RunGor(content, printTokens, printNodes, printVars, printVarsEachCycle)
			} else {
				fmt.Println("unknown command")
			}
		}
	}
}
