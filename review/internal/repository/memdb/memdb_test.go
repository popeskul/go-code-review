package memdb

import (
	"coupon_service/internal/service/entity"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepository_FindByCode(t *testing.T) {
	repo := New()

	repo.entries["VALIDCOUPON"] = entity.Coupon{
		Code:     "VALIDCOUPON",
		Discount: 10,
	}

	tests := []struct {
		name       string
		code       string
		wantError  error
		wantCoupon *entity.Coupon
	}{
		{
			name:      "find valid coupon",
			code:      "VALIDCOUPON",
			wantError: nil,
			wantCoupon: &entity.Coupon{
				Code:     "VALIDCOUPON",
				Discount: 10,
			},
		},
		{
			name:       "coupon not found",
			code:       "INVALIDCOUPON",
			wantError:  ErrCouponNotFound,
			wantCoupon: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCoupon, err := repo.FindByCode(tt.code)
			if tt.wantError != nil {
				assert.EqualError(t, err, tt.wantError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantCoupon, gotCoupon)
			}
		})
	}
}

func TestRepository_Save(t *testing.T) {
	repo := New()

	tests := []struct {
		name      string
		coupon    entity.Coupon
		wantError error
	}{
		{
			name: "save valid coupon",
			coupon: entity.Coupon{
				Code:     "NEWCOUPON",
				Discount: 20,
			},
			wantError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Save(tt.coupon)
			assert.NoError(t, err)
			savedCoupon, err := repo.FindByCode(tt.coupon.Code)
			assert.NoError(t, err)
			assert.Equal(t, &tt.coupon, savedCoupon)
		})
	}
}
