package metadata

import (
	"fmt"
	"math"
	"strconv"
)

func (r *CubeReport) GenerateSummaryCalc(summaryCalcs map[string][]*SummaryCalc) error {
	if len(summaryCalcs) == 0 {
		return nil
	}

	for k, calcs := range summaryCalcs {
		if _, ok := r.Summary[k]; !ok {
			continue
		}
		for _, v := range calcs {
			var ret_val float64
			var flag bool
			switch v.Op {
			case "/":
				if len(v.Params) < 2 {
					break
				}
				f1_val_str, ok := r.Summary[k][v.Params[0]]
				if !ok {
					break
				}
				f1_val, err := strconv.ParseFloat(f1_val_str, 64)
				if err != nil {
					break
				}
				//logger.Infof("v.Params[0] = %v, f1_val = %v", v.Params[0], f1_val)
				f2_val_str, ok := r.Summary[k][v.Params[1]]
				if !ok {
					break
				}
				f2_val, err := strconv.ParseFloat(f2_val_str, 64)
				if err != nil {
					break
				}
				//logger.Infof("v.Params[1] = %v, f2_val = %v", v.Params[1], f2_val)
				if math.Abs(f2_val) > 0.000001 {
					ret_val = f1_val / f2_val
				}
				//logger.Infof("op = %v, ret_val = %v", v.Op, ret_val)
				flag = true
			}
			if flag {
				if v.Multiple > 0 {
					ret_val = ret_val * v.Multiple
				}
				format := "%v"
				if v.Dp > 0 {
					format = fmt.Sprintf("%%0.%df", v.Dp)
				}
				r.Summary[k][v.Alias] = fmt.Sprintf(format, ret_val)
				//logger.Infof("k = %v, v.Alias = %v, r.Summary[k] = %v, r.Summary[k][v.Alias]=%v", k, v.Alias, r.Summary[k], r.Summary[k][v.Alias])
			}
		}
	}

	return nil
}
