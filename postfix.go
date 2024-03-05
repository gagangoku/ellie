package main

import (
	"errors"
	"fmt"
	"math"
	"strings"

	deque "github.com/gammazero/deque"
	"github.com/samber/lo"
)

const OP_ADD = "+"
const OP_SUB = "-"
const OP_MUL = "*"
const OP_DIV = "/"
const OP_POW = "^"
const OP_LT = "<"
const OP_GT = ">"
const OP_EQ = "=="
const OP_NE = "!="
const OP_LE = "<="
const OP_GE = ">="
const OP_IN = "IN"
const OP_CONTAINS = "CONTAINS"
const OP_BRACKET_START = "("
const OP_BRACKET_END = ")"

var ALL_OPS []string = []string{
	OP_ADD, OP_SUB, OP_MUL, OP_DIV, OP_POW,
}
var ALL_OPS_AND_BRACKETS []string = []string{
	OP_ADD, OP_SUB, OP_MUL, OP_DIV, OP_POW, OP_BRACKET_START, OP_BRACKET_END,
}

func precedence(op string) int {
	switch op {
	case OP_ADD, OP_SUB:
		return 1
	case OP_MUL, OP_DIV:
		return 2
	case OP_POW:
		return 3
	default:
		return 0
	}
}

func ParsePostfix(infix string, postfix []string) ([]string, error) {
	// Step 1 - ... "axsx...." .... "yyy...." .... => .... constants.c1 .... constants.c2 ...
	// Step 2 - parse token by token

	tokens := strings.Split(infix, " ")
	// Preallocating is faster than var stack []string
	stack := make([]string, 0, len(tokens)/2)

	for _, token := range tokens {
		if lo.IndexOf(ALL_OPS, token) >= 0 {
			for len(stack) > 0 && precedence(token) <= precedence(stack[len(stack)-1]) {
				postfix = append(postfix, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, token)
		} else if token == OP_BRACKET_START {
			stack = append(stack, token)
		} else if token == OP_BRACKET_END {
			for len(stack) > 0 && stack[len(stack)-1] != "(" {
				postfix = append(postfix, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			if len(stack) == 0 {
				return postfix, errors.New("invalid expression")
			}
			stack = stack[:len(stack)-1]
		} else {
			postfix = append(postfix, token)
		}
	}

	for len(stack) > 0 {
		if stack[len(stack)-1] == "(" {
			return postfix, errors.New("invalid expression")
		}
		postfix = append(postfix, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}

	return postfix, nil
}

type GenericValueType float64

func PostFixExecParts(parts []string, vals ...GenericValueType) GenericValueType {
	stack := deque.New[GenericValueType](10, 0)
	idx := 0
	for _, p := range parts {
		switch p {
		case OP_SUB:
			b := stack.PopBack()
			a := stack.PopBack()
			stack.PushBack(a - b)
		case OP_ADD:
			b := stack.PopBack()
			a := stack.PopBack()
			stack.PushBack(a + b)
		case OP_MUL:
			b := stack.PopBack()
			a := stack.PopBack()
			stack.PushBack(a * b)
		case OP_DIV:
			b := stack.PopBack()
			a := stack.PopBack()
			stack.PushBack(a / b)
		case OP_POW:
			b := stack.PopBack()
			a := stack.PopBack()
			stack.PushBack(GenericValueType(math.Pow(float64(a), float64(b))))
		default:
			stack.PushBack(vals[idx])
			idx++
		}
	}
	if stack.Len() != 1 {
		fmt.Println("Error: incorrect expression: ", parts)
	}
	return stack.PopBack()
}

func PostFixExec(expr string, vals ...GenericValueType) GenericValueType {
	parts := strings.Split(expr, " ")
	return PostFixExecParts(parts, vals...)
}
