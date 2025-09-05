package filters

import (
	"strconv"
	"strings"
)

type BoolFilter struct {
	Equals *bool
	IsNot  *bool
	IsNull *bool
	Or     *[]*BoolFilter
	And    *[]*BoolFilter
}

func (f *BoolFilter) SQL(columnKey string) (string, []any) {

	if f.Equals != nil {
		return columnKey + " = ?", []any{strconv.FormatBool(*f.Equals)}
	} else if f.IsNot != nil {
		return columnKey + " != ?", []any{strconv.FormatBool(*f.IsNot)}
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

func BoolEquals(value bool) *BoolFilter {
	return &BoolFilter{
		Equals: &value,
	}
}

func BoolIsNot(value bool) *BoolFilter {
	return &BoolFilter{
		IsNot: &value,
	}
}

func BoolIsNull(value bool) *BoolFilter {
	return &BoolFilter{
		IsNull: &value,
	}
}

func BoolAnd(values []*BoolFilter) *BoolFilter {
	return &BoolFilter{
		And: &values,
	}
}

func BoolOr(values []*BoolFilter) *BoolFilter {
	return &BoolFilter{
		Or: &values,
	}
}
