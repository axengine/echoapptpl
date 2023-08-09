package handler

import (
	"echoapptpl/http/bean"
	"echoapptpl/svc"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Handle struct {
	svc *svc.Service
}

func New(svc *svc.Service) *Handle {
	return &Handle{svc: svc}
}

// UserLoginHandle
// @Summary 用户登录
// @Description 用户登录
// @Tags USER
// @Accept json
// @Produce json
// @Param Request body bean.UserLoginRo true "request param"
// @Success 200 {object} bean.Resp "success"
// @Router /user/login [POST]
func (h *Handle) UserLoginHandle(c echo.Context) error {
	var ro bean.UserLoginRo
	if err := c.Bind(&ro); err != nil {
		return c.JSON(http.StatusOK, new(bean.Resp).FailMsg("invalid parameter"))
	}
	if err := c.Validate(ro); err != nil {
		return c.JSON(http.StatusOK, new(bean.Resp).FailMsg(err.Error()))
	}
	token, err := h.svc.UserLogin(c.Request().Context(), &ro, c.RealIP())
	if err != nil {
		return c.JSON(http.StatusOK, new(bean.Resp).FailErr(c, err))
	}

	return c.JSON(http.StatusOK, new(bean.Resp).Success(token))
}
