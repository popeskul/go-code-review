package entity

import "github.com/google/uuid"

type Coupon struct {
	ID             uuid.UUID `json:"id"`
	Code           string    `json:"code"`
	Discount       int       `json:"discount"`         // This is a percentage or a fixed amount
	MinBasketValue int       `json:"min_basket_value"` // Can this be negative?
}
