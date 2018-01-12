package metadata

import (
	"bytes"
	"errors"
	"fmt"
	"gitlab.xibao100.com/skyline/skyline/cubes/utils"
	"strings"
)

func (c *Cube) toSQLView(sqlViews *utils.MapData) (sqlview *SqlView, err error) {
	sqlview = NewSqlView(c)
	sqlview.SourceType = SOURCE_MYSQL
	sqlview.Fields, err = c.getReturnFields(sqlViews)
	if err != nil {
		return nil, err
	}

	if len(c.Union) > 0 && c.Source.Type == SOURCE_CUBE {
		sqlview.Sql, err = c.toUnionSQL(sqlViews)
		if err != nil {
			return nil, err
		}
		return sqlview, nil
	}
	if c.Sql != "" {
		sqlview.Sql = c.Sql
		sqlview.Fields = c.Dimensions
		return sqlview, nil
	}

	var buffer bytes.Buffer
	buffer.WriteString(`SELECT `)

	cnt := 0
	if len(c.Dimensions) > 0 {
		for _, v := range c.Dimensions {
			if cnt > 0 {
				buffer.WriteString(", ")
			}
			buffer.WriteString(format_field_name(v))
			cnt++
		}

		if len(c.Aggregates) > 0 {
			buffer.WriteString(", ")
		}
	}

	cnt = 0
	for f, fields := range c.Aggregates {
		for _, v := range fields {
			if cnt > 0 {
				buffer.WriteString(", ")
			}
			switch f {
			case "COUNT":
				if v.Expression == "*" || v.Expression == "1" {
					buffer.WriteString(fmt.Sprintf(" IFNULL(COUNT(1),0) AS `%s`", v.Alias))
				} else {
					buffer.WriteString(fmt.Sprintf(" IFNULL(COUNT(DISTINCT IFNULL(%s, 0)),0) AS `%s`", f, v.Expression, v.Alias))
				}
			default:
				buffer.WriteString(fmt.Sprintf(" IFNULL(%s(IFNULL(%s, 0)),0) AS `%s`", f, v.Expression, v.Alias))
			}
			cnt++
		}
	}
	if len(c.Dimensions) == 0 && len(c.Aggregates) == 0 {
		buffer.WriteString(fmt.Sprintf(" %s.* ", Sha1Name(TMP_TABLE_ALIAS)))
	}
	buffer.WriteString(` FROM ( `)
	buffer.WriteString(` SELECT `)
	cnt = 0
	base_select_fields, err := c.getDefaultFields(sqlViews)
	if err != nil {
		return nil, err
	}
	base_select_fields_map := make(map[string]struct{})
	for _, v := range base_select_fields {
		base_select_fields_map[utils.LowerTrim(v)] = struct{}{}
	}

	for _, v := range c.Mappings {
		if cnt > 0 {
			buffer.WriteString(", ")
		}
		buffer.WriteString(fmt.Sprintf("%s AS `%s`", v.Expression, v.Alias))
		if _, ok := base_select_fields_map[utils.LowerTrim(v.Alias)]; ok {
			delete(base_select_fields_map, utils.LowerTrim(v.Alias))
		}

		cnt++
	}
	if c.Source.Type == "mysql" && len(c.Tags) > 0 {
		cnt = 0
		for tagName, mappings := range c.Tags {
			if cnt > 0 || len(c.Mappings) > 0 {
				buffer.WriteString(", ")
			}
			buffer.WriteString(" CASE ")
			for _, m := range mappings {
				buffer.WriteString(fmt.Sprintf(` WHEN REPLACE(%s, " ", "") REGEXP "%s" `, format_field_name(m.Field), m.Val))
				if m.Val2 != "" {
					buffer.WriteString(fmt.Sprintf(` AND REPLACE(%s, " ", "") NOT REGEXP "%s" `, format_field_name(m.Field), m.Val2))
				}
				buffer.WriteString(fmt.Sprintf(` THEN "%s"`, m.TagVal))
			}
			buffer.WriteString(fmt.Sprintf(" ELSE \"\" END AS `%s`", tagName))

			if _, ok := base_select_fields_map[utils.LowerTrim(tagName)]; ok {
				delete(base_select_fields_map, utils.LowerTrim(tagName))
			}

			cnt++
		}
	}
	select_fields := []string{}
	for _, v := range base_select_fields {
		if _, ok := base_select_fields_map[utils.LowerTrim(v)]; ok {
			select_fields = append(select_fields, v)
		}
	}
	if len(c.Mappings) > 0 || (len(c.Tags) > 0 && c.Source.Type == "mysql") {
		buffer.WriteString(", ")
	}
	cnt = 0
	for _, v := range select_fields {
		if cnt > 0 {
			buffer.WriteString(", ")
		}
		buffer.WriteString(fmt.Sprintf(" %s.`%s` ", c.Store.Alias, v))
		cnt++
	}
	buffer.WriteString(` FROM `)
	if c.RunMode == SQLVIEW_MODE && c.Source.Type == SOURCE_CUBE && sqlViews != nil {
		v, err := get_sqlview_from_mapdata(c.Store.Sha1Name, sqlViews)
		if err != nil {
			return nil, err
		}
		buffer.WriteString(fmt.Sprintf("(%s) AS %s", v.Sql, c.Store.Alias))
	} else {
		buffer.WriteString(fmt.Sprintf("%s AS %s", c.Store.Sha1Name, c.Store.Alias))
	}

	if len(c.Join) > 0 {
		for _, join := range c.Join {
			if c.RunMode == SQLVIEW_MODE && c.Source.Type == SOURCE_CUBE && sqlViews != nil {
				logger.Infof("join.Store.Sha1Name = %s", join.Store.Sha1Name)
				v, err := get_sqlview_from_mapdata(join.Store.Sha1Name, sqlViews)
				if err != nil {
					return nil, err
				}

				buffer.WriteString(fmt.Sprintf(` %s (%s) AS %s `, join.Type, v.Sql, join.Store.Alias))
			} else {
				buffer.WriteString(fmt.Sprintf(` %s %s AS %s `, join.Type, join.Store.Sha1Name, join.Store.Alias))
			}

			buffer.WriteString(" ON (")
			cnt := 0
			for _, cond := range join.On {
				if cnt > 0 {
					buffer.WriteString(" AND ")
				}
				cond_val := fmt.Sprintf(`"%s"`, cond.Val)
				if len(strings.Split(cond.Val, ".")) > 1 {
					cond_val = cond.Val
				}
				switch cond.Op {
				case "BETWEEN":
					if len(strings.Split(cond.Val, "(")) > 1 || len(strings.Split(cond.Val2, "(")) > 1 {
						buffer.WriteString(fmt.Sprintf("%s BETWEEN %s AND %s", format_field_name(cond.Field), cond.Val, cond.Val2))
					} else {
						buffer.WriteString(fmt.Sprintf("%s BETWEEN \"%s\" AND \"%s\"", format_field_name(cond.Field), cond.Val, cond.Val2))
					}
				case "LIKE":
					buffer.WriteString(fmt.Sprintf("%s LIKE \"%s\"", format_field_name(cond.Field), cond_val))
				default:
					buffer.WriteString(fmt.Sprintf("%s %s %s", format_field_name(cond.Field), cond.Op, cond_val))
				}
				cnt++
			}
			if c.Source.Type == SOURCE_MYSQL && c.StoresLimit != nil && len(c.StoresLimit.FieldsSetting) > 0 {
				limitStore := c.StoresLimit.GetLimitStore(join.Store.Name)
				if limitStore != nil {
					for _, field := range limitStore.Fields {
						if v := c.StoresLimit.GetFieldSetting(field); v != nil {
							buffer.WriteString(fmt.Sprintf(" AND %s.`%s` = %v ", join.Store.Alias, field, v))
						}
					}
				}
			}

			buffer.WriteString(") ")
		}
	}

	buffer.WriteString(` WHERE 1=1 `)
	if c.Source.Type == SOURCE_MYSQL && c.StoresLimit != nil && len(c.StoresLimit.FieldsSetting) > 0 {
		limitStore := c.StoresLimit.GetLimitStore(c.Store.Name)
		if limitStore != nil {
			buffer.WriteString(` AND ( `)
			andCnt := 0
			for _, field := range limitStore.Fields {
				if v := c.StoresLimit.GetFieldSetting(field); v != nil {
					if andCnt > 0 {
						buffer.WriteString(" AND ")
					}
					buffer.WriteString(fmt.Sprintf(" %s.`%s` = %v ", c.Store.Alias, field, v))
					andCnt++
				}
			}
			buffer.WriteString(` ) `)
		}
	}
	if len(c.Filter) > 0 {
		buffer.WriteString(` AND ( `)

		cnt = 0
		for _, cond := range c.Filter {
			if cnt > 0 {
				buffer.WriteString(" AND ")
			}
			cond_val := fmt.Sprintf(`"%s"`, cond.Val)
			if len(strings.Split(cond.Val, ".")) > 1 {
				cond_val = cond.Val
			}
			switch cond.Op {
			case "BETWEEN":
				if len(strings.Split(cond.Val, "(")) > 1 || len(strings.Split(cond.Val2, "(")) > 1 {
					buffer.WriteString(fmt.Sprintf("%s BETWEEN %s AND %s", format_field_name(cond.Field), cond.Val, cond.Val2))
				} else {
					buffer.WriteString(fmt.Sprintf("%s BETWEEN \"%s\" AND \"%s\"", format_field_name(cond.Field), cond.Val, cond.Val2))
				}
			default:
				if cond_val == "" {
					buffer.WriteString(fmt.Sprintf("%s %s \"\"", format_field_name(cond.Field), cond.Op))
				} else {
					buffer.WriteString(fmt.Sprintf("%s %s %s", format_field_name(cond.Field), cond.Op, cond_val))
				}
			}
			cnt += 1
		}

		buffer.WriteString(` ) `)
	}
	buffer.WriteString(fmt.Sprintf(`) AS %s `, Sha1Name(TMP_TABLE_ALIAS)))

	if len(c.Aggregates) > 0 && len(c.Dimensions) > 0 {
		buffer.WriteString(` GROUP BY `)
		cnt = 0
		for _, v := range c.Dimensions {
			if cnt > 0 {
				buffer.WriteString(", ")
			}
			buffer.WriteString(format_field_name(v))
			cnt++
		}
	}

	if len(c.OrderBy) > 0 {
		buffer.WriteString(` ORDER BY `)
		cnt = 0
		for _, v := range c.OrderBy {
			if cnt > 0 {
				buffer.WriteString(", ")
			}
			buffer.WriteString(fmt.Sprintf("%s %s", format_field_name(v.Field), v.Order))
			cnt++
		}
	}
	if c.Limit != nil {
		buffer.WriteString(fmt.Sprintf(` LIMIT %d OFFSET %d `, c.Limit.Limit, c.Limit.Offset))
	}

	sqlview.Sql = buffer.String()
	return sqlview, nil
}

const TMP_TABLE_ALIAS = "XXXXXXXXXXXXXXXXXXXX_TMP"

func (c *Cube) toUnionSQL(sqlViews *utils.MapData) (string, error) {
	var buffer bytes.Buffer
	if len(c.Union) == 0 || c.Source.Type != SOURCE_CUBE {
		return "", errors.New("Only support cube union.")
	}

	fieldsCnt := 0
	cnt := 0
	for _, union := range c.Union {
		if cnt > 0 {
			buffer.WriteString(fmt.Sprintf(" %s ", union.UnionType))
		}

		buffer.WriteString(` SELECT `)

		var (
			fields []string
			err    error
		)
		if c.RunMode == SQLVIEW_MODE && sqlViews != nil {
			fields, err = c.getReturnFields(sqlViews)
		} else {
			fields, err = c.Sqlite.GetTableFields(union.Sha1Name)
		}
		if err != nil {
			logger.Error(err)
			return "", err
		}
		if fieldsCnt == 0 {
			fieldsCnt = len(fields)
		}
		if len(fields) != fieldsCnt {
			return "", errors.New(fmt.Sprintf("Cube union: %s, fields cnt not same.", union.Name))
		}

		fieldCnt := 0
		for _, v := range fields {
			if fieldCnt > 0 {
				buffer.WriteString(", ")
			}
			buffer.WriteString(format_field_name(v))
			fieldCnt++
		}

		if c.RunMode == SQLVIEW_MODE && sqlViews != nil {
			v, err := get_sqlview_from_mapdata(union.Sha1Name, sqlViews)
			if err != nil {
				logger.Error(err)
				return "", err
			}
			buffer.WriteString(fmt.Sprintf(" FROM (%s) AS %s", v.Sql, union.Sha1Name))
		} else {
			buffer.WriteString(fmt.Sprintf(" FROM %s AS %s", union.Sha1Name, union.Sha1Name))
		}

		cnt++
	}
	return buffer.String(), nil
}

func (c *Cube) getReturnFields(sqlViews *utils.MapData) ([]string, error) {
	if len(c.Union) > 0 && c.Source.Type == SOURCE_CUBE {
		fields := []string{}
		for _, v := range c.Dimensions {
			fields = append(fields, v)
		}
		if len(fields) == 0 {
			if c.RunMode == SQLVIEW_MODE && c.Source.Type == SOURCE_CUBE && sqlViews != nil {
				var storeSha1Name string
				if len(c.Union) > 0 {
					storeSha1Name = c.Union[0].Sha1Name
				} else {
					storeSha1Name = c.Store.Sha1Name
				}
				v, err := get_sqlview_from_mapdata(storeSha1Name, sqlViews)
				if err != nil {
					return nil, err
				}
				fields = v.Fields
			} else {
				tableFields, err := c.Sqlite.GetTableFields(c.Union[0].Sha1Name)
				if err != nil {
					logger.Error(err)
					return nil, err
				}
				fields = tableFields
			}
		}

		return fields, nil
	}

	fields := []string{}
	for _, v := range c.Dimensions {
		fields = append(fields, v)
	}
	for _, aggregate := range c.Aggregates {
		for _, v := range aggregate {
			fields = append(fields, v.Alias)
		}
	}

	if len(fields) == 0 {
		if c.RunMode == SQLVIEW_MODE && c.Source.Type == SOURCE_CUBE && sqlViews != nil {
			var storeSha1Name string
			if len(c.Union) > 0 {
				storeSha1Name = c.Union[0].Sha1Name
			} else {
				storeSha1Name = c.Store.Sha1Name
			}
			v, err := get_sqlview_from_mapdata(storeSha1Name, sqlViews)
			if err != nil {
				return nil, err
			}
			fields = v.Fields
		} else {
			switch c.Source.Type {
			case SOURCE_CUBE:
				fallthrough
			case SOURCE_CSV:
				fallthrough
			case SOURCE_JSON:
				if retFields, err := c.Sqlite.GetTableFields(c.Store.Sha1Name); err == nil {
					fields = retFields
				}
			case SOURCE_SQLITE:
				if retFields, err := c.Sqlite.GetTableFields(c.Store.Name); err == nil {
					fields = retFields
				}
			case SOURCE_MYSQL:
				if retFields, err := c.Mysql.GetTableFields(c.Store.Name); err == nil {
					fields = retFields
				}
			}
		}
	}
	return fields, nil
}

func get_sqlview_from_mapdata(sha1name string, sqlViews *utils.MapData) (*SqlView, error) {
	v := sqlViews.Get(sha1name)
	if v == nil {
		logger.Infof("sql view need wait: %s", sha1name)
		return nil, error_cube_need_wait
	}
	sqlview, ok := v.(*SqlView)
	if !ok {
		err := errors.New(fmt.Sprintf("SqlView[%v] not found.", v))
		logger.Error(err)
		return nil, err
	}

	return sqlview, nil
}

func (c *Cube) getDefaultFields(sqlViews *utils.MapData) ([]string, error) {
	if c.RunMode == SQLVIEW_MODE && c.Source.Type == SOURCE_CUBE && sqlViews != nil {
		v, err := get_sqlview_from_mapdata(c.Store.Sha1Name, sqlViews)
		if err != nil {
			return nil, err
		}

		return v.Fields, nil
	}

	var (
		fields []string
		err    error
	)
	switch c.Source.Type {
	case SOURCE_MYSQL:
		fields, err = c.Mysql.GetTableFields(c.Store.Name)
	case SOURCE_CUBE:
		fallthrough
	case SOURCE_CSV:
		fallthrough
	case SOURCE_JSON:
		fields, err = c.Sqlite.GetTableFields(c.Store.Sha1Name)
	case SOURCE_SQLITE:
		fields, err = c.Sqlite.GetTableFields(c.Store.Name)
	}
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	if len(fields) == 0 {
		return nil, utils.Errorf("No fields for %s", c.Store.Name)
	}

	return fields, nil
}

func format_field_name(field string) string {
	vals := strings.Split(field, ".")
	if len(vals) > 1 {
		return fmt.Sprintf("%s.`%s`", vals[0], vals[1])
	}

	return fmt.Sprintf("`%s`", field)
}
