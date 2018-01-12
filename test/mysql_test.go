package test

import (
	"fmt"
	"testing"
	"gitlab.xibao100.com/skyline/skyline/cubes/metadata"
)

const MYSQL_EXAMPLE_PATH = CUBES_BIN_DIR + "/example/mysql"

func TestMysql(t *testing.T) {
	tpls := []string{}
	json_tpls, err := list_dir(MYSQL_EXAMPLE_PATH, ".json")
	if err != nil {
		t.Error(err.Error())
	}
	tpls = append(tpls, json_tpls...)

	yaml_tpls, err := list_dir(MYSQL_EXAMPLE_PATH, ".yaml")
	if err != nil {
		t.Error(err.Error())
	}
	tpls = append(tpls, yaml_tpls...)

	if err := build_cubes_bin(); err != nil {
		t.Error(err.Error())
	}

	for _, v := range tpls {
		test_mysql_example(t, v)
	}
}

func test_mysql_example(t *testing.T, tpl string) {
	test_cubes_cmd(t, fmt.Sprintf("%s/%s", MYSQL_EXAMPLE_PATH, tpl), "", metadata.DEFAULT_MODE)
	test_cubes_cmd(t, fmt.Sprintf("%s/%s", MYSQL_EXAMPLE_PATH, tpl), "", metadata.SQLVIEW_MODE)
}
