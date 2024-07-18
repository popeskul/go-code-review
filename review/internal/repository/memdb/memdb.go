package memdb

import (
	"coupon_service/internal/service/entity"
	"errors"
)

type Repository struct {
	entries map[string]entity.Coupon
}

func New() *Repository {
	return &Repository{
		entries: make(map[string]entity.Coupon),
	}
}

var (
	ErrCouponNotFound = errors.New("coupon not found")
)

func (r *Repository) FindByCode(code string) (*entity.Coupon, error) {
	coupon, ok := r.entries[code]
	if !ok {
		return nil, ErrCouponNotFound
	}
	return &coupon, nil
}

func (r *Repository) Save(coupon entity.Coupon) error {
	r.entries[coupon.Code] = coupon
	return nil
}
