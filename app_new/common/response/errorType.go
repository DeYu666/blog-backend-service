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
}

var Errors = CustomErrors{
	TokenError: CustomError{66666, "登录授权失效"},

	DataBaseError: CustomError{10000, "其他数据与此关联"},
}
