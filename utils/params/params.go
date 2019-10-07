package paramsUtils

import (
	"mas/exception/http_err"
	reflectUtils "mas/utils/reflect"
	"reflect"
	"strings"
)

// 接口定义
type ParamsParser interface {
	Str(key string, desc string, config ...Config) (string, interface{})
	Int(key string, desc string, config ...Config) (int, interface{})
	Float(key string, desc string, config ...Config) (float64, interface{})
	Has(key string) bool
	Bool(key string, desc string, config ...Config) (bool, interface{})
	getRow(key string) reflect.Value
}

// 接口实体定义
type params struct {
	object interface{}
}

// 参数配置
type Config struct {
	Require      bool
	DefaultValue interface{}
}

// 创建params类
func NewParamsParser(object interface{}) ParamsParser {
	return &params{
		object: object,
	}
}

// 获取原始数据
func (this *params) getRow(key string) reflect.Value {

	// 获取值
	v := reflect.ValueOf(this.object)
	return v.FieldByName(key)
}

// 判断数据是否存在
func (this *params) Has(key string) bool {
	return reflectUtils.IsExist(this.getRow(key))
}

// 获取int类型数据
func (this *params) Int(key string, desc string, config ...Config) (int, interface{}) {

	value := this.getRow(key)

	// 判断类型是否正确
	if reflectUtils.IsExist(value) && strings.Contains(value.Type().String(), "int") {
		return int(value.Int()), nil
	}
	if len(config) > 0 && !config[0].Require {
		if v, ok := config[0].DefaultValue.(int); ok {
			return v, nil
		}
	}
	return -1, http_err.LackParams(desc)
}

// 获取float类型数据
func (this *params) Float(key string, desc string, config ...Config) (float64, interface{}) {
	value := this.getRow(key)

	// 判断类型是否正确
	if reflectUtils.IsExist(value) && strings.Contains(value.Type().String(), "float") {
		return value.Float(), nil
	}
	if len(config) > 0 && !config[0].Require {
		if v, ok := config[0].DefaultValue.(float64); ok {
			return float64(v), nil
		}
	}
	return -1, http_err.LackParams(desc)
}

// 获取bool类型数据
func (this *params) Bool(key string, desc string, config ...Config) (bool, interface{}) {
	value := this.getRow(key)

	// 判断类型是否正确
	if reflectUtils.IsExist(value) && strings.Contains(value.Type().String(), "bool") {
		return value.Bool(), nil
	}
	if len(config) > 0 && !config[0].Require {
		if v, ok := config[0].DefaultValue.(bool); ok {
			return v, nil
		}
	}
	return false, http_err.LackParams(desc)
}

// 获取string类型值
func (this *params) Str(key string, desc string, config ...Config) (string, interface{}) {

	value := this.getRow(key)

	// 判断类型是否正确
	if reflectUtils.IsExist(value) && strings.Contains(value.Type().String(), "string") {
		return value.String(), nil
	}
	if len(config) > 0 && !config[0].Require {
		if v, ok := config[0].DefaultValue.(string); ok {
			return v, nil
		}
	}
	return "", http_err.LackParams(desc)
}




