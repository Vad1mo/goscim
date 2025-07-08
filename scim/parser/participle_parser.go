package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer"
)

// Lexer definition
var scimLexer = lexer.Must(lexer.Regexp(`(\s+)` +
	`|(?P<string>"[^"]*")` +
	`|(?P<number>[+-]?\d*\.?\d+)` +
	`|(?P<boolean>true|false)` +
	`|(?P<and>(?i)and)` +
	`|(?P<or>(?i)or)` +
	`|(?P<not>(?i)not)` +
	`|(?P<present>(?i)pr)` +
	`|(?P<eq>(?i)eq)` +
	`|(?P<ne>(?i)ne)` +
	`|(?P<co>(?i)co)` +
	`|(?P<sw>(?i)sw)` +
	`|(?P<ew>(?i)ew)` +
	`|(?P<gt>(?i)gt)` +
	`|(?P<lt>(?i)lt)` +
	`|(?P<ge>(?i)ge)` +
	`|(?P<le>(?i)le)` +
	`|(?P<lparen>\()` +
	`|(?P<rparen>\))` +
	`|(?P<lbracket>\[)` +
	`|(?P<rbracket>\])` +
	`|(?P<attrname>[a-zA-Z_$][-_.a-zA-Z0-9:]*)`))

// Grammar structures

// Filter represents the complete filter
type Filter struct {
	Expression *OrExpression `@@`
}

// OrExpression handles OR operations
type OrExpression struct {
	Left  *AndExpression   `@@`
	Right []*AndExpression `( "or" @@ )*`
}

// AndExpression handles AND operations  
type AndExpression struct {
	Left  *NotExpression   `@@`
	Right []*NotExpression `( "and" @@ )*`
}

// NotExpression handles NOT operations
type NotExpression struct {
	Not  bool    `@"not"?`
	Term *Term   `@@`
}

// Term represents individual terms
type Term struct {
	Group     *OrExpression   `"lparen" @@ "rparen"`
	Attribute *AttributeExpr  `| @@`
	Bracket   *BracketExpr    `| @@`
}

// AttributeExpr handles attribute expressions
type AttributeExpr struct {
	Name     string `@"attrname"`
	Present  bool   `@"present"?`
	Operator string `@("eq" | "ne" | "co" | "sw" | "ew" | "gt" | "lt" | "ge" | "le")?`
	Value    *Value `@@?`
}

// BracketExpr handles attribute[expression] syntax
type BracketExpr struct {
	Name       string        `@"attrname"`
	Expression *OrExpression `"lbracket" @@ "rbracket"`
}

// Value represents values
type Value struct {
	String  *string  `@"string"`
	Number  *float64 `| @"number"`
	Boolean *bool    `| @"boolean"`
}

// Parser wrapper
type ParticleParser struct {
	parser *participle.Parser
}

// NewParticleParser creates a new particle parser
func NewParticleParser() (*ParticleParser, error) {
	parser, err := participle.Build(&Filter{},
		participle.Lexer(scimLexer),
	)
	if err != nil {
		return nil, err
	}
	return &ParticleParser{parser: parser}, nil
}

// Parse parses a SCIM filter string
func (p *ParticleParser) Parse(filter string) (*Filter, error) {
	if filter == "" {
		return &Filter{}, nil
	}
	f := &Filter{}
	err := p.parser.ParseString(filter, f)
	return f, err
}

// N1QL conversion methods

func (f *Filter) ToN1QL(resourceName string) (string, string) {
	if f.Expression == nil {
		return fmt.Sprintf("SELECT * FROM `%s`", resourceName), 
			   fmt.Sprintf("SELECT count(*) as count FROM `%s`", resourceName)
	}
	
	whereClause := f.Expression.ToN1QL()
	return fmt.Sprintf("SELECT * FROM `%s` WHERE %s", resourceName, whereClause),
		   fmt.Sprintf("SELECT count(*) as count FROM `%s` WHERE %s", resourceName, whereClause)
}

func (e *OrExpression) ToN1QL() string {
	if len(e.Right) == 0 {
		return e.Left.ToN1QL()
	}
	
	parts := []string{e.Left.ToN1QL()}
	for _, right := range e.Right {
		parts = append(parts, right.ToN1QL())
	}
	return strings.Join(parts, " or ")
}

func (e *AndExpression) ToN1QL() string {
	if len(e.Right) == 0 {
		return e.Left.ToN1QL()
	}
	
	parts := []string{e.Left.ToN1QL()}
	for _, right := range e.Right {
		parts = append(parts, right.ToN1QL())
	}
	return strings.Join(parts, " and ")
}

func (e *NotExpression) ToN1QL() string {
	result := e.Term.ToN1QL()
	if e.Not {
		return fmt.Sprintf("NOT (%s)", result)
	}
	return result
}

func (t *Term) ToN1QL() string {
	if t.Group != nil {
		return fmt.Sprintf("(%s)", t.Group.ToN1QL())
	}
	if t.Attribute != nil {
		return t.Attribute.ToN1QL()
	}
	if t.Bracket != nil {
		return t.Bracket.ToN1QL()
	}
	return ""
}

func (e *AttributeExpr) ToN1QL() string {
	quotedName := AddQuote(e.Name)
	
	if e.Present {
		return fmt.Sprintf("%s  IS NOT NULL", quotedName) // Note: double space to match ANTLR
	}
	
	if e.Value == nil {
		return ""
	}
	
	valueStr := e.Value.ToN1QL()
	
	switch strings.ToLower(e.Operator) {
	case "eq":
		return fmt.Sprintf("%s = %s", quotedName, valueStr)
	case "ne":
		return fmt.Sprintf("%s <> %s", quotedName, valueStr)
	case "co":
		unquoted := strings.Trim(valueStr, "\"")
		return fmt.Sprintf("%s LIKE \"%%%s%%\"", quotedName, unquoted)
	case "sw":
		unquoted := strings.Trim(valueStr, "\"")
		return fmt.Sprintf("%s LIKE \"%s%%\"", quotedName, unquoted)
	case "ew":
		unquoted := strings.Trim(valueStr, "\"")
		return fmt.Sprintf("%s LIKE \"%%%s\"", quotedName, unquoted)
	case "gt":
		return fmt.Sprintf("%s > %s", quotedName, valueStr)
	case "ge":
		return fmt.Sprintf("%s >= %s", quotedName, valueStr)
	case "lt":
		return fmt.Sprintf("%s < %s", quotedName, valueStr)
	case "le":
		return fmt.Sprintf("%s <= %s", quotedName, valueStr)
	default:
		return fmt.Sprintf("%s = %s", quotedName, valueStr)
	}
}

func (e *BracketExpr) ToN1QL() string {
	quotedName := AddQuote(e.Name)
	innerExpr := e.Expression.ToN1QL()
	return fmt.Sprintf("%s[%s]", quotedName, innerExpr)
}

func (v *Value) ToN1QL() string {
	if v.String != nil {
		return *v.String
	}
	if v.Number != nil {
		return strconv.FormatFloat(*v.Number, 'f', -1, 64)
	}
	if v.Boolean != nil {
		return strconv.FormatBool(*v.Boolean)
	}
	return ""
}

// ParticleFilterToN1QL is the new implementation using participle
func ParticleFilterToN1QL(resourceName string, filter string) (string, string) {
	parser, err := NewParticleParser()
	if err != nil {
		return fmt.Sprintf("SELECT * FROM `%s`", resourceName),
			   fmt.Sprintf("SELECT count(*) as count FROM `%s`", resourceName)
	}
	
	parsed, err := parser.Parse(filter)
	if err != nil {
		return fmt.Sprintf("SELECT * FROM `%s`", resourceName),
			   fmt.Sprintf("SELECT count(*) as count FROM `%s`", resourceName)
	}
	
	return parsed.ToN1QL(resourceName)
}