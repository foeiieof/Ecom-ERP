package demo

import (

	"go.mongodb.org/mongo-driver/v2/mongo"
)

//1.
// define interface for wrap up method
type DemoRepository interface {
  GetDemo(id string) (string,error)
}

//3. prepareing struct for using in method
type mongoRepo struct {
  db *mongo.Collection
}

//2. implement method or function

func (r *mongoRepo) GetDemo(id string) (string, error) {
  // if err := r.db.Ping(context.TODO(), readpref.Primary()); err != nil {
  //   return "", err
  // }
  // return "demo", nil
  // newShopeeModel := &model.ShopeeAuthModel{
  //   ShopeeID: id,
  // } 
  // r.db.InsertOne(context.TODO(), newShopeeModel)
  // if err := r.db.(context.TODO(), readpref.Primary()); err != nil {
  //   return "", nil
  // }
  return "shopee_auth", nil

}

//4. Create Constructor that return implement function
func NewDemoRepo(db *mongo.Collection) DemoRepository {
  return &mongoRepo{
    db: db,
  }
}
