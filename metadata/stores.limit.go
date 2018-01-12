package metadata

import (
	"encoding/json"
	"github.com/go-yaml/yaml"
	"gitlab.xibao100.com/skyline/skyline/cubes/utils"
	"sync"
)

type LimitStore struct {
	SourceType string   `json:"source_type,omitempty", yaml:"source_type,omitempty`
	Fields     []string `json:"fields,omitempty", yaml:"fields,omitempty`
}

type StoresLimit struct {
	LimitStores   map[string]*LimitStore `json:"limit_stores,omitempty", yaml:"limit_stores,omitempty`
	FieldsSetting map[string]interface{} `json:"fields_setting,omitempty", yaml:"fields_setting,omitempty`
	m             *sync.Mutex
}

func NewStoresLimit() *StoresLimit {
	return &StoresLimit{
		LimitStores:   make(map[string]*LimitStore),
		FieldsSetting: make(map[string]interface{}),
		m:             new(sync.Mutex),
	}
}

func NewStoresLimitFromStr(storesLimitStr string, format string) (*StoresLimit, error) {
	storesLimitTmp := NewStoresLimit()
	switch format {
	case TPL_JSON:
		if err := json.Unmarshal([]byte(storesLimitStr), storesLimitTmp); err != nil {
			logger.Errorf("ERROR Unmarshal: %v", err.Error())
			return nil, err
		}
	case TPL_YAML:
		if err := yaml.Unmarshal([]byte(storesLimitStr), storesLimitTmp); err != nil {
			logger.Errorf("ERROR Unmarshal: %v", err.Error())
			return nil, err
		}
	default:
		return nil, utils.Errorf("Unknown format:%s", format)
	}

	// KEY值统一大小写
	storesLimit := NewStoresLimit()
	for storeName, store := range storesLimitTmp.LimitStores {
		limitStore := &LimitStore{
			SourceType: utils.LowerTrim(store.SourceType),
			Fields:     []string{},
		}
		for _, v := range store.Fields {
			limitStore.Fields = append(limitStore.Fields, utils.LowerTrim(v))
		}
		storesLimit.LimitStores[utils.LowerTrim(storeName)] = limitStore
	}

	for k, v := range storesLimitTmp.FieldsSetting {
		storesLimit.FieldsSetting[utils.LowerTrim(k)] = v
	}

	logger.Infof("stores limit:%v", utils.Json(storesLimit))
	return storesLimit, nil
}

func (l *StoresLimit) GetLimitStore(storeName string) *LimitStore {
	storeName = utils.LowerTrim(storeName)
	if storeName == "" {
		return nil
	}

	l.m.Lock()
	v, ok := l.LimitStores[storeName]
	l.m.Unlock()
	if ok {
		return v
	}
	return nil
}

func (l *StoresLimit) SetFieldSetting(field string, val interface{}) {
	field = utils.LowerTrim(field)
	if field == "" {
		return
	}

	l.m.Lock()
	l.FieldsSetting[field] = val
	l.m.Unlock()
}

func (l *StoresLimit) GetFieldSetting(field string) interface{} {
	field = utils.LowerTrim(field)
	if field == "" {
		return nil
	}

	l.m.Lock()
	v, ok := l.FieldsSetting[field]
	l.m.Unlock()
	if ok {
		return v
	}
	return nil
}
