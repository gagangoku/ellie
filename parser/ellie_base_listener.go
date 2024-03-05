// Code generated from Ellie.g4 by ANTLR 4.13.1. DO NOT EDIT.

package parser // Ellie

import "github.com/antlr4-go/antlr/v4"

// BaseEllieListener is a complete listener for a parse tree produced by EllieParser.
type BaseEllieListener struct{}

var _ EllieListener = &BaseEllieListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseEllieListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseEllieListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseEllieListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseEllieListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterProg is called when production prog is entered.
func (s *BaseEllieListener) EnterProg(ctx *ProgContext) {}

// ExitProg is called when production prog is exited.
func (s *BaseEllieListener) ExitProg(ctx *ProgContext) {}

// EnterLhs is called when production lhs is entered.
func (s *BaseEllieListener) EnterLhs(ctx *LhsContext) {}

// ExitLhs is called when production lhs is exited.
func (s *BaseEllieListener) ExitLhs(ctx *LhsContext) {}

// EnterExpression is called when production expression is entered.
func (s *BaseEllieListener) EnterExpression(ctx *ExpressionContext) {}

// ExitExpression is called when production expression is exited.
func (s *BaseEllieListener) ExitExpression(ctx *ExpressionContext) {}

// EnterFunction is called when production function is entered.
func (s *BaseEllieListener) EnterFunction(ctx *FunctionContext) {}

// ExitFunction is called when production function is exited.
func (s *BaseEllieListener) ExitFunction(ctx *FunctionContext) {}

// EnterLogicalFn is called when production logicalFn is entered.
func (s *BaseEllieListener) EnterLogicalFn(ctx *LogicalFnContext) {}

// ExitLogicalFn is called when production logicalFn is exited.
func (s *BaseEllieListener) ExitLogicalFn(ctx *LogicalFnContext) {}

// EnterMathFn is called when production mathFn is entered.
func (s *BaseEllieListener) EnterMathFn(ctx *MathFnContext) {}

// ExitMathFn is called when production mathFn is exited.
func (s *BaseEllieListener) ExitMathFn(ctx *MathFnContext) {}

// EnterSetFn is called when production setFn is entered.
func (s *BaseEllieListener) EnterSetFn(ctx *SetFnContext) {}

// ExitSetFn is called when production setFn is exited.
func (s *BaseEllieListener) ExitSetFn(ctx *SetFnContext) {}

// EnterArguments is called when production arguments is entered.
func (s *BaseEllieListener) EnterArguments(ctx *ArgumentsContext) {}

// ExitArguments is called when production arguments is exited.
func (s *BaseEllieListener) ExitArguments(ctx *ArgumentsContext) {}

// EnterBool is called when production bool is entered.
func (s *BaseEllieListener) EnterBool(ctx *BoolContext) {}

// ExitBool is called when production bool is exited.
func (s *BaseEllieListener) ExitBool(ctx *BoolContext) {}
