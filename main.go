package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"slices"
	"strings"
)

/*
func ConvertVersionToInts(v string) ([3]int, error) {
	split := strings.Split(v, ".")
	str_1 := split[0]
	str_2 := split[1]
	str_3 := split[2]

	_1, err := strconv.Atoi(str_1)
	if err != nil {
		return [3]int{}, err
	}
	_2, err := strconv.Atoi(str_2)
	if err != nil {
		return [3]int{}, err
	}
	_3, err := strconv.Atoi(str_3)
	if err != nil {
		return [3]int{}, err
	}

	return [3]int{_1, _2, _3}, nil
}

/*
Returns 1 if v1 is bigger than v2, -1 if v1 is smaller than v2, and 0 if equal
* /
func CompareVersions(version_1, version_2 string) (int, error) {
	v1, err := ConvertVersionToInts(version_1)
	if err != nil {
		return 0, err
	}
	v2, err := ConvertVersionToInts(version_2)
	if err != nil {
		return 0, err
	}

	for i := range 3 {
		if v1[i] > v2[i] {
			return 1, nil
		} else if v1[i] < v2[i] {
			return -1, nil
		}
	}

	return 0, nil
}


func GetGorVersion() string {
	res, err := http.Get("https://raw.githubusercontent.com/voidwyrm-2/gor/main/gor_version.txt")
	if err != nil {
		log.Fatal(err.Error())
	}

	version, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err.Error())
	} else if string(version) == "404: Not Found" {
		log.Fatal("'https://raw.githubusercontent.com/voidwyrm-2/gor/main/gor_version.txt' not found")
	}

	return string(version)
}

func CheckCurrentGorVersion(localVersion string) {
	nonLocalVersion := GetGorVersion()
	cmp, err := CompareVersions(localVersion, nonLocalVersion)
	if err != nil {
		log.Fatal(err.Error())
	} else if cmp == -1 {
		fmt.Println("there's a new version of Gor availible!")
	}
}
*/

/*
func assertNoError[T any](v T, _ error) T {
	return v
}
*/

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

/*
func writeFile(filename string, data string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(data)
	if err != nil {
		return err
	}

	return nil
}
*/

func GorREPL(printTokens, printNodes, printVarsEachCycle bool) {
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
			_, err := RunGor(strings.Join(codeBuffer, "\n"), "<main>", false, printTokens, printNodes, true, printVarsEachCycle)
			if err != nil {
				fmt.Println(err.Error())
				codeBuffer = codeBuffer[:len(codeBuffer)-1]
			}
		}
	}
}

/*
reminder of how my versioning system works(because I'm forgetful):

[major version].[minor version].[sub minor version]

major version: increment when there are major changes

minor version: increment when there are minor changes

sub minor version: increment when there are very small changes
(e.g. a quick one-line grammar change, a change to a variable name)
*/
var GOR_VERSION = "0.4.1"

/*
func gorInit(allowRecursion bool) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err.Error())
	}

	onlineGorVersion := GetGorVersion()

	cachedGorVersion, err := readFile(path.Join(homeDir, "gor_version.txt"))
	if err != nil {
		if isFile404Err(err.Error()) && allowRecursion {
			writeFile(path.Join(homeDir, "gor_version.txt"), onlineGorVersion)
			gorInit(false)
			return
		}
		log.Fatal(err.Error())
	}

	if assertNoError(CompareVersions(strings.TrimSpace(onlineGorVersion), strings.TrimSpace(cachedGorVersion))) == 1 {
		err := writeFile(path.Join(homeDir, "gor_version.txt"), onlineGorVersion)
		if err != nil {
			log.Fatal(err.Error())
		} else if allowRecursion {
			gorInit(false)
			return
		}
	}

	GOR_VERSION = strings.TrimSpace(cachedGorVersion)
}
*/

func main() {
	//gorInit(true)

	printTokens := slices.Contains(os.Args, "-t")
	printNodes := slices.Contains(os.Args, "-n")
	printVars := slices.Contains(os.Args, "-v")
	printVarsEachCycle := slices.Contains(os.Args, "-cv")

	if slices.Contains(os.Args, "--version") {
		fmt.Println("Gor version " + GOR_VERSION)
		//CheckCurrentGorVersion(GOR_VERSION)
		return
	}

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
			GorREPL(printTokens, printNodes, printVarsEachCycle)
		case "help":
			fmt.Println(strings.Join([]string{
				"'help': shows this text",
				"'exit/quit': exits the interpreter",
				"'repl': starts the Gor repl",
				"'run [path]': runs a .gor file",
				"cli flags:",
				" --version: shows the current Gor version",
				" -t: print lexer tokens",
				" -n: print AST nodes",
				" -v: print variables after execution of all code",
				" -cv: print variables after execution of each node in the AST(this overrides -v)",
			}, "\n"))
		default:
			if len(command) > 4 && command[:4] == "run " {
				fileName := strings.TrimSpace(command[4:])

				if len(fileName) < 4 {
					if path.Ext(fileName) == "" {
						fileName += ".gor"
					} else if path.Ext(fileName) != ".gor" {
						fmt.Printf("file '%s' is not a .gor file\n", fileName)
						continue
					}
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

				RunGor(content, fileName, false, printTokens, printNodes, printVars, printVarsEachCycle)
			} else {
				fmt.Println("unknown command")
			}
		}
	}
}
