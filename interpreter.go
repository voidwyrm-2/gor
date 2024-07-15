package main

import (
	"errors"
	"fmt"
	"reflect"
)

/*
func SprintRS(str string, a ...any) string {
	interopCount := strings.Count(str, "{}")
	if interopCount != len(a) {
		panic(fmt.Sprintf("expected %d substituitions, but was given %d instead", interopCount, len(a)))
	}

	for {
		ind := strings.Index(str, "{}")
		if ind == -1 {
			return str
		}
	}
}
*/

func NewGorError(t Token, msg string) error {
	return fmt.Errorf("error on line %d, col %d-%d: %s", t.Ln, t.Start+1, t.End+2, msg)
}

func AssignVar(vars *map[string]any, funcs *map[string]any, identTok Token, value any) error {
	if _, ok := (*funcs)[identTok.Lit]; ok {
		return NewGorError(identTok, fmt.Sprintf("cannot assign value '%v' to function '%s'", value, identTok.Lit))
	} else if e, isErr := value.(error); isErr {
		return e
	}

	(*vars)[identTok.Lit] = value
	return nil
}

func CallFunc(vars *map[string]any, funcs *map[string]any, identTok Token, value any) error {
	return nil
}

func AddLabel(labels *map[string]uint, i uint, nameTok Token) error {
	if _, ok := (*labels)[nameTok.Lit]; ok {
		return NewGorError(nameTok, fmt.Sprintf("cannot create label '%v' as it already exists", nameTok.Lit))
	}

	(*labels)[nameTok.Lit] = i

	return nil
}

func LabelJump(i *uint, labelTok Token, labels map[string]uint) error {
	if _, ok := labels[labelTok.Lit]; !ok {
		return NewGorError(labelTok, fmt.Sprintf("cannot jump to label '%v' as it doesn't exist", labelTok.Lit))
	}

	*i = labels[labelTok.Lit]

	return nil
}

func Interpret(nodes []Node, printVars, printVarsEachCycle bool) error {
	var vars = make(map[string]any)
	var funcs = make(map[string]any)
	funcs["puts"] = func(a ...any) {
		fmt.Println(a...)
	}
	var labels = make(map[string]uint)

	for i, node := range nodes {
		if n, ok := node.(LabelNode); ok {
			err := AddLabel(&labels, uint(i), n.Name)
			if err != nil {
				return err
			}
		}
	}

	var i uint = 0
	for i < uint(len(nodes)) {
		fmt.Println(i)
		node := nodes[i]
		if n, ok := node.(AssignmentNode); ok {
			err := AssignVar(&vars, &funcs, n.Ident, n.Value.Generate(&vars, &funcs))
			if err != nil {
				return err
			}
			i++
		} else if n, ok := node.(FunccallNode); ok {
			err := CallFunc(&vars, &funcs, n.Ident, n.GenerateArgs(&vars, &funcs))
			if err != nil {
				return err
			}
			i++
		} else if _, ok := node.(LabelNode); ok {
			i++
		} else if n, ok := node.(JumptoNode); ok {
			err := LabelJump(&i, n.LabelIdent, labels)
			if err != nil {
				return err
			}
			i++
		} else {
			return errors.New("unknown node '" + reflect.TypeOf(node).Name() + "'")
		}

		if printVarsEachCycle {
			fmt.Println(vars)
			for vname, vval := range vars {
				fmt.Printf("'%s': %v, '%s'\n", vname, vval, reflect.TypeOf(vval).Name())
			}
			fmt.Println("")
		}
	}

	if printVars && !printVarsEachCycle {
		fmt.Println(vars)
		for vname, vval := range vars {
			fmt.Printf("'%s': %v, '%s'\n", vname, vval, reflect.TypeOf(vval).Name())
		}
	}

	return nil
}
