# validation


基于go实现的一个参数验证器，当前仅实现了int，string，是否为空，最大值，最小值判断几个最基础的

调用方式:

	rules := make([]map[string]interface{},0)
	rules = append(rules,map[string]interface{}{"field" : "name", "require" : true, "msg" : "姓名不能为空！",})
	rules = append(rules,map[string]interface{}{"field" : "name", "string" : true, "min" : 6, "max" : 12, "msg" :"姓名长度必须在6~12位之间",})
	rules = append(rules,map[string]interface{}{"field" : "age", "require" : true, "msg" : "年龄不能为空",})
	rules = append(rules,map[string]interface{}{"field" : "age", "int" : true, "min" : 18, "max" : 99, "msg" : "年龄必须在18~99岁之间",})
	rules = append(rules,map[string]interface{}{"field" : "desc", "string" : true,"max":100,"msg":"个人说明长度必须在100个字符以内"})

	data := make(map[string]interface{})
	data["name"] = "李三aaaaaaaaaa"
	data["age"] = 0
	data["desc"] = "过新年，穿新衣，晚上去打雪仗"


	validation := NewValidation(data,rules)

	validation.Check()
	if validation.IsFail(){
		fmt.Println("第一个错误:" + validation.FirstError().Error())
		fmt.Println("所有错误:" , validation.AllErrors())
	}else{
		fmt.Println("验证通过")
	}
