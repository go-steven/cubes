package metadata

import (
	"strings"
)

func GetTplVars(tpl []byte, tplCfg TplCfg) (TplCfg, error) {
	if tplCfg == nil || len(tplCfg) == 0 {
		return tplCfg, nil
	}

	tplFormat := GetTplFormat(string(tpl))
	rptTpl, err := ParseTpl(tpl, tplFormat)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	tplVars := rptTpl.TplVars
	loop := 0
	for {
		if len(tplVars) == 0 {
			break
		}
		logger.Infof("LOOP: %d", loop)
		if loop > MAX_LOOP {
			break
		}

		for cfg_k, cfg_v := range tplCfg {
			if len(tplVars) == 0 {
				break
			}
			for var_k, var_v := range tplVars {
				var_v = replace_tpl_cfg(var_v, cfg_k, cfg_v)
				tplVars[var_k] = var_v

				if !strings.Contains(var_v, TPL_SEP) {

					tplCfg[var_k] = var_v
					delete(tplVars, var_k)
				}
			}
		}
		loop++
	}

	return tplCfg, nil
}
