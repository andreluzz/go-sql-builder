package db

import (
	"database/sql"
	"encoding/json"
	"reflect"
)

//StructScan write rows data to struct
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
			if reflect.ValueOf(mapJSON[colName]).Kind() == reflect.Slice {
				mapJSON[colName] = string(mapJSON[colName].([]uint8))
			}
		}
		results = append(results, mapJSON)
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
