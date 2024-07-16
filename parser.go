package main

import (
	"fmt"
	"strconv"
	"strings"
)

type GorType interface {
	string | int | float32 | rune
}

type Node any

type JumptoNode struct {
	LabelIdent Token
}

type LabelNode struct {
	Name Token
}

type ModuleImportNode struct {
	PathIdent Token
}

type AssignableValue interface {
	Generate(*map[string]any, *map[string]any) any
}

type FunccallNode struct {
	Ident Token
	args  []AssignableValue
}

func (fn FunccallNode) GenerateArgs(vars *map[string]any, funcs *map[string]any) []any {
	var out []any
	for _, a := range fn.args {
		out = append(out, a.Generate(vars, funcs))
	}
	return out
}

type AssignmentNode struct {
	Ident Token
	Value AssignableValue
}

func RemoveNewlineTokens(tokens []Token) []Token {
	var out []Token
	for _, t := range tokens {
		if t.Type != NEWLINE {
			out = append(out, t)
		}
	}
	return out
}

func GenerateExpressionNodeFromTokens(tokens []Token) AssignableValue {
	tokens = RemoveNewlineTokens(tokens)
	//fmt.Println(tokens)
	if len(tokens) == 1 {
		return ValueNode{Val: tokens[0]}
	}
	if tokens[len(tokens)-1].Type == EOF {
		tokens = tokens[:len(tokens)-1]
	}
	/*
		if !tokens[int(len(tokens)/2)].IsMathOp() {
			tokens = append(tokens, NewToken(PLUS, "+", tokens[len(tokens)-1].Start+2, tokens[len(tokens)-1].End+2, tokens[len(tokens)-1].Ln))
			if tokens[0].Type == STRING {
				tokens = append(tokens, NewToken(STRING, "", tokens[len(tokens)-2].Start+2, tokens[len(tokens)-2].End+2, tokens[len(tokens)-1].Ln))
			} else {
				tokens = append(tokens, NewToken(NUMBER, "0", tokens[len(tokens)-2].Start+2, tokens[len(tokens)-2].End+2, tokens[len(tokens)-1].Ln))
			}
		}
	*/

	//fmt.Println(tokens)

	op := tokens[int(len(tokens)/2)]
	left := tokens[:int(len(tokens)/2)]
	right := tokens[int(len(tokens)/2)+1:]

	return ExpressionNode{Left: GenerateExpressionNodeFromTokens(left), Operand: op, Right: GenerateExpressionNodeFromTokens(right)}
}

type ExpressionNode struct {
	Left    AssignableValue
	Operand Token
	Right   AssignableValue
}

func (expr ExpressionNode) Generate(vars *map[string]any, funcs *map[string]any) any {
	left := expr.Left.Generate(vars, funcs)
	_, leftIsString := left.(string)
	_, leftIsInt := left.(int)
	_, leftIsFloat32 := left.(float32)
	_, leftIsError := left.(error)
	right := expr.Right.Generate(vars, funcs)
	_, rightIsString := right.(string)
	_, rightIsInt := right.(int)
	_, rightIsFloat32 := right.(float32)
	_, rightIsError := right.(error)

	if leftIsError {
		return left
	} else if rightIsError {
		return right
	}

	switch expr.Operand.Type {
	case PLUS:
		if leftIsString && rightIsString {
			return left.(string) + right.(string)
		} else if leftIsInt && rightIsInt {
			return left.(int) + right.(int)
		} else if leftIsFloat32 && rightIsFloat32 {
			return left.(float32) + right.(float32)
		}
	case HYPHEN:
		if leftIsInt && rightIsInt {
			return left.(int) - right.(int)
		} else if leftIsFloat32 && rightIsFloat32 {
			return left.(float32) - right.(float32)
		}
	case ASTERISK:
		if leftIsString && rightIsInt {
			acc := ""
			for range right.(int) {
				acc += left.(string)
			}
			return acc
		} else if leftIsInt && rightIsInt {
			return left.(int) * right.(int)
		} else if leftIsFloat32 && rightIsFloat32 {
			return left.(float32) * right.(float32)
		}
	case FORWARD_SLASH:
		if leftIsInt && rightIsInt {
			return left.(int) / right.(int)
		} else if leftIsFloat32 && rightIsFloat32 {
			return left.(float32) / right.(float32)
		}
	case PERCENT_SIGN:
		if leftIsInt && rightIsInt {
			return left.(int) % right.(int)
		}
	case EQUALS:
		return left == right
	case NOT_EQUALS:
		return left != right
	case GREATER_THAN:
		if leftIsInt && rightIsInt {
			return left.(int) > right.(int)
		} else if leftIsFloat32 && rightIsFloat32 {
			return left.(float32) > right.(float32)
		}
	case LESSER_THAN:
		if leftIsInt && rightIsInt {
			return left.(int) < right.(int)
		} else if leftIsFloat32 && rightIsFloat32 {
			return left.(float32) < right.(float32)
		}
	}
	return nil
}

type ValueNode struct {
	Val Token
}

func (v ValueNode) Generate(vars *map[string]any, funcs *map[string]any) any {
	switch v.Val.Type {
	case STRING:
		return v.Val.Lit
	case NUMBER:
		if strings.Contains(v.Val.Lit, ".") {
			res, _ := strconv.ParseFloat(v.Val.Lit, 32)
			return float32(res)
		}
		res, _ := strconv.Atoi(v.Val.Lit)
		return res
	case IDENT:
		if _, ok := (*vars)[v.Val.Lit]; !ok {
			return fmt.Errorf("unknown variable '%s'", v.Val.Lit)
		}
		return (*vars)[v.Val.Lit]
	}
	return nil
}

func Parse(tokens []Token) ([]Node, error) {
	if tokens[len(tokens)-1].Type == EOF {
		tokens = tokens[:len(tokens)-1]
	}

	var nodes []Node
	idx := 0
	for idx < len(tokens) {
		if tokens[idx].Istype(NEWLINE) {
			idx++
		} else if tokens[idx].Istype(IDENT) {
			ident := tokens[idx]
			idx++

			idx_2 := idx
			var ctokens []Token
			found := -1
			for idx_2 < len(tokens) {
				if tokens[idx_2].Type == NEWLINE {
					ctokens = tokens[:idx_2+1]
					found = idx_2
					break
				}
				idx_2++
			}
			if found == -1 {
				ctokens = tokens
			}

			if ctokens[idx].Istype(ASSIGN) {
				//fmt.Println("ctokens: ", ctokens[idx+1:])
				//fmt.Println(idx, idx+1, len(ctokens), idx+1 >= len(ctokens))
				if idx+1 < len(ctokens) {
					if idx+2 >= len(ctokens) {
						nodes = append(nodes, AssignmentNode{ident, ValueNode{ctokens[idx+1]}})
					} else {
						nodes = append(nodes, AssignmentNode{ident, GenerateExpressionNodeFromTokens(ctokens[idx+1:])})
					}
				} else {
					return []Node{}, NewGorError(ctokens[idx], fmt.Sprintf("expected expression, but found '%s' instead", string(ctokens[idx].Lit)))
				}
			} else {
				return []Node{}, NewGorError(ctokens[idx], fmt.Sprintf("expected assign glyph, but found '%s' instead", string(ctokens[idx].Lit)))
			}

			if found != -1 {
				idx = found
			} else {
				idx = len(tokens)
			}
		} else if tokens[idx].Istype(COLON) {
			if idx+1 < len(tokens) {
				if idx+2 < len(tokens) {
					if !tokens[idx+2].Istype(COLON) {
						return []Node{}, NewGorError(tokens[idx], fmt.Sprintf("expected colon(':'), but found '%s' instead", string(tokens[idx].Lit)))
					}
					nodes = append(nodes, LabelNode{Name: tokens[idx+1]})
					idx += 3
				} else {
					return []Node{}, NewGorError(tokens[idx], fmt.Sprintf("expected colon(':'), but found '%s' instead", string(tokens[idx].Lit)))
				}
			} else {
				return []Node{}, NewGorError(tokens[idx], fmt.Sprintf("expected identifier, but found '%s' instead", string(tokens[idx].Lit)))
			}
		} else if tokens[idx].Istype(KEYWORD) {
			switch tokens[idx].Lit {
			case "jumpto":
				if idx+1 < len(tokens) {
					if tokens[idx+1].Istype(IDENT) {
						nodes = append(nodes, JumptoNode{tokens[idx+1]})
						idx += 2
						continue
					}
				}
				return []Node{}, NewGorError(tokens[idx], fmt.Sprintf("expected identifier, but found '%s' instead", string(tokens[idx].Lit)))
			case "use":
				if idx+1 < len(tokens) {
					if tokens[idx+1].Istype(STRING) {
						nodes = append(nodes, ModuleImportNode{tokens[idx+1]})
						idx += 2
						continue
					}
				}
				return []Node{}, NewGorError(tokens[idx], fmt.Sprintf("expected string, but found '%s' instead", string(tokens[idx].Lit)))
			default:
				return []Node{}, NewGorError(tokens[idx], fmt.Sprintf("unknown keyword '%s'", tokens[idx].Lit))
			}
		} else {
			return []Node{}, NewGorError(tokens[idx], fmt.Sprintf("unexpected '%s'", tokens[idx].Lit))
		}
	}

	return nodes, nil
}
