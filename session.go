package mysqlib

import (
	"reflect"
)

// Session 会话实例的结构
type Session struct {
	builder    *Instance  //构建器内存地址
	modelInfo  *modelInfo //模型信息
	modelValue struct {   //模型的值
		Value   interface{}   //模型的实例
		isSlice bool          //模型的实例是Slice
		rValue  reflect.Value //模型实例的reflectValue
		rType   reflect.Type  //模型实例的reflectType
	}
	tableName string //临时作用于本次会话的表名
	//sql语句的结构
	stmt struct {
		action       string          //行为
		field        []*keyInterface //INSERT/UPDATE要从模型中取值的字段
		addValue     []*keyInterface //INSERT/UPDATE要写值到模型外的字段及其值
		where        []*whereCond    //where条件
		orders       []*orderBy
		limit        int
		offset       int
		resultString string        //最终生成的sql语句字符串
		resultValues []interface{} //最终汇总的参数值
	}
	err error //错误
}

type keyInterface struct {
	key   string
	value interface{}
}

//模型信息
type modelInfo struct {
	name       string                 //模型名称（用作缓存map中的key）
	tableName  string                 //sql表名
	fieldCount int                    //sql字段数（仅含标记信息的字段)
	fields     map[string]*modelField //sql字段信息key是sql字段名
}

//where条件结构
type whereCond struct {
	union    string      //连接符AND/OR
	field    string      //字段
	operator string      //运算符
	value    interface{} //字段值
}

//模型里的字段信息
type modelField struct {
	VarName string //模型变量名
	VarType string //模型变量类型
	SQLName string //数据表字段名
	//SQLType string //数据表字段类型
}

//排序规则
type orderBy struct {
	field     string
	direction string
}
