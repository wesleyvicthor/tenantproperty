package internal

import (
	"encoding/csv"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/wesleyvicthor/tenantproperty/pkg/domain"
)

type PropertyRepository struct {
	properties domain.Properties
}

func NewPropertyRepository() *PropertyRepository {
	f, e := os.Open("data/housess.csv")
	if e != nil {
		panic(e)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	// remove header
	reader.Read()
	// automatic count columns
	reader.FieldsPerRecord = 0
	var properties domain.Properties
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}

		propertyId, _ := strconv.Atoi(row[0])
		salary, _ := strconv.ParseFloat(row[5], 64)
		var criteria domain.Criteria
		criteria = strings.Split(row[7], ",")
		trimCriteria(criteria)
		properties = append(properties, &domain.Property{
			ID:           propertyId,
			Title:        row[1],
			Type:         row[2],
			Country:      row[3],
			City:         row[4],
			AveragePrice: salary,
			Zimmer:       row[6],
			Criteria:     criteria,
			Floor:        row[8],
			LivingSpace:  row[9],
		})
	}

	return &PropertyRepository{properties}
}

// OfID returns an existent property
func (r *PropertyRepository) OfID(id int) *domain.Property {
	for _, v := range r.properties {
		if id == v.ID {
			return v
		}
	}

	return nil
}

// All retrieve all available properties
func (r *PropertyRepository) All() domain.Properties {
	return r.SearchByAndSort("", "", "")
}

// SearchByAndSort return a list of properties and possible filtering
func (r *PropertyRepository) SearchByAndSort(country, city, sortBy string) domain.Properties {
	var result domain.Properties
	filter := func(h domain.Properties, s string, n func(*domain.Property) string) domain.Properties {
		var f domain.Properties
		for _, v := range h {
			if strings.EqualFold(n(v), s) {
				f = append(f, v)
			}
		}
		return f
	}
	if country != "" {
		result = filter(r.properties, country, func(v *domain.Property) string {
			return v.Country
		})
	}

	if result == nil {
		result = r.properties
	}

	if city != "" {
		result = filter(result, city, func(v *domain.Property) string {
			return v.City
		})
	}

	if sortBy == "price" {
		sort.Sort(ByAveragePrice{result})
	}

	if sortBy == "type" {
		sort.Sort(ByPropertyType{result})
	}

	if result == nil {
		return domain.Properties{}
	}

	return result
}

// ByAveragePrice allow sorting properties by salary
type ByAveragePrice struct{ domain.Properties }

// ByPropertyType allow sorting properties by seniority level
type ByPropertyType struct{ domain.Properties }

// Less sort by salary ascending
func (s ByAveragePrice) Less(i, j int) bool {
	return s.Properties[i].AveragePrice < s.Properties[j].AveragePrice
}

// Less sort by seniority level ascending
func (s ByPropertyType) Less(i, j int) bool {
	return s.Properties[i].Type < s.Properties[j].Type
}

func trimCriteria(criteria domain.Criteria) domain.Criteria {
	for i := 0; i < len(criteria); i++ {
		criteria[i] = strings.Trim(criteria[i], " ")
	}

	return criteria
}
