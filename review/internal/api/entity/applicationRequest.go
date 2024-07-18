package entity

import "coupon_service/internal/service/entity"

type ApplicationRequest struct {
	Code   string        `json:"code"`
	Basket entity.Basket `json:"basket"`
}

type ApplicationResponse struct {
	Basket entity.Basket `json:"basket"`
}
