package engine

import (
	"encoding/json"
	"fmt"
	"gitlab.xibao100.com/skyline/skyline/cubes/metadata"
	"gitlab.xibao100.com/skyline/skyline/cubes/source"
	"gitlab.xibao100.com/skyline/skyline/cubes/utils"
	"io/ioutil"
	"sort"
	"sync"
)

func create_table_from_json(sqlite *source.Sqlite, table, jsonFile string, tags map[string][]*metadata.TagMapping) error {
	logger.Infof("creating table[%s] from json[%s]", table, jsonFile)
	defer logger.Infof("created table[%s] from json[%s]", table, jsonFile)

	bytes, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		logger.Errorf("ERROR: failed to file :%v", err.Error())
		return err
	}

	rows := []map[string]interface{}{}
	if err := json.Unmarshal(bytes, &rows); err != nil {
		logger.Errorf("ERROR Unmarshal: %v", err.Error())
		return err
	}
	if len(rows) == 0 {
		logger.Errorf("ERROR: empty json file: %v", jsonFile)
		return err
	}
	fields := []string{}
	dmlFields := [][]string{}
	data := source.Rows{}
	var wg sync.WaitGroup
	for i, row := range rows {
		//logger.Infof("row = %v", row)
		if len(row) == 0 {
			continue
		}
		//logger.Infof("i = %d, len(record) = %d", i, len(row))
		if i == 0 {
			fields, dmlFields = generate_json_fields(row)
			for k, _ := range tags {
				fields = append(fields, k)
				dmlFields = append(dmlFields, []string{k, "TEXT"})
			}
			if err := sqlite.CreateTable(table, dmlFields, false); err != nil {
				return err
			}
		}
		//logger.Infof("len(fields) = %d, fields = %v", len(fields), fields)

		record := make(map[string]string)
		for k, v := range row {
			record[k] = fmt.Sprintf("%v", v)
		}
		record = tag_mapping(record, tags)
		if len(record) == 0 || (len(fields) > 0 && len(record) != len(fields)) {
			continue
		}

		data = append(data, record)
		if len(data) >= 2000 {
			wg.Add(1)
			go func(db *source.Sqlite, tableName string, rows source.Rows) {
				defer wg.Done()

				if err := db.InsertTable(tableName, rows); err != nil {
					logger.Error(err)
					return
				}
			}(sqlite, table, data)
			data = source.Rows{}
		}
	}
	if len(data) > 0 {
		wg.Add(1)
		go func(db *source.Sqlite, tableName string, rows source.Rows) {
			defer wg.Done()

			if err := db.InsertTable(tableName, rows); err != nil {
				logger.Error(err)
				return
			}
		}(sqlite, table, data)
	}
	wg.Wait()

	//查询数据
	totalResult, err := sqlite.Query(fmt.Sprintf("SELECT COUNT(1) AS total FROM %s", table), []string{"total"})
	if err != nil {
		logger.Error(err)
		return err
	}

	logger.Infof("totalResult = %v", utils.Json(totalResult))

	return nil
}

func generate_json_fields(row map[string]interface{}) (fields []string, dmlFields [][]string) {
	for k, _ := range row {
		fields = append(fields, k)
	}
	sort.Slice(fields, func(i, j int) bool { // less func
		return fields[i] < fields[j] // order by asc
	})
	for _, v := range fields {
		dmlFields = append(dmlFields, []string{v, "TEXT"})
	}
	return
}
