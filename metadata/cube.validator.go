package metadata

import (
	"gitlab.xibao100.com/skyline/skyline/cubes/utils"
)

func CheckCube(cube *Cube) (err error) {
	if err = check_cube_name(cube); err != nil {
		logger.Error(err)
		return
	}
	if err = check_cube_source(cube); err != nil {
		logger.Error(err)
		return
	}
	if err = check_cube_store(cube); err != nil {
		logger.Error(err)
		return
	}
	if err = check_cube_join(cube); err != nil {
		logger.Error(err)
		return
	}
	if err = check_cube_filter(cube); err != nil {
		logger.Error(err)
		return
	}
	if err = check_cube_mapping(cube); err != nil {
		logger.Error(err)
		return
	}
	if err = check_cube_dimension(cube); err != nil {
		logger.Error(err)
		return
	}
	if err = check_cube_orderby(cube); err != nil {
		logger.Error(err)
		return
	}
	if err = check_cube_limit(cube); err != nil {
		logger.Error(err)
		return
	}
	if err = check_cube_aggregate(cube); err != nil {
		logger.Error(err)
		return
	}
	if err = check_cube_tag(cube); err != nil {
		logger.Error(err)
		return
	}
	if err = check_cube_union(cube); err != nil {
		logger.Error(err)
		return
	}
	if err = check_cube_summary(cube); err != nil {
		logger.Error(err)
		return
	}
	return
}
func check_cube_name(cube *Cube) error {
	if cube.Name == "" || cube.Sha1Name == "" {
		return utils.Error("Empty cube name")
	}
	return nil
}

func check_cube_source(cube *Cube) error {
	if cube.Source == nil {
		return utils.Errorf("No source for cube: %s", cube.Name)
	}
	if !utils.InArrayStr(cube.Source.Type, []string{
		SOURCE_CUBE, SOURCE_MYSQL, SOURCE_SQLITE, SOURCE_JSON, SOURCE_CSV,
	}) {
		return utils.Errorf("Invalid source type[%s] for cube: %s", cube.Source.Type, cube.Name)
	}
	if cube.Source.Name == "" && utils.InArrayStr(cube.Source.Type, []string{
		SOURCE_SQLITE, SOURCE_JSON, SOURCE_CSV,
	}) {
		return utils.Errorf("Empty source name for cube: %s, source type:%s", cube.Name, cube.Source.Type)
	}
	return nil
}

func check_cube_store(cube *Cube) error {
	if cube.Source.Type == SOURCE_CUBE && cube.Union != nil && len(cube.Union) > 0 {
		return nil
	}

	if cube.Sql != "" && cube.Source.Type != SOURCE_CUBE {
		return nil
	}

	if cube.Store == nil {
		return utils.Errorf("No store for cube: %s", cube.Name)
	}
	if cube.Store.Name == "" || cube.Store.Sha1Name == "" {
		return utils.Error("Empty cube store name")
	}

	if cube.StoresLimit != nil && cube.Source.Type == SOURCE_MYSQL {
		if v := cube.StoresLimit.GetLimitStore(cube.Store.Name); v == nil {
			return utils.Errorf("Not allowed store: %s", cube.Store.Name)
		}
	}

	return nil
}

func check_cube_join(cube *Cube) error {
	if cube.Join != nil && len(cube.Join) > 0 {
		for _, join := range cube.Join {
			if !utils.InArrayStr(join.Type, []string{
				INNER_JOIN, LEFT_JOIN, RIGHT_JOIN,
			}) {
				return utils.Errorf("Invalid join type[%s] for cube: %s", join.Type, cube.Name)
			}

			if join.Store == nil {
				return utils.Errorf("No store for cube: %s join.", cube.Name)
			}
			if join.Store.Name == "" || join.Store.Sha1Name == "" {
				return utils.Errorf("Empty store name for cube: %s join.", cube.Name)
			}
			if cube.StoresLimit != nil && cube.Source.Type == SOURCE_MYSQL {
				if v := cube.StoresLimit.GetLimitStore(join.Store.Name); v == nil {
					return utils.Errorf("Not allowed store: %s", join.Store.Name)
				}
			}
			if len(join.On) == 0 {
				return utils.Errorf("No on conditions for cube: %s join.", cube.Name)
			}
			for _, cond := range join.On {
				if cond.Field == "" {
					return utils.Errorf("Empty on condition field for cube: %s join.", cube.Name)
				}
				if cond.Op == "" {
					return utils.Errorf("Empty on condition op for cube: %s join.", cube.Name)
				}
			}
			return nil
		}
	}
	return nil
}

func check_cube_filter(cube *Cube) error {
	if cube.Filter != nil && len(cube.Filter) > 0 {
		for _, cond := range cube.Filter {
			if cond.Field == "" {
				return utils.Errorf("Empty condition field for cube: %s filter.", cube.Name)
			}
			if cond.Op == "" {
				return utils.Errorf("Empty condition op for cube: %s filter.", cube.Name)
			}
		}
	}
	return nil
}

func check_cube_mapping(cube *Cube) error {
	if cube.Mappings != nil && len(cube.Mappings) > 0 {
		for _, mapping := range cube.Mappings {
			if mapping.Alias == "" {
				return utils.Errorf("Empty mapping alias for cube: %s mappings.", cube.Name)
			}
			if mapping.Expression == "" {
				return utils.Errorf("Empty mapping expression for cube: %s mappings.", cube.Name)
			}
		}
	}
	return nil
}

func check_cube_dimension(cube *Cube) error {
	if cube.Dimensions != nil && len(cube.Dimensions) > 0 {
		for _, v := range cube.Dimensions {
			if v == "" {
				return utils.Errorf("Has empty dimension for cube: %s.", cube.Name)
			}
		}
	}
	return nil
}

func check_cube_orderby(cube *Cube) error {
	if cube.OrderBy != nil && len(cube.OrderBy) > 0 {
		for _, orderby := range cube.OrderBy {
			if orderby.Field == "" {
				return utils.Errorf("Empty field for cube: %s order by.", cube.Name)
			}

			if !utils.InArrayStr(orderby.Order, []string{
				ORDER_ASC, ORDER_DESC,
			}) {
				return utils.Errorf("Invalid order[%s] for cube: %s order by.", orderby.Order, cube.Name)
			}
		}
	}
	return nil
}

func check_cube_limit(cube *Cube) error {
	if cube.Limit != nil {
		if cube.Limit.Limit == 0 {
			return utils.Errorf("Invalid limit[%d] for cube: %s", cube.Limit.Limit, cube.Name)
		}
	}
	return nil
}

func check_cube_aggregate(cube *Cube) error {
	if cube.Aggregates != nil && len(cube.Aggregates) > 0 {
		for f, mappings := range cube.Aggregates {
			if !utils.InArrayStr(f, []string{
				SUM, AVG, MIN, MAX, COUNT,
			}) {
				return utils.Errorf("Invalid aggregate function[%s] for cube: %s aggregate.", f, cube.Name)
			}
			if len(mappings) == 0 {
				return utils.Errorf("No aggregate fields for cube: %s aggregate, function:%s", cube.Name, f)
			}
			for _, mapping := range mappings {
				if mapping.Expression == "" {
					return utils.Errorf("Empty expression for cube: %s aggregate, function:%s", cube.Name, f)
				}
			}
		}
	}
	return nil
}

func check_cube_tag(cube *Cube) error {
	if cube.Tags != nil && len(cube.Tags) > 0 {
		for tag_name, mappings := range cube.Tags {
			if tag_name == "" {
				return utils.Errorf("Empty tag name for cube: %s tags", cube.Name)
			}
			if len(mappings) == 0 {
				return utils.Errorf("No tag mapping for cube: %s tags, tag name:%s", cube.Name, tag_name)
			}
			for _, mapping := range mappings {
				if mapping.TagVal == "" {
					return utils.Errorf("Empty tag value for cube: %s tags, tag name:%s", cube.Name, tag_name)
				}
				if mapping.Field == "" {
					return utils.Errorf("Empty mappinng field for cube: %s tags, tag name:%s", cube.Name, tag_name)
				}
				if mapping.Op != "REGEXP" {
					return utils.Errorf("Only REGEXP supported for cube: %s tags, tag name:%s", cube.Name, tag_name)
				}
				if mapping.Val == "" {
					return utils.Errorf("Empty include regexp for cube: %s tags, tag name:%s", cube.Name, tag_name)
				}
			}
		}
	}
	return nil
}

func check_cube_union(cube *Cube) error {
	if cube.Union != nil && len(cube.Union) > 0 {
		for _, union := range cube.Union {
			if union.Name == "" || union.Sha1Name == "" {
				return utils.Errorf("Empty union name for cube: %s union.", cube.Name)
			}
			if !utils.InArrayStr(union.UnionType, []string{
				UNION, UNION_ALL,
			}) {
				return utils.Errorf("Invalid union type[%s] for cube: %s union.", union.UnionType, cube.Name)
			}
		}
	}
	return nil
}

func check_cube_summary(cube *Cube) error {
	for _, v := range cube.Summary {
		if !utils.InArrayStr(v.Type, []string{
			SUM, CONTRAST,
		}) {
			return utils.Errorf("Invalid summary type[%s] for cube: %s summary.", v.Type, cube.Name)
		}
		if v.Alias == "" {
			return utils.Errorf("Empty alias for cube: %s summary.", cube.Name)
		}
	}
	return nil
}
