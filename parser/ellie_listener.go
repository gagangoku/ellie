// Code generated from Ellie.g4 by ANTLR 4.13.1. DO NOT EDIT.

package parser // Ellie

import "github.com/antlr4-go/antlr/v4"

// EllieListener is a complete listener for a parse tree produced by EllieParser.
type EllieListener interface {
	antlr.ParseTreeListener

	// EnterProg is called when entering the prog production.
	EnterProg(c *ProgContext)

	// EnterLhs is called when entering the lhs production.
	EnterLhs(c *LhsContext)

	// EnterExpression is called when entering the expression production.
	EnterExpression(c *ExpressionContext)

	// EnterFunction is called when entering the function production.
	EnterFunction(c *FunctionContext)

	// EnterLogicalFn is called when entering the logicalFn production.
	EnterLogicalFn(c *LogicalFnContext)

	// EnterMathFn is called when entering the mathFn production.
	EnterMathFn(c *MathFnContext)

	// EnterSetFn is called when entering the setFn production.
	EnterSetFn(c *SetFnContext)

	// EnterArguments is called when entering the arguments production.
	EnterArguments(c *ArgumentsContext)

	// EnterBool is called when entering the bool production.
	EnterBool(c *BoolContext)

	// ExitProg is called when exiting the prog production.
	ExitProg(c *ProgContext)

	// ExitLhs is called when exiting the lhs production.
	ExitLhs(c *LhsContext)

	// ExitExpression is called when exiting the expression production.
	ExitExpression(c *ExpressionContext)

	// ExitFunction is called when exiting the function production.
	ExitFunction(c *FunctionContext)

	// ExitLogicalFn is called when exiting the logicalFn production.
	ExitLogicalFn(c *LogicalFnContext)

	// ExitMathFn is called when exiting the mathFn production.
	ExitMathFn(c *MathFnContext)

	// ExitSetFn is called when exiting the setFn production.
	ExitSetFn(c *SetFnContext)

	// ExitArguments is called when exiting the arguments production.
	ExitArguments(c *ArgumentsContext)

	// ExitBool is called when exiting the bool production.
	ExitBool(c *BoolContext)
}
