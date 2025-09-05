package filters

import (
	"strconv"
	"strings"
)

type IntFilter struct {
	Equals               *int
	LessThan             *int
	GreaterThan          *int
	LessThanOrEqualTo    *int
	GreaterThanOrEqualTo *int
	IsNot                *int
	IsNull               *bool
	Or                   *[]*IntFilter
	And                  *[]*IntFilter
}

func (f *IntFilter) SQL(columnKey string) (string, []any) {
	if f.Equals != nil {
		return columnKey + " = ?", []any{strconv.Itoa(*f.Equals)}
	} else if f.LessThan != nil {
		return columnKey + " < ?", []any{strconv.Itoa(*f.LessThan)}
	} else if f.LessThanOrEqualTo != nil {
		return columnKey + " <= ?", []any{strconv.Itoa(*f.LessThanOrEqualTo)}
	} else if f.GreaterThan != nil {
		return columnKey + " > ?", []any{strconv.Itoa(*f.GreaterThan)}
	} else if f.GreaterThanOrEqualTo != nil {
		return columnKey + " >= ?", []any{strconv.Itoa(*f.GreaterThanOrEqualTo)}
	} else if f.IsNot != nil {
		return columnKey + " != ?", []any{strconv.Itoa(*f.IsNot)}
	} else if f.IsNull != nil {
		if *f.IsNull {
			return columnKey + " IS NULL", []any{}
		} else {
			return columnKey + " IS NOT NULL", []any{}
		}
	} else if f.And != nil {
		individualSQLStrings := []string{}
		individualParameters := []any{}
		for _, filter := range *f.And {
			sqlStr, parameters := filter.SQL(columnKey)
			if sqlStr != "" {
				individualSQLStrings = append(individualSQLStrings, sqlStr)
			}

			individualParameters = append(individualParameters, parameters...)

		}

		return "(" + strings.Join(individualSQLStrings, " AND ") + ")", individualParameters
	} else if f.Or != nil {
		individualSQLStrings := []string{}
		individualParameters := []any{}
		for _, filter := range *f.Or {
			sqlStr, parameters := filter.SQL(columnKey)
			if sqlStr != "" {
				individualSQLStrings = append(individualSQLStrings, sqlStr)
			}

			individualParameters = append(individualParameters, parameters...)
		}

		return "(" + strings.Join(individualSQLStrings, " OR ") + ")", individualParameters
	}

	return "", []any{}
}

func IntEquals(value int) *IntFilter {
	return &IntFilter{
		Equals: &value,
	}
}

func IntLessThan(value int) *IntFilter {
	return &IntFilter{
		LessThan: &value,
	}
}

func IntLessThanOrEqualTo(value int) *IntFilter {
	return &IntFilter{
		LessThanOrEqualTo: &value,
	}
}

func IntGreaterThan(value int) *IntFilter {
	return &IntFilter{
		GreaterThan: &value,
	}
}

func IntGreaterThanOrEqualTo(value int) *IntFilter {
	return &IntFilter{
		GreaterThanOrEqualTo: &value,
	}
}

func IntIsNot(value int) *IntFilter {
	return &IntFilter{
		IsNot: &value,
	}
}

func IntIsNull(value bool) *IntFilter {
	return &IntFilter{
		IsNull: &value,
	}
}

func IntAnd(values []*IntFilter) *IntFilter {
	return &IntFilter{
		And: &values,
	}
}

func IntOr(values []*IntFilter) *IntFilter {
	return &IntFilter{
		Or: &values,
	}
}
