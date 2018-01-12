package metadata

import (
	"encoding/json"
	"fmt"
	"github.com/go-yaml/yaml"
	"gitlab.xibao100.com/skyline/skyline/cubes/utils"
	"io/ioutil"
	"path"
	"strings"
)

const TPL_SEP = "@@@@@"

type TplCfg map[string]interface{}

func (t TplCfg) ReplaceTpl(reportTpl string) string {
	for k, v := range t {
		reportTpl = replace_tpl_cfg(reportTpl, k, v)
	}

	return reportTpl
}

func ParseTplCfg(tplCfg []byte, cfgFormat string) (TplCfg, error) {
	cfg := make(TplCfg)
	switch cfgFormat {
	case TPL_JSON:
		if err := json.Unmarshal(tplCfg, &cfg); err != nil {
			logger.Errorf("ERROR Unmarshal: %v", err.Error())
			return nil, err
		}
	case TPL_YAML:
		if err := yaml.Unmarshal(tplCfg, &cfg); err != nil {
			logger.Errorf("ERROR Unmarshal: %v", err.Error())
			return nil, err
		}
	default:
		return nil, utils.Errorf("Unknown tpl cfg format:%s", cfgFormat)
	}

	return cfg, nil
}

func LoadTplCfgFile(tplCfgFile string) (TplCfg, error) {
	if tplCfgFile == "" {
		return nil, nil
	}

	bytes, err := ioutil.ReadFile(tplCfgFile)
	if err != nil {
		logger.Errorf("ERROR: failed to read file[%s] :%v", tplCfgFile, err.Error())
		return nil, err
	}

	var cfgFormat string
	tplExt := utils.LowerTrim(path.Ext(tplCfgFile))
	switch tplExt {
	case EXT_JSON:
		cfgFormat = TPL_JSON
	case EXT_YML:
		fallthrough
	case EXT_YAML:
		cfgFormat = TPL_YAML
	default:
		return nil, utils.Errorf("Unknown ext:%s", tplExt)
	}
	cfg, err := ParseTplCfg(bytes, cfgFormat)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	logger.Infof("TplCfg:%v", utils.Json(cfg))
	return cfg, nil
}

func replace_tpl_cfg(tpl, cfg_k string, cfg_v interface{}) string {
	return strings.Replace(tpl, fmt.Sprintf("%s%s%s", TPL_SEP, utils.UpperTrim(cfg_k), TPL_SEP), fmt.Sprintf("%v", cfg_v), -1)
}
