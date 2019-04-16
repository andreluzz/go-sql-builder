package db

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/andreluzz/go-sql-builder/builder"
)

func parseSelectStruct(table, alias string, obj interface{}, embedded bool) ([]string, map[string]string, string, interface{}) {
	t := reflect.TypeOf(obj).Elem()

	fields := []string{}
	joins := make(map[string]string)
	pkField := "id"
	var pkValue interface{}
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag
		if tag.Get("sql") != "" && tag.Get("table") == "" {
			columnName := fmt.Sprintf("%s.%s %s", table, tag.Get("sql"), tag.Get("json"))
			if embedded {
				columnName = fmt.Sprintf("%s.%s %s__%s", table, tag.Get("sql"), alias, tag.Get("json"))
			}
			fields = append(fields, columnName)
			if tag.Get("pk") == "true" && !embedded {
				pkField = fmt.Sprintf("%s.%s", table, tag.Get("sql"))
				pkValue = reflect.ValueOf(obj).Elem().Field(i).Interface()
			}
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
			embeddedFields, embeddedJoins, _, _ := parseSelectStruct(tag.Get("alias"), tag.Get("json"), reflect.ValueOf(obj).Elem().Field(i).Interface(), true)
			fields = append(fields, embeddedFields...)
			for k, v := range embeddedJoins {
				joins[k] = v
			}
		}
	}

	return fields, joins, pkField, pkValue
}

//StructSelectQuery generates the select query based on the struct fields
func StructSelectQuery(table string, obj interface{}) (string, []interface{}) {
	fields, joins, pkField, pkValue := parseSelectStruct(table, "", obj, false)

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
		tag := t.Field(i).Tag
		if tag.Get("pk") != "true" && tag.Get("embedded") == "" && strings.Contains(updatableFields, tag.Get("sql")) {
			fields = append(fields, tag.Get("sql"))
			args = append(args, v.Field(i).Interface())
		}
		if tag.Get("pk") == "true" {
			pkField = tag.Get("sql")
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
