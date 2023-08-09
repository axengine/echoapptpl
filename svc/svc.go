package svc

import (
	"context"
	"echoapptpl/dbo"
	"echoapptpl/http/bean"
	"echoapptpl/middleware/jwt"
	"echoapptpl/model"
	"echoapptpl/types/errorx"
	"github.com/axengine/utils/hash"
	"gorm.io/gorm"
)

type Service struct {
	db *dbo.DBO
}

func New(db *dbo.DBO) *Service {
	return &Service{db: db}
}

func (svc *Service) UserLogin(ctx context.Context, ro *bean.UserLoginRo, ip string) (string, error) {
	// 使用密码登录
	var u model.User
	if err := svc.db.Where("login = ? ", ro.Login).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", errorx.NewError(errorx.CodeFailed, "用户名或密码错误")
		}
		return "", err
	}

	// 未审核用户不能登录
	if u.Unauthorized {
		return "", errorx.NewError(errorx.CodeFailed, "用户未审核")
	}

	// check password
	var loginErr error
	var token string
	hashedPswd := hash.Keccak256Base64([]byte(ro.Password + u.Salt))
	if hashedPswd != u.Password {
		loginErr = errorx.NewError(errorx.CodeFailed, "用户名或密码错误")
	}

	if loginErr == nil {
		token, loginErr = jwt.GenToken(u.ID, 0, u.Level)
	}

	return token, loginErr
}
