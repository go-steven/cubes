package metadata

import (
	"errors"
	"gitlab.xibao100.com/skyline/skyline/cubes/source"
	"gitlab.xibao100.com/skyline/skyline/cubes/utils"
	"regexp"
)

const (
	SOURCE_MYSQL  = "mysql"
	SOURCE_SQLITE = "sqlite"
	SOURCE_CSV    = "csv"
	SOURCE_JSON   = "json"
	SOURCE_CUBE   = "cube"
)

const (
	INNER_JOIN = "INNER JOIN"
	LEFT_JOIN  = "LEFT JOIN"
	RIGHT_JOIN = "RIGHT JOIN"
)

const (
	ORDER_ASC  = "ASC"
	ORDER_DESC = "DESC"
)

const (
	SUM      = "SUM"
	AVG      = "AVG"
	COUNT    = "COUNT"
	MIN      = "MIN"
	MAX      = "MAX"
	CONTRAST = "CONTRAST"
)

const (
	UNION     = "UNION"
	UNION_ALL = "UNION ALL"
)

const (
	DEFAULT_MODE = "default"
	SQLVIEW_MODE = "sqlview"
)

const ERROR_CUBE_NEED_WAIT = "ERROR: cube need wait."

var error_cube_need_wait = errors.New(ERROR_CUBE_NEED_WAIT)

type Report struct {
	Report []string       `json:"report,omitempty", yaml:"report,omitempty`
	Layout interface{}    `json:"layout,omitempty", yaml:"layout,omitempty`
	Cubes  *utils.MapData `json:"cubes,omitempty", yaml:"cubes,omitempty`

	RunMode string `json:"-", yaml:"-"`
}

type Cube struct {
	Name        string                    `json:"name,omitempty", yaml:"name,omitempty`
	Sha1Name    string                    `json:"sha1_name,omitempty", yaml:"sha1_name,omitempty`
	Display     interface{}               `json:"display,omitempty", yaml:"display,omitempty` // 在结果中原样返回，用户前段显示数据用
	Source      *Source                   `json:"source,omitempty", yaml:"source,omitempty`
	Store       *Store                    `json:"store,omitempty", yaml:"store,omitempty`
	Join        []*Join                   `json:"join,omitempty", yaml:"join,omitempty`
	Filter      []*Condition              `json:"filter,omitempty", yaml:"filter,omitempty`
	Mappings    []*Mapping                `json:"mappings,omitempty", yaml:"mappings,omitempty`
	Dimensions  []string                  `json:"dimensions,omitempty", yaml:"dimensions,omitempty`
	OrderBy     []*OrderBy                `json:"orderby,omitempty", yaml:"orderby,omitempty`
	Limit       *Limit                    `json:"limit,omitempty", yaml:"limit,omitempty`
	Aggregates  map[string][]*Mapping     `json:"aggregates,omitempty", yaml:"aggregates,omitempty`
	Tags        map[string][]*TagMapping  `json:"tags,omitempty", yaml:"tags,omitempty`
	Union       []*Union                  `json:"union,omitempty", yaml:"union,omitempty`
	Sql         string                    `json:"sql,omitempty", yaml:"sql,omitempty`
	Summary     []*Summary                `json:"summary,omitempty", yaml:"summary,omitempty`
	SummaryCalc map[string][]*SummaryCalc `json:"summary_calc,omitempty", yaml:"summary_calc,omitempty`

	StoresLimit *StoresLimit `json:"stores_limit,omitempty", yaml:"stores_limit,omitempty`

	TmpFields map[string]struct{} `json:"tmp_fields,omitempty", yaml:"tmp_fields,omitempty`

	RunMode string         `json:"-", yaml:"-"`
	ESqlite *source.Sqlite `json:"-", yaml:"-"`
	Sqlite  *source.Sqlite `json:"-", yaml:"-"`
	Mysql   *source.Mysql  `json:"-", yaml:"-"`
}

type Source struct {
	Name string `json:"name,omitempty", yaml:"name,omitempty`
	Type string `json:"type,omitempty", yaml:"type,omitempty`
}
type Mapping struct {
	Alias      string `json:"alias,omitempty", yaml:"alias,omitempty`
	Expression string `json:"expression,omitempty", yaml:"expression,omitempty`
}
type TagMapping struct {
	TagVal        string         `json:"tag_val,omitempty", yaml:"tag_val,omitempty`
	Field         string         `json:"field,omitempty", yaml:"field,omitempty`
	Op            string         `json:"op,omitempty", yaml:"op,omitempty`
	Val           string         `json:"val,omitempty", yaml:"val,omitempty`
	Val2          string         `json:"val2,omitempty", yaml:"val2,omitempty`
	IncludeRegexp *regexp.Regexp `json:"-", yaml:"-"`
	ExcludeRegexp *regexp.Regexp `json:"-", yaml:"-"`
}

type Condition struct {
	Field string `json:"field,omitempty", yaml:"field,omitempty`
	Op    string `json:"op,omitempty", yaml:"op,omitempty`
	Val   string `json:"val,omitempty", yaml:"val,omitempty`
	Val2  string `json:"val2,omitempty", yaml:"val2,omitempty`
}

type Store struct {
	Name     string `json:"name,omitempty", yaml:"name,omitempty`
	Sha1Name string `json:"sha1_name,omitempty", yaml:"sha1_name,omitempty`
	Alias    string `json:"alias,omitempty", yaml:"alias,omitempty`
}

type OrderBy struct {
	Field string `json:"field,omitempty", yaml:"field,omitempty`
	Order string `json:"order,omitempty", yaml:"order,omitempty`
}

type Limit struct {
	Limit  int `json:"limit,omitempty", yaml:"limit,omitempty`
	Offset int `json:"offset,omitempty", yaml:"offset,omitempty`
}

type Union struct {
	Name      string `json:"name,omitempty", yaml:"name,omitempty`
	Sha1Name  string `json:"sha1_name,omitempty", yaml:"sha1_name,omitempty`
	UnionType string `json:"union_type,omitempty", yaml:"union_type,omitempty`
}

type Summary struct {
	Type   string   `json:"type,omitempty", yaml:"type,omitempty`
	Alias  string   `json:"alias,omitempty", yaml:"alias,omitempty`
	Fields []string `json:"fields,omitempty", yaml:"fields,omitempty`
	Dp     uint8    `json:"dp,omitempty", yaml:"dp,omitempty` // decimal point
}

type SummaryCalc struct {
	Alias    string   `json:"alias,omitempty", yaml:"alias,omitempty`
	Op       string   `json:"op,omitempty", yaml:"op,omitempty`
	Params   []string `json:"params,omitempty", yaml:"params,omitempty`
	Dp       uint8    `json:"dp,omitempty", yaml:"dp,omitempty`             // decimal point
	Multiple float64  `json:"multiple,omitempty", yaml:"multiple,omitempty` // 乘以倍数
}

type Join struct {
	Type  string       `json:"type,omitempty", yaml:"type,omitempty`
	Store *Store       `json:"store,omitempty", yaml:"store,omitempty`
	On    []*Condition `json:"on,omitempty", yaml:"on,omitempty`
}

type CubeReport struct {
	Name    string                `json:"-", yaml:"-"`
	Display interface{}           `json:"display", yaml:"display,omitempty` // 在结果中原样返回，用户前段显示数据用
	Fields  []string              `json:"fields", yaml:"fields,omitempty`
	Data    source.Rows           `json:"data", yaml:"data,omitempty`
	Summary map[string]source.Row `json:"summary", yaml:"data,omitempty`
}

type ReportResult struct {
	Layout interface{}            `json:"layout", yaml:"layout,omitempty` // 在结果中原样返回，用户前段显示数据用
	Cubes  map[string]*CubeReport `json:"cubes", yaml:"cubes,omitempty`
}

type SqlView struct {
	Name        string                    `json:"name", yaml:"name,omitempty`
	Sha1Name    string                    `json:"sha1_name,omitempty", yaml:"sha1_name,omitempty`
	Display     interface{}               `json:"display,omitempty", yaml:"display,omitempty` // 在结果中原样返回，用户前段显示数据用
	SourceType  string                    `json:"source_type,omitempty", yaml:"source_type,omitempty`
	Fields      []string                  `json:"fields", yaml:"fields,omitempty`
	Sql         string                    `json:"sql", yaml:"sql,omitempty`
	Summary     []*Summary                `json:"summary,omitempty", yaml:"summary,omitempty`
	SummaryCalc map[string][]*SummaryCalc `json:"summary_calc,omitempty", yaml:"summary_calc,omitempty`
	TmpFields   map[string]struct{}       `json:"tmp_fields,omitempty", yaml:"tmp_fields,omitempty`

	ESqlite *source.Sqlite `json:"-", yaml:"-"`
	Sqlite  *source.Sqlite `json:"-", yaml:"-"`
	Mysql   *source.Mysql  `json:"-", yaml:"-"`
}
