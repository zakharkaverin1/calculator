// ast.go
package application

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type ASTNode struct {
	IsLeaf   bool
	Value    float64
	Operator string
	Left     *ASTNode
	Right    *ASTNode
}

type parser struct {
	input string
	pos   int
}

func ParseAST(expr string) (*ASTNode, error) {
	expr = strings.ReplaceAll(expr, " ", "")
	if expr == "" {
		return nil, fmt.Errorf("пустое выражение")
	}

	p := &parser{input: expr}
	node, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	if p.pos < len(p.input) {
		return nil, fmt.Errorf("неизвестный символ на позиции %d", p.pos)
	}
	return node, nil
}

func (p *parser) parseExpression() (*ASTNode, error) {
	return p.parseBinaryOp(p.parseTerm, []string{"+", "-"})
}

func (p *parser) parseTerm() (*ASTNode, error) {
	return p.parseBinaryOp(p.parseFactor, []string{"*", "/"})
}

func (p *parser) parseBinaryOp(next func() (*ASTNode, error), ops []string) (*ASTNode, error) {
	node, err := next()
	if err != nil {
		return nil, err
	}

	for {
		op := p.peekString(ops)
		if op == "" {
			break
		}
		p.pos++

		right, err := next()
		if err != nil {
			return nil, err
		}
		node = &ASTNode{Operator: op, Left: node, Right: right}
	}
	return node, nil
}

func (p *parser) parseFactor() (*ASTNode, error) {
	if p.peek() == '(' {
		p.pos++
		node, err := p.parseExpression()
		if err != nil || p.peek() != ')' {
			return nil, fmt.Errorf("незакрытая скобка")
		}
		p.pos++
		return node, nil
	}

	start := p.pos
	if p.peek() == '+' || p.peek() == '-' {
		p.pos++
	}
	for p.pos < len(p.input) && (unicode.IsDigit(rune(p.input[p.pos])) || p.input[p.pos] == '.') {
		p.pos++
	}

	numStr := p.input[start:p.pos]
	if numStr == "" {
		return nil, fmt.Errorf("ожидается число")
	}

	value, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return nil, fmt.Errorf("невалидное число: %s", numStr)
	}
	return &ASTNode{IsLeaf: true, Value: value}, nil
}

func (p *parser) peek() byte {
	if p.pos >= len(p.input) {
		return 0
	}
	return p.input[p.pos]
}

func (p *parser) peekString(ops []string) string {
	if p.pos >= len(p.input) {
		return ""
	}
	for _, op := range ops {
		if string(p.input[p.pos]) == op {
			return op
		}
	}
	return ""
}
