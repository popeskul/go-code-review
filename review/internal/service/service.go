package service

import (
	serviceEntity "coupon_service/internal/service/entity"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

//go:generate mockery --name Repository --output ./mocks --filename repository.go

var (
	ErrInvalidDiscountValue  = fmt.Errorf("invalid discount value")
	ErrInvalidCouponCode     = fmt.Errorf("invalid coupon code")
	ErrInvalidMinBasketValue = fmt.Errorf("invalid minimum basket value")
	ErrNegativeBasketValue   = fmt.Errorf("tried to apply discount to negative value")
	ErrCouponNotFound        = fmt.Errorf("coupon not found")
)

type Repository interface {
	FindByCode(string) (*serviceEntity.Coupon, error)
	Save(serviceEntity.Coupon) error
}

type Service struct {
	repo Repository
}

func New(repo Repository) Service {
	return Service{
		repo: repo,
	}
}

func (s Service) ApplyCoupon(basket serviceEntity.Basket, code string) (*serviceEntity.Basket, error) {
	b := &basket
	coupon, err := s.repo.FindByCode(code)
	if err != nil {
		return nil, err
	}

	switch {
	case b.Value > 0:
		b.AppliedDiscount = coupon.Discount
		b.ApplicationSuccessful = true
	case b.Value == 0:
		return b, nil
	default:
		return nil, ErrNegativeBasketValue
	}

	return b, nil
}

func (s Service) CreateCoupon(discount int, code string, minBasketValue int) (*serviceEntity.Coupon, error) {
	if discount <= 0 {
		return nil, fmt.Errorf("%w: %d", ErrInvalidDiscountValue, discount)
	}
	if code == "" {
		return nil, ErrInvalidCouponCode
	}
	if minBasketValue < 0 {
		return nil, fmt.Errorf("%w: %d", ErrInvalidMinBasketValue, minBasketValue)
	}

	coupon := serviceEntity.Coupon{
		Discount:       discount,
		Code:           code,
		MinBasketValue: minBasketValue,
		ID:             uuid.New(),
	}

	if err := s.repo.Save(coupon); err != nil {
		return nil, fmt.Errorf("failed to save coupon: %w", err)
	}
	return &coupon, nil
}

func (s Service) GetCoupons(codes []string) ([]serviceEntity.Coupon, error) {
	coupons := make([]serviceEntity.Coupon, 0, len(codes))
	var eg errgroup.Group
	var mu sync.Mutex
	var notFoundErrors []string

	for idx, code := range codes {
		idx, code := idx, code
		eg.Go(func() error {
			coupon, err := s.repo.FindByCode(code)
			if err != nil {
				mu.Lock()
				if err.Error() == ErrCouponNotFound.Error() {
					notFoundErrors = append(notFoundErrors, fmt.Sprintf("code: %s, index: %d", code, idx))
				} else {
					mu.Unlock()
					return err
				}
				mu.Unlock()
				return nil
			}

			mu.Lock()
			if coupon != nil {
				coupons = append(coupons, *coupon)
			}
			mu.Unlock()
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	if len(notFoundErrors) > 0 {
		return coupons, fmt.Errorf("coupons not found: %v", notFoundErrors)
	}

	return coupons, nil
}
