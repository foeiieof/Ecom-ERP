package shopee

import "time"

type ShopeeAuthEntity struct {
	ShopID        string
  Code          string
	AccessToken   string
	RefreshToken  string
	ExpiredAt     time.Time

	CreatedAt     time.Time
	CreatedBy     string

	MoidifiedAt   time.Time
	ModifiedBy    string
}

