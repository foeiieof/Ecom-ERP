package shopee

import "time"

type ShopeeAuthModel struct {
  PartnerID    string    `bson:"partner_id"`
	ShopID       string    `bson:"shop_id"`
  Code         string    `bson:"code"`
	AccessToken  string    `bson:"access_token"`
	RefreshToken string    `bson:"refresh_token"`
	ExpiredAt    time.Time `bson:"expired_at"`

	CreatedAt time.Time `bson:"created_at"`
	CreatedBy string    `bson:"created_by"`

	MoidifiedAt time.Time `bson:"modified_at"`
	ModifiedBy  string    `bson:"modified_by"`
}

type ShopeeAuthRequestModel struct {
	PartnerID   string    `bson:"partner_id"`
	PartnerKey  string    `bson:"partner_key"`
	PartnerName string    `bson:"partner_name"`
  GeneratedUrl string   `bson:"generated_url"`

	CreatedAt time.Time   `bson:"created_at"`
	CreatedBy string      `bson:"created_by"`

	MoidifiedAt time.Time `bson:"modified_at"`
	ModifiedBy  string    `bson:"modified_by"`
}

type ShopeePartnerModel struct {  
  PartnerID   string    `bson:"partner_id"`
	PartnerKey  string    `bson:"partner_key"`
	PartnerName string    `bson:"partner_name"`

  CreatedAt time.Time   `bson:"created_at"`
	CreatedBy string      `bson:"created_by"`

	MoidifiedAt time.Time `bson:"modified_at"`
	ModifiedBy  string    `bson:"modified_by"`
}

// ---------------------------------------------- Demo Template ------------------------------------
// type DemoModel Struct {
//   ... 
//   CreatedAt time.Time   `bson:"created_at"`
// 	CreatedBy string      `bson:"created_by"`
// 	MoidifiedAt time.Time `bson:"modified_at"`
// 	ModifiedBy  string    `bson:"modified_by"`
// }
// ---------------------------------------------- End Demo Template --------------------------------
