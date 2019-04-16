package db

import (
	"database/sql"
	"encoding/json"
	"reflect"
	"strings"
)

//StructScan write rows data to struct
func StructScan(rows *sql.Rows, obj interface{}) error {

	cols, _ := rows.Columns()

	results := []map[string]interface{}{}
	rowNumber := 0

	embeddedFields := getModelEmbeddedFields(obj)

	for rows.Next() {

		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return err
		}

		mapJSON := make(map[string]interface{})
		appendRow := false
		for i, colName := range cols {
			if !strings.Contains(colName, "__") {
				val := columnPointers[i].(*interface{})
				mapJSON[colName] = *val

				if rowNumber > 0 {
					if results[rowNumber-1][colName] != *val {
						appendRow = true
					}
				} else {
					appendRow = true
				}
			}
		}
		if appendRow {
			for _, f := range embeddedFields {
				mapJSON[f] = getRowEmbeddedObject(f, mapJSON["id"], rows)
			}
			results = append(results, mapJSON)
			rowNumber++
		}
	}
	rows.Close()

	if len(results) == 0 {
		return nil
	}

	var jsonMap []byte
	if reflect.Indirect(reflect.ValueOf(obj)).Kind() == reflect.Struct {
		jsonMap, _ = json.Marshal(results[0])
	} else {
		jsonMap, _ = json.Marshal(results)
	}
	json.Unmarshal(jsonMap, obj)

	return nil
}

func getRowEmbeddedObject(name string, id interface{}, rows *sql.Rows) []map[string]interface{} {
	cols, _ := rows.Columns()
	results := []map[string]interface{}{}
	rowNumber := 0
	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		rows.Scan(columnPointers...)

		mapJSON := make(map[string]interface{})
		appendRow := false

		//TODO get row pk id without using index 0
		rowID := columnPointers[0].(*interface{})
		if id != rowID {
			for i, colName := range cols {
				if strings.Contains(colName, name+"__") {
					index := strings.Index(colName, "__")
					val := columnPointers[i].(*interface{})
					embeddedColName := colName[index+2:]
					mapJSON[embeddedColName] = *val
					if rowNumber > 0 {
						if results[rowNumber-1][embeddedColName] != *val {
							appendRow = true
						}
					} else {
						appendRow = true
					}
				}
			}
		}

		if appendRow {
			results = append(results, mapJSON)
			rowNumber++
		}
	}
	return results
}

func getModelEmbeddedFields(obj interface{}) []string {
	model := obj
	if reflect.TypeOf(obj).Kind() == reflect.Slice {
		model = reflect.TypeOf(obj).Elem()
	}
	results := []string{}
	t := reflect.TypeOf(model).Elem()
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).Tag.Get("embedded") != "" {
			results = append(results, t.Field(i).Tag.Get("json"))
		}
	}
	return results
}
