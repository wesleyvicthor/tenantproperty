package internal

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/wesleyvicthor/tenantproperty/pkg/domain"
)

// Return infos of a Property
func (p PropertyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Path)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		res, _ := json.Marshal(map[string]string{"message": "Property ID must be a number"})
		w.Write(res)

		return
	}

	property := p.properties.OfID(id)
	if property == nil {
		w.WriteHeader(http.StatusNotFound)
		res, _ := json.Marshal(map[string]string{"message": "Property Not Found"})
		w.Write(res)

		return
	}

	response, _ := json.Marshal(property)
	w.Write(response)
}

type PropertyHandler struct {
	properties *PropertyRepository
}

func NewPropertyHandler(p *PropertyRepository) *PropertyHandler {
	return &PropertyHandler{p}
}

// List available Properties
func (p PropertiesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	city := query.Get("city")
	country := query.Get("country")
	sort := query.Get("sort")

	all := p.properties.SearchByAndSort(country, city, sort)

	response, _ := json.Marshal(all)
	w.Write(response)
}

type PropertiesHandler struct {
	properties *PropertyRepository
}

func NewPropertiesHandler(p *PropertyRepository) *PropertiesHandler {
	return &PropertiesHandler{p}
}

// Match a Tenant interest to available Property
func (pm PropertyMatchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		res, _ := json.Marshal(map[string]string{"message": "Method not allowed"})
		w.Write(res)
		return
	}

	// hard coded as an example, this could come from form submission
	tenant := domain.NewTenant("Penthouse", "WG,Garten,Balkon,Neubau")
	matcher := domain.NewPropertyMatcher(tenant, pm.properties.All())
	matcher.Match()

	response, _ := json.Marshal(matcher)
	w.Write(response)
}

type PropertyMatchHandler struct {
	properties *PropertyRepository
}

func NewPropertyMatchHandler(p *PropertyRepository) *PropertyMatchHandler {
	return &PropertyMatchHandler{p}
}
