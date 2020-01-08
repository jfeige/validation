package main

import (
	"errors"
	"fmt"
	"strconv"
)

/**
rules = append(rules,map[string]interface{}{"field" : "name", "require" : true, "msg" : "非法请求，数据类型错误！",})
rules = append(rules,map[string]interface{}{"field" : "name", "string" : true, "min" : 6, "max" : 12, "msg" :"姓名长度必须在6~12位之间",})
rules = append(rules,map[string]interface{}{"field" : "age", "require" : true, "msg" : "非法请求，数据类型错误！",})
rules = append(rules,map[string]interface{}{"field" : "age", "int" : true, "min" : 18, "max" : 99, "msg" : "年龄必须在18~99岁之间",})
rules = append(rules,map[string]interface{}{"field" : "desc", "string" : true,"max":100,"msg":"个人说明长度必须在100个字符以内"})
*/
type SingleRule struct {
	field      string
	value      interface{}
	data_value map[string]interface{}
	data_rule  map[string]interface{}
	isRequired bool
	isString   bool
	isInt      bool
	msg        string //该条规则错误信息
	err        error  //规则错误信息
}

func (this *Validate) handleRule(rule map[string]interface{}) {
	defer wg.Done()
	//field 为必传参数
	field, ok := rule["field"].(string)
	if !ok {
		this.errors = append(this.errors, ERR_REQUIRED_PARAMETER_MISSING)
		return
	}
	r := SingleRule{
		field:      field,
		data_value: this.data,
		data_rule:  rule,
	}
	r.is_msg()
	if r.err != nil {
		this.errors = append(this.errors, r.err)
		return
	}
	r.is_requried()
	if r.err != nil {
		this.errors = append(this.errors, r.err)
		return
	}
	r.is_string()
	if r.err != nil {
		this.errors = append(this.errors, r.err)
		return
	}
	r.is_int()
	if r.err != nil {
		this.errors = append(this.errors, r.err)
		return
	}
	r.is_min()
	if r.err != nil {
		this.errors = append(this.errors, r.err)
		return
	}
	r.is_max()
	if r.err != nil {
		this.errors = append(this.errors, r.err)
		return
	}
	return
}

func (this *SingleRule) is_msg() {
	defer func() {
		if err := recover(); err != nil {
			//类型转换错误
			this.err = fmt.Errorf("validation:field %s 's type is not string", this.field)
		}
	}()
	tmpMsg, ok := this.data_rule["msg"]
	if !ok {
		this.err = fmt.Errorf("validation:config msg is null.must be error message!")
		return
	}
	msg := tmpMsg.(string)
	if msg == "" {
		this.err = fmt.Errorf("validation:config msg is null!")
	}
	this.msg = msg
}

/**
required
如果规则里有required且设置为true，检测参数数组里，该字段是否存在
*/
func (this *SingleRule) is_requried() {
	tmpTequired, ok := this.data_rule["required"]
	if !ok {
		return
	}
	required, err := to_bool(tmpTequired)
	if err != nil {
		this.err = fmt.Errorf("validation:field:%s is wrong", this.field)
		return
	}
	if required {
		_, ok := this.data_value[this.field]
		if !ok {
			this.err = errors.New(this.msg)
		}
		this.isRequired = required
	}
}

/**
string
如果规则里有string且设置为true，检测该字段是否为string类型
*/
func (this *SingleRule) is_string() {
	tmpString, ok := this.data_rule["string"]
	if !ok {
		return
	}
	isString, err := to_bool(tmpString)
	if err != nil {
		this.err = fmt.Errorf("validation:config string'value is wrong,must be 'true' or 'false'")
		return
	}
	if isString {
		value, err := to_string(this.data_value[this.field])
		if err == nil {
			this.isString = true
			this.value = value
			return
		}
	}
	this.err = fmt.Errorf("validation:field %s 's type is not string", this.field)
}

/**
int
如果规则里有int且设置为true，检测该字段是否为string类型
*/
func (this *SingleRule) is_int() {
	defer func() {
		if err := recover(); err != nil {
			//类型转换错误
			this.err = fmt.Errorf("validation:field %s 's type is not int", this.field)
		}
	}()
	tmpInt, ok := this.data_rule["int"]
	if !ok {
		return
	}
	isInt, err := to_bool(tmpInt)
	if err != nil {
		this.err = fmt.Errorf("validation:config int'value is wrong,must be 'true' or 'false'")
		return
	}
	if isInt {
		value, err := to_int(this.data_value[this.field])
		if err == nil {
			this.isInt = true
			this.value = value
			return
		}
	}
	this.err = fmt.Errorf("validation:field %s 's type is not int", this.field)
}

/**
min
如果规则定义为string，则判断该字段值长度是否小于指定的min值
如果规则定义为int，则判断该字段值大小是否小于指定的min值
*/
func (this *SingleRule) is_min() {
	tmpMin, ok := this.data_rule["min"]
	if !ok {
		return
	}
	min, err := to_int(tmpMin)
	if err != nil {
		this.err = fmt.Errorf("validation:config min'value is not int type")
		return
	}
	if this.isString {
		//字符串计算长度
		if min < 0 {
			this.err = fmt.Errorf("validation:config min'value must >= 0")
		}
		value := this.value.(string)
		if len(value) < min {
			//业务错误
			this.err = errors.New(this.msg)
		}
	} else if this.isInt {
		//整型计算大小
		value := this.value.(int)
		if value < min {
			//业务错误
			this.err = errors.New(this.msg)
		}
	}
}

/**
max
如果规则定义为string，则判断该字段值长度是否大于指定的min值
如果规则定义为int，则判断该字段值大小是否大于指定的min值
*/
func (this *SingleRule) is_max() {
	tmpMax, ok := this.data_rule["max"]
	if !ok {
		return
	}
	max, err := to_int(tmpMax)
	if err != nil {
		this.err = fmt.Errorf("validation:config max'value is not int type")
	}
	if this.isString {
		//字符串计算长度
		if max < 0 {
			this.err = fmt.Errorf("validation:config min'value must >= 0")
		}
		value := this.value.(string)
		if len(value) > max {
			//业务错误
			this.err = errors.New(this.msg)
		}
	} else if this.isInt {
		//整型计算大小
		value := this.value.(int)
		if value > max {
			//业务错误
			this.err = errors.New(this.msg)
		}
	}
}

/**
interface{}转换bool型
*/
func to_bool(reply interface{}) (bool, error) {
	switch reply := reply.(type) {
	case int64:
		return reply != 0, nil
	case string:
		return strconv.ParseBool(reply)
	case nil:
		return false, ERR_FIELD_IS_NIL
	case Error:
		return false, ERR_BAD_FILED_VALUE
	case bool:
		return reply, nil
	}
	return false, fmt.Errorf("validation: unexpected type for Bool, got type %T", reply)
}

/**
interface{}转换string型
*/
func to_string(reply interface{}) (string, error) {
	switch reply := reply.(type) {
	case string:
		return reply, nil
	case int:
		return strconv.Itoa(reply), nil
	case int64:
		return strconv.FormatInt(reply, 10), nil
	case []byte:
		return string(reply), nil
	}
	return "", fmt.Errorf("validation: unexpected type for String, got type %T", reply)
}

func to_int(reply interface{}) (int, error) {
	switch reply := reply.(type) {
	case string:
		return strconv.Atoi(reply)
	case int:
		return reply, nil
	case int64:
		return int(reply), nil
	case []byte:
		return strconv.Atoi(string(reply))
	}
	return 0, fmt.Errorf("validation: unexpected type for Integer, got type %T", reply)
}
