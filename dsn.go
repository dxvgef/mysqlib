package mysqlib

import "errors"

// DSNOptions 数据库连接参数
type DSNOptions struct {
	Addr     string
	User     string
	Password string
	Database string
	Charset  string
}

// DSN 构建MySQL的DSN语句
func DSN(opt *DSNOptions) (string, error) {
	if opt.Database == "" {
		return "", errors.New("`Database`参数值不能为空")
	}
	if opt.Addr == "" {
		opt.Addr = "127.0.0.1:3306"
	}
	if opt.User == "" {
		opt.Addr = "root"
	}
	if opt.Charset == "" {
		opt.Charset = "utf8mb4,utf8"
	}
	result := opt.User + ":" + opt.Password + "@tcp(" + opt.Addr + ")/" + opt.Database + "?charset=" + opt.Charset
	return result, nil
}
