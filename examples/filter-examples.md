# SCIM Filter Examples - No ANTLR Knowledge Required

This document demonstrates that you can use complex SCIM filters in GoSCIM without any ANTLR knowledge.

## Basic Filter Examples

### Simple Equality
```http
GET /scim/v2/Users?filter=userName eq "admin"
```

### String Operations
```http
# Contains
GET /scim/v2/Users?filter=name.familyName co "Garcia"

# Starts with  
GET /scim/v2/Users?filter=userName sw "admin"

# Ends with
GET /scim/v2/Users?filter=emails ew "@company.com"
```

### Presence Check
```http
GET /scim/v2/Users?filter=phoneNumbers pr
```

### Numeric Comparisons
```http
GET /scim/v2/Users?filter=age gt 25
GET /scim/v2/Users?filter=age le 65
```

## Logical Operators

### AND Operations
```http
GET /scim/v2/Users?filter=active eq true and userType eq "Employee"
```

### OR Operations  
```http
GET /scim/v2/Users?filter=userName eq "admin" or userName eq "administrator"
```

### NOT Operations
```http
GET /scim/v2/Users?filter=not (active eq false)
```

## Complex Nested Expressions

### Parentheses Grouping
```http
GET /scim/v2/Users?filter=userType eq "Employee" and (emails co "company.com" or emails co "company.org")
```

### Multi-valued Attributes
```http
GET /scim/v2/Users?filter=emails[type eq "work" and value ew "@company.com"]
```

### Schema Extensions
```http
GET /scim/v2/Users?filter=urn:ietf:params:scim:schemas:extension:enterprise:2.0:User:department eq "Engineering"
```

## Real-World Query Examples

### Find Active Employees from Specific Domains
```http
GET /scim/v2/Users?filter=active eq true and userType eq "Employee" and (emails[type eq "work"].value ew "@company.com" or emails[type eq "work"].value ew "@subsidiary.com")
```

### Find Users Modified Since Date
```http
GET /scim/v2/Users?filter=meta.lastModified gt "2023-01-01T00:00:00Z"
```

### Find Users with Phone Numbers
```http
GET /scim/v2/Users?filter=phoneNumbers pr and phoneNumbers[type eq "mobile"]
```

### Complex Enterprise Query
```http
GET /scim/v2/Users?filter=active eq true and urn:ietf:params:scim:schemas:extension:enterprise:2.0:User:department eq "Engineering" and (title co "Senior" or title co "Lead")
```

## How These Work Internally (You Don't Need to Know This)

Behind the scenes, GoSCIM's ANTLR-generated parser converts these expressions to N1QL queries:

```
Input:  userName eq "admin" and active eq true
Output: SELECT * FROM `User` WHERE `userName` = "admin" AND `active` = true
```

**The key point**: You send SCIM syntax, get results back. No ANTLR knowledge required!

## Testing Filters

You can test any filter expression using the provided test files:

```bash
# Run parser tests to see examples
go test ./scim/parser -v

# Start the server and test via HTTP
go run main.go
curl "http://localhost:8080/scim/v2/Users?filter=userName%20eq%20%22admin%22"
```

## Conclusion

As these examples show, you can create sophisticated SCIM queries without knowing anything about ANTLR. The parser handles all the complexity for you automatically.