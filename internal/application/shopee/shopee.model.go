package shopee

import "time"

type ShopeeAuthModel struct {
	ShopeeID     string    `bson:"shop_id"`
	AccessToken  string    `bson:"token"`
	RefreshToken string    `bson:"refresh_token"`
	ExpiredAt    time.Time `bson:"expired_at"`

	CreatedAt time.Time `bson:"created_at"`
	CreatedBy string    `bson:"created_by"`

	MoidifiedAt time.Time `bson:"modified_at"`
	ModifiedBy  string    `bson:"modified_by"`
}

