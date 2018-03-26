package mysqlib

import (
	"testing"

	"log"
	//_ "github.com/go-sql-driver/mysql"
)

// 构建器实例
var mysql *Instance

// 数据库连接实例
//var db *sql.DB

// User 定义模型
type User struct {
	tableName struct{} `sql:"user"`
	ID        int64    `sql:"id"`
	Username  string   `sql:"username"`
	Password  string   `sql:"password"`
}

// TestInit 初始化
func TestInit(t *testing.T) {
	// 实例化一个构建器对象
	mysql = New(&Options{
		TableNameField:    "tableName", //标记表名的字段名
		TagName:           "sql",       //标记字符
		DisableModelCache: false,       //禁用模型缓存
	})
}

// TestConnect 测试构建DSN语句
func TestConnect(t *testing.T) {
	// 构建DSN语句
	dsn, err := DSN(&DSNOptions{
		Addr:     "127.0.0.1:3306",
		User:     "root",
		Password: "123456",
		Database: "mytest",
	})
	if err != nil {
		t.Error(err.Error())
		return
	}

	t.Log("构建的DSN语句：", dsn)

	// 连接数据库
	//db, err = sql.Open("mysql", dsn)
	//if err != nil {
	//	t.Error(err.Error())
	//	return
	//}
}

// TestInsert 测试构建INSERT语句
func TestInsert(t *testing.T) {
	// 实例化一个模型并赋值，用于写入到数据库
	var user User
	user.ID = 1
	user.Username = "dxvgef"
	user.Password = "123456"
	// 使用user实例做为模型，构建INSERT语句
	// 执行Build()时会构建最终的SQL语句，入参值为false时构建带?占位符的模板语句
	// Build会返回一个会话实例sqlSess
	sqlSess, err := mysql.Insert(&user).Build(false)
	if err != nil {
		t.Error(err.Error())
		return
	}
	// 用会话实例的GetStmt方法获得SQL语句
	stmt := sqlSess.GetStmt()
	t.Log("构建的SQL语句：", stmt)

	// 用会话实例的GetValues方法得到模板语句中所有要传入的参数
	// 注意：只有在Build()的入参是false时，此方法才会有值返回
	values := sqlSess.GetValues()
	t.Log("执行SQL语句所需要的参数：")
	t.Log(values)

	// 将生成的SQL语句和参数传入到db.Exec()中执行数据库查询
	//result, err := db.Exec(stmt, values...)
	//if err != nil {
	//	t.Error(err.Error())
	//	return
	//}
	// 受影响的行数
	//count, err := result.RowsAffected()
	//if err != nil {
	//	t.Error(err.Error())
	//	return
	//}
	//t.Log("受影响的行数：", count)
}

// 测试构建UPDATE语句
func TestUpdate(t *testing.T) {
	var user User
	user.Username = "abc"

	sqlSess, err := mysql.Update(&user).
		// 使用Column()指定要更新哪些字段
		// 注意：出于安全考虑，如果不指定要更新哪些字段，Build()不会构建SQL语句
		Column("username").
		Where("id", "=", 1).
		Build(false)
	if err != nil {
		log.Println(err.Error())
		return
	}

	stmt := sqlSess.GetStmt()
	t.Log("构建的SQL语句：", stmt)

	values := sqlSess.GetValues()
	t.Log("执行SQL语句所需要的参数：")
	t.Log(values)

	//result, err := db.Exec(stmt, values...)
	//if err != nil {
	//	t.Log(err.Error())
	//	return
	//}
	//
	//count, err := result.RowsAffected()
	//if err != nil {
	//	t.Log(err.Error())
	//	return
	//}
	//t.Log("受影响的行数：", count)
}

// TestSelect 测试构建SELECT语句
func TestSelect(t *testing.T) {
	var user User
	sqlSess, err := mysql.Select(&user).
		Column("username").
		Where("id", "=", 1).
		Limit(1).
		Build(true)
	if err != nil {
		t.Error(err.Error())
		return
	}

	stmt := sqlSess.GetStmt()
	t.Log("构建的SQL语句：", stmt)

	values := sqlSess.GetValues()
	t.Log("执行SQL语句所需要的参数：")
	t.Log(values)

	//rows, err := db.Query(stmt, values...)
	//if err != nil {
	//	t.Error(err.Error())
	//	return
	//}
	//
	//// 使用会话的ScanModel方法将单条记录集赋值到模型
	//err = sqlSess.ScanModel(rows)
	//if err != nil {
	//	t.Error(err.Error())
	//	return
	//}
	//t.Log("以下为读取出来的数据：")
	//t.Log(user.Username)
}

// TestSelectMultiple 测试通过模型slice查询多条记录
func TestSelectMultiple(t *testing.T) {
	// 注意这里定义的是模型的slice
	var users []User
	sqlSess, err := mysql.Select(&users).
		Build(false)
	if err != nil {
		t.Error(err.Error())
		return
	}

	stmt := sqlSess.GetStmt()
	t.Log("构建的SQL语句：", stmt)

	values := sqlSess.GetValues()
	t.Log("执行SQL语句所需要的参数：")
	t.Log(values)

	//rows, err := db.Query(stmt, values...)
	//if err != nil {
	//	t.Error(err.Error())
	//	return
	//}
	//
	//// 使用会话的ScanModelSlice 将多条记录集赋值到模型slice
	//err = sqlSess.ScanModelSlice(rows)
	//if err != nil {
	//	t.Error(err.Error())
	//	return
	//}
	//
	////循环遍历模型的值
	//t.Log("以下为读取出来的数据：")
	//for _, v := range users {
	//	t.Log(v.ID, v.Username, v.Password)
	//}
}

// TestDelete 测试构建DELETE语句
func TestDelete(t *testing.T) {
	// 使用空实例做模型
	sqlSess, err := mysql.Delete(&User{}).
		Where("id", "=", 1).
		Build(false)
	if err != nil {
		t.Error(err.Error())
		return
	}
	stmt := sqlSess.GetStmt()
	t.Log("构建的SQL语句：", stmt)

	values := sqlSess.GetValues()
	t.Log("执行SQL语句所需要的参数：")
	t.Log(values)

	//result, err := db.Exec(stmt, values...)
	//if err != nil {
	//	t.Error(err.Error())
	//	return
	//}
	//
	//count, err := result.RowsAffected()
	//if err != nil {
	//	t.Error(err.Error())
	//	return
	//}
	//
	//t.Log("受影响的行数：", count)
}
