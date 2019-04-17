package db

import (
	"database/sql"
	"encoding/json"
	"reflect"
)

//StructScan write rows data to struct array or struct
func StructScan(rows *sql.Rows, obj interface{}) error {

	cols, _ := rows.Columns()

	results := []map[string]interface{}{}

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
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			mapJSON[colName] = *val
		}
		results = append(results, mapJSON)
	}
	rows.Close()

	if len(results) == 0 {
		return nil
	}

	var jsonMap []byte
	if reflect.TypeOf(obj).Elem().Kind() == reflect.Struct {
		jsonMap, _ = json.Marshal(results[0])
	} else {
		jsonMap, _ = json.Marshal(results)
	}

	return json.Unmarshal(jsonMap, obj)
}
