package db

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/andreluzz/go-sql-builder/builder"
)

type join struct {
	table string
	on    string
}

func parseSelectStruct(table, alias string, obj interface{}, embedded bool) ([]string, []join, error) {
	element := reflect.TypeOf(obj)

	if embedded {
		element = obj.(reflect.Type)
	} else {
		if element.Kind() == reflect.Ptr {
			element = element.Elem()
		}
		if element.Kind() == reflect.Slice {
			element = element.Elem()
		}
		fmt.Println(element)
	}

	fields := []string{}
	joins := []join{}

	for i := 0; i < element.NumField(); i++ {
		tag := element.Field(i).Tag
		if tag.Get("sql") != "" && tag.Get("table") == "" {
			columnName := fmt.Sprintf("%s.%s as %s", table, tag.Get("sql"), tag.Get("json"))
			if embedded {
				columnName = fmt.Sprintf("%s.%s as %s__%s", alias, tag.Get("sql"), alias, tag.Get("json"))
			}
			fields = append(fields, columnName)
		} else if tag.Get("sql") != "" && tag.Get("table") != "" {
			columnName := fmt.Sprintf("%s.%s as %s", tag.Get("alias"), tag.Get("sql"), tag.Get("json"))
			if embedded {
				columnName = fmt.Sprintf("%s_%s.%s as %s__%s", alias, tag.Get("alias"), tag.Get("sql"), alias, tag.Get("json"))
			}
			fields = append(fields, columnName)
			joinTable := fmt.Sprintf("%s %s", tag.Get("table"), tag.Get("alias"))
			if embedded {
				joinTable = fmt.Sprintf("%s %s_%s", tag.Get("table"), alias, tag.Get("alias"))
				replacedOn := strings.Replace(tag.Get("on"), fmt.Sprintf("%s.", tag.Get("alias")), fmt.Sprintf("%s_%s.", alias, tag.Get("alias")), -1)
				replacedOn = strings.Replace(replacedOn, fmt.Sprintf("%s.", table), fmt.Sprintf("%s.", alias), -1)
				joins = append(joins, join{joinTable, replacedOn})
			} else {
				joins = append(joins, join{joinTable, tag.Get("on")})
			}
		} else if tag.Get("sql") == "" && tag.Get("table") != "" && !embedded {

			joinTable := fmt.Sprintf("%s %s", tag.Get("table"), tag.Get("alias"))
			joins = append(joins, join{joinTable, tag.Get("on")})
			structFields, structJoins, err := parseSelectStruct(tag.Get("table"), tag.Get("alias"), element.Field(i).Type.Elem(), true)
			if err != nil {
				return nil, nil, err
			}

			fields = append(fields, structFields...)
			joins = append(joins, structJoins...)
		}
	}

	return fields, joins, nil
}

//StructSelectQuery generates the select query based on the struct fields
//Object can be a poninter to an array or struct
func StructSelectQuery(table string, obj interface{}, conditions builder.Builder) (string, []interface{}, error) {
	fields, joins, err := parseSelectStruct(table, "", obj, false)
	if err != nil {
		return "", nil, err
	}

	statement := builder.Select(fields...).From(table)
	for _, j := range joins {
		statement.Join(j.table, j.on)
	}

	if conditions != nil {
		statement.Where(conditions)
	}

	query := builder.NewQuery()
	statement.Prepare(query)

	return query.String(), query.Value(), nil
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
		if tag.Get("sql") != "" && tag.Get("pk") != "true" && strings.Contains(updatableFields, tag.Get("sql")) {
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
