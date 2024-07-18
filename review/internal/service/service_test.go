package service

import (
	"coupon_service/internal/service/entity"
	"coupon_service/internal/service/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_ApplyCoupon(t *testing.T) {
	mockRepo := new(mocks.Repository)
	service := New(mockRepo)

	tests := []struct {
		name      string
		basket    entity.Basket
		code      string
		setupMock func()
		wantError error
		wantValue int
	}{
		{
			name: "apply valid coupon",
			basket: entity.Basket{
				Value: 100,
			},
			code: "VALIDCOUPON",
			setupMock: func() {
				mockRepo.On("FindByCode", "VALIDCOUPON").Return(&entity.Coupon{
					Code:     "VALIDCOUPON",
					Discount: 10,
				}, nil)
			},
			wantError: nil,
			wantValue: 10,
		},
		{
			name: "invalid coupon code",
			basket: entity.Basket{
				Value: 100,
			},
			code: "INVALIDCOUPON",
			setupMock: func() {
				mockRepo.On("FindByCode", "INVALIDCOUPON").Return(nil, ErrCouponNotFound)
			},
			wantError: ErrCouponNotFound,
			wantValue: 0,
		},
		{
			name: "negative basket value",
			basket: entity.Basket{
				Value: -10,
			},
			code: "VALIDCOUPON",
			setupMock: func() {
				mockRepo.On("FindByCode", "VALIDCOUPON").Return(&entity.Coupon{
					Code:     "VALIDCOUPON",
					Discount: 10,
				}, nil)
			},
			wantError: ErrNegativeBasketValue,
			wantValue: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.setupMock()
			gotBasket, err := service.ApplyCoupon(tt.basket, tt.code)
			if tt.wantError != nil {
				assert.EqualError(t, err, tt.wantError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantValue, gotBasket.AppliedDiscount)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_CreateCoupon(t *testing.T) {
	mockRepo := new(mocks.Repository)
	service := New(mockRepo)

	tests := []struct {
		name           string
		discount       int
		code           string
		minBasketValue int
		setupMock      func()
		wantError      error
	}{
		{
			name:           "valid coupon",
			discount:       10,
			code:           "NEWCOUPON",
			minBasketValue: 50,
			setupMock: func() {
				mockRepo.On("Save", mock.AnythingOfType("entity.Coupon")).Return(nil)
			},
			wantError: nil,
		},
		{
			name:           "invalid discount",
			discount:       0,
			code:           "INVALIDDISCOUNT",
			minBasketValue: 50,
			setupMock:      func() {},
			wantError:      ErrInvalidDiscountValue,
		},
		{
			name:           "empty code",
			discount:       10,
			code:           "",
			minBasketValue: 50,
			setupMock:      func() {},
			wantError:      ErrInvalidCouponCode,
		},
		{
			name:           "negative min basket value",
			discount:       10,
			code:           "NEGATIVEBASKET",
			minBasketValue: -10,
			setupMock:      func() {},
			wantError:      ErrInvalidMinBasketValue,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.setupMock()
			_, err := service.CreateCoupon(tt.discount, tt.code, tt.minBasketValue)
			if tt.wantError != nil {
				assert.ErrorIs(t, err, tt.wantError)
			} else {
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
