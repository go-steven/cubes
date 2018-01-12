package metadata

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-yaml/yaml"
	"gitlab.xibao100.com/skyline/skyline/cubes/utils"
	"io/ioutil"
	"path"
	"regexp"
	"strconv"
	"strings"
)

func ParseTpl(tpl []byte, tplFormat string) (*ReportTpl, error) {
	rptTpl := ReportTpl{}
	switch tplFormat {
	case TPL_JSON:
		if err := json.Unmarshal(tpl, &rptTpl); err != nil {
			logger.Errorf("ERROR Unmarshal: %v", err.Error())
			return nil, err
		}
	case TPL_YAML:
		if err := yaml.Unmarshal(tpl, &rptTpl); err != nil {
			logger.Errorf("ERROR Unmarshal: %v", err.Error())
			return nil, err
		}
	default:
		return nil, utils.Errorf("Unknown tpl format:%s", tplFormat)
	}

	return &rptTpl, nil
}

func LoadReportTplFile(tplFile string, tplCfgFile string) (*ReportTpl, error) {
	logger.Infof("LoadReportTplFile: tplFile = %s, tplCfgFile = %s", tplFile, tplCfgFile)
	bytes, err := ioutil.ReadFile(tplFile)
	if err != nil {
		logger.Errorf("ERROR: failed to read file[%s] :%v", tplFile, err.Error())
		return nil, err
	}
	content := string(bytes)
	logger.Infof("tpl content: %v", content)
	tplCfg, err := LoadTplCfgFile(tplCfgFile)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	if tplCfg != nil && len(tplCfg) > 0 {
		content = tplCfg.ReplaceTpl(content)
	}
	if strings.Contains(content, TPL_SEP) {
		return nil, errors.New("Report Tpl still has variables.")
	}

	var tplFormat string
	tplExt := utils.LowerTrim(path.Ext(tplFile))
	logger.Infof("tplExt = %s", tplExt)
	switch tplExt {
	case EXT_JSON:
		tplFormat = TPL_JSON
	case EXT_YML:
		fallthrough
	case EXT_YAML:
		tplFormat = TPL_YAML
	default:
		return nil, utils.Errorf("Unknown ext:%s", tplExt)
	}

	report, err := ParseTpl([]byte(content), tplFormat)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	logger.Infof("ReportTpl:%s", utils.Json(report))
	return report, nil
}

func (r *CubeTpl) Convert() (*Cube, error) {
	// Trim all string fields
	r.Name = utils.Trim(r.Name)
	r.Store = utils.Trim(r.Store)
	r.Dimensions = utils.Trim(r.Dimensions)
	r.Limit = utils.Trim(r.Limit)
	r.Source = utils.Trim(r.Source)
	if r.Source == "" {
		r.Source = SOURCE_CUBE
	}

	if r.Sql != "" {
		r.Sql = strings.Replace(r.Sql, "\n", " ", -1)
		r.Sql = strings.Replace(r.Sql, "\t", " ", -1)
	}
	cube := &Cube{
		Name:       r.Name,
		Display:    r.Display,
		Source:     &Source{},
		Filter:     []*Condition{},
		Mappings:   []*Mapping{},
		Dimensions: []string{},
		OrderBy:    []*OrderBy{},
		Aggregates: make(map[string][]*Mapping),
		Tags:       make(map[string][]*TagMapping),
		Union:      []*Union{},
		Join:       []*Join{},
		Sql:        utils.Trim(r.Sql),
		TmpFields:  make(map[string]struct{}),
	}
	cube.Sha1Name = Sha1Name(cube.Name)

	vals := strings.Split(r.Dimensions, SEP_LIST)
	for _, v := range vals {
		v = utils.Trim(v)
		if v != "" {
			cube.Dimensions = append(cube.Dimensions, v)
		}
	}

	for _, str := range r.OrderBy {
		orderBy := &OrderBy{
			Order: "ASC",
		}
		vals = strings.Split(str, SEP_LIST)
		if len(vals) == 0 {
			continue
		}
		orderBy.Field = utils.Trim(vals[0])
		if len(vals) >= 2 {
			orderBy.Order = utils.UpperTrim(vals[1])
		}
		cube.OrderBy = append(cube.OrderBy, orderBy)
	}

	if len(cube.OrderBy) == 0 && len(cube.Dimensions) > 0 {
		for _, v := range cube.Dimensions {
			cube.OrderBy = append(cube.OrderBy, &OrderBy{
				Field: v,
				Order: "ASC",
			})
		}
	}
	vals = strings.Split(r.Limit, SEP_LIST)
	if len(vals) > 1 {
		cube.Limit = &Limit{
			Limit:  1000,
			Offset: 0,
		}

		i, err := strconv.ParseInt(vals[0], 10, 64)
		if err == nil {
			cube.Limit.Limit = int(i)
		}

		if len(vals) >= 2 {
			i, err := strconv.ParseInt(vals[1], 10, 64)
			if err == nil {
				cube.Limit.Offset = int(i)
			}
		}
	}

	vals = strings.Split(r.Source, SEP_LIST)
	if len(vals) >= 1 {
		cube.Source.Type = utils.LowerTrim(vals[0])
	}
	if len(vals) >= 2 {
		cube.Source.Name = utils.Trim(vals[1])
	}

	//logger.Infof("r.Store = %s", r.Store)
	if r.Store != "" {
		vals = strings.Split(r.Store, SEP_LIST)
		if len(vals) > 0 {
			name := utils.Trim(vals[0])
			if name != "" {
				var alias string
				if len(vals) >= 2 {
					alias = utils.Trim(vals[1])
				}
				if alias == "" {
					alias = "t0"
				}
				cube.Store = &Store{
					Name:  name,
					Alias: alias,
				}
				if cube.Source.Type == SOURCE_MYSQL || cube.Source.Type == SOURCE_SQLITE {
					cube.Store.Sha1Name = cube.Store.Name
				} else {
					cube.Store.Sha1Name = Sha1Name(cube.Store.Name)
				}
			}
		}
	}
	//logger.Infof("cube.Store = %s", utils.Json(cube.Store))

	for _, andStr := range r.Filter {
		and := strings.Split(andStr, SEP_EXP)
		if len(and) < 2 {
			continue
		}
		cond := &Condition{
			Field: utils.Trim(and[0]),
			Op:    utils.UpperTrim(and[1]),
		}
		if len(and) >= 3 {
			cond.Val = utils.Trim(and[2])
		}
		if len(and) >= 4 {
			cond.Val2 = utils.Trim(and[3])
		}

		cube.Filter = append(cube.Filter, cond)
	}

	for _, v := range r.Mappings {
		vals := strings.Split(v, SEP_EXP)
		if len(vals) < 2 {
			continue
		}
		mapping := &Mapping{
			Alias:      utils.Trim(vals[0]),
			Expression: utils.Trim(vals[1]),
		}
		cube.Mappings = append(cube.Mappings, mapping)
	}

	for _, aggregate := range r.Aggregates {
		if len(aggregate) < 2 {
			continue
		}
		function := aggregate[0]
		mappings := []*Mapping{}
		for k, v := range aggregate {
			if k == 0 {
				continue
			}
			vals := strings.Split(v, SEP_EXP)
			if len(vals) == 0 {
				continue
			}
			mapping := &Mapping{
				Expression: utils.Trim(vals[0]),
			}
			if len(vals) >= 2 {
				mapping.Alias = utils.Trim(vals[1])
			}
			if mapping.Alias == "" {
				mapping.Alias = fmt.Sprintf("%s_%s", utils.LowerTrim(function), mapping.Expression)
			}
			mappings = append(mappings, mapping)
		}
		cube.Aggregates[function] = mappings
	}

	for tagName, mappings := range r.Tags {
		tagMappings := []*TagMapping{}
		for _, v := range mappings {
			vals := strings.Split(v, SEP_EXP)
			if len(vals) < 4 {
				continue
			}

			tagMapping := &TagMapping{
				TagVal: utils.Trim(vals[0]),
				Field:  utils.Trim(vals[1]),
				Op:     utils.UpperTrim(vals[2]),
				Val:    utils.Trim(vals[3]),
			}
			if len(vals) >= 5 {
				tagMapping.Val2 = utils.Trim(vals[4])
			}

			// TO DO: support more functions for tags
			if tagMapping.Op != "REGEXP" {
				return nil, errors.New("For now, only `regexp` supported for tags.")
			}
			reg, err := regexp.Compile(tagMapping.Val)
			if err != nil {
				logger.Error(err)
				return nil, err
			}
			tagMapping.IncludeRegexp = reg

			if tagMapping.Val2 != "" {
				reg, err = regexp.Compile(tagMapping.Val2)
				if err != nil {
					logger.Error(err)
					return nil, err
				}
				tagMapping.ExcludeRegexp = reg
			}
			tagMappings = append(tagMappings, tagMapping)
		}
		if len(tagMappings) > 0 {
			cube.Tags[tagName] = tagMappings
		}
	}

	for _, v := range r.Union {
		vals := strings.Split(v, SEP_LIST)
		if len(vals) >= 1 {
			union := &Union{
				Name:      utils.Trim(vals[0]),
				UnionType: UNION,
			}
			if len(vals) >= 2 {
				union.UnionType = utils.UpperTrim(vals[1])
			}
			if cube.Source.Type == SOURCE_MYSQL || cube.Source.Type == SOURCE_SQLITE {
				union.Sha1Name = union.Name
			} else {
				union.Sha1Name = Sha1Name(union.Name)
			}

			cube.Union = append(cube.Union, union)
		}
	}

	for _, v := range r.Join {
		join := v.Convert(cube.Source.Type)
		if join != nil {
			cube.Join = append(cube.Join, join)
		}
	}
	for _, summary_str := range r.Summary {
		vals := strings.Split(summary_str, SEP_EXP)
		if len(vals) <= 0 {
			continue
		}
		vals_1 := strings.Split(vals[0], SEP_LIST)
		if len(vals_1) <= 0 {
			continue
		}

		summary := &Summary{
			Type:   utils.UpperTrim(vals_1[0]),
			Fields: []string{},
		}

		if len(vals_1) >= 2 {
			summary.Alias = utils.Trim(vals_1[1])
		}
		if summary.Alias == "" {
			summary.Alias = get_default_summary_alias(summary.Type)
		}

		if len(vals) >= 2 {
			fields := strings.Split(vals[1], SEP_LIST)
			for _, v := range fields {
				val := utils.Trim(v)
				if val != "" {
					summary.Fields = append(summary.Fields, val)
				}
			}
		}
		if len(vals) >= 3 {
			dp, _ := utils.ParseUint(vals[2], 10, 64)
			summary.Dp = uint8(dp)
		}
		if summary.Dp <= 0 {
			summary.Dp = 2
		}
		cube.Summary = append(cube.Summary, summary)
	}

	for k, calc_list := range r.SummaryCalc {
		if _, ok := cube.SummaryCalc[k]; !ok {
			cube.SummaryCalc = make(map[string][]*SummaryCalc)
		}
		for _, calc_str := range calc_list {
			vals := strings.Split(calc_str, SEP_EXP)
			if len(vals) < 2 {
				continue
			}
			summaryCalc := &SummaryCalc{
				Alias:  vals[0],
				Params: []string{},
			}
			if len(vals) >= 3 {
				dp, _ := utils.ParseUint(vals[2], 10, 64)
				summaryCalc.Dp = uint8(dp)
			}
			if len(vals) >= 4 && vals[3] != "" {
				summaryCalc.Multiple, _ = strconv.ParseFloat(vals[3], 64)
			}

			op_vals := strings.Split(vals[1], SEP_LIST)
			if len(op_vals) < 2 {
				continue
			}
			summaryCalc.Op = op_vals[0]
			summaryCalc.Params = append(summaryCalc.Params, op_vals[1:]...)
			cube.SummaryCalc[k] = append(cube.SummaryCalc[k], summaryCalc)
		}
	}

	for _, v := range r.TmpFields {
		val := utils.Trim(v)
		if val != "" {
			cube.TmpFields[val] = struct{}{}
		}
	}
	//logger.Infof("cube = %s", utils.Json(cube))
	if err := CheckCube(cube); err != nil {
		logger.Error(err)
		return nil, err
	}

	return cube, nil
}

func (j *JoinTpl) Convert(sourceType string) *Join {
	vals := strings.Split(j.Store, SEP_LIST)
	if len(vals) == 0 {
		return nil
	}
	store := &Store{
		Name: utils.Trim(vals[0]),
	}
	if len(vals) >= 2 {
		store.Alias = utils.Trim(vals[1])
	}
	if store.Alias == "" {
		store.Alias = store.Name
	}
	if sourceType == SOURCE_MYSQL || sourceType == SOURCE_SQLITE {
		store.Sha1Name = store.Name
	} else {
		store.Sha1Name = Sha1Name(store.Name)
	}
	conds := []*Condition{}
	for _, v := range j.On {
		vals2 := strings.Split(v, SEP_EXP)
		if len(vals2) < 3 {
			continue
		}

		cond := &Condition{
			Field: utils.Trim(vals2[0]),
			Op:    utils.UpperTrim(vals2[1]),
			Val:   utils.Trim(vals2[2]),
		}
		if len(vals2) >= 4 {
			cond.Val2 = utils.Trim(vals2[3])
		}

		conds = append(conds, cond)
	}
	if len(conds) == 0 {
		return nil
	}

	join := &Join{
		Type:  utils.UpperTrim(j.Type),
		Store: store,
		On:    conds,
	}
	return join
}

func GetTplFormat(tpl string) (format string) {
	format = TPL_JSON
	if !strings.HasPrefix(tpl, "{") {
		format = TPL_YAML
	}
	return
}
