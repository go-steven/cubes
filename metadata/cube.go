package metadata

import (
	"bytes"
	"errors"
	"fmt"
	"gitlab.xibao100.com/skyline/skyline/cubes/source"
	"gitlab.xibao100.com/skyline/skyline/cubes/utils"
)

func Sha1Name(name string) string {
	return fmt.Sprintf("t_%s", utils.Sha1(name))
}

func NewCube() *Cube {
	return &Cube{}
}

func (c *Cube) Execute() error {
	c.RunMode = DEFAULT_MODE

	switch c.Source.Type {
	case SOURCE_MYSQL:
		if c.Mysql == nil {
			return errors.New("ERROR: No MYSQL connection.")
		}
		sqlview, err := c.toSQLView(nil)
		if err != nil {
			logger.Error(err)
			return err
		}
		logger.Infof("Run cube[%s] sql: %s", c.Name, sqlview.Sql)
		data, err := c.Mysql.Query(sqlview.Sql, sqlview.Fields)
		if err != nil {
			logger.Error(err)
			return err
		}
		if err := c.saveResultToSqlite(data); err != nil {
			logger.Error(err)
			return err
		}
	case SOURCE_SQLITE:
		if c.Sqlite == nil {
			return errors.New("ERROR: No SQLITE connection.")
		}
		sqlview, err := c.toSQLView(nil)
		if err != nil {
			logger.Error(err)
			return err
		}
		logger.Infof("Run cube[%s] sql: %s", c.Name, sqlview.Sql)
		data, err := c.Sqlite.Query(sqlview.Sql, sqlview.Fields)
		if err != nil {
			logger.Error(err)
			return err
		}
		if err := c.saveResultToSqlite(data); err != nil {
			logger.Error(err)
			return err
		}
	case SOURCE_CUBE:
		fallthrough
	case SOURCE_JSON:
		fallthrough
	case SOURCE_CSV:
		if c.Sqlite == nil {
			return errors.New("ERROR: No sqlite connection.")
		}
		sqlview, err := c.toSQLView(nil)
		if err != nil {
			logger.Error(err)
			return err
		}
		logger.Infof("Run cube[%s] sql: %s", c.Name, sqlview.Sql)
		data, err := c.ESqlite.Query(sqlview.Sql, sqlview.Fields)
		if err != nil {
			logger.Error(err)
			return err
		}
		if err := c.saveResultToSqlite(data); err != nil {
			logger.Error(err)
			return err
		}
	default:
		err := errors.New(fmt.Sprintf("Unknow source type:%s", c.Source.Type))
		logger.Error(err)
		return err
	}

	return nil
}

func (c *Cube) saveResultToSqlite(data source.Rows) error {
	dmlFields := [][]string{}
	returnFields, err := c.getReturnFields(nil)
	if err != nil {
		logger.Error(err)
		return err
	}
	for _, v := range returnFields {
		dmlFields = append(dmlFields, []string{v, "string"})
	}

	if err := c.ESqlite.CreateTable(c.Sha1Name, dmlFields, true); err != nil {
		logger.Error(err)
		return err
	}

	if err := c.ESqlite.InsertTable(c.Sha1Name, data); err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (c *Cube) GetReport() (*CubeReport, error) {
	fields, err := c.getReturnFields(nil)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	var sql bytes.Buffer
	sql.WriteString("SELECT ")
	cnt := 0
	for _, v := range fields {
		if cnt > 0 {
			sql.WriteString(fmt.Sprintf(", `%s`", v))
		} else {
			sql.WriteString(fmt.Sprintf("`%s`", v))
		}
		cnt++
	}
	sql.WriteString(fmt.Sprintf(" FROM %s", c.Sha1Name))
	rows, err := c.ESqlite.Query(sql.String(), fields)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	ret := &CubeReport{
		Name:    c.Name,
		Display: c.Display,
		Fields:  fields,
		Data:    rows,
	}
	if err := ret.GenerateSummary(c.Summary); err != nil {
		logger.Error(err)
		return nil, err
	}

	if err := ret.GenerateSummaryCalc(c.SummaryCalc); err != nil {
		logger.Error(err)
		return nil, err
	}

	// 处理不返回的字段列表
	if len(c.TmpFields) > 0 {
		newFields := []string{}
		for _, v := range ret.Fields {
			if _, ok := c.TmpFields[v]; !ok {
				newFields = append(newFields, v)
			}
		}
		ret.Fields = newFields

		newRows := source.Rows{}
		for _, row := range ret.Data {
			newRow := make(source.Row)
			for k, v := range row {
				if _, ok := c.TmpFields[k]; !ok {
					newRow[k] = v
				}
			}
			newRows = append(newRows, newRow)
		}
		ret.Data = newRows

		newSummary := make(map[string]source.Row)
		for k, row := range ret.Summary {
			newRow := make(source.Row)
			for k, v := range row {
				if _, ok := c.TmpFields[k]; !ok {
					newRow[k] = v
				}
			}
			newSummary[k] = newRow
		}
		ret.Summary = newSummary
	}

	return ret, nil
}
