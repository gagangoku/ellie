// Code generated from Ellie.g4 by ANTLR 4.13.1. DO NOT EDIT.

package parser // Ellie

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/antlr4-go/antlr/v4"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = strconv.Itoa
var _ = sync.Once{}

type EllieParser struct {
	*antlr.BaseParser
}

var EllieParserStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	LiteralNames           []string
	SymbolicNames          []string
	RuleNames              []string
	PredictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func ellieParserInit() {
	staticData := &EllieParserStaticData
	staticData.LiteralNames = []string{
		"", "'='", "'('", "')'", "'AND'", "'OR'", "'+'", "'-'", "'/'", "'*'",
		"'^'", "'<'", "'<='", "'>'", "'>='", "'=='", "'!='", "'CONTAINS'", "'IN'",
		"','",
	}
	staticData.SymbolicNames = []string{
		"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "",
		"", "", "", "TERM_TABLE", "TERM_TABLE_COL", "TERM_DB", "TERM_DB_COL",
		"TERM_CONST", "TERM_STAT", "TRUE", "FALSE", "NUMBER", "ID", "TEXT",
		"SPACE",
	}
	staticData.RuleNames = []string{
		"prog", "lhs", "expression", "function", "logicalFn", "mathFn", "setFn",
		"arguments", "bool",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 31, 84, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 1, 0, 1, 0, 1, 0, 1,
		0, 1, 0, 1, 1, 1, 1, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1,
		2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 3, 2, 42, 8, 2, 1, 2, 1, 2,
		1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 5, 2, 56, 8,
		2, 10, 2, 12, 2, 59, 9, 2, 1, 3, 1, 3, 1, 3, 3, 3, 64, 8, 3, 1, 3, 1, 3,
		1, 4, 1, 4, 1, 5, 1, 5, 1, 6, 1, 6, 1, 7, 1, 7, 1, 7, 5, 7, 77, 8, 7, 10,
		7, 12, 7, 80, 9, 7, 1, 8, 1, 8, 1, 8, 0, 1, 4, 9, 0, 2, 4, 6, 8, 10, 12,
		14, 16, 0, 5, 1, 0, 20, 25, 1, 0, 4, 5, 1, 0, 6, 16, 1, 0, 17, 18, 1, 0,
		26, 27, 90, 0, 18, 1, 0, 0, 0, 2, 23, 1, 0, 0, 0, 4, 41, 1, 0, 0, 0, 6,
		60, 1, 0, 0, 0, 8, 67, 1, 0, 0, 0, 10, 69, 1, 0, 0, 0, 12, 71, 1, 0, 0,
		0, 14, 73, 1, 0, 0, 0, 16, 81, 1, 0, 0, 0, 18, 19, 3, 2, 1, 0, 19, 20,
		5, 1, 0, 0, 20, 21, 3, 4, 2, 0, 21, 22, 5, 0, 0, 1, 22, 1, 1, 0, 0, 0,
		23, 24, 7, 0, 0, 0, 24, 3, 1, 0, 0, 0, 25, 26, 6, 2, -1, 0, 26, 42, 3,
		6, 3, 0, 27, 42, 3, 16, 8, 0, 28, 42, 5, 30, 0, 0, 29, 42, 5, 28, 0, 0,
		30, 42, 5, 21, 0, 0, 31, 42, 5, 20, 0, 0, 32, 42, 5, 23, 0, 0, 33, 42,
		5, 22, 0, 0, 34, 42, 5, 24, 0, 0, 35, 42, 5, 25, 0, 0, 36, 42, 5, 29, 0,
		0, 37, 38, 5, 2, 0, 0, 38, 39, 3, 4, 2, 0, 39, 40, 5, 3, 0, 0, 40, 42,
		1, 0, 0, 0, 41, 25, 1, 0, 0, 0, 41, 27, 1, 0, 0, 0, 41, 28, 1, 0, 0, 0,
		41, 29, 1, 0, 0, 0, 41, 30, 1, 0, 0, 0, 41, 31, 1, 0, 0, 0, 41, 32, 1,
		0, 0, 0, 41, 33, 1, 0, 0, 0, 41, 34, 1, 0, 0, 0, 41, 35, 1, 0, 0, 0, 41,
		36, 1, 0, 0, 0, 41, 37, 1, 0, 0, 0, 42, 57, 1, 0, 0, 0, 43, 44, 10, 15,
		0, 0, 44, 45, 3, 8, 4, 0, 45, 46, 3, 4, 2, 16, 46, 56, 1, 0, 0, 0, 47,
		48, 10, 14, 0, 0, 48, 49, 3, 10, 5, 0, 49, 50, 3, 4, 2, 15, 50, 56, 1,
		0, 0, 0, 51, 52, 10, 13, 0, 0, 52, 53, 3, 12, 6, 0, 53, 54, 3, 4, 2, 14,
		54, 56, 1, 0, 0, 0, 55, 43, 1, 0, 0, 0, 55, 47, 1, 0, 0, 0, 55, 51, 1,
		0, 0, 0, 56, 59, 1, 0, 0, 0, 57, 55, 1, 0, 0, 0, 57, 58, 1, 0, 0, 0, 58,
		5, 1, 0, 0, 0, 59, 57, 1, 0, 0, 0, 60, 61, 5, 29, 0, 0, 61, 63, 5, 2, 0,
		0, 62, 64, 3, 14, 7, 0, 63, 62, 1, 0, 0, 0, 63, 64, 1, 0, 0, 0, 64, 65,
		1, 0, 0, 0, 65, 66, 5, 3, 0, 0, 66, 7, 1, 0, 0, 0, 67, 68, 7, 1, 0, 0,
		68, 9, 1, 0, 0, 0, 69, 70, 7, 2, 0, 0, 70, 11, 1, 0, 0, 0, 71, 72, 7, 3,
		0, 0, 72, 13, 1, 0, 0, 0, 73, 78, 3, 4, 2, 0, 74, 75, 5, 19, 0, 0, 75,
		77, 3, 4, 2, 0, 76, 74, 1, 0, 0, 0, 77, 80, 1, 0, 0, 0, 78, 76, 1, 0, 0,
		0, 78, 79, 1, 0, 0, 0, 79, 15, 1, 0, 0, 0, 80, 78, 1, 0, 0, 0, 81, 82,
		7, 4, 0, 0, 82, 17, 1, 0, 0, 0, 5, 41, 55, 57, 63, 78,
	}
	deserializer := antlr.NewATNDeserializer(nil)
	staticData.atn = deserializer.Deserialize(staticData.serializedATN)
	atn := staticData.atn
	staticData.decisionToDFA = make([]*antlr.DFA, len(atn.DecisionToState))
	decisionToDFA := staticData.decisionToDFA
	for index, state := range atn.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(state, index)
	}
}

// EllieParserInit initializes any static state used to implement EllieParser. By default the
// static state used to implement the parser is lazily initialized during the first call to
// NewEllieParser(). You can call this function if you wish to initialize the static state ahead
// of time.
func EllieParserInit() {
	staticData := &EllieParserStaticData
	staticData.once.Do(ellieParserInit)
}

// NewEllieParser produces a new parser instance for the optional input antlr.TokenStream.
func NewEllieParser(input antlr.TokenStream) *EllieParser {
	EllieParserInit()
	this := new(EllieParser)
	this.BaseParser = antlr.NewBaseParser(input)
	staticData := &EllieParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	this.RuleNames = staticData.RuleNames
	this.LiteralNames = staticData.LiteralNames
	this.SymbolicNames = staticData.SymbolicNames
	this.GrammarFileName = "Ellie.g4"

	return this
}

// EllieParser tokens.
const (
	EllieParserEOF            = antlr.TokenEOF
	EllieParserT__0           = 1
	EllieParserT__1           = 2
	EllieParserT__2           = 3
	EllieParserT__3           = 4
	EllieParserT__4           = 5
	EllieParserT__5           = 6
	EllieParserT__6           = 7
	EllieParserT__7           = 8
	EllieParserT__8           = 9
	EllieParserT__9           = 10
	EllieParserT__10          = 11
	EllieParserT__11          = 12
	EllieParserT__12          = 13
	EllieParserT__13          = 14
	EllieParserT__14          = 15
	EllieParserT__15          = 16
	EllieParserT__16          = 17
	EllieParserT__17          = 18
	EllieParserT__18          = 19
	EllieParserTERM_TABLE     = 20
	EllieParserTERM_TABLE_COL = 21
	EllieParserTERM_DB        = 22
	EllieParserTERM_DB_COL    = 23
	EllieParserTERM_CONST     = 24
	EllieParserTERM_STAT      = 25
	EllieParserTRUE           = 26
	EllieParserFALSE          = 27
	EllieParserNUMBER         = 28
	EllieParserID             = 29
	EllieParserTEXT           = 30
	EllieParserSPACE          = 31
)

// EllieParser rules.
const (
	EllieParserRULE_prog       = 0
	EllieParserRULE_lhs        = 1
	EllieParserRULE_expression = 2
	EllieParserRULE_function   = 3
	EllieParserRULE_logicalFn  = 4
	EllieParserRULE_mathFn     = 5
	EllieParserRULE_setFn      = 6
	EllieParserRULE_arguments  = 7
	EllieParserRULE_bool       = 8
)

// IProgContext is an interface to support dynamic dispatch.
type IProgContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Lhs() ILhsContext
	Expression() IExpressionContext
	EOF() antlr.TerminalNode

	// IsProgContext differentiates from other interfaces.
	IsProgContext()
}

type ProgContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyProgContext() *ProgContext {
	var p = new(ProgContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EllieParserRULE_prog
	return p
}

func InitEmptyProgContext(p *ProgContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EllieParserRULE_prog
}

func (*ProgContext) IsProgContext() {}

func NewProgContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ProgContext {
	var p = new(ProgContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EllieParserRULE_prog

	return p
}

func (s *ProgContext) GetParser() antlr.Parser { return s.parser }

func (s *ProgContext) Lhs() ILhsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILhsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILhsContext)
}

func (s *ProgContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *ProgContext) EOF() antlr.TerminalNode {
	return s.GetToken(EllieParserEOF, 0)
}

func (s *ProgContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ProgContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ProgContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EllieListener); ok {
		listenerT.EnterProg(s)
	}
}

func (s *ProgContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EllieListener); ok {
		listenerT.ExitProg(s)
	}
}

func (p *EllieParser) Prog() (localctx IProgContext) {
	localctx = NewProgContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, EllieParserRULE_prog)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(18)
		p.Lhs()
	}
	{
		p.SetState(19)
		p.Match(EllieParserT__0)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(20)
		p.expression(0)
	}
	{
		p.SetState(21)
		p.Match(EllieParserEOF)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ILhsContext is an interface to support dynamic dispatch.
type ILhsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	TERM_TABLE() antlr.TerminalNode
	TERM_TABLE_COL() antlr.TerminalNode
	TERM_DB() antlr.TerminalNode
	TERM_DB_COL() antlr.TerminalNode
	TERM_CONST() antlr.TerminalNode
	TERM_STAT() antlr.TerminalNode

	// IsLhsContext differentiates from other interfaces.
	IsLhsContext()
}

type LhsContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLhsContext() *LhsContext {
	var p = new(LhsContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EllieParserRULE_lhs
	return p
}

func InitEmptyLhsContext(p *LhsContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EllieParserRULE_lhs
}

func (*LhsContext) IsLhsContext() {}

func NewLhsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LhsContext {
	var p = new(LhsContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EllieParserRULE_lhs

	return p
}

func (s *LhsContext) GetParser() antlr.Parser { return s.parser }

func (s *LhsContext) TERM_TABLE() antlr.TerminalNode {
	return s.GetToken(EllieParserTERM_TABLE, 0)
}

func (s *LhsContext) TERM_TABLE_COL() antlr.TerminalNode {
	return s.GetToken(EllieParserTERM_TABLE_COL, 0)
}

func (s *LhsContext) TERM_DB() antlr.TerminalNode {
	return s.GetToken(EllieParserTERM_DB, 0)
}

func (s *LhsContext) TERM_DB_COL() antlr.TerminalNode {
	return s.GetToken(EllieParserTERM_DB_COL, 0)
}

func (s *LhsContext) TERM_CONST() antlr.TerminalNode {
	return s.GetToken(EllieParserTERM_CONST, 0)
}

func (s *LhsContext) TERM_STAT() antlr.TerminalNode {
	return s.GetToken(EllieParserTERM_STAT, 0)
}

func (s *LhsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LhsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LhsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EllieListener); ok {
		listenerT.EnterLhs(s)
	}
}

func (s *LhsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EllieListener); ok {
		listenerT.ExitLhs(s)
	}
}

func (p *EllieParser) Lhs() (localctx ILhsContext) {
	localctx = NewLhsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, EllieParserRULE_lhs)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(23)
		_la = p.GetTokenStream().LA(1)

		if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&66060288) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IExpressionContext is an interface to support dynamic dispatch.
type IExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Function() IFunctionContext
	Bool_() IBoolContext
	TEXT() antlr.TerminalNode
	NUMBER() antlr.TerminalNode
	TERM_TABLE_COL() antlr.TerminalNode
	TERM_TABLE() antlr.TerminalNode
	TERM_DB_COL() antlr.TerminalNode
	TERM_DB() antlr.TerminalNode
	TERM_CONST() antlr.TerminalNode
	TERM_STAT() antlr.TerminalNode
	ID() antlr.TerminalNode
	AllExpression() []IExpressionContext
	Expression(i int) IExpressionContext
	LogicalFn() ILogicalFnContext
	MathFn() IMathFnContext
	SetFn() ISetFnContext

	// IsExpressionContext differentiates from other interfaces.
	IsExpressionContext()
}

type ExpressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExpressionContext() *ExpressionContext {
	var p = new(ExpressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EllieParserRULE_expression
	return p
}

func InitEmptyExpressionContext(p *ExpressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EllieParserRULE_expression
}

func (*ExpressionContext) IsExpressionContext() {}

func NewExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExpressionContext {
	var p = new(ExpressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EllieParserRULE_expression

	return p
}

func (s *ExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *ExpressionContext) Function() IFunctionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFunctionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFunctionContext)
}

func (s *ExpressionContext) Bool_() IBoolContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBoolContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBoolContext)
}

func (s *ExpressionContext) TEXT() antlr.TerminalNode {
	return s.GetToken(EllieParserTEXT, 0)
}

func (s *ExpressionContext) NUMBER() antlr.TerminalNode {
	return s.GetToken(EllieParserNUMBER, 0)
}

func (s *ExpressionContext) TERM_TABLE_COL() antlr.TerminalNode {
	return s.GetToken(EllieParserTERM_TABLE_COL, 0)
}

func (s *ExpressionContext) TERM_TABLE() antlr.TerminalNode {
	return s.GetToken(EllieParserTERM_TABLE, 0)
}

func (s *ExpressionContext) TERM_DB_COL() antlr.TerminalNode {
	return s.GetToken(EllieParserTERM_DB_COL, 0)
}

func (s *ExpressionContext) TERM_DB() antlr.TerminalNode {
	return s.GetToken(EllieParserTERM_DB, 0)
}

func (s *ExpressionContext) TERM_CONST() antlr.TerminalNode {
	return s.GetToken(EllieParserTERM_CONST, 0)
}

func (s *ExpressionContext) TERM_STAT() antlr.TerminalNode {
	return s.GetToken(EllieParserTERM_STAT, 0)
}

func (s *ExpressionContext) ID() antlr.TerminalNode {
	return s.GetToken(EllieParserID, 0)
}

func (s *ExpressionContext) AllExpression() []IExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpressionContext); ok {
			len++
		}
	}

	tst := make([]IExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpressionContext); ok {
			tst[i] = t.(IExpressionContext)
			i++
		}
	}

	return tst
}

func (s *ExpressionContext) Expression(i int) IExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *ExpressionContext) LogicalFn() ILogicalFnContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILogicalFnContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILogicalFnContext)
}

func (s *ExpressionContext) MathFn() IMathFnContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IMathFnContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IMathFnContext)
}

func (s *ExpressionContext) SetFn() ISetFnContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISetFnContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISetFnContext)
}

func (s *ExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EllieListener); ok {
		listenerT.EnterExpression(s)
	}
}

func (s *ExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EllieListener); ok {
		listenerT.ExitExpression(s)
	}
}

func (p *EllieParser) Expression() (localctx IExpressionContext) {
	return p.expression(0)
}

func (p *EllieParser) expression(_p int) (localctx IExpressionContext) {
	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()

	_parentState := p.GetState()
	localctx = NewExpressionContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IExpressionContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 4
	p.EnterRecursionRule(localctx, 4, EllieParserRULE_expression, _p)
	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(41)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 0, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(26)
			p.Function()
		}

	case 2:
		{
			p.SetState(27)
			p.Bool_()
		}

	case 3:
		{
			p.SetState(28)
			p.Match(EllieParserTEXT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 4:
		{
			p.SetState(29)
			p.Match(EllieParserNUMBER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 5:
		{
			p.SetState(30)
			p.Match(EllieParserTERM_TABLE_COL)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 6:
		{
			p.SetState(31)
			p.Match(EllieParserTERM_TABLE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 7:
		{
			p.SetState(32)
			p.Match(EllieParserTERM_DB_COL)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 8:
		{
			p.SetState(33)
			p.Match(EllieParserTERM_DB)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 9:
		{
			p.SetState(34)
			p.Match(EllieParserTERM_CONST)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 10:
		{
			p.SetState(35)
			p.Match(EllieParserTERM_STAT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 11:
		{
			p.SetState(36)
			p.Match(EllieParserID)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 12:
		{
			p.SetState(37)
			p.Match(EllieParserT__1)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(38)
			p.expression(0)
		}
		{
			p.SetState(39)
			p.Match(EllieParserT__2)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(57)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 2, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			p.SetState(55)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}

			switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 1, p.GetParserRuleContext()) {
			case 1:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, EllieParserRULE_expression)
				p.SetState(43)

				if !(p.Precpred(p.GetParserRuleContext(), 15)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 15)", ""))
					goto errorExit
				}
				{
					p.SetState(44)
					p.LogicalFn()
				}
				{
					p.SetState(45)
					p.expression(16)
				}

			case 2:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, EllieParserRULE_expression)
				p.SetState(47)

				if !(p.Precpred(p.GetParserRuleContext(), 14)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 14)", ""))
					goto errorExit
				}
				{
					p.SetState(48)
					p.MathFn()
				}
				{
					p.SetState(49)
					p.expression(15)
				}

			case 3:
				localctx = NewExpressionContext(p, _parentctx, _parentState)
				p.PushNewRecursionContext(localctx, _startState, EllieParserRULE_expression)
				p.SetState(51)

				if !(p.Precpred(p.GetParserRuleContext(), 13)) {
					p.SetError(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 13)", ""))
					goto errorExit
				}
				{
					p.SetState(52)
					p.SetFn()
				}
				{
					p.SetState(53)
					p.expression(14)
				}

			case antlr.ATNInvalidAltNumber:
				goto errorExit
			}

		}
		p.SetState(59)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 2, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.UnrollRecursionContexts(_parentctx)
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFunctionContext is an interface to support dynamic dispatch.
type IFunctionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ID() antlr.TerminalNode
	Arguments() IArgumentsContext

	// IsFunctionContext differentiates from other interfaces.
	IsFunctionContext()
}

type FunctionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFunctionContext() *FunctionContext {
	var p = new(FunctionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EllieParserRULE_function
	return p
}

func InitEmptyFunctionContext(p *FunctionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EllieParserRULE_function
}

func (*FunctionContext) IsFunctionContext() {}

func NewFunctionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FunctionContext {
	var p = new(FunctionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EllieParserRULE_function

	return p
}

func (s *FunctionContext) GetParser() antlr.Parser { return s.parser }

func (s *FunctionContext) ID() antlr.TerminalNode {
	return s.GetToken(EllieParserID, 0)
}

func (s *FunctionContext) Arguments() IArgumentsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IArgumentsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IArgumentsContext)
}

func (s *FunctionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FunctionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FunctionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EllieListener); ok {
		listenerT.EnterFunction(s)
	}
}

func (s *FunctionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EllieListener); ok {
		listenerT.ExitFunction(s)
	}
}

func (p *EllieParser) Function() (localctx IFunctionContext) {
	localctx = NewFunctionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, EllieParserRULE_function)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(60)
		p.Match(EllieParserID)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(61)
		p.Match(EllieParserT__1)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(63)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&2146435076) != 0 {
		{
			p.SetState(62)
			p.Arguments()
		}

	}
	{
		p.SetState(65)
		p.Match(EllieParserT__2)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ILogicalFnContext is an interface to support dynamic dispatch.
type ILogicalFnContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsLogicalFnContext differentiates from other interfaces.
	IsLogicalFnContext()
}

type LogicalFnContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLogicalFnContext() *LogicalFnContext {
	var p = new(LogicalFnContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EllieParserRULE_logicalFn
	return p
}

func InitEmptyLogicalFnContext(p *LogicalFnContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EllieParserRULE_logicalFn
}

func (*LogicalFnContext) IsLogicalFnContext() {}

func NewLogicalFnContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LogicalFnContext {
	var p = new(LogicalFnContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EllieParserRULE_logicalFn

	return p
}

func (s *LogicalFnContext) GetParser() antlr.Parser { return s.parser }
func (s *LogicalFnContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LogicalFnContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LogicalFnContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EllieListener); ok {
		listenerT.EnterLogicalFn(s)
	}
}

func (s *LogicalFnContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EllieListener); ok {
		listenerT.ExitLogicalFn(s)
	}
}

func (p *EllieParser) LogicalFn() (localctx ILogicalFnContext) {
	localctx = NewLogicalFnContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, EllieParserRULE_logicalFn)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(67)
		_la = p.GetTokenStream().LA(1)

		if !(_la == EllieParserT__3 || _la == EllieParserT__4) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IMathFnContext is an interface to support dynamic dispatch.
type IMathFnContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsMathFnContext differentiates from other interfaces.
	IsMathFnContext()
}

type MathFnContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMathFnContext() *MathFnContext {
	var p = new(MathFnContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EllieParserRULE_mathFn
	return p
}

func InitEmptyMathFnContext(p *MathFnContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EllieParserRULE_mathFn
}

func (*MathFnContext) IsMathFnContext() {}

func NewMathFnContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MathFnContext {
	var p = new(MathFnContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EllieParserRULE_mathFn

	return p
}

func (s *MathFnContext) GetParser() antlr.Parser { return s.parser }
func (s *MathFnContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MathFnContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *MathFnContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EllieListener); ok {
		listenerT.EnterMathFn(s)
	}
}

func (s *MathFnContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EllieListener); ok {
		listenerT.ExitMathFn(s)
	}
}

func (p *EllieParser) MathFn() (localctx IMathFnContext) {
	localctx = NewMathFnContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, EllieParserRULE_mathFn)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(69)
		_la = p.GetTokenStream().LA(1)

		if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&131008) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ISetFnContext is an interface to support dynamic dispatch.
type ISetFnContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser
	// IsSetFnContext differentiates from other interfaces.
	IsSetFnContext()
}

type SetFnContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySetFnContext() *SetFnContext {
	var p = new(SetFnContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EllieParserRULE_setFn
	return p
}

func InitEmptySetFnContext(p *SetFnContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EllieParserRULE_setFn
}

func (*SetFnContext) IsSetFnContext() {}

func NewSetFnContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SetFnContext {
	var p = new(SetFnContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EllieParserRULE_setFn

	return p
}

func (s *SetFnContext) GetParser() antlr.Parser { return s.parser }
func (s *SetFnContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SetFnContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SetFnContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EllieListener); ok {
		listenerT.EnterSetFn(s)
	}
}

func (s *SetFnContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EllieListener); ok {
		listenerT.ExitSetFn(s)
	}
}

func (p *EllieParser) SetFn() (localctx ISetFnContext) {
	localctx = NewSetFnContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, EllieParserRULE_setFn)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(71)
		_la = p.GetTokenStream().LA(1)

		if !(_la == EllieParserT__16 || _la == EllieParserT__17) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IArgumentsContext is an interface to support dynamic dispatch.
type IArgumentsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllExpression() []IExpressionContext
	Expression(i int) IExpressionContext

	// IsArgumentsContext differentiates from other interfaces.
	IsArgumentsContext()
}

type ArgumentsContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyArgumentsContext() *ArgumentsContext {
	var p = new(ArgumentsContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EllieParserRULE_arguments
	return p
}

func InitEmptyArgumentsContext(p *ArgumentsContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EllieParserRULE_arguments
}

func (*ArgumentsContext) IsArgumentsContext() {}

func NewArgumentsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ArgumentsContext {
	var p = new(ArgumentsContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EllieParserRULE_arguments

	return p
}

func (s *ArgumentsContext) GetParser() antlr.Parser { return s.parser }

func (s *ArgumentsContext) AllExpression() []IExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpressionContext); ok {
			len++
		}
	}

	tst := make([]IExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpressionContext); ok {
			tst[i] = t.(IExpressionContext)
			i++
		}
	}

	return tst
}

func (s *ArgumentsContext) Expression(i int) IExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *ArgumentsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ArgumentsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ArgumentsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EllieListener); ok {
		listenerT.EnterArguments(s)
	}
}

func (s *ArgumentsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EllieListener); ok {
		listenerT.ExitArguments(s)
	}
}

func (p *EllieParser) Arguments() (localctx IArgumentsContext) {
	localctx = NewArgumentsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, EllieParserRULE_arguments)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(73)
		p.expression(0)
	}
	p.SetState(78)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == EllieParserT__18 {
		{
			p.SetState(74)
			p.Match(EllieParserT__18)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(75)
			p.expression(0)
		}

		p.SetState(80)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IBoolContext is an interface to support dynamic dispatch.
type IBoolContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	TRUE() antlr.TerminalNode
	FALSE() antlr.TerminalNode

	// IsBoolContext differentiates from other interfaces.
	IsBoolContext()
}

type BoolContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBoolContext() *BoolContext {
	var p = new(BoolContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EllieParserRULE_bool
	return p
}

func InitEmptyBoolContext(p *BoolContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = EllieParserRULE_bool
}

func (*BoolContext) IsBoolContext() {}

func NewBoolContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BoolContext {
	var p = new(BoolContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = EllieParserRULE_bool

	return p
}

func (s *BoolContext) GetParser() antlr.Parser { return s.parser }

func (s *BoolContext) TRUE() antlr.TerminalNode {
	return s.GetToken(EllieParserTRUE, 0)
}

func (s *BoolContext) FALSE() antlr.TerminalNode {
	return s.GetToken(EllieParserFALSE, 0)
}

func (s *BoolContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BoolContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *BoolContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EllieListener); ok {
		listenerT.EnterBool(s)
	}
}

func (s *BoolContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(EllieListener); ok {
		listenerT.ExitBool(s)
	}
}

func (p *EllieParser) Bool_() (localctx IBoolContext) {
	localctx = NewBoolContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, EllieParserRULE_bool)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(81)
		_la = p.GetTokenStream().LA(1)

		if !(_la == EllieParserTRUE || _la == EllieParserFALSE) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

func (p *EllieParser) Sempred(localctx antlr.RuleContext, ruleIndex, predIndex int) bool {
	switch ruleIndex {
	case 2:
		var t *ExpressionContext = nil
		if localctx != nil {
			t = localctx.(*ExpressionContext)
		}
		return p.Expression_Sempred(t, predIndex)

	default:
		panic("No predicate with index: " + fmt.Sprint(ruleIndex))
	}
}

func (p *EllieParser) Expression_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	switch predIndex {
	case 0:
		return p.Precpred(p.GetParserRuleContext(), 15)

	case 1:
		return p.Precpred(p.GetParserRuleContext(), 14)

	case 2:
		return p.Precpred(p.GetParserRuleContext(), 13)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}
