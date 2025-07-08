package parser

import (
	"regexp"
	"strings"
)

// FilterToN1QL converts a SCIM filter to N1QL query using the new participle parser
func FilterToN1QL(resourceName string, filter string) (string, string) {
	return ParticleFilterToN1QL(resourceName, filter)
}

// AddQuote handles attribute name quoting for URN schemas and nested attributes
func AddQuote(value string) string {
	re := regexp.MustCompile(`^(urn[:\w\.\_]*)(:-*)?(:[\w]*)(\.)(.*)$`)
	urn := ""
	if re.MatchString(value) {
		urn = "`" + re.ReplaceAllString(value, `${1}${2}${3}`) + "`."
	}
	path := re.ReplaceAllString(value, `${5}`)
	path = urn + "`" + strings.Join(strings.Split(path, "."), "`.`") + "`"
	return path
}