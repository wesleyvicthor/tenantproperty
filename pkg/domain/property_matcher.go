package domain

import "strings"

type PropertyMatcher struct {
	Tenant          Tenant      `json:"tenant"`
	MostInteresting *Property   `json:"most_interesting"`
	Matches         []*Property `json:"matches"`
	properties      []*Property
}

func NewPropertyMatcher(c Tenant, v []*Property) PropertyMatcher {
	return PropertyMatcher{Tenant: c, properties: v}
}

func (c *PropertyMatcher) Match() {
	for _, v := range c.properties {
		if !strings.EqualFold(c.Tenant.PropertyType, v.Type) {
			continue
		}

		required := make(map[string]bool)
		for _, j := range v.Criteria {
			j = strings.ToLower(j)
			required[j] = true
		}
		skillC := 0
		for _, k := range c.Tenant.Criteria {
			k = strings.ToLower(k)
			if required[k] {
				skillC++
			}
		}

		pctage := skillC * 100 / len(v.Criteria)
		if pctage >= matchPercentage {
			c.Matches = append(c.Matches, v)
		}
	}

	if len(c.Matches) == 0 {
		return
	}

	bestMatch := c.Matches[0]
	c.Matches = c.Matches[1:]
	bestKey := 0
	for k, m := range c.Matches {
		if m.AveragePrice < bestMatch.AveragePrice {
			bestMatch = m
			bestKey = k
		}
	}

	c.Matches = append(c.Matches[:bestKey], c.Matches[bestKey+1:]...)
	c.MostInteresting = bestMatch
}
