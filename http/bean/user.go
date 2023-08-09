package bean

type UserLoginRo struct {
	// 登录名
	Login string `json:"login" validate:"required"`
	// md5(password)
	Password string `json:"password" validate:"required,md5"`
}
