package mysqlib

import (
	"reflect"
	"strings"
)

//解析模型结构
func (sess *Session) parseModel() {
	//将模型转为valueOf类型以便反射得到相关信息
	model := reflect.ValueOf(sess.modelValue.Value)

	//将转为reflect.Value类型的值存入session
	sess.modelValue.rValue = model.Elem()

	//得到模型的变量类型
	modelKind := sess.modelValue.rValue.Kind()

	//模型名称
	var modelName string

	//如果传进来的是结构体
	if modelKind == reflect.Slice {
		//保存模型的reflect.Type类型到session
		sess.modelValue.rType = sess.modelValue.rValue.Type().Elem()
		//标记是结构体
		sess.modelValue.isSlice = true
		//模型名
		modelName = sess.modelValue.rValue.Type().Elem().String()
	} else if modelKind == reflect.Struct {
		//保存模型的reflect.Type类型到session
		sess.modelValue.rType = model.Type().Elem()
		//模型名
		modelName = sess.modelValue.rValue.Type().String()
	}

	//如果没有禁用模型缓存
	if sess.builder.options.DisableModelCache == false {
		//从缓存中读取模型信息
		sess.modelInfo = sess.builder.modelCache[modelName]
	}
	//如果缓存中没有读到模型信息
	if sess.modelInfo == nil {
		//再从反射中获取模型信息
		sess.modelInfo = sess.reflectModel(modelName)
		//把反射出来的模型信息写入到缓存
		if sess.builder.options.DisableModelCache == false {
			sess.builder.modelCache[sess.modelInfo.name] = sess.modelInfo
		}
	}

	//会话的初始表名是从模型里的tag里读取的
	//如果没有取到，后面可以再用Table()方法指定
	sess.tableName = sess.modelInfo.tableName
}

//反射结构体得到模型信息
func (sess *Session) reflectModel(modelName string) *modelInfo {
	//创建一个模型信息
	var info modelInfo
	info.fields = make(map[string]*modelField)
	info.name = modelName

	//取得结构体所有字段的总数
	allFieldCount := sess.modelValue.rType.NumField()

	//如果字段总数>0才开始循环
	if allFieldCount > 0 {
		//遍历所有字段
		for i := 0; i < allFieldCount; i++ {
			var field modelField
			field.VarName = sess.modelValue.rType.Field(i).Name
			field.VarType = sess.modelValue.rType.Field(i).Type.Name()
			field.SQLName = sess.modelValue.rType.Field(i).Tag.Get(sess.builder.options.TagName)
			//如果存在标记
			if field.SQLName != "" {
				//如果是标记表名的字段
				if field.VarName == sess.builder.options.TableNameField {
					//赋值表名
					info.tableName = field.SQLName
				} else {
					//累加模型信息中的字段总数
					info.fieldCount++
					//写入字段信息
					info.fields[field.SQLName] = &field
				}
			}
		}
	}

	return &info
}
