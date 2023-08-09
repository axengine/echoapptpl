package jwt

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type Resp struct {
	ResCode int         `json:"resCode"`
	ResDesc string      `json:"resDesc"`
	Result  interface{} `json:"result"`
}

// JWTAuth 中间件，检查token
func JWTAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		authHeader := ctx.Request().Header.Get("Authorization")
		if authHeader == "" {
			return ctx.JSON(http.StatusOK, Resp{401, "Invalid Token", ""})
		}

		//按空格拆分
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			return ctx.JSON(http.StatusOK, Resp{401, "Invalid Token", ""})
		}

		//解析token包含的信息
		claims, err := ParseToken(parts[1])
		if err != nil {
			return ctx.JSON(http.StatusOK, Resp{401, "Invalid Token", ""})
		}

		// Path /admin need admin role
		if strings.HasPrefix(ctx.Request().RequestURI, "/admin") && claims.Role != 1 {
			return ctx.JSON(http.StatusOK, Resp{403, "Access Denied", ""})
		}

		// 将当前请求的claims信息保存到请求的上下文c上
		ctx.Set("claims", claims)
		return next(ctx) // 后续的处理函数可以用过ctx.Get("claims")来获取当前请求的用户信息
	}
}
