package mysqlib

// GetStmt 获得已经构建的SQL语句
func (sess *Session) GetStmt() string {
	return sess.stmt.resultString
}

// GetValues 获得SQL语句中?占位符对应的参数值
func (sess *Session) GetValues() []interface{} {
	return sess.stmt.resultValues
}
