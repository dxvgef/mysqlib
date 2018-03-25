package mysqlib

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
)

type keyString struct {
	key   string
	value string
}

// Build 开始构建语句，并赋值会话实例及错误消息
// 入参值为true时，构建可直接执行含有参数值的SQL语句
// 入参值为false时，构建含有?占位符的SQL语句，占位符对应的值通过GetValues()方法获得
// 注意：一个会话不要切换两种模式来构建
func (sess *Session) Build(final bool) (*Session, error) {
	if sess.err != nil {
		return nil, sess.err
	}

	var stmt bytes.Buffer

	//解析模型结构
	sess.parseModel()

	//如果没有定义表名
	if sess.tableName == "" {
		return nil, errors.New("没有定义表名")
	}

	//根据行为调用不同的解析方法
	switch sess.stmt.action {
	case "INSERT":
		stmt.WriteString(sess.buildInsert(final))
	case "UPDATE":
		stmt.WriteString(sess.buildUpdate(final))
		//拼接where语句
		stmt.WriteString(sess.buildWhere(final))
		//拼接order by语句
		stmt.WriteString(sess.buildOrderBy())
		//拼接limit语句
		stmt.WriteString(sess.buildLimit())
	case "SELECT":
		stmt.WriteString(sess.buildSelect(final))
		//拼接where语句
		stmt.WriteString(sess.buildWhere(final))
		//拼接order by语句
		stmt.WriteString(sess.buildOrderBy())
		//拼接limit语句
		stmt.WriteString(sess.buildLimit())
		//拼接offset语句
		stmt.WriteString(sess.buildOffset())
	case "DELETE":
		stmt.WriteString("DELETE FROM `")
		stmt.WriteString(sess.tableName)
		stmt.WriteString("`")
		//拼接where语句
		stmt.WriteString(sess.buildWhere(final))
		//拼接order by语句
		stmt.WriteString(sess.buildOrderBy())
		//拼接limit语句
		stmt.WriteString(sess.buildLimit())
	default:
		return nil, errors.New("未知的行为")
	}

	sess.stmt.resultString = stmt.String()

	return sess, nil
}

//拼接INSERT语句
func (sess *Session) buildInsert(final bool) string {
	var stmt bytes.Buffer
	stmt.WriteString("INSERT INTO `")
	stmt.WriteString(sess.tableName)
	stmt.WriteString("` (")

	var allField []keyString

	// ------------------ 拼接column部分 ----------------------------
	// 如果没有用Column()指定field
	var field keyString
	if len(sess.stmt.field) == 0 {
		//把模型里所有的字段及其值写入到allField里
		for _, v := range sess.modelInfo.fields {
			value := sess.modelValue.rValue.FieldByName(v.VarName).Interface()
			if final == false {
				sess.stmt.resultValues = append(sess.stmt.resultValues, value)
			}
			field.key = v.SQLName
			field.value = interfaceToString(value)
			//sess.stmt.field = append(sess.stmt.field, &field)
			//将k/v插入到allField
			allField = append(allField, field)
		}
	} else {
		//如果用Column()指定了field
		for _, sqlName := range sess.stmt.field {
			value := sess.modelValue.rValue.FieldByName(sess.modelInfo.fields[sqlName.key].VarName).Interface()
			if final == false {
				sess.stmt.resultValues = append(sess.stmt.resultValues, value)
			}
			field.key = sqlName.key
			field.value = interfaceToString(value)
			allField = append(allField, field)
		}
	}

	//把额外添加的字段也汇总到allField
	for _, v := range sess.stmt.extraField {
		if final == false {
			sess.stmt.resultValues = append(sess.stmt.resultValues, v.value)
		}
		field.key = v.key
		field.value = interfaceToString(v.value)
		allField = append(allField, field)
	}

	// 拼接column部分
	for k, v := range allField {

		// 拼接语句
		if k == 0 {
			stmt.WriteString("`")
			stmt.WriteString(v.key)
			stmt.WriteString("`")
		} else {
			stmt.WriteString(", `")
			stmt.WriteString(v.key)
			stmt.WriteString("`")
		}
	}
	//VALUES前面的拼接完成
	stmt.WriteString(") VALUES (")

	//遍历开始拼接参数值
	if final == true {
		for k, v := range allField {
			if k == 0 {
				stmt.WriteString(v.value)
			} else {
				stmt.WriteString(", ")
				stmt.WriteString(v.value)
			}
		}
	} else {
		fieldCount := len(allField)
		for i := 0; i < fieldCount; i++ {
			if i == 0 {
				stmt.WriteString("?")
			} else {
				stmt.WriteString(", ?")
			}
		}
	}

	stmt.WriteString(");")

	return stmt.String()
}

//拼接UPDATE语句
func (sess *Session) buildUpdate(final bool) string {
	var stmt bytes.Buffer
	stmt.WriteString("UPDATE `")
	stmt.WriteString(sess.tableName)
	stmt.WriteString("` SET ")

	var allField []keyString

	// ------------------ 拼接SET部分 ----------------------------
	var field keyString
	//遍历Column
	for _, sqlName := range sess.stmt.field {
		value := sess.modelValue.rValue.FieldByName(sess.modelInfo.fields[sqlName.key].VarName).Interface()
		if final == false {
			sess.stmt.resultValues = append(sess.stmt.resultValues, value)
		}
		field.key = sqlName.key
		field.value = interfaceToString(value)
		allField = append(allField, field)
	}

	//把额外添加的字段也汇总到allField
	for _, v := range sess.stmt.extraField {
		if final == false {
			sess.stmt.resultValues = append(sess.stmt.resultValues, v.value)
		}
		field.key = v.key
		field.value = interfaceToString(v.value)
		allField = append(allField, field)
	}
	// 拼接set语句
	for k, v := range allField {
		if final == true {
			if k == 0 {
				stmt.WriteString("`")
				stmt.WriteString(v.key)
				stmt.WriteString("`=")
				stmt.WriteString(v.value)
			} else {
				stmt.WriteString(", `")
				stmt.WriteString(v.key)
				stmt.WriteString("`=")
				stmt.WriteString(v.value)
			}
		} else {
			if k == 0 {
				stmt.WriteString("`")
				stmt.WriteString(v.key)
				stmt.WriteString("`=?")
			} else {
				stmt.WriteString(", `")
				stmt.WriteString(v.key)
				stmt.WriteString("`=?")
			}
		}
	}

	return stmt.String()
}

//拼接SELECT语句
func (sess *Session) buildSelect(final bool) string {
	var stmt bytes.Buffer
	stmt.WriteString("SELECT ")

	var allField []keyString

	// ------------------ 拼接column部分 ----------------------------
	// 如果没有用Column()指定field也没有用ColumnRaw()指定过字段
	var field keyString
	if len(sess.stmt.field) == 0 && len(sess.stmt.extraField) == 0 {
		//把模型里所有的字段及其值写入到allField里
		for _, v := range sess.modelInfo.fields {
			field.key = v.SQLName
			if sess.modelValue.isSlice == false {
				field.value = interfaceToString(sess.modelValue.rValue.FieldByName(v.VarName).Interface())
			}
			sess.stmt.field = append(sess.stmt.field, &keyInterface{
				key:   field.key,
				value: field.value,
			})
			//将k/v插入到allField
			allField = append(allField, field)
		}
	} else {
		//如果用Column()指定了field
		for _, sqlName := range sess.stmt.field {
			field.key = sqlName.key
			if sess.modelValue.isSlice == false {
				field.value = interfaceToString(sess.modelValue.rValue.FieldByName(sess.modelInfo.fields[sqlName.key].VarName).Interface())
			}
			allField = append(allField, field)
		}
	}

	//把额外添加的字段也汇总到allField
	for _, v := range sess.stmt.extraField {
		field.key = v.key
		field.value = interfaceToString(v.value)
		allField = append(allField, field)
		sess.stmt.field = append(sess.stmt.field, &keyInterface{
			key:   field.key,
			value: field.value,
		})
	}

	// 拼接column部分
	for k, v := range allField {
		if k == 0 {
			if v.value == "" {
				stmt.WriteString(v.key)
			} else {
				stmt.WriteString(v.key)
				stmt.WriteString("`")
			}
		} else {
			if v.value == "" {
				stmt.WriteString(", ")
				stmt.WriteString(v.key)
			} else {
				stmt.WriteString(", `")
				stmt.WriteString(v.key)
				stmt.WriteString("`")
			}
		}
	}

	//VALUES前面的拼接完成
	stmt.WriteString(" FROM `")
	stmt.WriteString(sess.tableName)
	stmt.WriteString("`")

	return stmt.String()
}

//构建where语句
func (sess *Session) buildWhere(final bool) string {
	whereCount := len(sess.stmt.where)
	if whereCount == 0 {
		return ""
	}
	var stmt bytes.Buffer
	stmt.WriteString(" WHERE ")

	//遍历where条件
	var value string
	for i := 0; i < whereCount; i++ {
		value = ""
		//如果是IN或NOT IN
		if sess.stmt.where[i].operator == "IN" || sess.stmt.where[i].operator == "NOT IN" {
			if final == true {
				tmpValues := sliceToString(sess.stmt.where[i].value)
				value += strings.Join(tmpValues, ", ")
				if i == 0 {
					stmt.WriteString("(`")
					stmt.WriteString(sess.stmt.where[i].field)
					stmt.WriteString("` ")
					stmt.WriteString(sess.stmt.where[i].operator)
					stmt.WriteString(" (")
					stmt.WriteString(value)
					stmt.WriteString("))")
				} else {
					stmt.WriteString(" ")
					stmt.WriteString(sess.stmt.where[i].union)
					stmt.WriteString(" (`")
					stmt.WriteString(sess.stmt.where[i].field)
					stmt.WriteString("` ")
					stmt.WriteString(sess.stmt.where[i].operator)
					stmt.WriteString(" (")
					stmt.WriteString(value)
					stmt.WriteString("))")
				}
			} else {
				if i == 0 {
					stmt.WriteString("(`")
					stmt.WriteString(sess.stmt.where[i].field)
					stmt.WriteString("` ")
					stmt.WriteString(sess.stmt.where[i].operator)
					stmt.WriteString(" (?))")
				} else {
					stmt.WriteString(" ")
					stmt.WriteString(sess.stmt.where[i].union)
					stmt.WriteString(" (`")
					stmt.WriteString(sess.stmt.where[i].field)
					stmt.WriteString("` ")
					stmt.WriteString(sess.stmt.where[i].operator)
					stmt.WriteString(" (?))")
				}
				sess.stmt.resultValues = append(sess.stmt.resultValues, sess.stmt.where[i].value)
			}
		} else if sess.stmt.where[i].operator == "[!RAW!]" {
			if i == 0 {
				stmt.WriteString(sess.stmt.where[i].field)
			} else {
				stmt.WriteString(" ")
				stmt.WriteString(sess.stmt.where[i].union)
				stmt.WriteString(" ")
				stmt.WriteString(sess.stmt.where[i].field)
			}
		} else {
			if final == true {
				value = interfaceToString(sess.stmt.where[i].value)
				if i == 0 {
					stmt.WriteString("(`")
					stmt.WriteString(sess.stmt.where[i].field)
					stmt.WriteString("`")
					stmt.WriteString(sess.stmt.where[i].operator)
					stmt.WriteString(value)
					stmt.WriteString(")")
				} else {
					stmt.WriteString(" ")
					stmt.WriteString(sess.stmt.where[i].union)
					stmt.WriteString(" (`")
					stmt.WriteString(sess.stmt.where[i].field)
					stmt.WriteString("`")
					stmt.WriteString(sess.stmt.where[i].operator)
					stmt.WriteString(value)
					stmt.WriteString(")")
				}
			} else {
				if i == 0 {
					stmt.WriteString("(`")
					stmt.WriteString(sess.stmt.where[i].field)
					stmt.WriteString("`")
					stmt.WriteString(sess.stmt.where[i].operator)
					stmt.WriteString("?)")
				} else {
					stmt.WriteString(" ")
					stmt.WriteString(sess.stmt.where[i].union)
					stmt.WriteString(" (`")
					stmt.WriteString(sess.stmt.where[i].field)
					stmt.WriteString("`")
					stmt.WriteString(sess.stmt.where[i].operator)
					stmt.WriteString("?)")
				}
				sess.stmt.resultValues = append(sess.stmt.resultValues, sess.stmt.where[i].value)
			}
		}
	}
	return stmt.String()
}

//构建ORDER BY语句
func (sess *Session) buildOrderBy() string {
	len := len(sess.stmt.orders)
	if len == 0 {
		return ""
	}
	var stmt bytes.Buffer
	stmt.WriteString(" ORDER BY ")
	for i := 0; i < len; i++ {
		if i == 0 {
			stmt.WriteString("`")
			stmt.WriteString(sess.stmt.orders[i].field)
			stmt.WriteString("` ")
			stmt.WriteString(sess.stmt.orders[i].direction)
		} else {
			stmt.WriteString(", `")
			stmt.WriteString(sess.stmt.orders[i].field)
			stmt.WriteString("` ")
			stmt.WriteString(sess.stmt.orders[i].direction)
		}
	}

	return stmt.String()
}

//构建LIMIT语句
func (sess *Session) buildLimit() string {
	var stmt bytes.Buffer
	if sess.stmt.limit > 0 {
		stmt.WriteString(" LIMIT ")
		stmt.WriteString(strconv.Itoa(sess.stmt.limit))
	}
	return stmt.String()
}

//构建拼接OFFSET语句
func (sess *Session) buildOffset() string {
	var stmt bytes.Buffer
	if sess.stmt.offset > 0 {
		stmt.WriteString(" OFFSET ")
		stmt.WriteString(strconv.Itoa(sess.stmt.offset))
	}
	return stmt.String()
}
