package domain

type Property struct {
	ID           int      `json:"id"`
	Title        string   `json:"title"`
	Type         string   `json:"type"`
	Country      string   `json:"country"`
	City         string   `json:"city"`
	AveragePrice float64  `json:"average_price"`
	Zimmer       string   `json:"zimmer"`
	Criteria     Criteria `json:"criteria"`
	Floor        string   `json:"floor"`
	LivingSpace  string   `json:"living_space"`
}

type Criteria []string

type Properties []*Property

type PropertyRepository interface {
	OfID(id int) *Property
	SearchByAndSort(country, city, sort string) Properties
	All() Properties
}

// Len used by sort interface
func (p Properties) Len() int { return len(p) }

// Swap used by sort interface
func (p Properties) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
