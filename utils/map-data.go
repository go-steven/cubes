package utils

import (
	"sync"
)

type MapData struct {
	data map[interface{}]interface{}
	m    *sync.RWMutex
}

func NewMapData() *MapData {
	return &MapData{
		data: make(map[interface{}]interface{}),
		m:    new(sync.RWMutex),
	}
}

func (this *MapData) Get(k interface{}) interface{} {
	this.m.RLock()
	v, ok := this.data[k]
	this.m.RUnlock()
	if !ok {
		return nil
	}

	return v
}

func (this *MapData) Set(k, v interface{}) {
	this.m.Lock()
	this.data[k] = v
	this.m.Unlock()
}

func (this *MapData) Del(k interface{}) {
	this.m.Lock()
	if _, ok := this.data[k]; ok {
		delete(this.data, k)
	}
	this.m.Unlock()
}

func (this *MapData) Keys() []interface{} {
	result := []interface{}{}

	this.m.RLock()
	for k, _ := range this.data {
		result = append(result, k)
	}
	this.m.RUnlock()

	return result
}

func (this *MapData) Values() []interface{} {
	result := []interface{}{}

	this.m.RLock()
	for _, v := range this.data {
		result = append(result, v)
	}
	this.m.RUnlock()

	return result
}

func (this *MapData) Copy() map[interface{}]interface{} {
	result := make(map[interface{}]interface{})

	this.m.RLock()
	for k, v := range this.data {
		result[k] = v
	}
	this.m.RUnlock()

	return result
}

func (this *MapData) Len() (ret int) {
	this.m.RLock()
	ret = len(this.data)
	this.m.RUnlock()

	return
}
