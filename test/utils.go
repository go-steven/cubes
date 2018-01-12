package test

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.xibao100.com/skyline/skyline/cubes/metadata"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func run_shell_cmd(cmd string) (string, error) {
	fmt.Println("Run shell cmd: ", cmd)
	f, err := exec.Command("/bin/sh", "-c", cmd).Output()
	if err != nil {
		return "", err
	}
	return string(f), nil
}

func build_cubes_bin() error {
	ret, err := run_shell_cmd(fmt.Sprintf("cd %s && go build && cd %s", CUBES_BIN_DIR, TEST_PATH))
	if err != nil {
		return errors.New(fmt.Sprintf("build cubes bin failed: %v", err))
	}

	if len(ret) != 0 {
		return errors.New(fmt.Sprintf("build cubes bin failed: %v", ret))
	}

	return nil
}

func output_file_name() string {
	return fmt.Sprintf("/tmp/test%d.output", time.Now().UnixNano())
}

func read_report_from_output(output_file string) (map[string]*metadata.CubeReport, error) {
	bytes, err := ioutil.ReadFile(output_file)
	if err != nil {
		return nil, err
	}

	report := make(map[string]*metadata.CubeReport)
	if err := json.Unmarshal(bytes, &report); err != nil {
		return nil, errors.New(fmt.Sprintf("ERROR Unmarshal: %v", err.Error()))
	}

	//println("xxxxx: ", utils.Json(report))
	return report, nil
}

func tail_lines(s string, num int) (ret []string) {
	lines := strings.Split(s, "\n")
	total_lines := len(lines)
	if total_lines <= num {
		return lines
	}

	for i := total_lines - num; i < total_lines; i++ {
		ret = append(ret, lines[i])
	}
	return
}
func test_cubes_cmd(t *testing.T, tpl, tplcfg, sqlmode string) {
	outputFile := output_file_name()
	ret, err := run_shell_cmd(fmt.Sprintf("%s --tpl=%s --tplcfg=%s --sqlmode=%s --output=%s", CUBES_BIN, tpl, tplcfg, sqlmode, outputFile))
	if err != nil {
		t.Error(err.Error())
	}
	report, err := read_report_from_output(outputFile)
	if err != nil {
		t.Error(strings.Join(tail_lines(ret, 5), "\n"))
		t.Error(err.Error())
	}
	if report == nil {
		t.Error(strings.Join(tail_lines(ret, 5), "\n"))
		t.Errorf("No report for tpl: %s.", tpl)
	}
	os.Remove(outputFile)
}

func list_dir(dirPth string, suffix string) (files []string, err error) {
	files = make([]string, 0, 10)
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}
	suffix = strings.ToLower(suffix) //忽略后缀匹配的大小写
	for _, fi := range dir {
		if fi.IsDir() { // 忽略目录
			continue
		}
		if strings.HasSuffix(strings.ToLower(fi.Name()), suffix) { //匹配文件
			files = append(files, fi.Name())
		}
	}
	return files, nil
}
