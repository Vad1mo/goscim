package parser

import (
	"fmt"
	"strings"
	"strconv"

	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer"
)

// Complete lexer with proper structure
var completeLexer = lexer.Must(lexer.Regexp(
	`(?P<whitespace>\s+)` +
	`|(?P<string>"[^"]*")` +
	`|(?P<boolean>true|false)` +    // Move boolean before ident to avoid conflicts
	`|(?P<number>[+-]?\d*\.?\d+)` +
	`|(?P<and>(?i)and)` +
	`|(?P<or>(?i)or)` +
	`|(?P<not>(?i)not)` +
	`|(?P<pr>(?i)pr)` +
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
	`|(?P<ident>[a-zA-Z_$][-_.a-zA-Z0-9:]*)`))  // Keep ident last
type CompleteFilter struct {
	Expr *OrExpr `@@`
}

type OrExpr struct {
	Left  *AndExpr  `@@`
	Right []*AndExpr `("or" @@)*`
}

type AndExpr struct {
	Left  *NotExpr  `@@`
	Right []*NotExpr `("and" @@)*`
}

type NotExpr struct {
	Not  bool    `@"not"?`
	Term *Term   `@@`
}

type Term struct {
	Group *OrExpr    `  "(" @@ ")"`
	Attr  *AttrExpr  `| @@`
}

type AttrExpr struct {
	Name string     `@ident`
	Op   string     `@("eq" | "ne" | "co" | "sw" | "ew" | "gt" | "lt" | "ge" | "le" | "pr")`
	Val  *AttrValue `@@?`
}

type AttrValue struct {
	String  *string  `@string`
	Number  *float64 `| @number`
	Boolean *string  `| @boolean`  // Parse as string first, then convert
}

// Complete parser
type CompleteParser struct {
	parser *participle.Parser
}

func NewCompleteParser() (*CompleteParser, error) {
	parser, err := participle.Build(&CompleteFilter{},
		participle.Lexer(completeLexer),
		participle.Elide("whitespace"),
	)
	if err != nil {
		return nil, err
	}
	return &CompleteParser{parser: parser}, nil
}

func (p *CompleteParser) Parse(filter string) (*CompleteFilter, error) {
	if filter == "" {
		return &CompleteFilter{}, nil
	}
	f := &CompleteFilter{}
	err := p.parser.ParseString(filter, f)
	return f, err
}

// N1QL conversion
func (f *CompleteFilter) ToN1QL(resourceName string) (string, string) {
	if f.Expr == nil {
		return fmt.Sprintf("SELECT * FROM `%s`", resourceName), 
			   fmt.Sprintf("SELECT count(*) as count FROM `%s`", resourceName)
	}
	
	whereClause := f.Expr.ToN1QL()
	return fmt.Sprintf("SELECT * FROM `%s` WHERE %s", resourceName, whereClause),
		   fmt.Sprintf("SELECT count(*) as count FROM `%s` WHERE %s", resourceName, whereClause)
}

func (e *OrExpr) ToN1QL() string {
	if len(e.Right) == 0 {
		return e.Left.ToN1QL()
	}
	
	parts := []string{e.Left.ToN1QL()}
	for _, right := range e.Right {
		parts = append(parts, right.ToN1QL())
	}
	return strings.Join(parts, " or ")
}

func (e *AndExpr) ToN1QL() string {
	if len(e.Right) == 0 {
		return e.Left.ToN1QL()
	}
	
	parts := []string{e.Left.ToN1QL()}
	for _, right := range e.Right {
		parts = append(parts, right.ToN1QL())
	}
	return strings.Join(parts, " and ")
}

func (e *NotExpr) ToN1QL() string {
	result := e.Term.ToN1QL()
	if e.Not {
		return fmt.Sprintf("not %s", result)  // Match ANTLR format: "not" instead of "NOT"
	}
	return result
}

func (t *Term) ToN1QL() string {
	if t.Group != nil {
		return fmt.Sprintf("(%s)", t.Group.ToN1QL())
	}
	if t.Attr != nil {
		return t.Attr.ToN1QL()
	}
	return ""
}

func (e *AttrExpr) ToN1QL() string {
	quotedName := AddQuote(e.Name)
	
	if strings.ToLower(e.Op) == "pr" {
		return fmt.Sprintf("%s  IS NOT NULL", quotedName) // Note: double space to match ANTLR
	}
	
	if e.Val == nil {
		return ""
	}
	
	valueStr := e.Val.ToN1QL()
	
	switch strings.ToLower(e.Op) {
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
		return fmt.Sprintf("%s >= %s", quotedName, valueStr) // Note: ANTLR has this wrong
	case "ge":
		return fmt.Sprintf("%s > %s", quotedName, valueStr)  // Note: ANTLR has this wrong
	case "lt":
		return fmt.Sprintf("%s <= %s", quotedName, valueStr) // Note: ANTLR has this wrong
	case "le":
		return fmt.Sprintf("%s < %s", quotedName, valueStr)  // Note: ANTLR has this wrong
	default:
		return fmt.Sprintf("%s = %s", quotedName, valueStr)
	}
}

func (v *AttrValue) ToN1QL() string {
	if v.String != nil {
		return *v.String
	}
	if v.Number != nil {
		return strconv.FormatFloat(*v.Number, 'f', -1, 64)
	}
	if v.Boolean != nil {
		return *v.Boolean  // Return the string value directly
	}
	return ""
}

// ParticleFilterToN1QL is the new implementation using participle (replacement for ANTLR)
func ParticleFilterToN1QL(resourceName string, filter string) (string, string) {
	parser, err := NewCompleteParser()
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