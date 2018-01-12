package metadata

import (
	"errors"
	"fmt"
	"gitlab.xibao100.com/skyline/skyline/cubes/source"
	"gitlab.xibao100.com/skyline/skyline/cubes/utils"
)

func NewSqlView(cube *Cube) *SqlView {
	return &SqlView{
		Name:        cube.Name,
		Sha1Name:    cube.Sha1Name,
		Display:     cube.Display,
		SourceType:  SOURCE_MYSQL,
		Fields:      []string{},
		Summary:     cube.Summary,
		SummaryCalc: cube.SummaryCalc,
		TmpFields:   cube.TmpFields,

		ESqlite: cube.ESqlite,
		Sqlite:  cube.Sqlite,
		Mysql:   cube.Mysql,
	}
}
func (v *SqlView) Execute() (*CubeReport, error) {
	if v.Sql == "" || len(v.Fields) == 0 {
		return nil, utils.Errorf("ERROR: Invalid sql view: %s", v.Name)
	}

	var (
		data source.Rows
		err  error
	)
	switch v.SourceType {
	case SOURCE_MYSQL:
		if v.Mysql == nil {
			return nil, utils.Error("ERROR: No MYSQL connection.")
		}

		logger.Infof("Run sqlview[%s] sql: %s", v.Name, v.Sql)
		data, err = v.Mysql.Query(v.Sql, v.Fields)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
	case SOURCE_SQLITE:
		if v.Sqlite == nil {
			return nil, errors.New("ERROR: No SQLITE connection.")
		}
		logger.Infof("Run sqlview[%s] sql: %s", v.Name, v.Sql)
		data, err = v.Sqlite.Query(v.Sql, v.Fields)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
	default:
		return nil, errors.New(fmt.Sprintf("Unknow source type:%s", v.SourceType))
	}

	ret := &CubeReport{
		Name:    v.Name,
		Display: v.Display,
		Fields:  v.Fields,
		Data:    data,
	}

	if err := ret.GenerateSummary(v.Summary); err != nil {
		logger.Error(err)
		return nil, err
	}

	if err := ret.GenerateSummaryCalc(v.SummaryCalc); err != nil {
		logger.Error(err)
		return nil, err
	}

	// 处理不返回的字段列表
	if len(v.TmpFields) > 0 {
		newFields := []string{}
		for _, field := range ret.Fields {
			if _, ok := v.TmpFields[field]; !ok {
				newFields = append(newFields, field)
			}
		}
		ret.Fields = newFields

		newRows := source.Rows{}
		for _, row := range ret.Data {
			newRow := make(source.Row)
			for k, val := range row {
				if _, ok := v.TmpFields[k]; !ok {
					newRow[k] = val
				}
			}
			newRows = append(newRows, newRow)
		}
		ret.Data = newRows

		newSummary := make(map[string]source.Row)
		for k, row := range ret.Summary {
			newRow := make(source.Row)
			for k, val := range row {
				if _, ok := v.TmpFields[k]; !ok {
					newRow[k] = val
				}
			}
			newSummary[k] = newRow
		}
		ret.Summary = newSummary
	}

	return ret, nil
}
