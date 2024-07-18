package entity

type CouponRequest struct {
	Codes []string `json:"codes"`
}

type CouponResponse struct {
	Coupons []Coupon `json:"coupons"`
}
