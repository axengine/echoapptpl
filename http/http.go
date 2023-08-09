package http

import (
	"context"
	"echoapptpl/handler"
	"echoapptpl/middleware/jwt"
	"github.com/axengine/utils/log"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/zap"
	"sync"

	_ "echoapptpl/docs"
)

type HttpServer struct {
	e    *echo.Echo
	h    *handler.Handle
	exit chan struct{}
}

func New(h *handler.Handle) *HttpServer {
	return &HttpServer{h: h,
		exit: make(chan struct{}, 1)}
}

type customValidator struct {
	once     sync.Once
	validate *validator.Validate
}

func (c *customValidator) Validate(i interface{}) error {
	c.lazyInit()
	return c.validate.Struct(i)
}

func (c *customValidator) lazyInit() {
	c.once.Do(func() {
		c.validate = validator.New()
	})
}

func (s *HttpServer) Start(ctx context.Context, listen string) {
	s.e = echo.New()

	// set CORS
	s.e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.DefaultCORSConfig))

	// init jwt middleware
	jwt.Init(viper.GetString("JWT_KEY"), "hudson")

	// enable docs
	s.e.GET("/docs/*any", echoSwagger.WrapHandler)

	s.e.Validator = &customValidator{}
	lora := s.e.Group("/user")
	{
		lora.POST("/login", s.h.UserLoginHandle)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				if err := s.e.Shutdown(ctx); err != nil {
					log.Logger.Error("Shutdown", zap.Error(err))
				}
				s.exit <- struct{}{}
			}
		}
	}()

	go func() {
		if err := s.e.Start(listen); err != nil {
			log.Logger.Error("Start", zap.Error(err))
		}
	}()
}

func (s *HttpServer) Stop(ctx context.Context) {
	<-s.exit
}
