package domain

import (
	"strings"
)

const matchPercentage = 50

type Tenant struct {
	PropertyType string   `json:"property_type"`
	Criteria     Criteria `json:"criteria"`
}

func NewTenant(propertyType string, criteria string) Tenant {
	return Tenant{propertyType, strings.Split(criteria, ",")}
}
