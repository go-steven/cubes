package metadata

import (
	"fmt"
	"gitlab.xibao100.com/skyline/skyline/cubes/source"
	"math"
	"strconv"
	"strings"
)

func (r *CubeReport) GenerateSummary(summary []*Summary) error {
	r.Summary = make(map[string]source.Row)
	if summary == nil {
		return nil
	}

	var (
		summaryRet source.Row
		err        error
	)
	for _, s := range summary {
		switch s.Type {
		case SUM:
			summaryRet, err = summary_sum(r.Data, r.Fields, s)
		case CONTRAST:
			summaryRet, err = summary_constrast(r.Data, r.Fields, s)
		}
		if err != nil {
			logger.Error(err)
			return err
		}
		r.Summary[s.Alias] = summaryRet
	}

	return nil
}

func summary_sum(data source.Rows, dataFields []string, summary *Summary) (source.Row, error) {
	ret := make(source.Row)
	for k, v := range dataFields {
		if k == 0 {
			ret[v] = summary.Alias
		} else {
			ret[v] = ""
		}
	}
	stats := make(map[string]float64)
	for _, v := range summary.Fields {
		stats[v] = 0.0
		ret[v] = trim_zero(fmt.Sprintf("%f", stats[v]))
	}
	if len(data) == 0 {
		return ret, nil
	}

	for _, row := range data {
		for k, v := range row {
			if _, ok := stats[k]; ok {
				floatVal, err := strconv.ParseFloat(v, 64)
				if err != nil {
					logger.Error(err)
					return nil, err
				}
				stats[k] += floatVal
			}
		}
	}

	for k, v := range stats {
		ret[k] = trim_zero(fmt.Sprintf("%f", v))
	}

	return ret, nil
}

func trim_zero(v string) string {
	vals := strings.Split(v, ".")
	if len(vals) < 2 {
		return v
	}

	str := vals[1]
	i := len(str) - 1
	for ; i >= 0; i-- {
		if string(str[i]) != "0" {
			break
		}
	}

	if i <= 0 {
		return vals[0]
	} else {
		s := ""
		for j := 0; j < i+1; j++ {
			s += string(str[j])
		}
		return fmt.Sprintf("%s.%s", vals[0], s)
	}
}
func summary_constrast(data source.Rows, dataFields []string, summary *Summary) (source.Row, error) {
	ret := make(source.Row)
	for k, v := range dataFields {
		if k == 0 {
			ret[v] = summary.Alias
		} else {
			ret[v] = ""
		}
	}

	//format := "%." + fmt.Sprintf("%d", summary.Dp) + "f%%"
	format := "%.0f%%"
	stats := make(map[string]float64)
	for _, v := range summary.Fields {
		stats[v] = 0.0
		ret[v] = trim_zero(fmt.Sprintf(format, stats[v]))
	}
	rowNum := len(data)
	if rowNum < 2 {
		return ret, nil
	}

	rowA := data[rowNum-1]
	rowB := data[rowNum-2]

	for k, v_a := range rowA {
		v_b := rowB[k]
		if _, ok := stats[k]; ok {
			var (
				floatValA, floatValB float64
				err                  error
			)
			if strings.Contains(v_a, "%") {
				v_a = strings.Replace(v_a, "%", "", -1)
				floatValA, err = strconv.ParseFloat(v_a, 64)
				if err != nil {
					logger.Error(err)
					return nil, err
				}
				floatValA = floatValA / 100.0
			} else {
				floatValA, err = strconv.ParseFloat(v_a, 64)
				if err != nil {
					logger.Error(err)
					return nil, err
				}
			}
			if strings.Contains(v_b, "%") {
				v_b = strings.Replace(v_b, "%", "", -1)
				floatValB, err = strconv.ParseFloat(v_b, 64)
				if err != nil {
					logger.Error(err)
					return nil, err
				}
				floatValB = floatValB / 100.0
			} else {
				floatValB, err = strconv.ParseFloat(v_b, 64)
				if err != nil {
					logger.Error(err)
					return nil, err
				}
			}
			if math.Abs(floatValB) > 0.00000001 {
				stats[k] += floatValA/floatValB - 1
			}
		}
	}

	for k, v := range stats {
		ret[k] = trim_zero(fmt.Sprintf(format, v*100))
	}

	return ret, nil
}

func get_default_summary_alias(summaryType string) (ret string) {
	switch summaryType {
	case SUM:
		ret = "汇总"
	case CONTRAST:
		ret = "环比"
	}
	return
}
