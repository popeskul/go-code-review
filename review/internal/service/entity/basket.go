package entity

type Basket struct {
	Value                 int  `json:"value"`
	AppliedDiscount       int  `json:"applied_discount"`       // % or number
	ApplicationSuccessful bool `json:"application_successful"` // % or number
}
