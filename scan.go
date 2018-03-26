package mysqlib

import (
	"database/sql"
	"reflect"
)

// ScanModelSlice 将到多条记录赋值到模型Slice
func (sess *Session) ScanModelSlice(rows *sql.Rows) (err error) {
	defer rows.Close()
	//遍历数据库返回的记录集
	for rows.Next() {
		//一行记录的载体
		row := make([]interface{}, len(sess.stmt.field))
		//根据模型的类型，动态创建一个结构体，用于存储一条记录
		newRow := reflect.New(sess.modelValue.rType).Elem()
		//遍历要输出的字段
		for i, sqlName := range sess.stmt.field {
			row[i] = newRow.FieldByName(sess.modelInfo.fields[sqlName.key].VarName).Addr().Interface()
		}
		//获取记录集
		err = rows.Scan(row...)
		if err != nil {
			if err.Error() == "sql: Rows are closed" {
				err = sql.ErrNoRows
			}
			return
		}
		//把newRow结构体append到模型中
		sess.modelValue.rValue.Set(reflect.Append(sess.modelValue.rValue, newRow))
	}

	return
}

// ScanModel 将单条记录赋值到模型
func (sess *Session) ScanModel(rows *sql.Rows) (err error) {
	defer rows.Close()
	//一行记录的载体
	row := make([]interface{}, len(sess.stmt.field))
	//遍历要输出的字段
	for i, sqlName := range sess.stmt.field {
		//将模型字段的内存地址赋值给oneRow记录的载体
		row[i] = sess.modelValue.rValue.FieldByName(sess.modelInfo.fields[sqlName.key].VarName).Addr().Interface()
	}

	//获取记录集
	rows.Next()
	err = rows.Scan(row...)
	if err != nil {
		if err.Error() == "sql: Rows are closed" {
			err = sql.ErrNoRows
		}
		return
	}

	return
}
