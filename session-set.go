package mysqlib

import (
	"errors"
	"strings"
)

// Table 设置本次会话的表名，优先级于结构体中定义的表名，但仅影响本次会话
func (sess *Session) Table(tableName string) *Session {
	//设置临时表名
	sess.tableName = tableName
	return sess
}

// Column 要影响的字段，作用于以下操作：
// INSERT：只插入哪些字段
// UPDATE：只更新哪些字段
// SELECT：只返回哪些字段
func (sess *Session) Column(fields ...string) *Session {
	for _, v := range fields {
		var field keyInterface
		field.key = v
		sess.stmt.field = append(sess.stmt.field, &field)
	}
	return sess
}

// ColumnRaw 用原生SQL语句设置要影响的字段，用于对Column()方法的补充
// 主要用于影响模型没有定义的额外字段及其值，或者执行数据库函数
// 如果不传入value参数，会被当作MySQL的函数来看待，不会为字段名加上引号
func (sess *Session) ColumnRaw(fieldName string, value ...interface{}) *Session {
	var field keyInterface
	field.key = fieldName
	if len(value) > 0 {
		field.value = value[0]
	}
	sess.stmt.extraField = append(sess.stmt.extraField, &field)
	return sess
}

// Where 设置AND WHERE条件，作用跟AndWhere()一样
func (sess *Session) Where(field, operator string, value interface{}) *Session {
	return sess.whereHandle("AND", field, operator, value)
}

// AndWhere 设置AND WHERE条件
func (sess *Session) AndWhere(field, operator string, value interface{}) *Session {
	return sess.whereHandle("AND", field, operator, value)
}

// OrWhere 设置OR WHERE
func (sess *Session) OrWhere(field, operator string, value interface{}) *Session {
	return sess.whereHandle("OR", field, operator, value)
}

// WhereRaw 传入原生WHERE语句
// 传入的字符串会直接拼接，务必注意安全
func (sess *Session) WhereRaw(stmt string) *Session {
	return sess.whereHandle("AND", stmt, "[!RAW!]", nil)
}

// AndWhereRaw 传入原生AND WHERE语句
// 传入的字符串会直接拼接，务必注意安全
func (sess *Session) AndWhereRaw(stmt string) *Session {
	return sess.whereHandle("AND", stmt, "[!RAW!]", nil)
}

// OrWhereRaw 传入原生OR WHERE语句
// 传入的字符串会直接拼接，务必注意安全
func (sess *Session) OrWhereRaw(stmt string) *Session {
	return sess.whereHandle("OR", stmt, "[!RAW!]", nil)
}

// WhereIn 传入WHERE IN语句，作用和AndWhereIn()一样
func (sess *Session) WhereIn(field string, value interface{}) *Session {
	return sess.whereHandle("AND", field, "IN", value)
}

// AndWhereIn 传入AND IN语句
func (sess *Session) AndWhereIn(field string, value interface{}) *Session {
	return sess.whereHandle("AND", field, "IN", value)
}

// OrWhereIn 传入OR IN语句
func (sess *Session) OrWhereIn(field string, value interface{}) *Session {
	return sess.whereHandle("OR", field, "IN", value)
}

// WhereNotIn 传入AND NOT IN语句，作用和AndWhereNotIn()一样
func (sess *Session) WhereNotIn(field string, value interface{}) *Session {
	return sess.whereHandle("AND", field, "NOT IN", value)
}

// AndWhereNotIn 传入AND NOT IN语句
func (sess *Session) AndWhereNotIn(field string, value interface{}) *Session {
	return sess.whereHandle("AND", field, "NOT IN", value)
}

// OrWhereNotIn 传入OR NOT IN语句
func (sess *Session) OrWhereNotIn(field string, value interface{}) *Session {
	return sess.whereHandle("OR", field, "NOT IN", value)
}

func (sess *Session) whereHandle(union, field, operator string, value interface{}) *Session {
	var cond whereCond
	cond.union = union
	cond.field = field
	cond.operator = operator
	cond.value = value
	sess.stmt.where = append(sess.stmt.where, &cond)
	return sess
}

//OrderBy 排序语句
func (sess *Session) OrderBy(field, direction string) *Session {
	direction = strings.ToUpper(direction)
	if direction != "ASC" && direction != "DESC" {
		sess.err = errors.New("排序规则只能是`ASC`或者`DESC`")
		return sess
	}
	var order orderBy
	order.field = field
	order.direction = direction
	sess.stmt.orders = append(sess.stmt.orders, &order)
	return sess
}

//Limit 限制返回记录条数
func (sess *Session) Limit(value int) *Session {
	sess.stmt.limit = value
	return sess
}

//Offset 游标偏移数量
func (sess *Session) Offset(value int) *Session {
	sess.stmt.offset = value
	return sess
}
