package metadata

import (
	"bytes"
	"errors"
	"fmt"
	"gitlab.xibao100.com/skyline/skyline/cubes/utils"
	"sync"
)

func NewReport() *Report {
	return &Report{
		Report: []string{},
		Cubes:  utils.NewMapData(),
	}
}

const MAX_LOOP = 200

func (r *Report) Execute() (map[string]*CubeReport, error) {
	if r.RunMode == SQLVIEW_MODE {
		return r.sqlview_execute()
	}

	return r.default_execute()
}

func (r *Report) default_execute() (map[string]*CubeReport, error) {
	cubes_result := utils.NewMapData()

	cube_tasks := make(map[string]*Cube)
	for k, v := range r.Cubes.Copy() {
		cube, ok := v.(*Cube)
		if !ok {
			err := errors.New("Map data should return cube.")
			logger.Error(err)
			return nil, err
		}
		cubeSha1Name, ok := k.(string)
		if !ok {
			err := errors.New("Map key should return string.")
			logger.Error(err)
			return nil, err
		}
		cube_tasks[cubeSha1Name] = cube
	}

	loop := 0

	for {
		if len(cube_tasks) == 0 {
			break
		}

		logger.Infof("LOOP: %d", loop)
		loop++
		if loop >= MAX_LOOP {
			err := errors.New(fmt.Sprintf("Has invalid cubes: %v", utils.Json(cube_tasks)))
			logger.Error(err)
			return nil, err
		}

		batch_tasks := []*Cube{}
		for _, c := range cube_tasks {
			if v := cubes_result.Get(c.Sha1Name); v != nil {
				delete(cube_tasks, c.Sha1Name)
				continue
			}

			logger.Infof("======================processing cube: %s", c.Name)
			needWait := false
			if c.Source.Type == SOURCE_CUBE {
				if c.Store != nil {
					if v := cubes_result.Get(c.Store.Sha1Name); v == nil {
						needWait = true
					}
				}
				for _, union := range c.Union {
					if v := cubes_result.Get(union.Sha1Name); v == nil {
						needWait = true
					}
				}

				for _, join := range c.Join {
					if v := cubes_result.Get(join.Store.Sha1Name); v == nil {
						needWait = true
					}
				}
			}
			if needWait {
				logger.Infof("======================cube: %s need wait, cube.store=%s, cube.union=%s.", c.Name, utils.Json(c.Store), utils.Json(c.Union))
				continue
			}

			batch_tasks = append(batch_tasks, c)
		}

		var wg sync.WaitGroup
		errorsMap := utils.NewMapData()
		for _, c := range batch_tasks {
			wg.Add(1)
			go func(cube *Cube) {
				defer wg.Done()

				if err := cube.Execute(); err != nil {
					errorsMap.Set(cube.Name, err)
					logger.Error(err)
					return
				}
				report, err := cube.GetReport()
				if err != nil {
					errorsMap.Set(cube.Name, err)
					logger.Error(err)
					return
				}

				cubes_result.Set(cube.Sha1Name, report)
			}(c)
		}
		wg.Wait()
		batch_tasks = []*Cube{}

		if errorsMap.Len() > 0 {
			var buffer bytes.Buffer
			for k, v := range errorsMap.Copy() {
				buffer.WriteString(fmt.Sprintf("%v: %v", k, v))
			}
			err := errors.New(fmt.Sprintf("Errors found in cube execution, errors:%v", buffer.String()))
			logger.Error(err)
			return nil, err
		}
	}

	ret := make(map[string]*CubeReport)
	for _, v := range r.Report {
		if report := cubes_result.Get(Sha1Name(v)); report != nil {
			cube_result, ok := report.(*CubeReport)
			if !ok {
				err := errors.New(fmt.Sprintf("Report cube[%s] not found.", v))
				logger.Error(err)
				return nil, err
			}
			ret[cube_result.Name] = cube_result
		}
	}
	logger.Infof("report ret: %v", utils.Json(ret))
	return ret, nil
}

func (r *Report) sqlview_execute() (map[string]*CubeReport, error) {
	sqlviews := utils.NewMapData()

	sqlview_tasks := make(map[string]*Cube)
	for k, v := range r.Cubes.Copy() {
		cube, ok := v.(*Cube)
		if !ok {
			err := errors.New("Map data should return cube.")
			logger.Error(err)
			return nil, err
		}
		cubeSha1Name, ok := k.(string)
		if !ok {
			err := errors.New("Map key should return string.")
			logger.Error(err)
			return nil, err
		}
		sqlview_tasks[cubeSha1Name] = cube
	}

	loop := 0

	for {
		if len(sqlview_tasks) == 0 {
			break
		}

		logger.Infof("LOOP: %d", loop)
		loop++
		if loop >= MAX_LOOP {
			err := errors.New(fmt.Sprintf("Has invalid cubes: %v", utils.Json(sqlview_tasks)))
			logger.Error(err)
			return nil, err
		}

		for _, c := range sqlview_tasks {
			if v := sqlviews.Get(c.Sha1Name); v != nil {
				delete(sqlview_tasks, c.Sha1Name)
				continue
			}

			logger.Infof("======================processing cube: %s, sourceType = %s", c.Name, c.Source.Type)
			needWait := false
			if c.Source.Type == SOURCE_CUBE {
				if c.Store != nil {
					if v := sqlviews.Get(c.Store.Sha1Name); v == nil {
						logger.Infof("Need wait: name = %s, sha1name = %s", c.Store.Name, c.Store.Sha1Name)
						needWait = true
					}
				}
				for _, union := range c.Union {
					if v := sqlviews.Get(union.Sha1Name); v == nil {
						logger.Infof("Need wait: name = %s, sha1name = %s", union.Name, union.Sha1Name)
						needWait = true
					}
				}
				for _, join := range c.Join {
					if v := sqlviews.Get(join.Store.Sha1Name); v == nil {
						logger.Infof("Need wait: name = %s, sha1name = %s", join.Store.Name, join.Store.Sha1Name)
						needWait = true
					}
				}
			}
			if needWait {
				logger.Infof("======================cube: %s need wait, cube.store=%s, cube.union=%s.", c.Name, utils.Json(c.Store), utils.Json(c.Union))
				continue
			}

			sqlview, err := c.toSQLView(sqlviews)
			if err != nil {
				if err.Error() == ERROR_CUBE_NEED_WAIT {
					logger.Infof("to sql view, return: need wait.")
					needWait = true
				} else {
					logger.Error(err)
					return nil, err
				}
			}

			if needWait {
				logger.Infof("======================cube: %s need wait, cube.store=%s, cube.union=%s.", c.Name, utils.Json(c.Store), utils.Json(c.Union))
				continue
			}

			logger.Infof("sqlview:%s", utils.Json(sqlview))
			logger.Infof("set sql view: name = %s, sha1name=%s", c.Name, c.Sha1Name)
			sqlviews.Set(c.Sha1Name, sqlview)
		}
	}

	var wg sync.WaitGroup
	errorsMap := utils.NewMapData()
	sqlview_result := utils.NewMapData()
	for _, rpt := range r.Report {
		v := sqlviews.Get(Sha1Name(rpt))
		if v == nil {
			err := utils.Errorf("No rpt sqlview for %s.", rpt)
			logger.Error(err)
			return nil, err
		}
		sqlview, ok := v.(*SqlView)
		if !ok {
			err := errors.New(fmt.Sprintf("Sqlview[%v] not found.", v))
			logger.Error(err)
			return nil, err
		}
		wg.Add(1)
		go func(v *SqlView) {
			defer wg.Done()

			cube_result, err := v.Execute()
			if err != nil {
				errorsMap.Set(v.Name, err)
				logger.Error(err)
				return
			}
			sqlview_result.Set(v.Name, cube_result)
		}(sqlview)
	}
	wg.Wait()
	if errorsMap.Len() > 0 {
		var buffer bytes.Buffer
		for k, v := range errorsMap.Copy() {
			buffer.WriteString(fmt.Sprintf("%v: %v", k, v))
		}
		err := errors.New(fmt.Sprintf("Errors found in cube execution, errors:%v", buffer.String()))
		logger.Error(err)
		return nil, err
	}

	ret := make(map[string]*CubeReport)
	for k, v := range sqlview_result.Copy() {
		name, _ := k.(string)
		rpt, _ := v.(*CubeReport)
		ret[name] = rpt
	}
	logger.Infof("report ret: %v", utils.Json(ret))
	return ret, nil
}
