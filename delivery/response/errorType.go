package response

type CustomError struct {
	ErrorCode int
	ErrorMsg  string
}

type CustomErrors struct {
	BusinessError CustomError
	ValidateError CustomError
	DataBaseError CustomError
	TokenError    CustomError
	ServeError    CustomError
}

/*
	code 规则制定
	1. 1-999 为系统错误
	2. 1000-1999 为业务错误
	3. 2000-2999 为验证错误
	4. 3000-3999 为数据库错误
	5. 4000-4999 为服务错误
*/

var Errors = CustomErrors{
	BusinessError: CustomError{1000, "业务错误"},
	TokenError:    CustomError{2000, "登录授权失效"},
	ValidateError: CustomError{2001, "验证错误"},
	DataBaseError: CustomError{3000, "其他数据与此关联"},
	ServeError:    CustomError{4000, "服务器错误"},
}
