package util

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

type NullString sql.NullString

func (x *NullString) MarshalJSON() ([]byte, error) {
	if !x.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(x.String)
}

type FilterSQLer interface {
	SQL(columnKey string) (string, []any)
}

var structColumnsMap = make(map[string][]string)

func GetColumnsFromModel(obj any) ([]string, error) {
	rValue := reflect.ValueOf(obj)

	if rValue.Kind() != reflect.Pointer {
		return nil, errors.New("struct is not pointer")
	}

	rElem := rValue.Elem()
	rType := rElem.Type()

	structTypeKey := rType.String()

	if structColumnsMap[structTypeKey] != nil {
		return structColumnsMap[structTypeKey], nil
	}

	numFields := rElem.NumField()
	columnNames := []string{}

	for i := 0; i < numFields; i++ {
		field := rElem.Field(i)
		if field.CanSet() {
			fieldT := rType.Field(i)
			if strings.Contains(fieldT.Tag.Get("orm"), "ignore") || fieldT.Type.Kind() == reflect.Struct || (fieldT.Type.Kind() == reflect.Pointer && (fieldT.Type.Elem().Kind() == reflect.Struct || fieldT.Type.Elem().Kind() == reflect.Slice || fieldT.Type.Elem().Kind() == reflect.Array)) {
				continue
			}
			// convert pascal case to lower camel case
			lowerCamelCase := []rune{}
			for j, r := range rType.Field(i).Name {
				if unicode.IsUpper(r) && j > 0 {
					lowerCamelCase = append(lowerCamelCase, r)
				} else {
					lowerCamelCase = append(lowerCamelCase, unicode.ToLower(r))
				}
			}

			columnNames = append(columnNames, string(lowerCamelCase))
		}
	}

	structColumnsMap[structTypeKey] = columnNames
	return columnNames, nil
}

func GetCreateSQL(data any) ([]string, []any, error) {
	columns := []string{}
	params := []any{}

	rValue := reflect.ValueOf(data)

	if rValue.Kind() == reflect.Pointer {
		return nil, nil, errors.New("object is a pointer")
	}

	rType := rValue.Type()

	numFields := rValue.NumField()

	for i := 0; i < numFields; i++ {
		field := rValue.Field(i)
		fieldT := rType.Field(i)

		if strings.Contains(fieldT.Tag.Get("orm"), "ignore") || fieldT.Type.Kind() == reflect.Struct || (fieldT.Type.Kind() == reflect.Pointer && (fieldT.Type.Elem().Kind() == reflect.Struct || fieldT.Type.Elem().Kind() == reflect.Slice || fieldT.Type.Elem().Kind() == reflect.Array)) {
			continue
		}
		// convert pascal case to lower camel case
		lowerCamelCase := []rune{}
		for j, r := range rType.Field(i).Name {
			if unicode.IsUpper(r) && j > 0 {
				lowerCamelCase = append(lowerCamelCase, r)
			} else {
				lowerCamelCase = append(lowerCamelCase, unicode.ToLower(r))
			}
		}
		columns = append(columns, string(lowerCamelCase))
		params = append(params, field.Interface())

	}

	return columns, params, nil
}

func GetUpdateSQL(data any) (string, []any, error) {
	sql := []string{}
	params := []any{}

	rValue := reflect.ValueOf(data)

	if rValue.Kind() == reflect.Pointer {
		return "", nil, errors.New("object is a pointer")
	}

	rType := rValue.Type()

	numFields := rValue.NumField()

	for i := 0; i < numFields; i++ {
		field := rValue.Field(i)
		fieldT := rType.Field(i)

		if strings.Contains(fieldT.Tag.Get("orm"), "ignore") || fieldT.Type.Kind() == reflect.Struct || (fieldT.Type.Kind() == reflect.Pointer && (fieldT.Type.Elem().Kind() == reflect.Struct || fieldT.Type.Elem().Kind() == reflect.Slice || fieldT.Type.Elem().Kind() == reflect.Array)) {
			continue
		}
		// convert pascal case to lower camel case
		lowerCamelCase := []rune{}
		for j, r := range rType.Field(i).Name {
			if unicode.IsUpper(r) && j > 0 {
				lowerCamelCase = append(lowerCamelCase, r)
			} else {
				lowerCamelCase = append(lowerCamelCase, unicode.ToLower(r))
			}
		}
		sql = append(sql, string(lowerCamelCase)+"=?")
		params = append(params, field.Interface())

	}

	return strings.Join(sql, ", "), params, nil
}

func GetWhereSQL(data any) (string, []any, error) {
	sql := []string{}
	params := []any{}

	rValue := reflect.ValueOf(data)

	if rValue.Kind() == reflect.Pointer {
		return "", nil, errors.New("object is a pointer")
	}

	rType := rValue.Type()

	numFields := rValue.NumField()

	for i := 0; i < numFields; i++ {
		field := rValue.Field(i)
		fieldT := rType.Field(i)

		if strings.Contains(fieldT.Tag.Get("orm"), "ignore") {
			continue
		}

		// convert pascal case to lower camel case
		lowerCamelCase := []rune{}
		for j, r := range rType.Field(i).Name {
			if unicode.IsUpper(r) && j > 0 {
				lowerCamelCase = append(lowerCamelCase, r)
			} else {
				lowerCamelCase = append(lowerCamelCase, unicode.ToLower(r))
			}
		}

		if fieldT.Type.Kind() == reflect.Pointer {

			if field.IsNil() {
				continue
			}

			if field.Type().Implements(reflect.TypeOf((*FilterSQLer)(nil)).Elem()) {
				filterSql, filterParams := field.Interface().(FilterSQLer).SQL(string(lowerCamelCase))

				sql = append(sql, filterSql)
				params = append(params, filterParams...)
			} else {
				if !(fieldT.Type.Elem().Kind() == reflect.Struct || fieldT.Type.Elem().Kind() == reflect.Slice || fieldT.Type.Elem().Kind() == reflect.Array) {
					sql = append(sql, string(lowerCamelCase)+"=?")
					params = append(params, field.Interface())
				}
			}
		} else if !(field.Kind() == reflect.Struct || field.Kind() == reflect.Slice || field.Kind() == reflect.Array) {
			sql = append(sql, string(lowerCamelCase)+"=?")
			params = append(params, field.Interface())
		}
	}

	return strings.Join(sql, ", "), params, nil
}

func ScanRow(db *sql.DB, dest any, query string, args ...any) error {
	rValue := reflect.ValueOf(dest)

	if rValue.Kind() != reflect.Pointer {
		return errors.New("destination is not a pointer")
	}

	rElem := rValue.Elem()
	rElemType := rElem.Type()
	if rElem.Kind() != reflect.Struct {
		return errors.New("destination is not a pointer to a struct")
	}

	row := db.QueryRow(query, args...)

	numFields := rElem.NumField()
	scanArgs := []any{}

	for i := 0; i < numFields; i++ {
		field := rElem.Field(i)
		if field.CanSet() {
			fieldT := rElemType.Field(i)
			if strings.Contains(fieldT.Tag.Get("orm"), "ignore") || fieldT.Type.Kind() == reflect.Struct || (fieldT.Type.Kind() == reflect.Pointer && (fieldT.Type.Elem().Kind() == reflect.Struct || fieldT.Type.Elem().Kind() == reflect.Slice || fieldT.Type.Elem().Kind() == reflect.Array)) {
				continue
			}

			scanArgs = append(scanArgs, field.Addr().Interface())
		}
	}

	err := row.Scan(scanArgs...)
	if err != nil {
		return err
	}

	return nil
}

func ScanRows(db *sql.DB, dest any, query string, args ...any) error {
	rValue := reflect.ValueOf(dest)

	if rValue.Kind() != reflect.Pointer {
		return errors.New("destination is not a pointer")
	}

	rElem := rValue.Elem()
	if rElem.Kind() != reflect.Slice {
		return errors.New("destination is not a pointer to a slice: ")
	}

	rElemType := rElem.Type().Elem()
	if rElemType.Kind() != reflect.Struct {
		return errors.New("destination is not a pointer to an slice of structs")
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		newElem := reflect.New(rElemType).Elem()

		numFields := newElem.NumField()
		scanArgs := []any{}

		for i := 0; i < numFields; i++ {
			field := newElem.Field(i)

			if field.CanSet() {
				fieldT := rElemType.Field(i)
				if strings.Contains(fieldT.Tag.Get("orm"), "ignore") || fieldT.Type.Kind() == reflect.Struct || (fieldT.Type.Kind() == reflect.Pointer && (fieldT.Type.Elem().Kind() == reflect.Struct || fieldT.Type.Elem().Kind() == reflect.Slice || fieldT.Type.Elem().Kind() == reflect.Array)) {
					continue
				}

				scanArgs = append(scanArgs, field.Addr().Interface())
			}
		}

		err := rows.Scan(scanArgs...)
		if err != nil {
			return err
		}

		rElem.Set(reflect.Append(rElem, newElem))
	}

	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}

func FormatDebugSQLStringWithParameters(sqlString string, parameters []any) string {
	adjustedSqlString := sqlString
	for _, parameter := range parameters {
		adjustedSqlString = strings.Replace(adjustedSqlString, "?", fmt.Sprintf("%v", parameter), 1)
	}
	return adjustedSqlString
}
