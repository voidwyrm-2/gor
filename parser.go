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

type IfStatementNode struct {
	Expr  AssignableValue
	Nodes []Node
}

type ElsifStatementNode struct {
	Expr  AssignableValue
	Nodes []Node
}

type ElseStatementNode struct {
	Nodes []Node
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

func CheckTokenType(tokens []Token, index int, _type tokType) bool {
	if index < len(tokens) {
		return tokens[index].Istype(_type)
	}
	return false
}

func CollectUntilToken(tokens []Token, _type tokType, oposing_type tokType) ([]Token, bool) {
	nest := 0
	var out []Token
	gotBroken := false
	for _, t := range tokens {
		if t.Istype(_type) {
			if nest > 0 {
				nest--
				continue
			}
			gotBroken = true
			break
		} else if t.Istype(oposing_type) && oposing_type != NULLTOKEN {
			nest++
		}
		out = append(out, t)
	}
	return out, gotBroken
}

func IndexTokens(tokens []Token, _type tokType) int {
	for i, t := range tokens {
		if t.Istype(_type) {
			return i
		}
	}
	return -1
}

func IndexTokensWithCascadeFailsafe(tokens []Token, types []tokType) int {
	i := 0
	index := -1
	for i < len(types) {
		index = IndexTokens(tokens, types[i])
		if index == -1 {
			i++
			continue
		}
		break
	}

	return index
}

func GenerateExpressionNodeFromTokens(tokens []Token) (AssignableValue, error) {
	tokens = RemoveNewlineTokens(tokens)
	if len(tokens) == 0 {
		panic("length of tokens is 0")
	} else if len(tokens) == 1 {
		return ValueNode{Val: tokens[0]}, nil
	}

	index := IndexTokensWithCascadeFailsafe(tokens, []tokType{AND, OR, EQUALS, NOT_EQUALS, GREATER_THAN, LESSER_THAN, FORWARD_SLASH, PERCENT_SIGN, ASTERISK, HYPHEN, PLUS})
	if index == -1 {
		return ExpressionNode{}, fmt.Errorf("invalid tokens for expression: %v", tokens)
	}

	op := tokens[index]
	left := tokens[:index]
	right := tokens[index+1:]

	genLeft, leftErr := GenerateExpressionNodeFromTokens(left)
	if leftErr != nil {
		return ExpressionNode{}, leftErr
	}
	genRight, rightErr := GenerateExpressionNodeFromTokens(right)
	if rightErr != nil {
		return ExpressionNode{}, rightErr
	}
	return ExpressionNode{Left: genLeft, Operand: op, Right: genRight}, nil
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
	if len(tokens) == 0 {
		return []Node{}, nil
	}
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

			if tokens[idx].Istype(ASSIGN) {
				idx++
				//fmt.Println(tokens)
				exprToks, _ := CollectUntilToken(tokens[idx:], SEMICOLON, NULLTOKEN)
				if len(exprToks) == 0 {
					return []Node{}, NewGorError(tokens[idx], fmt.Sprintf("expected expression, but found '%s' instead", string(tokens[idx].Lit)))
				}
				gen, err := GenerateExpressionNodeFromTokens(exprToks)
				if err != nil {
					return []Node{}, err
				}
				nodes = append(nodes, AssignmentNode{Ident: ident, Value: gen})

				if !CheckTokenType(tokens, idx+len(exprToks), SEMICOLON) {
					return []Node{}, NewGorError(tokens[idx], "expected ';'")
				}
				idx += len(exprToks) + 1
			} else {
				return []Node{}, NewGorError(tokens[idx], fmt.Sprintf("expected assign glyph, but found '%s' instead", string(tokens[idx].Lit)))
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
			case "if", "elsif", "else":
				orig := tokens[idx].Lit
				idx++
				ifExprToks, ok := CollectUntilToken(tokens[idx:], LBRACE, NULLTOKEN)
				if !ok {
					return []Node{}, NewGorError(tokens[idx-1], "expected '{'")
				}
				idx += len(ifExprToks) + 1

				ifBodyToks, ok := CollectUntilToken(tokens[idx:], RBRACE, LBRACE)
				if !ok {
					return []Node{}, NewGorError(tokens[idx-1], "expected '}'")
				}
				//fmt.Println("body", ifBodyToks)

				ifBodyNodes, ifParseErr := Parse(ifBodyToks)
				if ifParseErr != nil {
					return []Node{}, ifParseErr
				}

				if orig == "elsif" {
					gen, err := GenerateExpressionNodeFromTokens(ifExprToks)
					if err != nil {
						return []Node{}, err
					}
					nodes = append(nodes, ElsifStatementNode{Expr: gen, Nodes: ifBodyNodes})
				} else if orig == "else" {
					nodes = append(nodes, ElseStatementNode{Nodes: ifBodyNodes})
				} else {
					gen, err := GenerateExpressionNodeFromTokens(ifExprToks)
					if err != nil {
						return []Node{}, err
					}
					nodes = append(nodes, IfStatementNode{Expr: gen, Nodes: ifBodyNodes})
				}
				idx += len(ifBodyToks) + 1
			case "jumpto":
				if CheckTokenType(tokens, idx+1, IDENT) {
					if !CheckTokenType(tokens, idx+2, SEMICOLON) {
						return []Node{}, NewGorError(tokens[idx+1], "expected ';'")
					}
					nodes = append(nodes, JumptoNode{tokens[idx+1]})
					idx += 3
					continue
				}
				return []Node{}, NewGorError(tokens[idx], fmt.Sprintf("expected identifier, but found '%s' instead", string(tokens[idx].Lit)))
			case "use":
				if CheckTokenType(tokens, idx+1, STRING) {
					nodes = append(nodes, ModuleImportNode{tokens[idx+1]})
					idx += 2
					continue
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
