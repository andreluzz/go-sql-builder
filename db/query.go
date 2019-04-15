package db

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/andreluzz/go-sql-builder/builder"
)

//StructSelectQuery generates the select query based on the struct fields
func StructSelectQuery(table string, obj interface{}) (string, []interface{}) {
	v := reflect.ValueOf(obj).Elem()
	t := reflect.TypeOf(obj).Elem()

	fields := []string{}
	joins := make(map[string]string)
	pkField := "id"
	var pkValue interface{}
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag
		if tag.Get("sql") != "" && tag.Get("table") == "" {
			columnName := fmt.Sprintf("%s.%s %s", table, tag.Get("sql"), tag.Get("json"))
			fields = append(fields, columnName)

			if tag.Get("pk") == "true" {
				pkField = fmt.Sprintf("%s.%s", table, tag.Get("sql"))
				pkValue = v.Field(i).Interface()
			}
		} else if tag.Get("table") != "" {
			columnName := fmt.Sprintf("%s.%s %s", tag.Get("alias"), tag.Get("sql"), tag.Get("json"))
			fields = append(fields, columnName)
			table := fmt.Sprintf("%s %s", tag.Get("table"), tag.Get("alias"))
			joins[table] = tag.Get("on")
		}
	}

	statement := builder.Select(fields...).From(table)
	for t, on := range joins {
		statement.Join(t, on)
	}

	if pkValue != nil {
		statement.Where(builder.Equal(pkField, pkValue))
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
		if t.Field(i).Tag.Get("pk") != "true" && t.Field(i).Tag.Get("table") == "" {
			fields = append(fields, t.Field(i).Tag.Get("sql"))
			args = append(args, v.Field(i).Interface())
		}
		if t.Field(i).Tag.Get("pk") == "true" {
			pkField = t.Field(i).Tag.Get("sql")
		}
	}

	statement := builder.Insert(table, fields...).Values(args...).Return(pkField)
	query := builder.NewQuery()
	statement.Prepare(query)

	return query.String(), query.Value()
}

//StructUpdateQuery generates the update query based on the struct fields
func StructUpdateQuery(table string, obj interface{}, updatableFields string) (string, []interface{}, error) {
	v := reflect.ValueOf(obj).Elem()
	t := reflect.TypeOf(obj).Elem()

	fields := []string{}
	args := []interface{}{}
	pkField := ""
	var pkValue interface{}
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).Tag.Get("pk") != "true" && strings.Contains(updatableFields, t.Field(i).Tag.Get("sql")) {
			fields = append(fields, t.Field(i).Tag.Get("sql"))
			args = append(args, v.Field(i).Interface())
		}
		if t.Field(i).Tag.Get("pk") == "true" {
			pkField = t.Field(i).Tag.Get("sql")
			pkValue = v.Field(i).Interface()
		}
	}

	if pkValue == nil || pkField == "" {
		return "", nil, errors.New("invalid update pk value")
	}

	statement := builder.Update(table, fields...).Values(args...).Where(builder.Equal(pkField, pkValue))
	query := builder.NewQuery()
	statement.Prepare(query)

	return query.String(), query.Value(), nil
}

//StructDeleteQuery generates the delete query based on the struct fields
func StructDeleteQuery(table string, obj interface{}) (string, []interface{}, error) {
	v := reflect.ValueOf(obj).Elem()
	t := reflect.TypeOf(obj).Elem()
	pkField := ""
	var pkValue interface{}
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).Tag.Get("pk") == "true" {
			pkField = t.Field(i).Tag.Get("sql")
			pkValue = v.Field(i).Interface()
		}
	}

	if pkValue == nil || pkField == "" {
		return "", nil, errors.New("invalid delete pk value")
	}

	statement := builder.Delete(table).Where(builder.Equal(pkField, pkValue))
	query := builder.NewQuery()
	statement.Prepare(query)

	return query.String(), query.Value(), nil
}
