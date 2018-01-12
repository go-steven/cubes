package engine

import (
	"errors"
	"fmt"
	"gitlab.xibao100.com/skyline/skyline/cubes/metadata"
	"gitlab.xibao100.com/skyline/skyline/cubes/source"
	"gitlab.xibao100.com/skyline/skyline/cubes/utils"
	"strings"
)

func tag_mapping(row map[string]string, tags map[string][]*metadata.TagMapping) map[string]string {
	for tagName, mappings := range tags {
		tagVal := ""
		for _, m := range mappings {
			fieldVal, ok := row[m.Field]
			if !ok {
				continue
			}
			if m.IncludeRegexp.MatchString(fieldVal) && (m.ExcludeRegexp == nil || !m.ExcludeRegexp.MatchString(fieldVal)) {
				tagVal = m.TagVal
				break
			}
		}
		row[tagName] = tagVal
	}
	return row
}

func sqlite_tags_mapping(sqlite *source.Sqlite, table string, tags map[string][]*metadata.TagMapping) error {
	fields, err := sqlite.GetTableFields(table)
	if err != nil {
		logger.Error(err)
		return err
	}
	if len(fields) == 0 {
		return errors.New(fmt.Sprintf("Table[%s] not found.", table))
	}
	logger.Infof("fields : %v", fields)
	sql := fmt.Sprintf("SELECT %s FROM %s ORDER BY %s", strings.Join(fields, ", "), table, fields[0])
	rows, err := sqlite.Query(sql, fields)
	if err != nil {
		logger.Error(err)
		return err
	}

	for k, _ := range tags {
		fields = append(fields, k)
	}
	dmlFields := [][]string{}
	for _, v := range fields {
		dmlFields = append(dmlFields, []string{v, "TEXT"})
	}
	newTable := metadata.Sha1Name(fmt.Sprintf("%s_with_tags", table))
	if err := sqlite.CreateTable(newTable, dmlFields, true); err != nil {
		return err
	}
	data := source.Rows{}
	for _, row := range rows {
		record := tag_mapping(row, tags)
		data = append(data, record)
		if len(data) > 2000 {
			if err := sqlite.InsertTable(newTable, data); err != nil {
				logger.Error(err)
				return err
			}
			data = source.Rows{}
		}
	}
	if len(data) > 0 {
		if err := sqlite.InsertTable(newTable, data); err != nil {
			logger.Error(err)
			return err
		}
	}

	//查询数据
	totalResult, err := sqlite.Query(fmt.Sprintf("SELECT COUNT(1) AS total FROM %s", newTable), []string{"total"})
	if err != nil {
		logger.Error(err)
		return err
	}

	logger.Infof("totalResult = %v", utils.Json(totalResult))

	return nil
}
