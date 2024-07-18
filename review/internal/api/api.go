package api

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"coupon_service/internal/service/entity"
	"github.com/gin-gonic/gin"
)

type Service interface {
	ApplyCoupon(basket entity.Basket, couponCode string) (*entity.Basket, error)
	CreateCoupon(discount int, code string, expiration int) (*entity.Coupon, error)
	GetCoupons(codes []string) ([]entity.Coupon, error)
}

type Config struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type API struct {
	srv *http.Server
	mux *gin.Engine
	svc Service
	cfg Config
}

func New(cfg Config, svc Service) *API {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	api := &API{
		mux: r,
		cfg: cfg,
		svc: svc,
	}

	api.srv = &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: api.mux,
	}

	return api
}

func (a *API) registerRoutes() {
	apiGroup := a.mux.Group("/api")
	apiGroup.POST("/apply", a.ApplyCoupon)
	apiGroup.POST("/create", a.CreateCoupon)
	apiGroup.GET("/coupons", a.GetCoupons)
}

func (a *API) Start() error {
	a.registerRoutes()
	return a.srv.ListenAndServe()
}

func (a *API) Close(ctx context.Context) error {
	if err := a.srv.Shutdown(ctx); err != nil {
		log.Println("Server forced to shutdown:", err)
		return err
	}
	log.Println("Server exiting")
	return nil
}
