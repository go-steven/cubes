package test

import (
	"fmt"
	"testing"
	"gitlab.xibao100.com/skyline/skyline/cubes/metadata"
)

const REPORT_PATH = CUBES_BIN_DIR + "/reports"

func TestReport(t *testing.T) {
	test_cubes_cmd(t, fmt.Sprintf("%s/%s", REPORT_PATH, "client_daily_report_tpl.json"), fmt.Sprintf("%s/%s", REPORT_PATH, "config.json"), metadata.SQLVIEW_MODE)
	test_cubes_cmd(t, fmt.Sprintf("%s/%s", REPORT_PATH, "client_daily_report_tpl.yaml"), fmt.Sprintf("%s/%s", REPORT_PATH, "config.yaml"), metadata.SQLVIEW_MODE)
	test_cubes_cmd(t, fmt.Sprintf("%s/%s", REPORT_PATH, "client_weekly_report_tpl.json"), fmt.Sprintf("%s/%s", REPORT_PATH, "config.json"), metadata.SQLVIEW_MODE)
	test_cubes_cmd(t, fmt.Sprintf("%s/%s", REPORT_PATH, "client_weekly_report_tpl.yaml"), fmt.Sprintf("%s/%s", REPORT_PATH, "config.yaml"), metadata.SQLVIEW_MODE)
}