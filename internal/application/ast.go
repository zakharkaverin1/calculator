package application

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type ASTNode struct {
	IsLeaf        bool
	Value         float64
	Operator      string
	Left, Right   *ASTNode
	TaskScheduled bool
}

func ParseAST(expression string) (*ASTNode, error) {
	expr := strings.ReplaceAll(expression, " ", "")
	if expr == "" {
		return nil, fmt.Errorf("пустое")
	}
	p := &parser{input: expr, pos: 0}
	node, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	if p.pos < len(p.input) {
		return nil, fmt.Errorf("неизвестный токен на позиции ", p.pos)
	}
	return node, nil
}

type parser struct {
	input string
	pos   int
}

func (p *parser) peek() rune {
	if p.pos < len(p.input) {
		return rune(p.input[p.pos])
	}
	return 0
}

func (p *parser) get() rune {
	ch := p.peek()
	p.pos++
	return ch
}

func (p *parser) parseExpression() (*ASTNode, error) {
	node, err := p.parseTerm()
	if err != nil {
		return nil, err
	}
	for {
		ch := p.peek()
		if ch == '+' || ch == '-' {
			op := string(p.get())
			right, err := p.parseTerm()
			if err != nil {
				return nil, err
			}
			node = &ASTNode{
				IsLeaf:   false,
				Operator: op,
				Left:     node,
				Right:    right,
			}
		} else {
			break
		}
	}
	return node, nil
}

func (p *parser) parseTerm() (*ASTNode, error) {
	node, err := p.parseFactor()
	if err != nil {
		return nil, err
	}
	for {
		ch := p.peek()
		if ch == '*' || ch == '/' {
			op := string(p.get())
			right, err := p.parseFactor()
			if err != nil {
				return nil, err
			}
			node = &ASTNode{
				IsLeaf:   false,
				Operator: op,
				Left:     node,
				Right:    right,
			}
		} else {
			break
		}
	}
	return node, nil
}

func (p *parser) parseFactor() (*ASTNode, error) {
	ch := p.peek()
	if ch == '(' {
		p.get()
		node, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if p.peek() != ')' {
			return nil, fmt.Errorf("blablabla")
		}
		p.get()
		return node, nil
	}
	start := p.pos
	if ch == '+' || ch == '-' {
		p.get()
	}
	for {
		ch = p.peek()
		if unicode.IsDigit(ch) || ch == '.' {
			p.get()
		} else {
			break
		}
	}
	token := p.input[start:p.pos]
	if token == "" {
		return nil, fmt.Errorf("ожидалось иное число ", start)
	}
	value, err := strconv.ParseFloat(token, 64)
	if err != nil {
		return nil, fmt.Errorf("Невалидные данные ", token)
	}
	return &ASTNode{
		IsLeaf: true,
		Value:  value,
	}, nil
}