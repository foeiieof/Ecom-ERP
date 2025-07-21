package shopee

import "time"

type ShopeeAuth struct {
	ShopID        string
	AccessToken   string
	RefreshToken  string
	ExpiredAt     time.Time

	CreatedAt     time.Time
	CreatedBy     string

	MoidifiedAt   time.Time
	ModifiedBy    string
}
