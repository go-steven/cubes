package metadata

const (
	SEP_LIST = ","
	SEP_EXP  = ";"
)

const (
	TPL_JSON = "json"
	TPL_YAML = "yaml"
)

const (
	EXT_CSV  = ".csv"
	EXT_JSON = ".json"
	EXT_YAML = ".yaml"
	EXT_YML  = ".yml"
)

type ReportTpl struct {
	Report  []string          `json:"report,omitempty", yaml:"report,omitempty"`
	Layout  interface{}       `json:"layout,omitempty", yaml:"layout,omitempty"`
	Cubes   []*CubeTpl        `json:"cubes,omitempty", yaml:"cubes,omitempty"`
	TplVars map[string]string `json:"tpl_vars,omitempty", yaml:"tpl_vars,omitempty"`
}

type CubeTpl struct {
	Name        string              `json:"name,omitempty", yaml:"name,omitempty"`
	Display     interface{}         `json:"display,omitempty", yaml:"display,omitempty"` // 在结果中原样返回，用户前段显示数据用
	Source      string              `json:"source,omitempty", yaml:"source,omitempty"`
	Store       string              `json:"store,omitempty", yaml:"store,omitempty"`
	Join        []*JoinTpl          `json:"join,omitempty", yaml:"join,omitempty"`
	Filter      []string            `json:"filter,omitempty", yaml:"filter,omitempty"`
	Mappings    []string            `json:"mappings,omitempty", yaml:"mappings,omitempty"`
	Dimensions  string              `json:"dimensions,omitempty", yaml:"dimensions,omitempty"`
	OrderBy     []string            `json:"orderby,omitempty", yaml:"orderby,omitempty"`
	Limit       string              `json:"limit,omitempty", yaml:"limit,omitempty"`
	Aggregates  [][]string          `json:"aggregates,omitempty", yaml:"aggregates,omitempty"`
	Tags        map[string][]string `json:"tags,omitempty", yaml:"tags,omitempty"`
	Union       []string            `json:"Union,omitempty", yaml:"Union,omitempty"`
	Sql         string              `json:"sql,omitempty", yaml:"sql,omitempty"`
	Summary     []string            `json:"summary,omitempty", yaml:"summary,omitempty`
	SummaryCalc map[string][]string `json:"summary_calc,omitempty", yaml:"summary_calc,omitempty`
	TmpFields   []string            `json:"tmp_fields,omitempty", yaml:"tmp_fields,omitempty`
}

type JoinTpl struct {
	Type  string   `json:"type,omitempty", yaml:"type,omitempty"`
	Store string   `json:"store,omitempty", yaml:"store,omitempty"`
	On    []string `json:"on,omitempty", yaml:"on,omitempty"`
}
