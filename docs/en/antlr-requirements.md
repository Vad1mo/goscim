# ANTLR Knowledge Requirements for GoSCIM

This document answers the question: **"How much ANTLR knowledge is needed to use this project?"**

## TL;DR - Quick Answer

- **Regular Users**: **NO ANTLR knowledge required** - just use the project as-is
- **Most Contributors**: **NO ANTLR knowledge required** - for most development tasks
- **Grammar Modifiers**: **Intermediate ANTLR knowledge** - only if changing filter syntax
- **Advanced Customization**: **Advanced ANTLR knowledge** - for complex parser modifications

## Detailed Breakdown by User Type

### 1. 🚀 Regular Users (Just Using GoSCIM)

**ANTLR Knowledge Required: NONE ❌**

If you're simply using GoSCIM as a SCIM server, you need **zero ANTLR knowledge**.

**What you need instead:**
- Go 1.16+ 
- Couchbase Server 6.0+
- Basic understanding of SCIM 2.0 protocol

**Why no ANTLR knowledge needed:**
- All parser files are **pre-generated** and included in the repository
- The SCIM filter parsing happens automatically behind the scenes
- You just send SCIM filter expressions and they work

**Example - You can use complex filters without knowing ANTLR:**
```http
GET /scim/v2/Users?filter=userName eq "admin" and (emails co "company.com" or active eq true)
```

### 2. 🛠️ Contributors (Adding Features)

**ANTLR Knowledge Required: NONE ❌** (for most tasks)

Most contribution scenarios don't require touching the parser:

**Tasks requiring NO ANTLR knowledge:**
- ✅ Adding new SCIM resource types (Users, Groups, etc.)
- ✅ Creating custom schema extensions  
- ✅ Adding new API endpoints
- ✅ Database integration improvements
- ✅ Authentication/authorization features
- ✅ Performance optimizations
- ✅ Bug fixes in non-parser code
- ✅ Documentation improvements
- ✅ Adding tests

**Example - Adding a new resource type:**
```bash
# No ANTLR needed - just JSON configuration
mkdir config/schemas/
echo '{"id": "urn:custom:Employee", ...}' > config/schemas/Employee.json
mkdir config/resourceType/
echo '{"schemas": [...], "endpoint": "/Employees"}' > config/resourceType/Employee.json
```

### 3. 🔧 Grammar Modifiers (Changing Filter Syntax)

**ANTLR Knowledge Required: INTERMEDIATE ⚠️**

Only needed if you want to modify how SCIM filters are parsed.

**When you need ANTLR knowledge:**
- Adding new SCIM filter operators
- Changing filter syntax rules
- Supporting custom query extensions
- Modifying the grammar in `ScimFilter.g4`

**ANTLR concepts you need to understand:**
- Grammar syntax and rules
- Tokens vs. parser rules
- Basic regular expressions
- How to regenerate parsers

**Example modification process:**
```bash
# 1. Modify the grammar file
vim ScimFilter.g4

# 2. Regenerate parser (requires ANTLR installation)
wget http://www.antlr.org/download/antlr-4.7-complete.jar
alias antlr='java -jar $PWD/antlr-4.7-complete.jar'
antlr -Dlanguage=Go -o scim/parser ScimFilter.g4

# 3. Test your changes
go test ./scim/parser -v
```

**Skills needed:**
- Understanding ANTLR grammar files (`.g4` format)
- Basic knowledge of lexical analysis and parsing
- Ability to write and debug grammar rules
- Understanding of how tokens are defined

### 4. 🎓 Advanced Parser Customization

**ANTLR Knowledge Required: ADVANCED 🔥**

For complex modifications to the parsing pipeline.

**Advanced scenarios:**
- Custom visitor/listener implementations
- Complex AST transformations
- Performance optimizations in parser
- Error handling customization
- Integration with other query languages

**Advanced ANTLR concepts needed:**
- Visitor vs. Listener patterns
- Parse tree manipulation
- ANTLR runtime API
- Memory management in parsers
- Error recovery strategies

**Example advanced customization:**
```go
// Custom visitor implementation
type CustomFilterVisitor struct {
    *BaseScimFilterVisitor
    // Custom state
}

func (v *CustomFilterVisitor) VisitATTR_OPER_CRITERIA(ctx *ATTR_OPER_CRITERIAContext) interface{} {
    // Custom parsing logic
    return v.VisitChildren(ctx)
}
```

## Current ANTLR Setup in GoSCIM

### ✅ What's Already Provided

1. **Pre-generated Parser Files:**
   ```
   scim/parser/
   ├── scimfilter_parser.go      # Generated parser
   ├── scimfilter_lexer.go       # Generated lexer  
   ├── scimfilter_listener.go    # Generated listener interface
   ├── scimfilter_base_listener.go # Base listener implementation
   └── scimfilter_listener_implement.go # Custom N1QL conversion logic
   ```

2. **Ready-to-use ANTLR JAR:** `antlr-4.7-complete.jar` (included in repo)

3. **Working Grammar:** `ScimFilter.g4` defines complete SCIM filter syntax

4. **Test Suite:** `scim/parser/parser_test.go` validates parser functionality

### 🔄 When Regeneration is Needed

You only need to regenerate parser files if you modify `ScimFilter.g4`:

```bash
# Check if grammar changed
git status ScimFilter.g4

# If changed, regenerate (requires Java + ANTLR)
antlr -Dlanguage=Go -o scim/parser ScimFilter.g4

# Test changes
go test ./scim/parser -v
```

## Supported SCIM Filter Features (No ANTLR Knowledge Needed)

The current parser supports all standard SCIM 2.0 filter expressions:

### Comparison Operators
```
userName eq "admin"           # equals
name.familyName ne "Smith"    # not equals  
userName co "admin"           # contains
userName sw "admin"           # starts with
userName ew "admin"           # ends with
age gt 25                     # greater than
age ge 25                     # greater than or equal
age lt 65                     # less than
age le 65                     # less than or equal
active pr                     # present (has value)
```

### Logical Operators
```
userName eq "admin" and active eq true
userName eq "admin" or userName eq "user"
not (userName eq "admin")
```

### Complex Expressions
```
userType eq "Employee" and (emails co "company.com" or emails co "company.org")
emails[type eq "work" and value ew "@company.com"]
```

### Attribute References
```
name.familyName               # nested attributes
emails[primary eq true].value # multi-valued attributes
urn:ietf:params:scim:schemas:extension:enterprise:2.0:User:department eq "Engineering"
```

## Getting Help

### If You're a Regular User
- 📖 Read the [getting started guide](getting-started.md)
- 💬 Ask questions in [GitHub Discussions](https://github.com/Vad1mo/goscim/discussions)
- 🐛 Report issues if filters don't work as expected

### If You Need to Modify the Grammar
- 📚 Study the [ANTLR documentation](https://github.com/antlr/antlr4/blob/master/doc/index.md)
- 🔍 Examine the existing `ScimFilter.g4` file
- 🧪 Look at the test cases in `scim/parser/parser_test.go`
- 💡 Start with small grammar changes and test frequently

### Learning Resources
- [ANTLR 4 Documentation](https://github.com/antlr/antlr4/blob/master/doc/index.md)
- [The Definitive ANTLR 4 Reference](https://pragprog.com/titles/tpantlr2/the-definitive-antlr-4-reference/) (book)
- [ANTLR 4 Tutorial](https://tomassetti.me/antlr-mega-tutorial/)

## Summary

**99% of users and contributors will never need to learn ANTLR.** The project is designed to work out-of-the-box with pre-generated parser files. ANTLR knowledge is only required for the specialized use case of modifying the SCIM filter grammar itself.

Choose your path:
- **Just using GoSCIM?** → Skip ANTLR entirely
- **Contributing features?** → Skip ANTLR entirely  
- **Modifying filter syntax?** → Learn intermediate ANTLR
- **Advanced parser hacking?** → Master ANTLR deeply