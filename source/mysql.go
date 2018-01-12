package source

import (
	"errors"
	"fmt"
	"github.com/bububa/mymysql/autorc"
	_ "github.com/bububa/mymysql/thrsafe"
	"regexp"
	"strings"
	"gitlab.xibao100.com/skyline/skyline/cubes/utils"
)

type Mysql struct {
	db        *autorc.Conn
	sqlForbid *regexp.Regexp
}

func NewMysql(db *autorc.Conn) *Mysql {
	return &Mysql{
		db:        db,
		sqlForbid: SqlForbidRegexp(),
	}
}
func (m *Mysql) Query(sql string, fields []string) (Rows, error) {
	logger.Infof("MYSQL run SQL: %s, fields:%v\n", sql, fields)
	if len(fields) == 0 {
		return nil, errors.New("ERROR: no return fields given.")
	}
	if m.sqlForbid != nil && SqlForbid(sql, m.sqlForbid) {
		return nil, errors.New("ERROR: has not-allowed keywords in SQL.")
	}
	rows, res, err := m.db.Query(sql)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	logger.Infof("fields: %s", utils.Json(fields))
	for k, v := range fields {
		if v == "" {
			continue
		}
		if v[0] == '`' && v[len(v) - 1] == '`' {
			v = v[1:len(v)-1]
		}
		fields[k] = v
	}
	logger.Infof("new fields: %s", utils.Json(fields))
	
	ret := Rows{}
	for _, row := range rows {
		//logger.Infof("row: %s", utils.Json(row))
		retRow := make(Row)
		for _, v := range fields {
			retRow[v] = row.Str(res.Map(v))
		}

		ret = append(ret, retRow)
	}
	logger.Infof("SQL result: %d", len(ret))
	return ret, nil
}

func (m *Mysql) GetTableFields(table string) ([]string, error) {
	var tableName, schema string
	vals := strings.Split(table, ".")
	if len(vals) == 1 {
		tableName = strings.TrimSpace(vals[0])
	}
	if len(vals) == 2 {
		schema = strings.TrimSpace(vals[0])
		tableName = strings.TrimSpace(vals[1])
	}

	tableInfo, err := m.Query(fmt.Sprintf(`SELECT column_name FROM information_schema.columns WHERE ('%s' = '' OR table_schema = '%s') and table_name = '%s' ORDER BY ordinal_position`, schema, schema, tableName), []string{
		"column_name",
	})
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	if len(tableInfo) == 0 {
		return nil, errors.New(fmt.Sprintf("Table[%s] not found.", table))
	}
	fields := []string{}
	for _, row := range tableInfo {
		fields = append(fields, strings.ToLower(row["column_name"]))
	}

	return fields, nil
}
