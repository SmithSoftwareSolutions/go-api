package filters

import "strings"

type StringFilter struct {
	Equals   *string
	Contains *string
	IsNot    *string
	IsNull   *bool
	Or       *[]*StringFilter
	And      *[]*StringFilter
}

func (f *StringFilter) SQL(columnKey string) (string, []any) {
	if f.Equals != nil {
		return columnKey + " = ?", []any{*f.Equals}
	} else if f.Contains != nil {
		return columnKey + " LIKE ?", []any{"%" + *f.Contains + "%"}
	} else if f.IsNot != nil {
		return columnKey + " != ?", []any{*f.IsNot}
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

func StrEquals(value string) *StringFilter {
	return &StringFilter{
		Equals: &value,
	}
}

func StrContains(value string) *StringFilter {
	return &StringFilter{
		Contains: &value,
	}
}

func StrIsNot(value string) *StringFilter {
	return &StringFilter{
		IsNot: &value,
	}
}

func StrIsNull(value bool) *StringFilter {
	return &StringFilter{
		IsNull: &value,
	}
}

func StrAnd(values []*StringFilter) *StringFilter {
	return &StringFilter{
		And: &values,
	}
}

func StrOr(values []*StringFilter) *StringFilter {
	return &StringFilter{
		Or: &values,
	}
}
