package service

import (
	"sort"
	"strings"

	"github.com/shopspring/decimal"
)

// Sorting of hall rows (batch B5). Up to three `field[:asc|desc]` rules;
// rows whose metric is not published (state ∉ {ok, stale}) always sort last
// regardless of direction, and group_id ascending breaks every tie.

const ProviderHallMaxSortRules = 3

var providerHallSortFields = map[string]bool{
	"display_order":    true,
	"rate":             true,
	"historical_price": true,
	"predicted_rate":   true,
	"ttft_fast95":      true,
	"cache_rate":       true,
	"success_rate":     true,
}

type ProviderHallSortRule struct {
	Field string
	Desc  bool
}

// ParseProviderHallSort parses the query form. An empty string yields the
// default order (display_order asc). Unknown fields, bad directions, more
// than three rules or repeated fields are rejected.
func ParseProviderHallSort(raw string) ([]ProviderHallSortRule, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return []ProviderHallSortRule{{Field: "display_order"}}, true
	}
	parts := strings.Split(raw, ",")
	if len(parts) > ProviderHallMaxSortRules {
		return nil, false
	}
	seen := map[string]bool{}
	rules := make([]ProviderHallSortRule, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		field, dir, hasDir := strings.Cut(part, ":")
		field = strings.TrimSpace(field)
		if !providerHallSortFields[field] || seen[field] {
			return nil, false
		}
		rule := ProviderHallSortRule{Field: field}
		if hasDir {
			switch strings.ToLower(strings.TrimSpace(dir)) {
			case "asc":
			case "desc":
				rule.Desc = true
			default:
				return nil, false
			}
		}
		seen[field] = true
		rules = append(rules, rule)
	}
	return rules, true
}

// providerHallSortKey is one row's sortable projection. A nil value means
// "not published" and sorts after every published value.
type providerHallSortKey struct {
	groupID         int64
	displayOrder    int
	rate            decimal.Decimal
	historicalPrice *decimal.Decimal
	predictedRate   *decimal.Decimal
	ttftFast95      *float64
	cacheRate       *float64
	successRate     *float64
}

func providerHallSortDecimal[T any](m ProviderHallMetric[T], parse func(T) *decimal.Decimal) *decimal.Decimal {
	if m.Value == nil || (m.State != ProviderHallMetricOK && m.State != ProviderHallMetricStale) {
		return nil
	}
	return parse(*m.Value)
}

func providerHallSortFloat(m ProviderHallMetric[float64]) *float64 {
	if m.Value == nil || (m.State != ProviderHallMetricOK && m.State != ProviderHallMetricStale) {
		return nil
	}
	v := *m.Value
	return &v
}

func providerHallParseDecimal(s string) *decimal.Decimal {
	d, err := decimal.NewFromString(s)
	if err != nil {
		return nil
	}
	return &d
}

// compare returns -1/0/1 for the given field honouring the missing-last rule.
func (k providerHallSortKey) compare(other providerHallSortKey, rule ProviderHallSortRule) int {
	cmpDecimalPtr := func(a, b *decimal.Decimal) int {
		switch {
		case a == nil && b == nil:
			return 0
		case a == nil:
			return 1
		case b == nil:
			return -1
		}
		c := a.Cmp(*b)
		if rule.Desc {
			c = -c
		}
		return c
	}
	cmpFloatPtr := func(a, b *float64) int {
		switch {
		case a == nil && b == nil:
			return 0
		case a == nil:
			return 1
		case b == nil:
			return -1
		case *a < *b:
			if rule.Desc {
				return 1
			}
			return -1
		case *a > *b:
			if rule.Desc {
				return -1
			}
			return 1
		}
		return 0
	}
	switch rule.Field {
	case "display_order":
		c := 0
		if k.displayOrder < other.displayOrder {
			c = -1
		} else if k.displayOrder > other.displayOrder {
			c = 1
		}
		if rule.Desc {
			c = -c
		}
		return c
	case "rate":
		c := k.rate.Cmp(other.rate)
		if rule.Desc {
			c = -c
		}
		return c
	case "historical_price":
		return cmpDecimalPtr(k.historicalPrice, other.historicalPrice)
	case "predicted_rate":
		return cmpDecimalPtr(k.predictedRate, other.predictedRate)
	case "ttft_fast95":
		return cmpFloatPtr(k.ttftFast95, other.ttftFast95)
	case "cache_rate":
		return cmpFloatPtr(k.cacheRate, other.cacheRate)
	case "success_rate":
		return cmpFloatPtr(k.successRate, other.successRate)
	}
	return 0
}

// providerHallSortRows orders rows in place by the rules with group_id asc as
// the final tiebreak. It is stable so equal keys keep their input order.
func providerHallSortRows(rows []*ProviderHallRowResult, rules []ProviderHallSortRule) {
	sort.SliceStable(rows, func(i, j int) bool {
		a, b := rows[i].sortKey(), rows[j].sortKey()
		for _, rule := range rules {
			if c := a.compare(b, rule); c != 0 {
				return c < 0
			}
		}
		return a.groupID < b.groupID
	})
}
