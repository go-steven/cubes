package engine

import (
	"encoding/csv"
	"fmt"
	"gitlab.xibao100.com/skyline/skyline/cubes/metadata"
	"gitlab.xibao100.com/skyline/skyline/cubes/source"
	"gitlab.xibao100.com/skyline/skyline/cubes/utils"
	"io"
	"os"
	"sync"
)

func create_table_from_csv(sqlite *source.Sqlite, table, csvFile string, ignoreHeader bool, tags map[string][]*metadata.TagMapping) error {
	logger.Infof("creating table[%s] from csv[%s]", table, csvFile)
	defer logger.Infof("created table[%s] from csv[%s]", table, csvFile)

	file, err := os.Open(csvFile)
	if err != nil {
		logger.Error(err)
		return err
	}
	defer file.Close()

	fields := [][]string{}
	data := source.Rows{}
	reader := csv.NewReader(file)
	cnt := 0
	var wg sync.WaitGroup
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			logger.Error(err)
			return err
		}
		//logger.Infof("record = %v", record)
		cnt++

		if len(record) == 0 {
			continue
		}
		//logger.Infof("cnt = %d, len(record) = %d", cnt, len(record))
		if cnt == 1 {
			fields = generate_csv_fields(len(record))
			for k, _ := range tags {
				fields = append(fields, []string{k, "TEXT"})
			}
			if err := sqlite.CreateTable(table, fields, false); err != nil {
				return err
			}
			if ignoreHeader {
				continue
			}
		}
		//logger.Infof("len(fields) = %d, fields = %v", len(fields), fields)

		row := make(map[string]string)
		for k, v := range record {
			row[fields[k][0]] = v
		}
		row = tag_mapping(row, tags)
		if len(row) == 0 || (len(fields) > 0 && len(row) != len(fields)) {
			continue
		}

		data = append(data, row)
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

func generate_csv_fields(fieldsLen int) [][]string {
	ret := [][]string{}
	for i := 0; i < fieldsLen; i++ {
		ret = append(ret, []string{fmt.Sprintf("f%d", i+1), "TEXT"})
	}
	return ret
}
