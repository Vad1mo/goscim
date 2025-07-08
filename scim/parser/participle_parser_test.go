package parser_test

import (
	"testing"

	"github.com/arturoeanton/goscim/scim/parser"
)

// Test cases for the new participle parser implementation
func TestParticleParser(t *testing.T) {
	tests := []struct {
		name           string
		filter         string
		expectedQuery  string
		expectedCount  string
	}{
		{
			name:   "Empty filter",
			filter: "",
			expectedQuery: "SELECT * FROM `User`",
			expectedCount: "SELECT count(*) as count FROM `User`",
		},
		{
			name:   "Simple equality",
			filter: `userName eq "bjensen"`,
			expectedQuery: "SELECT * FROM `User` WHERE `userName` = \"bjensen\"",
			expectedCount: "SELECT count(*) as count FROM `User` WHERE `userName` = \"bjensen\"",
		},
		{
			name:   "Present operator",
			filter: "title pr",
			expectedQuery: "SELECT * FROM `User` WHERE `title`  IS NOT NULL",
			expectedCount: "SELECT count(*) as count FROM `User` WHERE `title`  IS NOT NULL",
		},
		{
			name:   "Contains operator",
			filter: `name.familyName co "O'Malley"`,
			expectedQuery: "SELECT * FROM `User` WHERE `name`.`familyName` LIKE \"%O'Malley%\"",
			expectedCount: "SELECT count(*) as count FROM `User` WHERE `name`.`familyName` LIKE \"%O'Malley%\"",
		},
		{
			name:   "Starts with operator",
			filter: `userName sw "J"`,
			expectedQuery: "SELECT * FROM `User` WHERE `userName` LIKE \"J%\"",
			expectedCount: "SELECT count(*) as count FROM `User` WHERE `userName` LIKE \"J%\"",
		},
		{
			name:   "Ends with operator",
			filter: `email ew "@example.com"`,
			expectedQuery: "SELECT * FROM `User` WHERE `email` LIKE \"%@example.com\"",
			expectedCount: "SELECT count(*) as count FROM `User` WHERE `email` LIKE \"%@example.com\"",
		},
		{
			name:   "Greater than",
			filter: `meta.lastModified gt "2011-05-13T04:42:34Z"`,
			expectedQuery: "SELECT * FROM `User` WHERE `meta`.`lastModified` >= \"2011-05-13T04:42:34Z\"",  // ANTLR bug: GT mapped to >=
			expectedCount: "SELECT count(*) as count FROM `User` WHERE `meta`.`lastModified` >= \"2011-05-13T04:42:34Z\"",
		},
		{
			name:   "Greater than or equal",
			filter: `meta.lastModified ge "2011-05-13T04:42:34Z"`,
			expectedQuery: "SELECT * FROM `User` WHERE `meta`.`lastModified` > \"2011-05-13T04:42:34Z\"",   // ANTLR bug: GE mapped to >
			expectedCount: "SELECT count(*) as count FROM `User` WHERE `meta`.`lastModified` > \"2011-05-13T04:42:34Z\"",
		},
		{
			name:   "Less than",
			filter: `meta.lastModified lt "2011-05-13T04:42:34Z"`,
			expectedQuery: "SELECT * FROM `User` WHERE `meta`.`lastModified` <= \"2011-05-13T04:42:34Z\"",  // ANTLR bug: LT mapped to <=
			expectedCount: "SELECT count(*) as count FROM `User` WHERE `meta`.`lastModified` <= \"2011-05-13T04:42:34Z\"",
		},
		{
			name:   "Less than or equal",
			filter: `meta.lastModified le "2011-05-13T04:42:34Z"`,
			expectedQuery: "SELECT * FROM `User` WHERE `meta`.`lastModified` < \"2011-05-13T04:42:34Z\"",   // ANTLR bug: LE mapped to <
			expectedCount: "SELECT count(*) as count FROM `User` WHERE `meta`.`lastModified` < \"2011-05-13T04:42:34Z\"",
		},
		{
			name:   "Not equal",
			filter: `userType ne "Employee"`,
			expectedQuery: "SELECT * FROM `User` WHERE `userType` <> \"Employee\"",
			expectedCount: "SELECT count(*) as count FROM `User` WHERE `userType` <> \"Employee\"",
		},
		{
			name:   "Boolean true",
			filter: `active eq true`,
			expectedQuery: "SELECT * FROM `User` WHERE `active` = true",
			expectedCount: "SELECT count(*) as count FROM `User` WHERE `active` = true",
		},
		{
			name:   "Boolean false",
			filter: `active eq false`,
			expectedQuery: "SELECT * FROM `User` WHERE `active` = false",
			expectedCount: "SELECT count(*) as count FROM `User` WHERE `active` = false",
		},
		{
			name:   "Number value",
			filter: `age eq 25`,
			expectedQuery: "SELECT * FROM `User` WHERE `age` = 25",
			expectedCount: "SELECT count(*) as count FROM `User` WHERE `age` = 25",
		},
		{
			name:   "AND operation",
			filter: `title pr and userType eq "Employee"`,
			expectedQuery: "SELECT * FROM `User` WHERE `title`  IS NOT NULL and `userType` = \"Employee\"",
			expectedCount: "SELECT count(*) as count FROM `User` WHERE `title`  IS NOT NULL and `userType` = \"Employee\"",
		},
		{
			name:   "OR operation",
			filter: `title pr or userType eq "Intern"`,
			expectedQuery: "SELECT * FROM `User` WHERE `title`  IS NOT NULL or `userType` = \"Intern\"",
			expectedCount: "SELECT count(*) as count FROM `User` WHERE `title`  IS NOT NULL or `userType` = \"Intern\"",
		},
		{
			name:   "Complex expression with parentheses",
			filter: `userType eq "Employee" and (emails co "example.com" or emails co "example.org")`,
			expectedQuery: "SELECT * FROM `User` WHERE `userType` = \"Employee\" and (`emails` LIKE \"%example.com%\" or `emails` LIKE \"%example.org%\")",
			expectedCount: "SELECT count(*) as count FROM `User` WHERE `userType` = \"Employee\" and (`emails` LIKE \"%example.com%\" or `emails` LIKE \"%example.org%\")",
		},
		{
			name:   "URN schema attribute",
			filter: `urn:ietf:params:scim:schemas:extension:one:2.0:Element.boolean eq true`,
			expectedQuery: "SELECT * FROM `User` WHERE `urn:ietf:params:scim:schemas:extension:one:2.0:Element`.`boolean` = true",
			expectedCount: "SELECT count(*) as count FROM `User` WHERE `urn:ietf:params:scim:schemas:extension:one:2.0:Element`.`boolean` = true",
		},
		{
			name:   "Complex URN expression",
			filter: `urn:ietf:params:scim:schemas:extension:one:2.0:User.userType eq "Employee" and (emails sw "example.com" or a.a.emails sw "example.org")`,
			expectedQuery: "SELECT * FROM `User` WHERE `urn:ietf:params:scim:schemas:extension:one:2.0:User`.`userType` = \"Employee\" and (`emails` LIKE \"example.com%\" or `a`.`a`.`emails` LIKE \"example.org%\")",
			expectedCount: "SELECT count(*) as count FROM `User` WHERE `urn:ietf:params:scim:schemas:extension:one:2.0:User`.`userType` = \"Employee\" and (`emails` LIKE \"example.com%\" or `a`.`a`.`emails` LIKE \"example.org%\")",
		},
		{
			name:   "NOT operation",
			filter: `not (userType eq "Employee")`,
			expectedQuery: "SELECT * FROM `User` WHERE not (`userType` = \"Employee\")",
			expectedCount: "SELECT count(*) as count FROM `User` WHERE not (`userType` = \"Employee\")",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query, count := parser.ParticleFilterToN1QL("User", tt.filter)
			if query != tt.expectedQuery {
				t.Errorf("Query mismatch:\nGot:      %s\nExpected: %s", query, tt.expectedQuery)
			}
			if count != tt.expectedCount {
				t.Errorf("Count query mismatch:\nGot:      %s\nExpected: %s", count, tt.expectedCount)
			}
		})
	}
}

// TestParticleParserAgainstAntlr compares participle results with ANTLR results
func TestParticleParserAgainstAntlr(t *testing.T) {
	testFilters := []string{
		`userName eq "bjensen"`,
		`name.familyName co "O'Malley"`,
		`userName sw "J"`,
		`title pr`,
		`meta.lastModified gt "2011-05-13T04:42:34Z"`,
		`title pr and userType eq "Employee"`,
		`title pr or userType eq "Intern"`,
		`userType eq "Employee" and (emails co "example.com" or emails co "example.org")`,
		`urn:ietf:params:scim:schemas:extension:one:2.0:Element.boolean eq true`,
		`active eq true`,
		`active eq false`,
		`age eq 25`,
	}

	for _, filter := range testFilters {
		t.Run(filter, func(t *testing.T) {
			// Get ANTLR results
			antlrQuery, antlrCount := parser.FilterToN1QL("User", filter)
			
			// Get Participle results
			particleQuery, particleCount := parser.ParticleFilterToN1QL("User", filter)
			
			// Compare results
			if antlrQuery != particleQuery {
				t.Errorf("Query mismatch for filter '%s':\nANTLR:     %s\nParticiple: %s", filter, antlrQuery, particleQuery)
			}
			if antlrCount != particleCount {
				t.Errorf("Count query mismatch for filter '%s':\nANTLR:     %s\nParticiple: %s", filter, antlrCount, particleCount)
			}
		})
	}
}

// TestParticleParserEdgeCases tests edge cases and error conditions
func TestParticleParserEdgeCases(t *testing.T) {
	tests := []struct {
		name   string
		filter string
		shouldNotPanic bool
	}{
		{
			name:   "Empty string",
			filter: "",
			shouldNotPanic: true,
		},
		{
			name:   "Only whitespace",
			filter: "   ",
			shouldNotPanic: true,
		},
		{
			name:   "Invalid syntax",
			filter: "userName eq",
			shouldNotPanic: true,
		},
		{
			name:   "Unclosed quotes",
			filter: `userName eq "test`,
			shouldNotPanic: true,
		},
		{
			name:   "Unmatched parentheses",
			filter: "userName eq \"test\" and (title pr",
			shouldNotPanic: true,
		},
		{
			name:   "Invalid operator",
			filter: "userName invalid \"test\"",
			shouldNotPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil && tt.shouldNotPanic {
					t.Errorf("Function panicked when it shouldn't: %v", r)
				}
			}()
			
			query, count := parser.ParticleFilterToN1QL("User", tt.filter)
			
			// Should at least return valid fallback queries
			if query == "" || count == "" {
				t.Errorf("Empty query returned for filter '%s'", tt.filter)
			}
		})
	}
}