package strs

import (
	"github.com/gertd/go-pluralize"
)

// PluralizeClient represents a client to support a pluralization.
type PluralizeClient struct {
	client *pluralize.Client
}

// NewPluralizeClient creates a new client.
func NewPluralizeClient() *PluralizeClient {
	c := &PluralizeClient{
		client: pluralize.NewClient(),
	}
	c.AddPluralRule("(?i)uri$", "uris")
	c.AddSingularRule("(?i)uris$", "uri")
	c.AddUncountableRule("(?i)info$")
	return c
}

// ToPlural converts the given string to its plural name.
func (c *PluralizeClient) ToPlural(s string) string {
	return c.client.Plural(c.client.Singular(s))
}

// AddPluralRule adds a pluralization rule to the collection.
func (c *PluralizeClient) AddPluralRule(rule string, replacement string) {
	c.client.AddPluralRule(rule, replacement)
}

// AddSingularRule adds a singularization rule to the collection.
func (c *PluralizeClient) AddSingularRule(rule string, replacement string) {
	c.client.AddSingularRule(rule, replacement)
}

// AddUncountableRule adds an uncountable word rule.
func (c *PluralizeClient) AddUncountableRule(word string) {
	c.client.AddUncountableRule(word)
}

// AddIrregularRule adds an irregular word definition.
func (c *PluralizeClient) AddIrregularRule(single string, plural string) {
	c.client.AddIrregularRule(single, plural)
}
