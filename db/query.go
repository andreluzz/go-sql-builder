package db

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/andreluzz/go-sql-builder/builder"
)

func parseSelectStruct(table, alias string, obj interface{}, embedded bool) ([]string, map[string]string) {
	t := reflect.TypeOf(obj).Elem()

	fields := []string{}
	joins := make(map[string]string)

	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag
		if tag.Get("sql") != "" && tag.Get("table") == "" {
			columnName := fmt.Sprintf("%s.%s %s", table, tag.Get("sql"), tag.Get("json"))
			if embedded {
				columnName = fmt.Sprintf("%s.%s %s__%s", table, tag.Get("sql"), alias, tag.Get("json"))
			}
			fields = append(fields, columnName)
		} else if tag.Get("table") != "" && tag.Get("embedded") == "" {
			columnName := fmt.Sprintf("%s.%s %s", tag.Get("alias"), tag.Get("sql"), tag.Get("json"))
			fields = append(fields, columnName)
			joinTable := fmt.Sprintf("%s %s", tag.Get("table"), tag.Get("alias"))
			joins[joinTable] = tag.Get("on")
		} else if tag.Get("embedded") == "slice" {
			if tag.Get("relation_table") != "" {
				joinTable := fmt.Sprintf("%s %s", tag.Get("relation_table"), tag.Get("relation_alias"))
				joins[joinTable] = tag.Get("relation_on")
			}
			joinTable := fmt.Sprintf("%s %s", tag.Get("table"), tag.Get("alias"))
			joins[joinTable] = tag.Get("on")
			embeddedFields, embeddedJoins := parseSelectStruct(tag.Get("alias"), tag.Get("json"), reflect.ValueOf(obj).Elem().Field(i).Interface(), true)
			fields = append(fields, embeddedFields...)
			for k, v := range embeddedJoins {
				joins[k] = v
			}
		}
	}

	return fields, joins
}

//StructSelectQuery generates the select query based on the struct fields
func StructSelectQuery(table string, obj interface{}, conditions builder.Builder) (string, []interface{}) {
	fields, joins := parseSelectStruct(table, "", obj, false)

	statement := builder.Select(fields...).From(table)
	for t, on := range joins {
		statement.Join(t, on)
	}

	if conditions != nil {
		statement.Where(conditions)
	}

	query := builder.NewQuery()
	statement.Prepare(query)

	return query.String(), query.Value()
}

//StructInsertQuery generates the insert query based on the struct fields
func StructInsertQuery(table string, obj interface{}) (string, []interface{}) {
	v := reflect.ValueOf(obj).Elem()
	t := reflect.TypeOf(obj).Elem()

	fields := []string{}
	args := []interface{}{}
	pkField := "id"
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag
		if tag.Get("sql") != "" && tag.Get("pk") != "true" && tag.Get("table") == "" {
			fields = append(fields, tag.Get("sql"))
			args = append(args, v.Field(i).Interface())
		}
		if tag.Get("pk") == "true" {
			pkField = tag.Get("sql")
		}
	}

	statement := builder.Insert(table, fields...).Values(args...).Return(pkField)
	query := builder.NewQuery()
	statement.Prepare(query)

	return query.String(), query.Value()
}

//StructMultipleInsertQuery generates the insert query based on the array of structs
func StructMultipleInsertQuery(table string, obj interface{}) (string, []interface{}) {
	t := reflect.TypeOf(obj).Elem()
	fields := []string{}
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag
		if tag.Get("sql") != "" && tag.Get("pk") != "true" && tag.Get("table") == "" {
			fields = append(fields, tag.Get("sql"))
		}
	}

	statement := builder.Insert(table, fields...)

	switch reflect.TypeOf(obj).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(obj)
		for i := 0; i < s.Len(); i++ {
			valueStruct := s.Index(i)
			args := []interface{}{}
			for i := 0; i < valueStruct.Type().NumField(); i++ {
				tag := valueStruct.Type().Field(i).Tag
				if tag.Get("sql") != "" && tag.Get("pk") != "true" && tag.Get("table") == "" {
					args = append(args, valueStruct.Field(i).Interface())
				}
			}
			statement.Values(args...)
		}
	}

	query := builder.NewQuery()
	statement.Prepare(query)

	return query.String(), query.Value()
}

//StructUpdateQuery generates the update query based on the struct fields
func StructUpdateQuery(table string, obj interface{}, updatableFields string, conditions builder.Builder) (string, []interface{}, error) {
	v := reflect.ValueOf(obj).Elem()
	t := reflect.TypeOf(obj).Elem()

	fields := []string{}
	args := []interface{}{}

	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag
		if tag.Get("sql") != "" && tag.Get("pk") != "true" && tag.Get("embedded") == "" && strings.Contains(updatableFields, tag.Get("sql")) {
			fields = append(fields, tag.Get("sql"))
			args = append(args, v.Field(i).Interface())
		}
	}

	if conditions == nil {
		return "", nil, errors.New("invalida where conditions")
	}

	statement := builder.Update(table, fields...).Values(args...).Where(conditions)
	query := builder.NewQuery()
	statement.Prepare(query)

	return query.String(), query.Value(), nil
}

//StructDeleteQuery generates the delete query based on the struct fields
func StructDeleteQuery(table string, conditions builder.Builder) (string, []interface{}, error) {
	statement := builder.Delete(table).Where(conditions)
	query := builder.NewQuery()
	statement.Prepare(query)

	return query.String(), query.Value(), nil
}
