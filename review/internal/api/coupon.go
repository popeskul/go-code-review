package api

import (
	apiEntity "coupon_service/internal/api/entity"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *API) ApplyCoupon(c *gin.Context) {
	apiReq := apiEntity.ApplicationRequest{}
	if err := c.ShouldBindJSON(&apiReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	basket, err := a.svc.ApplyCoupon(apiReq.Basket, apiReq.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	apiResp := apiEntity.ApplicationResponse{
		Basket: *basket,
	}
	c.JSON(http.StatusOK, apiResp)
}

func (a *API) GetCoupons(c *gin.Context) {
	apiReq := apiEntity.CouponRequest{}
	if err := c.ShouldBindJSON(&apiReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	coupons, err := a.svc.GetCoupons(apiReq.Codes)
	if err != nil {
		if coupons != nil {
			var apiCoupons []apiEntity.Coupon
			for _, coupon := range coupons {
				apiCoupons = append(apiCoupons, apiEntity.Coupon{
					ID:             coupon.ID,
					Code:           coupon.Code,
					Discount:       coupon.Discount,
					MinBasketValue: coupon.MinBasketValue,
				})
			}
			apiResp := apiEntity.CouponResponse{
				Coupons: apiCoupons,
			}
			c.JSON(http.StatusPartialContent, gin.H{"error": err.Error(), "coupons": apiResp.Coupons})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	var apiCoupons []apiEntity.Coupon
	for _, coupon := range coupons {
		apiCoupons = append(apiCoupons, apiEntity.Coupon{
			ID:             coupon.ID,
			Code:           coupon.Code,
			Discount:       coupon.Discount,
			MinBasketValue: coupon.MinBasketValue,
		})
	}

	apiResp := apiEntity.CouponResponse{
		Coupons: apiCoupons,
	}
	c.JSON(http.StatusOK, apiResp)
}

func (a *API) CreateCoupon(c *gin.Context) {
	apiReq := apiEntity.Coupon{}
	if err := c.ShouldBindJSON(&apiReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	_, err := a.svc.CreateCoupon(apiReq.Discount, apiReq.Code, apiReq.MinBasketValue)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}
